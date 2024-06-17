package model

import "github.com/korneliushen/handlel.app/scraper/lib"

// produkter som skal bli til json data/lagt inn i database
type Products []Product

type Product struct {
	ObjectID           string              `json:"objectID"`
	Id                 string              `json:"id"`
	Title              string              `json:"title"`
	SubTitle           string              `json:"subtitle"`
	Category           string              `json:"category"`
	SubCategory        string              `json:"subcategory"`
	Prices             []Price             `json:"prices"`
	OnSale             bool                `json:"onsale"`
	ImageLink          string              `json:"imagelink"`
	Description        string              `json:"description"`
	Duration           string              `json:"duration"`
	Unit               string              `json:"unit"`
	UnitType           string              `json:"unittype"`
	Size               string              `json:"size"`
	Vendor             string              `json:"vendor"`
	Brand              string              `json:"brand"`
	Ingredients        string              `json:"ingredients"`
	Storage            string              `json:"storage"`
	OriginCountry      string              `json:"origincountry"`
	Features           string              `json:"features"`
	Allergens          string              `json:"allergens"`
	MayContainTracesOf string              `json:"maycontaintracesof"`
	Weight             string              `json:"weight"`
  Notes              string              `json:"notes"`
	NutritionalContent *NutritionalContent `json:"nutritionalcontent"`
}

type Price struct {
	Store         string  `json:"store"`
	Price         float64 `json:"price"`
	OriginalPrice float64 `json:"originalprice"`
	UnitPrice     float64 `json:"unitprice"`
	Url           string  `json:"url"`
}

// fields er på norsk her for å kunne matche de med det som kommer fra databasen med reflect
type NutritionalContent struct {
	Energi          string `json:"energy"`
	Kalorier        string `json:"calories"`
	Natrium         string `json:"sodium"`
	Fett            string `json:"fat"`
	MettetFett      string `json:"saturatedfat"`
	EnumettetFett   string `json:"monounsaturatedfat"`
	FlerumettetFett string `json:"polyunsaturatedfat"`
	Karbohydrater   string `json:"carbohydrates"`
	Sukkerarter     string `json:"sugars"`
	Stivelse        string `json:"starch"`
	Kostfiber       string `json:"dietaryfiber"`
	Protein         string `json:"protein"`
	Salt            string `json:"salt"`
}

func (products *Products) Format(apiProducts ApiProducts) {
	// lagrer alle produkter som allerede har blitt sjekket i et array
	// da blir det ikke duplicates og vi kan returnere tidlig om produktet
	// allerede er ferdig
	var checkedGtins []string

	// mapper over alle produkter vi har fått fra databasen og formatterer
	// dataen i egne structs
	for _, firstProduct := range apiProducts {
		gtin := firstProduct.Data.Ean

		// om produktet allerede er sjekket, skip dette produktet
		if lib.IsIn(gtin, checkedGtins) {
			continue
		}
		checkedGtins = append(checkedGtins, gtin)

		// finner andre produkter med samme gtin og legger til i et array
		sameProduct := []ApiProduct{firstProduct}

		for _, secondProduct := range apiProducts {
			if gtin == secondProduct.Data.Ean &&
				firstProduct.Store != secondProduct.Store {
				// legger til produktet i sameProduct array, alle produkter i dette
				// arrayet sjekkes nå priser legges inn
				sameProduct = append(sameProduct, secondProduct)
			}
		}

		// formaterer dataen til alle produkter med samme gtin
		firstProduct.FormatData(sameProduct, products)
	}
}
