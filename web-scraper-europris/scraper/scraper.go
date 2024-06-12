package scraper

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

func Scrape() {
	ignoredCategories := []string{"Batterier", "Kontorrekvisita", "Papir", "Engangs", "Poser & folie", "Brygge selv", "Lyspærer", "Rengjøring"}

	c := colly.NewCollector()

	categories := Kategorier{}


	// besøker oda og henter alle kategorier
	// for hver kategori, kjører getUnderCategories med kategori linken man får
	c.OnHTML("span.category-navigation__button", func(e *colly.HTMLElement) {
		start := time.Now()
		category := &Kategori{}

		categoryName := e.ChildText("a")

		if categoryName == "" || isIgnored(categoryName, ignoredCategories) {
			return
		}

		category.Navn = categoryName

		fmt.Println("Getting data for: ", categoryName)

		categoryLink := e.ChildAttr("a", "href")
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

	c.Visit("https://www.europris.no/dagligvarer/")

	err := WriteData(categories, "../data.json")
	if err != nil {
		fmt.Printf("Error writing data to json: %v\n", err)
	}
}


func isIgnored(categoryName string, ignoredCategories []string) bool {
    for _, ignoredCategory := range ignoredCategories {
        if categoryName == ignoredCategory {
            return true
        }
    }
    return false
}

func getUnderCategories(categoryName string, categoryLink string, category *Kategori) {
    c := colly.NewCollector()

    c.OnHTML("div.Kategori:nth-child(2)", func(e *colly.HTMLElement) {
        // Find all subcategory elements
        e.ForEach("li.not-search > a", func(_ int, elem *colly.HTMLElement) {
            // Extract subcategory name and link
            trimmedName := strings.TrimSpace(elem.Text)
            if trimmedName == "" {
                return
            }

            underCategory := &Underkategori{Navn: trimmedName}
            fmt.Println("Getting data for undercategory:", trimmedName, "in category", categoryName)

            underCategoryLink := elem.Attr("href") // Retrieve the href attribute of the current element
            getProducts(underCategoryLink, underCategory)

            category.Underkategorier = append(category.Underkategorier, *underCategory)
        })
    })

    link := fmt.Sprintf("%s", categoryLink)
    c.Visit(link)
}



// får mengden sider for å vite hvor mange sider som må besøkes
// kjører så getProductInfo() på mengden sider -> får alle produkter på alle sider
func getProducts(underCategoryLink string, underCategory *Underkategori) error {
    c := colly.NewCollector()
    link := fmt.Sprintf("%s", underCategoryLink)

    // Define a function to handle the "next page" button
    var visitNextPage func(string)
    visitNextPage = func(pageLink string) {
        // Visit the current page and extract data
        c.OnHTML("div.pages", func(e *colly.HTMLElement) {
            fmt.Println("visiting page:", pageLink)

            // Extract page information here if needed

            // Call getProductInfo for each page if needed
            // This is where you can call getProductInfo or perform other actions

            // Find the "next page" button
            nextPageButton := e.DOM.Find("a.action.next")
            if nextPageButton.Length() > 0 {
                // If the "next page" button exists, visit the next page
                nextPageURL, _ := nextPageButton.Attr("href")
                visitNextPage(nextPageURL)
            }
        })
        c.Visit(pageLink)
    }

    // Start the pagination process
    visitNextPage(link)

    return nil
}



// henter data til hvert produkt på siden og legger data inn i underCategory struct instansen
func GetProductInfo(link string, cursor int, underCategory *Underkategori) error {
	// lager en err variabel, som kan få en value i onHTML
	// om den får en value, blir den returnert, om den forblir nil, blir det returnert
	var err error

	// lager en instans av Produkter
	c := colly.NewCollector()

	c.OnHTML("div.product-info-main", func(e *colly.HTMLElement) {
		// lager en instans av Produkt, som har all infoen om produktet
		productInfo := Produkt{}
		// bruker en referanse verdi her for at det skal funke med reflect
		// contents := &Innhold{}
		// nutritionalContents := &Næringsinnhold{}

		title := e.ChildText("h1.page-title > span")

		if title == "" {
			err = errors.New("No title, skipping")
			return
		}

		fmt.Println("Getting data for", title)

		// legg til alle på produkt structen
		productInfo.Tittel = title
		// productInfo.Undertittel
		productInfo.FullNavn = productInfo.Tittel
		// productInfo.Pris = e.ChildText("div > div > div > span")
		// productInfo.KiloPris = e.ChildText("div > div > div > p")
		productInfo.BildeLink = e.ChildAttr("div > img", "src")

		// link for side med næringsinnhold
		// productLink := e.ChildAttr("div > div a", "href")
		// dataLink := fmt.Sprintf("https://oda.com%s", productLink)

		// // henter data fra produkt siden
		// n := colly.NewCollector()

		// n.OnHTML("#read-more-product-description-nb", func(h *colly.HTMLElement) {
		// 	// henter beskrivelsen
		// 	description := h.ChildText("p")
		// 	contents.Beskrivelse = description
		// })

		// // henter innhold og næringsinnhold (om det eksisterer)
		// n.OnHTML("div.k-grid.k-pt-3.k-pb-6", func(h *colly.HTMLElement) {
		// 	key := h.ChildText("div > span")
		// 	value := h.ChildText("div > p")

		// 	// oppbevaring har annerledes oppsett enn annet innhold, basically bare en extra div og så en span
		// 	if strings.Contains(key, "Oppbevaring") {
		// 		key = "Oppbevaring"
		// 		value = h.ChildText("div > div span")
		// 	}

		// 	// ingredienser er satt opp annerledes avhengig av hvor mange ingredienser det er, så det accountes for her
		// 	if strings.Contains(key, "Ingredienser") {
		// 		key = "Ingredienser"
		// 		value = getIngredients(h, value)
		// 	}

		// 	setFieldValue(contents, nutritionalContents, key, value)
		// })

		// n.Visit(dataLink)

		// // legger næringsinnhold til i innhold
		// contents.Næringsinnhold = *nutritionalContents
		// // legger innhold til i produktinfo
		// productInfo.Innhold = *contents
		// // setter produktet inn i produkter arrayet
		// underCategory.Produkter = append(underCategory.Produkter, productInfo)
	})

	if err != nil {
		return err
	}

	visitLink := fmt.Sprintf("%s&cursor=%v", link, cursor)
	c.Visit(visitLink)

	return nil
}
