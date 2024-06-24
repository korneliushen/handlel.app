package model

import (
	"cmp"
	"fmt"
	"reflect"
	"slices"
	"strings"
)

// Bare no extra som trengs til ngdata (har ikke giddet å flytte ennå)
type ApiResponse struct {
	TimedOut bool `json:"timed_out"`
	Hits     Hits `json:"hits"`
}

type Hits struct {
	AmountOfProducts int          `json:"total"`
	Products         []BaseProduct `json:"hits"`
}

// type alias for an array of BaseProduct
type BaseProducts []BaseProduct

type BaseProduct struct {
	Store        string         `json:"store"`
	BaseUrl      string         `json:"base_url"`
	BaseImageUrl string         `json:"base_img_url"`
	Type         string         `json:"_type"`
	ApiId        string         `json:"_id"`
	Data         BaseProductData `json:"_source"`
}

type BaseProductData struct {
	Notes                 string                  `json:"notes"`
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
	ImageLinkSmall        string                  `json:"imagePathSmall"`
	ImageLinkMedium       string                  `json:"imagePathMedium"`
	ImageLinkLarge        string                  `json:"imagePathLarge"`
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
	Allergens             []BaseAllergens          `json:"allergens"`
	NutritionalContent    []BaseNutritionalContent `json:"nutritionalContent"`
	Associated            BaseAssociated           `json:"associated"`
}

type BaseAllergens struct {
	Name string `json:"displayName"`
	Code string `json:"code"`
}

type BaseNutritionalContent struct {
	Name   string  `json:"displayName"`
	Amount float32 `json:"amount"`
	Unit   string  `json:"unit"`
}

type BaseAssociated struct {
	Products []string `json:"slg"`
}

// oppdaterer data i BaseProduct
func (product BaseProduct) Extend(store, baseUrl string) BaseProduct {
	product.Store = store
	product.BaseUrl = baseUrl
	return product
}

// legger en method til i BaseProduct struct så vi kan accesse baseProduct
func (baseProduct *BaseProduct) FormatData(productData []BaseProduct, products *Products) {
	product := Product{}

	// for algolia
	product.ObjectID = baseProduct.Data.Ean

	product.Id = baseProduct.Data.Ean

	// legger til alle fields fra BaseProduct som har samme navn som Product
	vDest := reflect.ValueOf(&product).Elem()
	vSrc := reflect.ValueOf(&baseProduct.Data).Elem()
	for i := range vDest.NumField() {
		fieldDest := vDest.Field(i)
		fieldSrc := vSrc.FieldByName(vDest.Type().Field(i).Name)

		if fieldSrc.IsValid() && fieldSrc.Type() == fieldDest.Type() {
			fieldDest.Set(fieldSrc)
		}
	}

	// lager en string for vekt, med value og unit
	product.Weight = fmt.Sprintf("%v%s",
		baseProduct.Data.Weight, baseProduct.Data.WeightMeasurementType)

	// lager hele url-en for bildelinker for ulike størrelser
	product.Images.Small = baseProduct.Data.ImageLinkSmall
	product.Images.Medium = baseProduct.Data.ImageLinkMedium
	product.Images.Large = baseProduct.Data.ImageLinkLarge

	// Fikser kategori navn om kategori ikke er en empty string
	// (bruker hard-coda kategori navn for å gjøre ting til samme kategori)
	if product.Category != "" {
		product.Category = getCorrectCategoryName(baseProduct.Data.Category)
	}

	// Lager et array av priser, å gjøre det på denne måten gjør det lettere
	// når dataen skal sendes til database
	var prices []Price
	storeMap := map[string]bool{}
	// Sjekker at prisen ikke er 0, om den er det er det ikke vits å sende til
	// databasen
	for _, item := range productData {
    // Gjør også en sjekk for å sjekke om produktet er på salg fra en av
    // butikkene, så selv om ikke det første produktet er på salg, kan vi
    // den fortsatt markeres som onSale om andre produkter er det
    if item.Data.OnSale {
      product.OnSale = true
    }

		if _, exists := storeMap[item.Store]; exists {
			continue
		}
		storeMap[item.Store] = true
		prices = append(prices, Price{
			Store:         item.Store,
			Price:         item.Data.Price,
			OriginalPrice: item.Data.OriginalPrice,
			UnitPrice:     item.Data.ComparePricePerUnit,
			Url:           fmt.Sprintf("%s%s", item.BaseUrl, item.Data.Slug),
		})
	}

	// Sorterer basert på pris, så det første elementet i arrayet vil være det
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
	for _, allergen := range baseProduct.Data.Allergens {
		if allergen.Code == "JA" {
			allergens = append(allergens, allergen.Name)
		} else if allergen.Code == "KAN" {
			mayContainTracesOf = append(mayContainTracesOf, allergen.Name)
		}
	}
	product.Allergens = strings.Join(allergens, ", ")
	product.MayContainTracesOf = strings.Join(mayContainTracesOf, ", ")

	// næringsinnhold
	nutritionalContentData := baseProduct.Data.NutritionalContent

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
