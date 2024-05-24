package main

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
	Id            int
	Gtin          string
	FullNavn      string
	FørsteNavn    string
	AndreNavn     string
	Kategori      string
	Underkategori string
	Priser        Priser
	BildeLink     string
	Innhold       Innhold
}

type Priser struct {
	Joker         float32
	Meny          float32
	Spar          float32
	JokerOriginal float32
	MenyOriginal  float32
	SparOriginal  float32
	JokerKilo     float32
	MenyKilo      float32
	SparKilo      float32
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
	Allergener         []Allergens
	KanInneholdeSporAv string
	Vekt               string
	Bruksområde        string
	Næringsinnhold     []NutritionalContent
}

// currently ikke i bruk
type Næringsinnhold struct {
	Energi                      string
	Natrium                     string
	Fett                        string
	HvoravMettedeFettsyrer      string
	HvoravEnumettedeFettsyrer   string
	HvoravFlerumettedeFettsyrer string
	Karbohydrater               string
	HvoravSukkerarter           string
	HvoravPolyoler              string
	HvoravStivelse              string
	Kostfiber                   string
	Protein                     string
	Salt                        string
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
	ImageLink             string               `json:"imagePath"` // https://bilder.ngdata.no/BildeLink/medium.jpg (eller small)
	WeightMeasurementType string               `json:"measurementType"`
	Weight                float32              `json:"weight"`
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
	Code string `json:"code"`
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
