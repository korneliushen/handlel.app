package main

// TODO: rewrite hele denne dritten (ta alt på engelsk med bedre naming)

// NOTE
// alt som er data fra api-en er på engelsk, alt som er laget selv og skal sendes til databasen, er på norsk
// NOTE

// brukt for å få alle kategorier som skal hentes produkter fra
type Kategorier struct {
	Kategorier []Kategori
}

type Kategori struct {
	Navn            string
	Link            string
	Underkategorier []Underkategori
}

type Underkategori struct {
	Navn string
	Link string
}

// produkter som skal bli til json data/lagt inn i database
type Produkter struct {
	Produkter []Produkt
}

type Produkt struct {
	Gtin          string
	Tittel        string
	Undertittel   string
	Kategori      string
	Underkategori string
	Priser        Priser
	BildeLink     string
	Innhold       Innhold
}

type Priser struct {
	Priser []Pris
}

type Pris struct {
	Butikk       string
	Pris         float32
	OriginalPris float32
	EnhetsPris   float32
	EnhetsType   string
	Url          string
}

type Innhold struct {
	Beskrivelse        string
	Holdbarhetsgaranti string
	Enhet              string
	Størrelse          string
	Leverandør         string
	Ingredienser       string
	Oppbevaring        string
	Opprinnelsesland   string
	Opphavssted        string
	Egenskaper         string
	Allergener         string
	KanInneholdeSporAv string
	Vekt               string
	Bruksområde        string
	Næringsinnhold     Næringsinnhold
}

type Næringsinnhold struct {
	Energi          string
	Kalorier        string
	Natrium         string
	Fett            string
	MettetFett      string
	EnumettetFett   string
	FlerumettetFett string
	Karbohydrater   string
	Sukkerarter     string
	Stivelse        string
	Kostfiber       string
	Protein         string
	Salt            string
}

// setter data fra api inn i eget struct, gjør det lettere å assigne data senere
type ApiResponse struct {
	TimedOut bool `json:"timed_out"`
	Hits     Hits `json:"hits"`
}

type Hits struct {
	AmountOfProducts int       `json:"total"`
	Products         []Product `json:"hits"`
}

type Product struct {
	Type  string      `json:"_type"`
	ApiId string      `json:"_id"`
	Data  ProductData `json:"_source"`
}

type ProductData struct {
	Ean                   string               `json:"ean"`
	Title                 string               `json:"title"`
	Subtitle              string               `json:"subtitle"`
	Slug                  string               `json:"slugifiedUrl"`
	Description           string               `json:"description"`
	Category              string               `json:"categoryName"`
	SubCategory           string               `json:"shoppingListGroupName"`
	Price                 float32              `json:"pricePerUnit"`
	OriginalPrice         float32              `json:"pricePerUnitOriginal"`
	CalcPricePerUnit      float32              `json:"calcPricePerUnit"`
	CalcUnit              string               `json:"calcUnit"`
	ImageLink             string               `json:"imagePath"` // https://bilder.ngdata.no/BildeLink/medium.jpg (eller small)
	WeightMeasurementType string               `json:"measurementType"`
	Weight                float32              `json:"measurementValue"`
	Unit                  string               `json:"unit"`
	Size                  string               `json:"packageSize"`
	Ingredients           string               `json:"ingredients"`
	AllergyDeclaration    string               `json:"allergyDeclaration"`
	Vendor                string               `json:"vendor"`
	OnSale                bool                 `json:"isOffer"`
	OriginCountry         string               `json:"countryOfOrigin"`
	Allergens             []Allergens          `json:"allergens"`
	NutritionalContent    []NutritionalContent `json:"nutritionalContent"`
	Associated            Associated           `json:"associated"`
}

type Allergens struct {
	Name string `json:"displayName"`
}

type NutritionalContent struct {
	Id     string  `json:"name"`
	Name   string  `json:"displayName"`
	Amount float32 `json:"amount"`
	Unit   string  `json:"unit"`
}

type Associated struct {
	Products []string `json:"slg"`
}
