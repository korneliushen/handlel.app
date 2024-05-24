package main

import (
	"encoding/json"
	"fmt"
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
	product.BildeLink = fmt.Sprintf("%s%s%s", "https://bilder.ngdata.no/", menyData.Data.ImageLink, "/medium.jpg")

	product.Priser.Meny = menyData.Data.Price
	product.Priser.MenyOriginal = menyData.Data.OriginalPrice
	product.Priser.MenyEnhet = menyData.Data.KgPrice

	product.Priser.Joker = jokerData.Data.Price
	product.Priser.JokerOriginal = jokerData.Data.OriginalPrice
	product.Priser.JokerEnhet = jokerData.Data.KgPrice

	product.Priser.Spar = sparData.Data.Price
	product.Priser.SparOriginal = sparData.Data.OriginalPrice
	product.Priser.SparEnhet = jokerData.Data.KgPrice

	// innhold
	product.Innhold.Vekt = fmt.Sprintf("%v%s", menyData.Data.Weight, menyData.Data.WeightMeasurementType)
	product.Innhold.Beskrivelse = menyData.Data.Description
	product.Innhold.Enhet = menyData.Data.Unit
	product.Innhold.Størrelse = menyData.Data.Size
	product.Innhold.Leverandør = menyData.Data.Vendor
	product.Innhold.Opprinnelsesland = menyData.Data.OriginCountry

	var allergens []string
	for i := range menyData.Data.Allergens {
		allergens = append(allergens, menyData.Data.Allergens[i].Name)
	}

	product.Innhold.Allergener = strings.Join(allergens, ", ")
	product.Innhold.Ingredienser = menyData.Data.Ingredients
	product.Innhold.KanInneholdeSporAv = menyData.Data.AllergyDeclaration

	// TODO FOR MEG IMRGN: fikse næringsinnhold, må gjøre det om til et json object av typen til Næringsinnhold structen +
	// fikse prices field (tror jeg først må inserte i prices og så referere til den)

	// næringsinnhold
	nutritionalContent := Næringsinnhold{}
	product.Innhold.Næringsinnhold = nutritionalContent

	// bytt ut med faktiske data
	products.Produkter = append(products.Produkter, product)
}

func insertData(product Produkt) error {
	client := db()
	defer client.Close()

	// gjør om næringsinnhold (type Næringsinnhold struct) til json
	json, err := json.Marshal(product.Innhold.Næringsinnhold)

	rows, err := client.Query(`
		INSERT INTO 
		products (id, title, subtitle, imagelink, category, subcategory, description, weight, origincountry, ingredients, vendor, size, unit, allergens, allergydeclaration, nutritionalcontent) 
		VALUES 
		($1, $2, $3, $4, $5, $6 , $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)`,
		product.Gtin, product.Tittel, product.Undertittel, product.BildeLink, product.Kategori, product.Underkategori, product.Innhold.Beskrivelse, product.Innhold.Vekt, product.Innhold.Opprinnelsesland, product.Innhold.Ingredienser, product.Innhold.Leverandør, product.Innhold.Størrelse, product.Innhold.Enhet, product.Innhold.Allergener, product.Innhold.KanInneholdeSporAv, json,
	)
	if err != nil {
		return err
	}
	defer rows.Close()

	return nil
}
