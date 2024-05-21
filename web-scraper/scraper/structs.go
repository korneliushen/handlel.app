package scraper

type Kategorier struct {
	Kategorier []Kategori
}

type Kategori struct {
	Navn            string
	Underkategorier []Underkategori
}

type Underkategori struct {
	Navn      string
	Produkter []Produkt
}

// TODO: gjøre alt som kan bli int/float til det
type Produkt struct {
	Tittel    string
	Pris      string
	KiloPris  string
	BildeLink string
	Innhold   Innhold
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

// type Næringsinnhold struct {}
