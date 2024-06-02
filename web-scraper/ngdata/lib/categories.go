package lib

import (
	"github.com/gocolly/colly"
)

func getCategories(store string, categories *[]Category) {
	c := colly.NewCollector()

	categoriesStarted := false

	c.OnHTML(storeData[store].targetClass, func(e *colly.HTMLElement) {
		categoryName := e.ChildText("a span")

		if categoryName == storeData[store].firstCategory {
			categoriesStarted = true
		}

		if categoriesStarted {
			// lager instans av kategori med alle verdier jeg har til nå
			*categories = append(*categories, Category{Name: categoryName, Store: store})
		}
	})

	c.Visit(storeData[store].url)
}
