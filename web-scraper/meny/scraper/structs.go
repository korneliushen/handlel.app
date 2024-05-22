package scraper

type Kategorier struct {
	Kategorier []Kategori
}

type Kategori struct {
	Navn            string
	Link            string
	QueryFormatNavn string
	Underkategorier []Underkategori
}

type Underkategori struct {
	Navn            string
	Link            string
	QueryFormatNavn string
}
