package scraper

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly"
)

func Scrape() {
	c := colly.NewCollector()

	c.OnHTML("li.cw-categories__item", func(e *colly.HTMLElement) {
		// henter alle kategorier
		category := e.ChildAttr("a", "href")
		if strings.Contains(category, "/?c=") {
			fmt.Println("Kategori:", category)

			n := colly.NewCollector()

			// henter underkategorier
			fmt.Println("Underkategorier:")
			n.OnHTML(".cw-categories__item--active ul li", func(h *colly.HTMLElement) {
				subCategory := h.ChildAttr("a.cw-categories__title", "href")
				fmt.Println(subCategory)
			})

			link := fmt.Sprintf("https://meny.no%s", category)
			n.Visit(link)
		}
	})

	c.Visit("https://meny.no/varer/")
}
