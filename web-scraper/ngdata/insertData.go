package main

import (
	"cmp"
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
	"slices"
	"strings"
)

func getPrices(gtin string, jokerData ApiResponse, sparData ApiResponse) (Product, Product) {
	jokerProduct := Product{}
	sparProduct := Product{}

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
func formatData(menyData Product, jokerData Product, sparData Product, products *Produkter) {
	product := Produkt{}

	product.Gtin = menyData.Data.Ean
	product.Tittel = menyData.Data.Title
	product.Undertittel = menyData.Data.Subtitle
	product.Kategori = menyData.Data.Category
	product.Underkategori = menyData.Data.SubCategory
	// lager hele url-en for bildelink
	product.BildeLink = fmt.Sprintf("%s%s%s", "https://bilder.ngdata.no/", menyData.Data.ImageLink, "/medium.jpg")

	// lager et array av priser, å gjøre det på denne måten gjør det lettere når dataen skal sendes til database
	prices := Priser{}
	// sjekker at prisen ikke er 0, om den er det er det ikke vits å sende til databasen
	if menyData.Data.Price != 0 {
		prices.Priser = append(prices.Priser, Pris{Butikk: "meny", Pris: menyData.Data.Price, OriginalPris: menyData.Data.OriginalPrice, EnhetsPris: menyData.Data.CalcPricePerUnit, EnhetsType: menyData.Data.CalcUnit, Url: fmt.Sprintf("%s%s", "https://meny.no/varer", menyData.Data.Slug)})
	}
	if jokerData.Data.Price != 0 {
		prices.Priser = append(prices.Priser, Pris{Butikk: "joker", Pris: jokerData.Data.Price, OriginalPris: jokerData.Data.OriginalPrice, EnhetsPris: jokerData.Data.CalcPricePerUnit, EnhetsType: jokerData.Data.CalcUnit, Url: fmt.Sprintf("%s%s", "https://joker.no/nettbutikk/varer", jokerData.Data.Slug)})
	}
	if sparData.Data.Price != 0 {
		prices.Priser = append(prices.Priser, Pris{Butikk: "spar", Pris: sparData.Data.Price, OriginalPris: sparData.Data.OriginalPrice, EnhetsPris: sparData.Data.CalcPricePerUnit, EnhetsType: sparData.Data.CalcUnit, Url: fmt.Sprintf("%s%s", "https://spar.no/nettbutikk/varer", sparData.Data.Slug)})
	}

	// sorterer basert på pris, så det første elementet i arrayet vil være det billigste
	priceCmp := func(a, b Pris) int {
		return cmp.Compare(a.Pris, b.Pris)
	}
	slices.SortFunc(prices.Priser, priceCmp)
	product.Priser = prices

	// innhold
	// vekt kombinerer vekten og typen (g, kg, osv.)
	product.Innhold.Vekt = fmt.Sprintf("%v%s", menyData.Data.Weight, menyData.Data.WeightMeasurementType)
	product.Innhold.Beskrivelse = menyData.Data.Description
	product.Innhold.Enhet = menyData.Data.Unit
	product.Innhold.Størrelse = menyData.Data.Size
	product.Innhold.Leverandør = menyData.Data.Vendor
	product.Innhold.Opprinnelsesland = menyData.Data.OriginCountry

	// mapper over allergener array som vi fikk fra databasen, og legger til navnet i struct så det er av typen []string
	var allergens []string
	for i := range menyData.Data.Allergens {
		allergens = append(allergens, menyData.Data.Allergens[i].Name)
	}
	product.Innhold.Allergener = strings.Join(allergens, ", ")

	product.Innhold.Ingredienser = menyData.Data.Ingredients
	product.Innhold.KanInneholdeSporAv = menyData.Data.AllergyDeclaration

	// næringsinnhold
	nutritionalContent := Næringsinnhold{}
	nutritionalContentData := menyData.Data.NutritionalContent
	nutritionalContentType := reflect.TypeOf(nutritionalContent)
	nutritionalContentValue := reflect.ValueOf(&nutritionalContent).Elem()

	// legger til næringsinnhold data i fields med navn som matcher dataen fra api-en (reflect)
	for i := 0; i < len(nutritionalContentData) && i < nutritionalContentType.NumField(); i++ {
		field := nutritionalContentValue.Field(i)
		if field.CanSet() {
			field.SetString(fmt.Sprintf("%v%s", nutritionalContentData[i].Amount, nutritionalContentData[i].Unit))
		}
	}

	product.Innhold.Næringsinnhold = nutritionalContent

	products.Produkter = append(products.Produkter, product)
}

func insertData(product Produkt, db *sql.DB) error {
	fmt.Println("Legger inn data for:", product.Tittel)

	// gjør om næringsinnhold (type Næringsinnhold struct) til nutritionalContentJson
	nutritionalContentJson, err := json.Marshal(product.Innhold.Næringsinnhold)
	if err != nil {
		return err
	}

	// lager json objekt med priser
	pricesJson, err := json.Marshal(product.Priser.Priser)
	if err != nil {
		return err
	}

	// legger til en rad i Products table i databasen. om en rad med samme id (gtin) allerede eksisterer, blir den replaced
	// her gjører bare queryen klart, uten dette blir goroutinene helt fked up og overlapper
	productsStmt, err := db.Prepare(`
		INSERT INTO products (id, title, subtitle, imagelink, category, subcategory, description, weight, origincountry, ingredients, vendor, size, unit, allergens, allergydeclaration, nutritionalcontent, prices)
		VALUES ($1, $2, $3, $4, $5, $6 , $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
		ON CONFLICT (id)
		DO UPDATE SET
			title = EXCLUDED.title,
			subtitle = EXCLUDED.subtitle,
			imagelink = EXCLUDED.imagelink,
			category = EXCLUDED.category,
			subcategory = EXCLUDED.subcategory,
			description = EXCLUDED.description,
			weight = EXCLUDED.weight,
			origincountry = EXCLUDED.origincountry,
			ingredients = EXCLUDED.ingredients,
			vendor = EXCLUDED.vendor,
			size = EXCLUDED.size,
			unit = EXCLUDED.unit,
			allergens = EXCLUDED.allergens,
			allergydeclaration = EXCLUDED.allergydeclaration,
			nutritionalcontent = EXCLUDED.nutritionalcontent,
			prices = EXCLUDED.prices
		`)
	if err != nil {
		return err
	}
	defer productsStmt.Close()

	// queryen executes
	_, err = productsStmt.Exec(product.Gtin, product.Tittel, product.Undertittel, product.BildeLink, product.Kategori, product.Underkategori, product.Innhold.Beskrivelse, product.Innhold.Vekt, product.Innhold.Opprinnelsesland, product.Innhold.Ingredienser, product.Innhold.Leverandør, product.Innhold.Størrelse, product.Innhold.Enhet, product.Innhold.Allergener, product.Innhold.KanInneholdeSporAv, nutritionalContentJson, pricesJson)
	if err != nil {
		return err
	}

	return nil
}
