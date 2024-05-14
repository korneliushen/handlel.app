package main

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

// func getCategories() {
// 	c := colly.NewCollector()

// 	c.OnRequest(func(r *colly.Request) {
// 		fmt.Println("Visiting: ", r.URL)
// 	})

// 	c.OnError(func(_ *colly.Response, err error) {
// 		log.Println("Something went wrong: ", err)
// 	})

// 	c.OnResponse(func(r *colly.Response) {
// 		fmt.Println("Page visited: ", r.Request.URL)
// 	})

// 	c.OnHTML(".k-p-4", func(e *colly.HTMLElement) {
// 		fmt.Println("categoryLink: ", e.ChildAttr("div > a", "href"))
// 		categoryLink := e.ChildAttr("div > a", "href")
// 		getUnderCategories(categoryLink)
// 	})

// 	c.OnScraped(func(r *colly.Response) {
// 		fmt.Println(r.Request.URL, " scraped!")
// 	})

// 	c.Visit("https://oda.com/no/products/")
// }

// func getUnderCategories(categoryLink string) {
// 	c := colly.NewCollector()

// 	c.OnRequest(func(r *colly.Request) {
// 		fmt.Println("Visiting: ", r.URL)
// 	})

// 	c.OnError(func(_ *colly.Response, err error) {
// 		log.Println("Something went wrong: ", err)
// 	})

// 	c.OnResponse(func(r *colly.Response) {
// 		fmt.Println("Page visited: ", r.Request.URL)
// 	})

// 	c.OnHTML("section", func(e *colly.HTMLElement) {
// 		fmt.Println("underCategoryLink: ", e.ChildAttr("section > a", "href"))
// 		underCategoryLink := e.ChildAttr("section > a", "href")
// 		getProductInfo(underCategoryLink)
// 	})

// 	c.OnScraped(func(r *colly.Response) {
// 		fmt.Println(r.Request.URL, " scraped!")
// 	})

// 	link := fmt.Sprintf("https://oda.com%s", categoryLink)

// 	c.Visit(link)
// }

// func getProductInfo(underCategoryLink string) {
// 	c := colly.NewCollector()

// 	c.OnRequest(func(r *colly.Request) {
// 		fmt.Println("Visiting: ", r.URL)
// 	})

// 	c.OnError(func(_ *colly.Response, err error) {
// 		log.Println("Something went wrong: ", err)
// 	})

// 	c.OnResponse(func(r *colly.Response) {
// 		fmt.Println("Page visited: ", r.Request.URL)
// 	})

// 	c.OnHTML("article", func(e *colly.HTMLElement) {
// 		fmt.Println("Title: ", e.ChildText("div > div h2"))
// 		fmt.Println("Price: ", e.ChildText("div > div > div > span"))
// 		fmt.Println("Kilopris: ", e.ChildText("div > div > div > p"))
// 		fmt.Println("Bildelenke: ", e.ChildAttr("div > img", "src"))
// 	})

// 	c.OnScraped(func(r *colly.Response) {
// 		fmt.Println(r.Request.URL, " scraped!")
// 	})

// 	n := colly.NewCollector()

// 	n.OnHTML(".k-pill--extra-small", func(e *colly.HTMLElement) {
// 		fmt.Println("Title: ", e.ChildText("main > div > div > div > span > div > div > span"))
// 	})

// 	link := fmt.Sprintf("https://oda.com%s&cursor=", underCategoryLink)

// 	c.Visit(link)
// }

func getProductInfo() {
	c := colly.NewCollector()

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting: ", r.URL)
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong: ", err)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Page visited: ", r.Request.URL)
	})

	c.OnHTML("article", func(e *colly.HTMLElement) {
		fmt.Println("Title: ", e.ChildText("div > div h2"))
		fmt.Println("Price: ", e.ChildText("div > div > div > span"))
		fmt.Println("Kilopris: ", e.ChildText("div > div > div > p"))
		fmt.Println("Bildelenke: ", e.ChildAttr("div > img", "src"))
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println(r.Request.URL, " scraped!")
	})

	n := colly.NewCollector()

	n.OnHTML("main", func(e *colly.HTMLElement) {
		fmt.Println("Antall produkter:  ", e.ChildText("main > div > div > div > span > div > div > a.k-choice-chip--selected.k-choice-chip--primary > span.k-pill--extra-small"))
		tall := strings.Split(e.ChildText("main > div > div > div > span > div > div > a.k-choice-chip--selected.k-choice-chip--primary > span.k-pill--extra-small"), "")
		tallLengde := len(tall) / 2
		var s []string
		for i := 0; i < tallLengde; i++ {
			s = append(s, tall[i])
		}
		str := strings.Join(s, "")
		antallVarer, err := strconv.ParseFloat(str, 64)
		if err != nil {
			fmt.Println("Oh no")
		}
		underCategoryLink := "/no/categories/1135-bakeri-og-konditori/1136-handverksbakeri/?filters="
		antallVarerDelt := antallVarer / 24
		antallSider := int(math.Ceil(float64(antallVarerDelt)))
		fmt.Println(antallSider)
		link := fmt.Sprintf("https://oda.com%s&cursor=%v", underCategoryLink, antallSider)
		c.Visit(link)
	})
	n.Visit("https://oda.com/no/categories/1135-bakeri-og-konditori/1136-handverksbakeri/?filters=&cursor=2")
}
