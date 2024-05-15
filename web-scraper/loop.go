package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

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

func getProductInfo(link string, cursor int) {
	c := colly.NewCollector()

	// henter data for hvert element (printer bare for nå)
	c.OnHTML("article", func(e *colly.HTMLElement) {
		fmt.Println("Title: ", e.ChildText("div > div h2"))
		fmt.Println("Price: ", e.ChildText("div > div > div > span"))
		fmt.Println("Kilopris: ", e.ChildText("div > div > div > p"))
		fmt.Println("Bildelenke: ", e.ChildAttr("div > img", "src"))
	})

	visitLink := fmt.Sprintf("%s&cursor=%v", link, cursor)
	c.Visit(visitLink)
}
