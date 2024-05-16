package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

func getProducts() {
	c := colly.NewCollector()

	categories := Kategorier{}

	// besøker oda og henter alle kategorier
	// for hver kategori, kjører getUnderCategories med kategori linken man får
	c.OnHTML("article.k-p-4", func(e *colly.HTMLElement) {
		start := time.Now()
		category := Kategori{}

		categoryName := e.ChildText("div > a > h1")
		category.Navn = categoryName

		fmt.Println("Getting data for: ", categoryName)

		categoryLink := e.ChildAttr("div > a", "href")
		getUnderCategories(categoryName, categoryLink, category)

		categories.Kategorier = append(categories.Kategorier, category)

		// temp code just for testing
		jsonData, err := json.MarshalIndent(categories, "", "    ")
		if err != nil {
			fmt.Printf("Error marshalling to JSON: %v\n", err)
			return
		}

		fmt.Println("Data: ", string(jsonData))

		end := time.Now()
		elapsed := end.Sub(start)
		fmt.Println("Found data for category: ", categoryName, "in: ", elapsed, "seconds")
	})

	c.Visit("https://oda.com/no/products/")

	jsonData, err := json.MarshalIndent(categories, "", "    ")
	if err != nil {
		fmt.Printf("Error marshalling to JSON: %v\n", err)
		return
	}

	err = os.WriteFile("./data.json", jsonData, 0666)
	if err != nil {
		fmt.Printf("Error writing json data to file %v\n", err)
		return
	}

	fmt.Println("Data: ", string(jsonData))
}

func getUnderCategories(categoryName string, categoryLink string, category Kategori) {
	c := colly.NewCollector()

	// finner hver under kategori
	// Eksempel: oda.com/no/categories/20-frukt-og-gront og så Frukt, Bær, osv.
	// for hver underkategori, henter mengden sider den underkategorien har (cursor)
	c.OnHTML("section", func(e *colly.HTMLElement) {
		// lager instans av under kategori
		underCategory := Underkategori{}

		// gir instansen et navn
		underCategoryName := e.ChildText("section > a > h2")
		underCategory.Navn = underCategoryName

		fmt.Println("Getting data for undercategory: ", underCategoryName, "in category", categoryName)

		underCategoryLink := e.ChildAttr("section > a", "href")
		getPageCount(underCategoryName, underCategoryLink, underCategory)

		category.Underkategorier = append(category.Underkategorier, underCategory)
	})

	link := fmt.Sprintf("https://oda.com%s", categoryLink)

	c.Visit(link)
}

func getPageCount(underCategoryName string, underCategoryLink string, underCategory Underkategori) {
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
			getProductInfo(link, i+1, underCategory)
		}
	})
	link := fmt.Sprintf("https://oda.com%s", underCategoryLink)
	c.Visit(link)
}

// Sammenligner key til innholdet (navnet), med et field i Innhold structen
// om den finner en key som matcher en field, legges det til i instansen av Innhold
func setFieldValue(in *Innhold, key string, value string) {
	v := reflect.ValueOf(in).Elem()

	// denne koden gjør at alt med hvorav funker som det skal
	key = strings.Title(key)
	key = strings.ReplaceAll(key, " ", "")

	// finner fieldet som matcher key
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

func getProductInfo(link string, cursor int, underCategory Underkategori) {
	// lager en instans av Produkter
	c := colly.NewCollector()

	// henter data for hvert element (printer bare for nå)
	c.OnHTML("article", func(e *colly.HTMLElement) {
		// lager en instans av Produkt, som har all infoen om produktet
		productInfo := Produkt{}
		// bruker en referanse verdi her for at det skal funke med reflect
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

		// legger innhold til i produktinfo
		productInfo.Innhold = *contents
		// setter produktet inn i produkter arrayet
		underCategory.Produkter = append(underCategory.Produkter, productInfo)
	})

	// besøker linken med underkategori og cursor (side nr)
	// visitLink := fmt.Sprintf("%s&cursor=%v", link, cursor)
	visitLink := fmt.Sprintf("https://oda.com/no/categories/20-frukt-og-gront/21-frukt/?filters=&cursor=3")
	c.Visit(visitLink)
}
