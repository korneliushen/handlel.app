package lib

// brukt for å få alle kategorier som skal hentes produkter fra

type StoreData struct {
	Store       string
	ApiRes      ApiResponse
	Category    string
	SubCategory string
}

type Categories struct {
	Categories []Category
}

type Category struct {
	Name          string
	Store         string
	SubCategories []string
}

// produkter som skal bli til json data/lagt inn i database
type Products struct {
	Products []Product
}

type Product struct {
	Gtin        string
	Title       string
	SubTitle    string
	Category    string
	SubCategory string
	Prices      Prices
	OnSale      bool
	Content     Content
	Images      Images
}

// ulik størrelse på bilder
type Images struct {
	ImageLinkXSmall string
	ImageLinkSmall  string
	ImageLinkMedium string
	ImageLinkLarge  string
	ImageLinkXLarge string
}

type Prices struct {
	Prices []Price
}

type Price struct {
	Store         string  `json:"store"`
	Price         float64 `json:"price"`
	OriginalPrice float64 `json:"originalprice"`
	UnitPrice     float64 `json:"unitprice"`
	Url           string  `json:"url"`
}

type Content struct {
	Description        string
	Duration           string
	Unit               string
	UnitType           string
	Size               string
	Vendor             string
	Brand              string
	Ingredients        string
	Storage            string
	OriginCountry      string
	Features           string
	Allergens          string
	MayContainTracesOf string
	Weight             string
	NutritionalContent *NutritionalContent
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

// setter data fra api inn i eget struct, gjør det lettere å assigne data senere
type ApiResponse struct {
	TimedOut bool `json:"timed_out"`
	Hits     Hits `json:"hits"`
}

type Hits struct {
	AmountOfProducts int          `json:"total"`
	Products         []ApiProduct `json:"hits"`
}

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
	CompareUnit           string                  `json:"compareUnit"`
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
