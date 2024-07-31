package model

import (
	"cmp"
	"slices"

	"github.com/korneliushen/handlel.app/scraper/lib"
)

// produkter som skal bli til json data/lagt inn i database
type Products []Product

type Product struct {
  Store              string              `json:"store"`
  BaseUrl            string              `json:"baseurl"`
  Slug               string              `json:"slug"`
	Id                 string              `json:"id"`
  ObjectID           string              `json:"objectID"`
	Title              string              `json:"title"`
	SubTitle           string              `json:"subtitle"`
	Category           string              `json:"category"`
	SubCategory        string              `json:"subcategory"`
	Prices             Price               `json:"prices"`
	Images             Images              `json:"images"`
	OnSale             bool                `json:"onsale"`
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
  Today   []PriceToday `json:"today"`
  History []PriceHistory `json:"history"`
}

type PriceToday struct {
	Store         string  `json:"store"`
	Price         float64 `json:"price"`
	OriginalPrice float64 `json:"originalprice"`
	UnitPrice     float64 `json:"unitprice"`
	Url           string  `json:"url"`
}

type PriceHistory struct {
  Date string `json:"date"`
  Price float64 `json:"price"`
  Store string `json:"store"`
}

type Images struct {
	Small  string `json:"small"`
	Medium string `json:"medium"`
	Large  string `json:"large"`
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

// TODO: endre på verdein vi får som products i stedet for å returnere
func (products Products) Format() Products {
  newProducts := &Products{}

	// lagrer alle produkter som allerede har blitt sjekket i et array
	// da blir det ikke duplicates og vi kan returnere tidlig om produktet
	// allerede er ferdig
	var checkedGtins []string

	// mapper over alle produkter vi har fått fra databasen og formatterer
	// dataen i egne structs
	for _, firstProduct := range products {
		gtin := firstProduct.Id

		// om produktet allerede er sjekket, skip dette produktet
		if lib.IsIn(gtin, checkedGtins) {
			continue
		}
		checkedGtins = append(checkedGtins, gtin)

		// finner andre produkter med samme gtin og legger til i et array
		sameProduct := []Product{firstProduct}

		for _, secondProduct := range products {
			if gtin == secondProduct.Id &&
				firstProduct.Store != secondProduct.Store {
				// legger til produktet i sameProduct array, alle produkter i dette
				// arrayet sjekkes nå priser legges inn
				sameProduct = append(sameProduct, secondProduct)
			}
		}

		// Formaterer dataen til alle produkter med samme gtin
		firstProduct.AggregateProductData(sameProduct, newProducts)
	}

  // Assigner newProducts til products så vi har det nye products arrayet med
  // aggregated data
  return *newProducts
}

func (product *Product) AggregateProductData(productData Products, products *Products) {
  // For algolia
  product.ObjectID = product.Id

	// Lager et array av priser, å gjøre det på denne måten gjør det lettere
	// når dataen skal sendes til database
	var prices Price
	storeMap := map[string]bool{}
	// Sjekker at prisen ikke er 0, om den er det er det ikke vits å sende til
	// databasen
	for _, item := range productData {
    // Gjør også en sjekk for å sjekke om produktet er på salg fra en av
    // butikkene, så selv om ikke det første produktet er på salg, kan vi
    // den fortsatt markeres som onSale om andre produkter er det
    if item.OnSale {
      product.OnSale = true
    }

		if _, exists := storeMap[item.Store]; exists {
			continue
		}
		storeMap[item.Store] = true
		prices.Today = append(prices.Today, PriceToday{
			Store:         item.Store,
			Price:         item.Prices.Today[0].Price,
			OriginalPrice: item.Prices.Today[0].OriginalPrice,
			UnitPrice:     item.Prices.Today[0].UnitPrice,
			Url:           item.Prices.Today[0].Url,
		})

    // TODO: pris historikk
    //
    //
	}

	// Sorterer basert på pris, så det første elementet i arrayet vil være det
	// billigste
	priceCmp := func(a, b PriceToday) int {
		return cmp.Compare(a.Price, b.Price)
	}
	slices.SortFunc(prices.Today, priceCmp)
	product.Prices.Today = prices.Today

  if product.Category != "" {
		product.Category = getCorrectCategoryName(product.Category)
	}

  *products = append(*products, *product)
}

// Temporary løsning til jeg finner ut hvordan vi kan automatisere det
// Lager et array med kategorinavn som er forskjellige, som så blir brukt
// etterpå for å gjøre alt i samme kategori til samme kategori navn
// MENY, JOKER, SPAR
var categoryNames = [...][3]string{
	// Ting som basically kan automatiseres
	{"Frukt & grønt", "Frukt/Grønt", "Frukt og grønt"},
	{"Fisk & skalldyr", "Fisk/Skalldyr", "Fisk og skalldyr"},
	{"Personlige artikler", "Personlige Artikler"},

	// Ting som er mobbing (må kjøre noe contains greier ig)
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

// Bruker slicen av slices over til å gjøre alt til samme kategori navn
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

