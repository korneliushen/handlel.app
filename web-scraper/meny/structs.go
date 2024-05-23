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
	Pris          string
	KiloPris      string
	BildeLink     string
	Innhold       Innhold
}

type Innhold struct {
	Salgsinformasjon   string
	Beskrivelse        string
	Holdbarhetsgaranti string
	Utleveringsdager   string
	Størrelse          string
	Leverandør         string
	Ingredienser       string
	Oppbevaring        string
	Opprinnelsesland   string
	Opphavssted        string
	Egenskaper         string
	Inneholder         string
	Allergener         string
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
