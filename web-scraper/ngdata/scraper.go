package main

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly"
)

// TODO: legg til en timeout på scraper
const TIMEOUT = 10

func getCategories() Kategorier {
	fmt.Println("Getting categories")

	// kategorier instans
	categories := Kategorier{}

	c := colly.NewCollector()

	c.OnHTML("li.cw-categories__item", func(e *colly.HTMLElement) {
		// henter alle kategorier
		categoryLink := e.ChildAttr("a", "href")
		// bruker bare kategorien om linken inneholder /?c (den får ellers med noe annet på siden som ikke trengs)
		if strings.Contains(categoryLink, "/?c=") {
			categoryName := e.ChildText("a span")
			categoryLink = fmt.Sprintf("https://meny.no%s", categoryLink)

			// lager instans av kategori med alle verdier jeg har til nå
			category := Kategori{Navn: categoryName, Link: categoryLink}

			n := colly.NewCollector()

			// henter underkategorier
			n.OnHTML(".cw-categories__item--active ul li", func(h *colly.HTMLElement) {
				subCategoryName := h.ChildText("a span")
				subCategoryLink := fmt.Sprintf("htts://meny.no%s", h.ChildAttr("a.cw-categories__title", "href"))

				// lager instans av underkategori
				subCategory := Underkategori{Navn: subCategoryName, Link: subCategoryLink}

				// legger til underkategorien i Underkategorier feltet til kategori instansen vi lagde over
				category.Underkategorier = append(category.Underkategorier, subCategory)
			})

			n.Visit(categoryLink)

			// legger til kategori instansen i Kategorier instansen
			categories.Kategorier = append(categories.Kategorier, category)
		}
	})

	c.Visit("https://meny.no/varer/")

	return categories
}
