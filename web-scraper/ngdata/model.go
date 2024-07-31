package ngdata

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/korneliushen/handlel.app/scraper/model"
)

// Bare no extra som trengs til ngdata (har ikke giddet å flytte ennå)
type ApiResponse struct {
	TimedOut bool `json:"timed_out"`
	Hits     Hits `json:"hits"`
}

type Hits struct {
	AmountOfProducts int          `json:"total"`
	Products         []ApiProduct `json:"hits"`
}

// type alias for an array of BaseProduct
type ApiProducts []ApiProduct

type ApiProduct struct {
	Store        string          `json:"store"`
	BaseUrl      string          `json:"base_url"`
	BaseImageUrl string          `json:"base_img_url"`
	Type         string          `json:"_type"`
	ApiId        string          `json:"_id"`
	Data         ApiProductData `json:"_source"`
}

type ApiProductData struct {
	Notes                 string                  `json:"notes"`
	Ean                   string                  `json:"ean"`
	Title                 string                  `json:"title"`
	SubTitle              string                  `json:"subtitle"`
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
	Allergens             []ApiAllergens          `json:"allergens"`
	NutritionalContent    []ApiNutritionalContent `json:"nutritionalContent"`
	Associated            ApiAssociated           `json:"associated"`
}

type ApiAllergens struct {
	Name string `json:"displayName"`
	Code string `json:"code"`
}

type ApiNutritionalContent struct {
	Name   string  `json:"displayName"`
	Amount float32 `json:"amount"`
	Unit   string  `json:"unit"`
}

type ApiAssociated struct {
	Products []string `json:"slg"`
}

// oppdaterer data i BaseProduct
func (product ApiProduct) Extend(store, baseUrl string) ApiProduct {
	product.Store = store
	product.BaseUrl = baseUrl
	return product
}

// Formaterer dataen fra databasen til formatet databasen og algolia vil ha
func (baseProduct ApiProduct) Format() model.Product {
  product := model.Product{}

	// Legger til alle fields fra BaseProduct som har samme navn som Product
	vDest := reflect.ValueOf(&product).Elem()
	vSrc := reflect.ValueOf(&baseProduct.Data).Elem()
	for i := range vDest.NumField() {
		fieldDest := vDest.Field(i)
		fieldSrc := vSrc.FieldByName(vDest.Type().Field(i).Name)

		if fieldSrc.IsValid() && fieldSrc.Type() == fieldDest.Type() {
			fieldDest.Set(fieldSrc)
		}
	}

	product.Id = baseProduct.Data.Ean

  product.Store = baseProduct.Store

  price := model.PriceToday{}
  price.Store = baseProduct.Store
  price.Price = baseProduct.Data.Price
  price.OriginalPrice = baseProduct.Data.OriginalPrice
  price.UnitPrice = baseProduct.Data.ComparePricePerUnit

  // ComparePricePerUnit er noen ganger null av en eller annen grunn.
  // Har bare sett dette skje på produkter som selger et stk, så dette vil
  // være accurate for de tilfellene
  if baseProduct.Data.ComparePricePerUnit == 0 {
    price.UnitPrice = price.Price
  }

  price.Url = fmt.Sprintf("%s%s", baseProduct.BaseUrl, baseProduct.Data.Slug)
  product.Prices.Today = append(product.Prices.Today, price)

	// lager en string for vekt, med value og unit
	product.Weight = fmt.Sprintf("%v%s",
		baseProduct.Data.Weight, baseProduct.Data.WeightMeasurementType)

  // Base url for bilder
  baseImgUrl := "https://bilder.ngdata.no/"
  // Lager small, medium og large versjoner av image. Lager hele url-en
	product.Images.Small = fmt.Sprintf("%s%s%s", baseImgUrl, baseProduct.Data.ImageLink, "/small.jpg")
  product.Images.Medium = fmt.Sprintf("%s%s%s", baseImgUrl, baseProduct.Data.ImageLink, "/medium.jpg")
  product.Images.Large = fmt.Sprintf("%s%s%s", baseImgUrl, baseProduct.Data.ImageLink, "/large.jpg")

	// Mapper over allergener array som vi fikk fra databasen
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
		nutritionalContent := model.NutritionalContent{}
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

  return product
}
