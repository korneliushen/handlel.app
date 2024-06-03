package lib

import (
	"cmp"
	"database/sql"
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"slices"
	"strings"
	"sync"

	_ "github.com/lib/pq"
)

// temporary løsning til jeg finner ut hvordan vi kan automatisere det
// lager et array med kategorinavn som er forskjellige, som så blir brukt etterpå for å gjøre alt i samme kategori til samme kategori navn
// MENY, JOKER, SPAR
var categoryNames = [...][3]string{
	// ting som basically kan automatiseres
	{"Frukt & grønt", "Frukt/Grønt", "Frukt og grønt"},
	{"Fisk & skalldyr", "Fisk/Skalldyr", "Fisk og skalldyr"},
	{"Personlige artikler", "Personlige Artikler"},

	// ting som er mobbing (må kjøre noe contains greier ig)
	{"Bakeri", "Bakerivarer", "Brød og bakervarer"},
	{"Meieri & egg", "Meieriprodukter", "Meieri og egg"},
	{"Pålegg & frokost", "Frokost/Pålegg", "Pålegg og frokost"},
	{"Drikke", "Drikkevarer"},
	{"Dessert og iskrem", "Dessert"},
	{"Bakevarer og kjeks", "Kaker/Bakevarer", "Bakeartikler og kjeks"},
	{"Barneprodukter", "Barn"},
	{"Hus & hjem", "Hus/Hjem Artikler", "Hus og hjem"},
	{"Snacks & godteri", "Snacks/Godteri", "Snacks og godteri"},
	{"Dyr", "Dyreprodukter"},
}

// bruker slicen av slices over til å gjøre alt til samme kategori navn
func getCorrectCategoryName(category string) string {
	for i := range categoryNames {
		for j := range categoryNames[i] {
			if categoryNames[i][j] == category {
				return categoryNames[i][0]
			}
		}
	}
	return category
}

// lager instanser av egne structs med dataen fra fetchProducts
func formatData(productData []ApiProduct, products *[]Product) {
	product := Product{}

	primaryData := productData[0]

	product.Gtin = primaryData.Data.Ean
	product.Title = primaryData.Data.Title
	product.SubTitle = primaryData.Data.Subtitle
	product.Category = getCorrectCategoryName(primaryData.Data.Category)
	product.SubCategory = primaryData.Data.SubCategory
	product.OnSale = primaryData.Data.OnSale
	// lager hele url-en for bildelinker for ulike størrelser
	product.Images.ImageLinkXSmall = fmt.Sprintf("%s%s%s", "https://bilder.ngdata.no/", primaryData.Data.ImageLink, "/xsmall.jpg")
	product.Images.ImageLinkSmall = fmt.Sprintf("%s%s%s", "https://bilder.ngdata.no/", primaryData.Data.ImageLink, "/small.jpg")
	product.Images.ImageLinkMedium = fmt.Sprintf("%s%s%s", "https://bilder.ngdata.no/", primaryData.Data.ImageLink, "/medium.jpg")
	product.Images.ImageLinkLarge = fmt.Sprintf("%s%s%s", "https://bilder.ngdata.no/", primaryData.Data.ImageLink, "/large.jpg")
	product.Images.ImageLinkXLarge = fmt.Sprintf("%s%s%s", "https://bilder.ngdata.no/", primaryData.Data.ImageLink, "/xlarge.jpg")

	// lager et array av priser, å gjøre det på denne måten gjør det lettere når dataen skal sendes til database
	prices := Prices{}
	storeMap := map[string]bool{}
	// sjekker at prisen ikke er 0, om den er det er det ikke vits å sende til databasen
	for _, product := range productData {
		if _, exists := storeMap[product.Store]; exists {
			continue
		}
		storeMap[product.Store] = true
		prices.Prices = append(prices.Prices, Price{Store: product.Store, Price: math.Round(product.Data.Price), OriginalPrice: math.Round(product.Data.OriginalPrice), UnitPrice: math.Round(product.Data.ComparePricePerUnit), Url: fmt.Sprintf("%s%s", product.BaseUrl, product.Data.Slug)})
	}

	// sorterer basert på pris, så det første elementet i arrayet vil være det billigste
	priceCmp := func(a, b Price) int {
		return cmp.Compare(a.Price, b.Price)
	}
	slices.SortFunc(prices.Prices, priceCmp)
	product.Prices = prices

	// innhold
	// vekt kombinerer vekten og typen (g, kg, osv.)
	product.Content.Weight = fmt.Sprintf("%v%s", primaryData.Data.Weight, primaryData.Data.WeightMeasurementType)
	product.Content.Description = primaryData.Data.Description
	product.Content.Unit = primaryData.Data.Unit
	product.Content.UnitType = primaryData.Data.CompareUnit
	product.Content.Size = primaryData.Data.Size
	product.Content.Vendor = primaryData.Data.Vendor
	product.Content.Brand = primaryData.Data.Brand
	product.Content.OriginCountry = primaryData.Data.OriginCountry
	product.Content.Ingredients = primaryData.Data.Ingredients

	// mapper over allergener array som vi fikk fra databasen
	// i databasen så bestemmer koden hva itemet i arrayen betyr for produktet
	// om koden er JA, blir det lagt til i allergens, om det er kan blir det lagt til i mayContainTracesOf
	var allergens []string
	var mayContainTracesOf []string
	for _, allergen := range primaryData.Data.Allergens {
		if allergen.Code == "JA" {
			allergens = append(allergens, allergen.Name)
		} else if allergen.Code == "KAN" {
			mayContainTracesOf = append(mayContainTracesOf, allergen.Name)
		}
	}
	product.Content.Allergens = strings.Join(allergens, ", ")
	product.Content.MayContainTracesOf = strings.Join(mayContainTracesOf, ", ")

	// næringsinnhold
	nutritionalContentData := primaryData.Data.NutritionalContent

	// om det ikke er noe næringsinnhold
	if len(nutritionalContentData) == 0 {
		product.Content.NutritionalContent = nil
	} else {
		nutritionalContent := NutritionalContent{}
		v := reflect.ValueOf(&nutritionalContent).Elem()
		// legger til næringsinnhold data i fields med navn som matcher dataen fra api-en (reflect)
		for i := range len(nutritionalContentData) {
			field := v.FieldByName(nutritionalContentData[i].Name)
			if field.CanSet() {
				field.SetString(fmt.Sprintf("%v%s", nutritionalContentData[i].Amount, nutritionalContentData[i].Unit))
			}
		}
		product.Content.NutritionalContent = &nutritionalContent
	}

	*products = append(*products, product)
}

func insertData(products *[]Product) {
	db := db()
	defer db.Close()

	// lager en waitgroup, som venter på goroutines for å bli ferdig før den starter en ny
	var wg sync.WaitGroup
	// limiter hvor mange go routines som kan kjøre om om gangen
	sem := make(chan struct{}, 4)

	for i := range *products {
		// legger til et item i wait groupen
		wg.Add(1)
		sem <- struct{}{}

		go func(product Product) {
			// når funksjonen er ferdig, blir waitgroup instansen ferdig + sem (det som keeper track av hvor mange ting som kan kjøre om gangen) blir oppdatert
			defer wg.Done()
			defer func() { <-sem }()

			// legger til data i databasen
			if err := query(product, db); err != nil {
				fmt.Printf("Error inserting data for %s: %v", product.Title, err)
			}
		}((*products)[i])
	}

	wg.Wait()
}

func query(product Product, db *sql.DB) error {
	fmt.Println("Legger inn data for:", product.Title)

	// gjør om næringsinnhold (type Næringsinnhold struct) til nutritionalContentJson
	nutritionalContentJson, err := json.Marshal(product.Content.NutritionalContent)
	if err != nil {
		return err
	}

	// lager json objekt med priser
	pricesJson, err := json.Marshal(product.Prices.Prices)
	if err != nil {
		return err
	}

	// legger til en rad i Products table i databasen. om en rad med samme id (gtin) allerede eksisterer, blir den replaced
	// her gjører bare queryen klart, uten dette blir goroutinene helt fked up og overlapper
	productsStmt, err := db.Prepare(`
		INSERT INTO products (id, title, subtitle, imagelinkxsmall, imagelinksmall, imagelinkmedium, imagelinklarge, imagelinkxlarge, category, subcategory, onsale, description, weight, origincountry, ingredients, vendor, brand, size, unit, unittype, allergens, mayContainTracesOf, nutritionalcontent, prices)
		VALUES ($1, $2, $3, $4, $5, $6 , $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24)
		ON CONFLICT (id)
		DO UPDATE SET
			title = EXCLUDED.title,
			subtitle = EXCLUDED.subtitle,
			imagelinkxsmall = EXCLUDED.imagelinkxsmall,
			imagelinksmall = EXCLUDED.imagelinksmall,
			imagelinkmedium = EXCLUDED.imagelinkmedium,
			imagelinklarge = EXCLUDED.imagelinklarge,
			imagelinkxlarge = EXCLUDED.imagelinkxlarge,
			category = EXCLUDED.category,
			subcategory = EXCLUDED.subcategory,
			onsale = EXCLUDED.onsale,
			description = EXCLUDED.description,
			weight = EXCLUDED.weight,
			origincountry = EXCLUDED.origincountry,
			ingredients = EXCLUDED.ingredients,
			vendor = EXCLUDED.vendor,
			brand = EXCLUDED.brand,
			size = EXCLUDED.size,
			unit = EXCLUDED.unit,
			unittype = EXCLUDED.unittype,
			allergens = EXCLUDED.allergens,
			mayContainTracesOf = EXCLUDED.mayContainTracesOf,
			nutritionalcontent = EXCLUDED.nutritionalcontent,
			prices = EXCLUDED.prices
		`)
	if err != nil {
		return err
	}
	defer productsStmt.Close()

	// queryen executes
	_, err = productsStmt.Exec(product.Gtin, product.Title, product.SubTitle, product.Images.ImageLinkXSmall, product.Images.ImageLinkSmall, product.Images.ImageLinkMedium, product.Images.ImageLinkLarge, product.Images.ImageLinkXLarge, product.Category, product.SubCategory, product.OnSale, product.Content.Description, product.Content.Weight, product.Content.OriginCountry, product.Content.Ingredients, product.Content.Vendor, product.Content.Brand, product.Content.Size, product.Content.Unit, product.Content.UnitType, product.Content.Allergens, product.Content.MayContainTracesOf, nutritionalContentJson, pricesJson)
	if err != nil {
		return err
	}

	return nil
}
