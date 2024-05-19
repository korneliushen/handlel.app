package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

func scrape() {
	c := colly.NewCollector()

	categories := Kategorier{}

	// besøker oda og henter alle kategorier
	// for hver kategori, kjører getUnderCategories med kategori linken man får
	c.OnHTML("article.k-p-4", func(e *colly.HTMLElement) {
		start := time.Now()
		category := &Kategori{}

		categoryName := e.ChildText("div > a > h1")

		if categoryName == "" {
			return
		}

		category.Navn = categoryName

		fmt.Println("Getting data for: ", categoryName)

		categoryLink := e.ChildAttr("div > a", "href")
		getUnderCategories(categoryName, categoryLink, category)

		categories.Kategorier = append(categories.Kategorier, *category)

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

	writeData(categories)
}

func getUnderCategories(categoryName string, categoryLink string, category *Kategori) {
	c := colly.NewCollector()

	// finner hver under kategori
	// Eksempel: oda.com/no/categories/20-frukt-og-gront og så Frukt, Bær, osv.
	// for hver underkategori, henter mengden sider den underkategorien har (cursor)
	c.OnHTML("section", func(e *colly.HTMLElement) {
		// lager instans av under kategori
		underCategory := &Underkategori{} // Use a pointer to Underkategori

		// gir instansen et navn
		underCategoryName := e.ChildText("section > a > h2")

		if underCategoryName == "" {
			return
		}

		underCategory.Navn = underCategoryName

		fmt.Println("Getting data for undercategory: ", underCategoryName, "in category", categoryName)

		underCategoryLink := e.ChildAttr("section > a", "href")
		getProducts(underCategoryLink, underCategory)

		category.Underkategorier = append(category.Underkategorier, *underCategory)
	})

	link := fmt.Sprintf("https://oda.com%s", categoryLink)

	c.Visit(link)
}

// får mengden sider for å vite hvor mange sider som må besøkes
// kjører så getProductInfo() på mengden sider -> får alle produkter på alle sider
func getProducts(underCategoryLink string, underCategory *Underkategori) {
	c := colly.NewCollector()

	c.OnHTML("main", func(e *colly.HTMLElement) {
		// finner mengden sider
		pageCountSlice := strings.Split(e.ChildText("main > div > div > div > span > div > div > a.k-choice-chip--selected.k-choice-chip--primary > span.k-pill--extra-small"), "")
		pageCount := getPageCount(pageCountSlice)

		// lager en link for underkategorien
		link := fmt.Sprintf("https://oda.com%s", underCategoryLink)

		// for hver side, hent produktinfo for alle produktene på siden
		for i := 0; i < pageCount; i++ {
			getProductInfo(link, i+1, underCategory)
		}
	})
	link := fmt.Sprintf("https://oda.com%s", underCategoryLink)
	c.Visit(link)
}

func getProductInfo(link string, cursor int, underCategory *Underkategori) {
	// lager en instans av Produkter
	c := colly.NewCollector()

	// henter data for hvert element (printer bare for nå)
	c.OnHTML("article", func(e *colly.HTMLElement) {
		// lager en instans av Produkt, som har all infoen om produktet
		productInfo := Produkt{}
		// bruker en referanse verdi her for at det skal funke med reflect
		contents := &Innhold{}
		nutritionalContents := &Næringsinnhold{}

		title := e.ChildText("div > div h2")

		if title == "" {
			return
		}

		fmt.Println("Getting data for", title)

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

			// oppbevaring har annerledes oppsett enn annet innhold, basically bare en extra div og så en span
			if strings.Contains(key, "Oppbevaring") {
				key = "Oppbevaring"
				value = h.ChildText("div > div span")
			}

			// ingredienser er satt opp annerledes avhengig av hvor mange ingredienser det er, så det accountes for her
			if strings.Contains(key, "Ingredienser") {
				key = "Ingredienser"
				value = getIngredients(h, value)
			}

			setFieldValue(contents, nutritionalContents, key, value)
		})

		n.Visit(dataLink)

		// legger næringsinnhold til i innhold
		contents.Næringsinnhold = *nutritionalContents
		// legger innhold til i produktinfo
		productInfo.Innhold = *contents
		// setter produktet inn i produkter arrayet
		underCategory.Produkter = append(underCategory.Produkter, productInfo)
	})

	visitLink := fmt.Sprintf("%s&cursor=%v", link, cursor)
	c.Visit(visitLink)
}
