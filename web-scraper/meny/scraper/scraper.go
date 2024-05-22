package scraper

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/gocolly/colly"
)

func GetCategories() Kategorier {
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

			// gjør om navnet til et format som kan brukes i query for api kall
			categoryNameQueryFormat := url.QueryEscape(categoryName)

			fmt.Println("Henter underkategorier for:", categoryName)

			// lager instans av kategori med alle verdier jeg har til nå
			category := Kategori{Navn: categoryName, Link: categoryLink, QueryFormatNavn: categoryNameQueryFormat}

			n := colly.NewCollector()

			// henter underkategorier
			n.OnHTML(".cw-categories__item--active ul li", func(h *colly.HTMLElement) {
				subCategoryName := h.ChildText("a span")
				subCategoryLink := fmt.Sprintf("htts://meny.no%s", h.ChildAttr("a.cw-categories__title", "href"))

				// gjør om til format som kan brukes i query
				subCategoryNameQueryFormat := url.QueryEscape(subCategoryName)

				// lager instans av underkategori
				subCategory := Underkategori{Navn: subCategoryName, Link: subCategoryLink, QueryFormatNavn: subCategoryNameQueryFormat}

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
