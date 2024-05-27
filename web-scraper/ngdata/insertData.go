package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
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

	// TODO: remake dette for å gjøre at man bare kan compare når requesten sendes og slipper dette, det funker tho
	product.Gtin = menyData.Data.Ean
	product.Tittel = menyData.Data.Title
	product.Undertittel = menyData.Data.Subtitle
	product.Kategori = menyData.Data.Category
	product.Underkategori = menyData.Data.SubCategory
	// lager hele url-en for bildelink
	product.BildeLink = fmt.Sprintf("%s%s%s", "https://bilder.ngdata.no/", menyData.Data.ImageLink, "/medium.jpg")

	// lager et array av priser, å gjøre det på denne måten gjør det lettere når dataen skal sendes til database
	prices := Priser{}
	prices.Priser = append(prices.Priser, Pris{Gtin: menyData.Data.Ean, Butikk: "meny", Pris: menyData.Data.Price, OriginalPris: menyData.Data.OriginalPrice, EnhetsPris: menyData.Data.CalcPricePerUnit, EnhetsType: menyData.Data.CalcUnit, Url: fmt.Sprintf("%s%s", "https://meny.no/varer", menyData.Data.Slug)})
	prices.Priser = append(prices.Priser, Pris{Gtin: jokerData.Data.Ean, Butikk: "joker", Pris: jokerData.Data.Price, OriginalPris: jokerData.Data.OriginalPrice, EnhetsPris: jokerData.Data.CalcPricePerUnit, EnhetsType: jokerData.Data.CalcUnit, Url: fmt.Sprintf("%s%s", "https://joker.no/nettbutikk/varer", jokerData.Data.Slug)})
	prices.Priser = append(prices.Priser, Pris{Gtin: sparData.Data.Ean, Butikk: "spar", Pris: sparData.Data.Price, OriginalPris: sparData.Data.OriginalPrice, EnhetsPris: sparData.Data.CalcPricePerUnit, EnhetsType: sparData.Data.CalcUnit, Url: fmt.Sprintf("%s%s", "https://spar.no/nettbutikk/varer", sparData.Data.Slug)})
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

	// gjør om næringsinnhold (type Næringsinnhold struct) til json
	json, err := json.Marshal(product.Innhold.Næringsinnhold)

	// legger til en rad i Products table i databasen. om en rad med samme id (gtin) allerede eksisterer, blir den replaced
	productsStmt, err := db.Prepare(`
		INSERT INTO products (id, title, subtitle, imagelink, category, subcategory, description, weight, origincountry, ingredients, vendor, size, unit, allergens, allergydeclaration, nutritionalcontent) 
		VALUES ($1, $2, $3, $4, $5, $6 , $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
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
			nutritionalcontent = EXCLUDED.nutritionalcontent
		`)
	if err != nil {
		return err
	}
	defer productsStmt.Close()

	_, err = productsStmt.Exec(product.Gtin, product.Tittel, product.Undertittel, product.BildeLink, product.Kategori, product.Underkategori, product.Innhold.Beskrivelse, product.Innhold.Vekt, product.Innhold.Opprinnelsesland, product.Innhold.Ingredienser, product.Innhold.Leverandør, product.Innhold.Størrelse, product.Innhold.Enhet, product.Innhold.Allergener, product.Innhold.KanInneholdeSporAv, json)
	if err != nil {
		return err
	}

	pricesStmt, err := db.Prepare(`
		INSERT INTO prices (id, gtin, store, price, priceoriginal, priceunit, unittype, url, product_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (id)
		DO UPDATE SET
			gtin = EXCLUDED.gtin,
			store = EXCLUDED.store,
			price = EXCLUDED.price,
			priceoriginal = EXCLUDED.priceoriginal,
			priceunit = EXCLUDED.priceunit,
			unittype = EXCLUDED.unittype,
			url = EXCLUDED.url,
			product_id = EXCLUDED.product_id
		`,
	)
	if err != nil {
		return err
	}
	defer pricesStmt.Close()

	// om prisen ikke er 0 (da er det ikke noe vits å putte inn en row), legger til eller replacer pris basert på id
	// id er bygd opp av gtin og butikk, noe som vil være unikt
	for i := range product.Priser.Priser {
		price := product.Priser.Priser[i]

		if price.Pris == 0 {
			continue
		}

		id := fmt.Sprintf("%s/%s", product.Gtin, price.Butikk)

		pricesStmt.Exec(id, product.Gtin, price.Butikk, price.Pris, price.OriginalPris, price.EnhetsPris, price.EnhetsType, price.Url, product.Gtin)
	}

	return nil
}
