package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

func getCategories() {
	c := colly.NewCollector()

	c.OnHTML(".k-p-4", func(e *colly.HTMLElement) {
		categoryLink := e.ChildAttr("div > a", "href")
		getUnderCategories(categoryLink)
	})

	c.Visit("https://oda.com/no/products/")
}

func getUnderCategories(categoryLink string) {
	c := colly.NewCollector()

	fmt.Println("going to: ", categoryLink)

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
		fmt.Println(antallSider)
		underCategoryLink := "/no/categories/1135-bakeri-og-konditori/1136-handverksbakeri/?filters="
		link := fmt.Sprintf("https://oda.com%s", underCategoryLink)
		for i := range antallSider {
			getProductInfo(link, i+1)
		}
	})
	link := fmt.Sprintf("https://oda.com%s", underCategoryLink)
	c.Visit(link)
}

func getProductInfo(link string, cursor int) {
	c := colly.NewCollector()

	c.OnHTML("article", func(e *colly.HTMLElement) {
		fmt.Println("Title: ", e.ChildText("div > div h2"))
		fmt.Println("Price: ", e.ChildText("div > div > div > span"))
		fmt.Println("Kilopris: ", e.ChildText("div > div > div > p"))
		fmt.Println("Bildelenke: ", e.ChildAttr("div > img", "src"))
	})

	visitLink := fmt.Sprintf("%s&cursor=%v", link, cursor)
	fmt.Println("Link", visitLink)
	c.Visit(visitLink)
}
