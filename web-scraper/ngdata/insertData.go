package main

import "fmt"

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
func insertData(menyData Product, jokerData Product, sparData Product, products *Produkter) {
	product := Produkt{}

	// TODO: remake dette for å gjøre at man bare kan compare når requesten sendes og slipper dette, det funker tho
	product.Gtin = menyData.Data.Ean
	product.Tittel = menyData.Data.Title
	product.Undertittel = menyData.Data.Subtitle
	product.Kategori = menyData.Data.Category
	product.Underkategori = menyData.Data.SubCategory
	product.BildeLink = fmt.Sprintf("%s%s%s", "https://bilder.ngdata.no/", menyData.Data.ImageLink, "/medium.jpg")

	// TODO: priser (må endre på dette systemet senere for å få med alle priser)
	product.Priser.Meny = menyData.Data.Price
	product.Priser.MenyOriginal = menyData.Data.OriginalPrice
	product.Priser.Joker = jokerData.Data.Price
	product.Priser.JokerOriginal = jokerData.Data.OriginalPrice
	product.Priser.Spar = sparData.Data.Price
	product.Priser.SparOriginal = sparData.Data.OriginalPrice

	// innhold
	product.Innhold.Vekt = fmt.Sprintf("%v%s", menyData.Data.Weight, menyData.Data.WeightMeasurementType)
	product.Innhold.Beskrivelse = menyData.Data.Description
	product.Innhold.Enhet = menyData.Data.Unit
	product.Innhold.Størrelse = menyData.Data.Size
	product.Innhold.Leverandør = menyData.Data.Vendor
	product.Innhold.Opprinnelsesland = menyData.Data.OriginCountry
	product.Innhold.Allergener = menyData.Data.Allergens
	product.Innhold.Ingredienser = menyData.Data.Ingredients
	product.Innhold.KanInneholdeSporAv = menyData.Data.AllergyDeclaration

	// næringsinnhold
	product.Innhold.Næringsinnhold = menyData.Data.NutritionalContent

	// bytt ut med faktiske data
	products.Produkter = append(products.Produkter, product)
}
