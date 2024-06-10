package model

import (
	"cmp"
	"fmt"
	"reflect"
	"slices"
	"strings"
)

// TODO: Rename alt her til noe annet enn bare Api... (kanskje shared eller no)

// Bare no extra som trengs til ngdata (har ikke giddet å flytte ennå)
type ApiResponse struct {
	TimedOut bool `json:"timed_out"`
	Hits     Hits `json:"hits"`
}

type Hits struct {
	AmountOfProducts int          `json:"total"`
	Products         []ApiProduct `json:"hits"`
}

// type alias for an array of ApiProduct
type ApiProducts []ApiProduct

type ApiProduct struct {
	Store   string         `json:"store"`
	BaseUrl string         `json:"base_url"`
	Type    string         `json:"_type"`
	ApiId   string         `json:"_id"`
	Data    ApiProductData `json:"_source"`
}

type ApiProductData struct {
	Ean                   string                  `json:"ean"`
	Title                 string                  `json:"title"`
	Subtitle              string                  `json:"subtitle"`
	Slug                  string                  `json:"slugifiedUrl"`
	Description           string                  `json:"description"`
	Category              string                  `json:"categoryName"`
	SubCategory           string                  `json:"shoppingListGroupName"`
	Price                 float64                 `json:"pricePerUnit"`
	OriginalPrice         float64                 `json:"pricePerUnitOriginal"`
	ComparePricePerUnit   float64                 `json:"comparePricePerUnit"`
	UnitType              string                  `json:"compareUnit"`
	ImageLink             string                  `json:"imagePath"`
	WeightMeasurementType string                  `json:"measurementType"`
	Weight                float64                 `json:"measurementValue"`
	Unit                  string                  `json:"unit"`
	Size                  string                  `json:"packageSize"`
	Ingredients           string                  `json:"ingredients"`
	AllergyDeclaration    string                  `json:"allergyDeclaration"`
	Vendor                string                  `json:"vendor"`
	Brand                 string                  `json:"brand"`
	OnSale                bool                    `json:"isOffer"`
	OriginCountry         string                  `json:"countryOfOrigin"`
	Allergens             []ApiAllergens          `json:"allergens"`
	NutritionalContent    []ApiNutritionalContent `json:"nutritionalContent"`
	Associated            ApiAssociated           `json:"associated"`
}

type ApiAllergens struct {
	Name string `json:"displayName"`
	Code string `json:"code"`
}

type ApiNutritionalContent struct {
	Id     string  `json:"name"`
	Name   string  `json:"displayName"`
	Amount float32 `json:"amount"`
	Unit   string  `json:"unit"`
}

type ApiAssociated struct {
	Products []string `json:"slg"`
}

// oppdaterer data i ApiProduct
func (product ApiProduct) Extend(store, baseUrl string) ApiProduct {
	product.Store = store
	product.BaseUrl = baseUrl
	return product
}

// legger en method til i ApiProduct struct så vi kan accesse apiProduct
func (apiProduct *ApiProduct) FormatData(productData []ApiProduct, products *Products) {
	product := Product{}

	// for algolia
	product.ObjectID = apiProduct.Data.Ean

	product.Id = apiProduct.Data.Ean

	// legger til alle fields fra ApiProduct som har samme navn som Product
	vDest := reflect.ValueOf(&product).Elem()
	vSrc := reflect.ValueOf(&apiProduct.Data).Elem()
	for i := range vDest.NumField() {
		fieldDest := vDest.Field(i)
		fieldSrc := vSrc.FieldByName(vDest.Type().Field(i).Name)

		if fieldSrc.IsValid() && fieldSrc.Type() == fieldDest.Type() {
			fieldDest.Set(fieldSrc)
		}
	}

	// lager en string for vekt, med value og unit
	product.Weight = fmt.Sprintf("%v%s",
		apiProduct.Data.Weight, apiProduct.Data.WeightMeasurementType)

	// lager hele url-en for bildelinker for ulike størrelser
	product.ImageLink = fmt.Sprintf("%s%s",
		"https://bilder.ngdata.no/", apiProduct.Data.ImageLink)

	// fikser kategori navn
	// (bruker hard-coda kategori navn for å gjøre ting til samme kategori)
	product.Category = getCorrectCategoryName(apiProduct.Data.Category)

	// lager et array av priser, å gjøre det på denne måten gjør det lettere
	// når dataen skal sendes til database
	var prices []Price
	storeMap := map[string]bool{}
	// sjekker at prisen ikke er 0, om den er det er det ikke vits å sende til
	// databasen
	for _, product := range productData {
		if _, exists := storeMap[product.Store]; exists {
			continue
		}
		storeMap[product.Store] = true
		prices = append(prices, Price{
			Store:         product.Store,
			Price:         product.Data.Price,
			OriginalPrice: product.Data.OriginalPrice,
			UnitPrice:     product.Data.ComparePricePerUnit,
			Url:           fmt.Sprintf("%s%s", product.BaseUrl, product.Data.Slug),
		})
	}

	// sorterer basert på pris, så det første elementet i arrayet vil være det
	// billigste
	priceCmp := func(a, b Price) int {
		return cmp.Compare(a.Price, b.Price)
	}
	slices.SortFunc(prices, priceCmp)
	product.Prices = prices

	// mapper over allergener array som vi fikk fra databasen
	// i databasen så bestemmer koden hva itemet i arrayen betyr for produktet
	// om koden er JA, blir det lagt til i allergens, om det er kan blir det
	// lagt til i mayContainTracesOf
	var allergens []string
	var mayContainTracesOf []string
	for _, allergen := range apiProduct.Data.Allergens {
		if allergen.Code == "JA" {
			allergens = append(allergens, allergen.Name)
		} else if allergen.Code == "KAN" {
			mayContainTracesOf = append(mayContainTracesOf, allergen.Name)
		}
	}
	product.Allergens = strings.Join(allergens, ", ")
	product.MayContainTracesOf = strings.Join(mayContainTracesOf, ", ")

	// næringsinnhold
	nutritionalContentData := apiProduct.Data.NutritionalContent

	// om det ikke er noe næringsinnhold
	if len(nutritionalContentData) == 0 {
		product.NutritionalContent = nil
	} else {
		nutritionalContent := NutritionalContent{}
		v := reflect.ValueOf(&nutritionalContent).Elem()
		// legger til næringsinnhold data i fields med navn som matcher dataen
		// fra api-en (reflect)
		for i := range len(nutritionalContentData) {
			field := v.FieldByName(nutritionalContentData[i].Name)
			if field.CanSet() {
				// value er bygd opp av mengde og enhet (xxxg)
				field.SetString(fmt.Sprintf("%v%s",
					nutritionalContentData[i].Amount, nutritionalContentData[i].Unit))
			}
		}
		product.NutritionalContent = &nutritionalContent
	}

	// legger til produktet i products array
	*products = append(*products, product)
}

// temporary løsning til jeg finner ut hvordan vi kan automatisere det
// lager et array med kategorinavn som er forskjellige, som så blir brukt
// etterpå for å gjøre alt i samme kategori til samme kategori navn
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
