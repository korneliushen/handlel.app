package scraper

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
