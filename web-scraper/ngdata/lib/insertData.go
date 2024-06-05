package lib

import (
	"cmp"
	"database/sql"
	"encoding/json"
	"fmt"
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

<<<<<<< Updated upstream
	product.Gtin = primaryData.Data.Ean
=======
	// for algolia
	product.ObjectID = primaryData.Data.Ean

	product.Id = primaryData.Data.Ean
>>>>>>> Stashed changes
	product.Title = primaryData.Data.Title
	product.SubTitle = primaryData.Data.Subtitle
	product.Category = getCorrectCategoryName(primaryData.Data.Category)
	product.SubCategory = primaryData.Data.SubCategory
	product.OnSale = primaryData.Data.OnSale
	// lager hele url-en for bildelinker for ulike størrelser
	product.ImageLink = fmt.Sprintf("%s%s", "https://bilder.ngdata.no/", primaryData.Data.ImageLink)

	// lager et array av priser, å gjøre det på denne måten gjør det lettere når dataen skal sendes til database
	var prices []Price
	storeMap := map[string]bool{}
	// sjekker at prisen ikke er 0, om den er det er det ikke vits å sende til databasen
	for _, product := range productData {
		if _, exists := storeMap[product.Store]; exists {
			continue
		}
		storeMap[product.Store] = true
		prices = append(prices, Price{Store: product.Store, Price: product.Data.Price, OriginalPrice: product.Data.OriginalPrice, UnitPrice: product.Data.ComparePricePerUnit, Url: fmt.Sprintf("%s%s", product.BaseUrl, product.Data.Slug)})
	}

	// sorterer basert på pris, så det første elementet i arrayet vil være det billigste
	priceCmp := func(a, b Price) int {
		return cmp.Compare(a.Price, b.Price)
	}
	slices.SortFunc(prices, priceCmp)
	product.Prices = prices

	// innhold
	// vekt kombinerer vekten og typen (g, kg, osv.)
	product.Weight = fmt.Sprintf("%v%s", primaryData.Data.Weight, primaryData.Data.WeightMeasurementType)
	product.Description = primaryData.Data.Description
	product.Unit = primaryData.Data.Unit
	product.UnitType = primaryData.Data.CompareUnit
	product.Size = primaryData.Data.Size
	product.Vendor = primaryData.Data.Vendor
	product.Brand = primaryData.Data.Brand
	product.OriginCountry = primaryData.Data.OriginCountry
	product.Ingredients = primaryData.Data.Ingredients

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
	product.Allergens = strings.Join(allergens, ", ")
	product.MayContainTracesOf = strings.Join(mayContainTracesOf, ", ")

	// næringsinnhold
	nutritionalContentData := primaryData.Data.NutritionalContent

	// om det ikke er noe næringsinnhold
	if len(nutritionalContentData) == 0 {
		product.NutritionalContent = nil
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
		product.NutritionalContent = &nutritionalContent
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
	nutritionalContentJson, err := json.Marshal(product.NutritionalContent)
	if err != nil {
		return err
	}

	// lager json objekt med priser
	pricesJson, err := json.Marshal(product.Prices)
	if err != nil {
		return err
	}

	// legger til en rad i Products table i databasen. om en rad med samme id (gtin) allerede eksisterer, blir den replaced
	// her gjører bare queryen klart, uten dette blir goroutinene helt fked up og overlapper
	productsStmt, err := db.Prepare(`
		INSERT INTO products (id, title, subtitle, imagelink, category, subcategory, onsale, description, weight, origincountry, ingredients, vendor, brand, size, unit, unittype, allergens, mayContainTracesOf, nutritionalcontent, prices)
		VALUES ($1, $2, $3, $4, $5, $6 , $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20)
		ON CONFLICT (id)
		DO UPDATE SET
			title = EXCLUDED.title,
			subtitle = EXCLUDED.subtitle,
			imagelink = EXCLUDED.imagelink,
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
	_, err = productsStmt.Exec(product.Id, product.Title, product.SubTitle, product.ImageLink, product.Category, product.SubCategory, product.OnSale, product.Description, product.Weight, product.OriginCountry, product.Ingredients, product.Vendor, product.Brand, product.Size, product.Unit, product.UnitType, product.Allergens, product.MayContainTracesOf, nutritionalContentJson, pricesJson)
	if err != nil {
		return err
	}

	return nil
}
