package lib

import (
	"cmp"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"reflect"
	"slices"
	"strings"

	_ "github.com/lib/pq"
)

func db() *sql.DB {
	connStr := os.Getenv("NEON_SECRET")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	// limiter open og idle connections for å ikke med et uhell overloade databasen
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return db
}

func getPrices(gtin string, jokerData ApiResponse, sparData ApiResponse) (ApiProduct, ApiProduct) {
	jokerProduct := ApiProduct{}
	sparProduct := ApiProduct{}

	// finner produkt fra joker med samme gtin
	for l := range jokerData.Hits.Products {
		if jokerData.Hits.Products[l].Data.Ean == gtin {
			jokerProduct = jokerData.Hits.Products[l]
		}
	}

	// finner produkt fra spar med samme gtin
	for l := range sparData.Hits.Products {
		if sparData.Hits.Products[l].Data.Ean == gtin {
			sparProduct = sparData.Hits.Products[l]
		}
	}

	return jokerProduct, sparProduct
}

// lager instanser av egne structs med dataen fra fetchProducts
func formatData(menyData ApiProduct, jokerData ApiProduct, sparData ApiProduct, products *Products) {
	product := Product{}

	product.Gtin = menyData.Data.Ean
	product.Title = menyData.Data.Title
	product.SubTitle = menyData.Data.Subtitle
	product.Category = menyData.Data.Category
	product.SubCategory = menyData.Data.SubCategory
	product.OnSale = menyData.Data.OnSale
	// lager hele url-en for bildelinker for ulike størrelser
	product.Images.ImageLinkXSmall = fmt.Sprintf("%s%s%s", "https://bilder.ngdata.no/", menyData.Data.ImageLink, "/xsmall.jpg")
	product.Images.ImageLinkSmall = fmt.Sprintf("%s%s%s", "https://bilder.ngdata.no/", menyData.Data.ImageLink, "/small.jpg")
	product.Images.ImageLinkMedium = fmt.Sprintf("%s%s%s", "https://bilder.ngdata.no/", menyData.Data.ImageLink, "/medium.jpg")
	product.Images.ImageLinkLarge = fmt.Sprintf("%s%s%s", "https://bilder.ngdata.no/", menyData.Data.ImageLink, "/large.jpg")
	product.Images.ImageLinkXLarge = fmt.Sprintf("%s%s%s", "https://bilder.ngdata.no/", menyData.Data.ImageLink, "/xlarge.jpg")

	// lager et array av priser, å gjøre det på denne måten gjør det lettere når dataen skal sendes til database
	prices := Prices{}
	// sjekker at prisen ikke er 0, om den er det er det ikke vits å sende til databasen
	if menyData.Data.Price != 0 {
		prices.Prices = append(prices.Prices, Price{Store: "meny", Price: math.Round(menyData.Data.Price), OriginalPrice: math.Round(menyData.Data.OriginalPrice), UnitPrice: math.Round(menyData.Data.ComparePricePerUnit), Url: fmt.Sprintf("%s%s", "https://meny.no/varer", menyData.Data.Slug)})
	}
	if jokerData.Data.Price != 0 {
		prices.Prices = append(prices.Prices, Price{Store: "joker", Price: math.Round(jokerData.Data.Price), OriginalPrice: math.Round(jokerData.Data.OriginalPrice), UnitPrice: math.Round(jokerData.Data.ComparePricePerUnit), Url: fmt.Sprintf("%s%s", "https://joker.no/nettbutikk/varer", jokerData.Data.Slug)})
	}
	if sparData.Data.Price != 0 {
		prices.Prices = append(prices.Prices, Price{Store: "spar", Price: math.Round(sparData.Data.Price), OriginalPrice: math.Round(sparData.Data.OriginalPrice), UnitPrice: math.Round(sparData.Data.ComparePricePerUnit), Url: fmt.Sprintf("%s%s", "https://spar.no/nettbutikk/varer", sparData.Data.Slug)})
	}

	// sorterer basert på pris, så det første elementet i arrayet vil være det billigste
	priceCmp := func(a, b Price) int {
		return cmp.Compare(a.Price, b.Price)
	}
	slices.SortFunc(prices.Prices, priceCmp)
	product.Prices = prices

	// innhold
	// vekt kombinerer vekten og typen (g, kg, osv.)
	product.Content.Weight = fmt.Sprintf("%v%s", menyData.Data.Weight, menyData.Data.WeightMeasurementType)
	product.Content.Description = menyData.Data.Description
	product.Content.Unit = menyData.Data.Unit
	product.Content.UnitType = menyData.Data.CompareUnit
	product.Content.Size = menyData.Data.Size
	product.Content.Vendor = menyData.Data.Vendor
	product.Content.Brand = menyData.Data.Brand
	product.Content.OriginCountry = menyData.Data.OriginCountry
	product.Content.Ingredients = menyData.Data.Ingredients

	// mapper over allergener array som vi fikk fra databasen
	// i databasen så bestemmer koden hva itemet i arrayen betyr for produktet
	// om koden er JA, blir det lagt til i allergens, om det er kan blir det lagt til i mayContainTracesOf
	var allergens []string
	var mayContainTracesOf []string
	for _, allergen := range menyData.Data.Allergens {
		if allergen.Code == "JA" {
			allergens = append(allergens, allergen.Name)
		} else if allergen.Code == "KAN" {
			mayContainTracesOf = append(mayContainTracesOf, allergen.Name)
		}
	}
	product.Content.Allergens = strings.Join(allergens, ", ")
	product.Content.MayContainTracesOf = strings.Join(mayContainTracesOf, ", ")

	// næringsinnhold
	nutritionalContentData := menyData.Data.NutritionalContent

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

	products.Products = append(products.Products, product)
}

func insertData(product Product, db *sql.DB) error {
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
