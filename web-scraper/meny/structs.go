package main

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

type Produkter struct {
	Produkter []Produkt
}

type Produkt struct {
	Id            int
	Gtin          int
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
	Joker     float32
	Meny      float32
	Spar      float32
	JokerKilo float32
	MenyKilo  float32
	SparKilo  float32
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
	Allergener         []string
	KanInneholdeSporAv string
	Vekt               string
	Bruksområde        string
	Næringsinnhold     Næringsinnhold
}

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

type ApiResponse struct {
	TidsAvbrudd bool     `json:"timed_out"`
	Resultat    Resultat `json:"hits"`
}

type Resultat struct {
	Antall    int          `json:"total"`
	Produkter []ApiProdukt `json:"hits"`
}

type ApiProdukt struct {
	Type  string `json:"_type"`
	ApiId string `json:"_id"`
}
