package main

import (
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

type Produkter struct {
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
	Ingredienser       []string
	Oppbevaring        string
	// næringsinnhold (gjør til egen struct kanskje senere)
	Energi                      string
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

func getProducts() {
	c := colly.NewCollector()

	// besøker oda og henter alle kategorier
	// for hver kategori, kjører getUnderCategories med kategori linken man får
	c.OnHTML(".k-p-4", func(e *colly.HTMLElement) {
		categoryLink := e.ChildAttr("div > a", "href")
		getUnderCategories(categoryLink)
	})

	c.Visit("https://oda.com/no/products/")
}

func getUnderCategories(categoryLink string) {
	c := colly.NewCollector()

	// finner hver under kategori
	// Eksempel: oda.com/no/categories/20-frukt-og-gront og så Frukt, Bær, osv.
	// for hver underkategori, henter mengden sider den underkategorien har (cursor)
	c.OnHTML("section", func(e *colly.HTMLElement) {
		underCategoryLink := e.ChildAttr("section > a", "href")
		getPageCount(underCategoryLink)
	})

	link := fmt.Sprintf("https://oda.com%s", categoryLink)

	c.Visit(link)
}

func getPageCount(underCategoryLink string) {
	c := colly.NewCollector()

	c.OnHTML("main", func(e *colly.HTMLElement) {
		// finner mengden sider
		tall := strings.Split(e.ChildText("main > div > div > div > span > div > div > a.k-choice-chip--selected.k-choice-chip--primary > span.k-pill--extra-small"), "")
		tallLengde := len(tall) / 2
		var s []string
		for i := 0; i < tallLengde; i++ {
			s = append(s, tall[i])
		}
		str := strings.Join(s, "")
		antallVarer, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return
		}
		antallVarerDelt := antallVarer / 24
		antallSider := int(math.Ceil(float64(antallVarerDelt)))
		// lager en link for underkategorien
		link := fmt.Sprintf("https://oda.com%s", underCategoryLink)
		// for hver side, hent produktinfo for alle produktene på siden
		for i := range antallSider {
			getProductInfo(link, i+1)
		}
	})
	link := fmt.Sprintf("https://oda.com%s", underCategoryLink)
	c.Visit(link)
}

// TODO: rewrite
// convert ord som ikke matcher i struct
// legg til cases for float og int (for gram osv??)
func setFieldValue(in *Innhold, key string, value string) {
	v := reflect.ValueOf(in).Elem()

	key = strings.Title(key)

	key = strings.ReplaceAll(key, " ", "")

	fmt.Println(key)

	field := v.FieldByName(key)

	if !field.IsValid() || !field.CanSet() {
		fmt.Printf("Cannot set field %s\n", key)
		return
	}

	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Slice:
		field.Set(reflect.ValueOf([]string{value}))
	default:
		fmt.Printf("Unsupported kind %s\n", field.Kind())
	}
}

func getProductInfo(link string, cursor int) {
	products := Produkter{}

	c := colly.NewCollector()

	// henter data for hvert element (printer bare for nå)
	c.OnHTML("article", func(e *colly.HTMLElement) {
		productInfo := Produkt{}
		contents := &Innhold{}

		title := e.ChildText("div > div h2")

		if title == "" {
			return
		}

		productInfo.Tittel = title
		productInfo.Pris = e.ChildText("div > div > div > span")
		productInfo.KiloPris = e.ChildText("div > div > div > p")
		productInfo.BildeLink = e.ChildAttr("div > img", "src")

		// link for side med næringsinnhold
		productLink := e.ChildAttr("div > div a", "href")
		dataLink := fmt.Sprintf("https://oda.com%s", productLink)

		// henter data fra produkt siden
		n := colly.NewCollector()

		n.OnHTML("#read-more-product-description-nb", func(h *colly.HTMLElement) {
			// henter beskrivelsen
			description := h.ChildText("p")
			contents.Beskrivelse = description
		})

		// henter innhold og næringsinnhold (om det eksisterer)
		n.OnHTML("div.k-grid.k-pt-3.k-pb-6", func(h *colly.HTMLElement) {
			key := h.ChildText("div > span")
			value := h.ChildText("div > p")
			if strings.Contains(key, "Oppbevaring") {
				key = "Oppbevaring"
				value = h.ChildText("div > div span")
			}
			setFieldValue(contents, key, value)
		})

		n.Visit(dataLink)

		productInfo.Innhold = *contents
		products.Produkter = append(products.Produkter, productInfo)
	})

	// besøker linken med underkategori og cursor (side nr)
	// visitLink := fmt.Sprintf("%s&cursor=%v", link, cursor)
	visitLink := fmt.Sprintf("https://oda.com/no/categories/20-frukt-og-gront/21-frukt/?filters=&cursor=3")
	c.Visit(visitLink)

	jsonData, err := json.Marshal(products)
	if err != nil {
		fmt.Println("Error marshalling to JSON:", err)
		return
	}

	fmt.Println(string(jsonData))
}
