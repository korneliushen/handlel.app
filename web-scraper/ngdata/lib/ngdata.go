package lib

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/gocolly/colly"
)

// separerer kategori og sub-kategori
const BASE_URL = "https://platform-rest-prod.ngdata.no/api/products"
const OPTIONS = "?page=1&page_size=1000&full_response=true&fieldset=maximal&facets=Category%2CAllergen&facet=Categories%3A"
const SEPERATOR = "%3BShoppingListGroups%3A"

var stores = [3]string{"meny", "joker", "spar"}

// enum med id til ngdata butikkene
var storeIds = map[string]string{
	"meny":  "/1300/7080001150488",
	"joker": "/1220/7080001395933",
	"spar":  "/1210/7080001097950",
	// "kiwi": "/1100/",
}

func getCategories() Categories {
	// kategorier instans
	categories := Categories{}

	c := colly.NewCollector()

	c.OnHTML("li.cw-categories__item", func(e *colly.HTMLElement) {
		// henter alle kategorier
		categoryLink := e.ChildAttr("a", "href")
		// bruker bare kategorien om linken inneholder /?c (den får ellers med noe annet på siden som ikke trengs)
		if strings.Contains(categoryLink, "/?c=") {
			categoryName := e.ChildText("a span")
			categoryLink = fmt.Sprintf("https://meny.no%s", categoryLink)

			// lager instans av kategori med alle verdier jeg har til nå
			category := Category{Name: categoryName, Link: categoryLink}

			n := colly.NewCollector()

			// henter underkategorier
			n.OnHTML(".cw-categories__item--active ul li", func(h *colly.HTMLElement) {
				subCategoryName := h.ChildText("a span")
				subCategoryLink := fmt.Sprintf("htts://meny.no%s", h.ChildAttr("a.cw-categories__title", "href"))

				// lager instans av underkategori
				subCategory := SubCategories{Name: subCategoryName, Link: subCategoryLink}

				// legger til underkategorien i Underkategorier feltet til kategori instansen vi lagde over
				category.SubCategories = append(category.SubCategories, subCategory)
			})

			n.Visit(categoryLink)

			// legger til kategori instansen i Kategorier instansen
			categories.Category = append(categories.Category, category)
		}
	})

	c.Visit("https://meny.no/varer/")

	return categories
}

func getProducts(store string, category string, subCategory string) (ApiResponse, error) {
	// bare meny funker helt for nå
	url := getUrl(store, category, subCategory)

	// får data om produkter fra api-en
	data, err := fetchProducts(url)
	if err != nil {
		return ApiResponse{}, err
	}

	return data, nil
}

// genererer en url for norgesgruppen api med butikk, kategori og underkategori
func getUrl(store, category, subCategory string) string {
	category = getCategoryName(store, category)

	// joker trenger stor bokstav på starten av alle ord i kategori eller underkategori
	if store == "joker" {
		category = strings.Title(category)
		subCategory = strings.Title(subCategory)
	}

	// gjør kategori og underkategori mulig å putte inn i query
	category = url.QueryEscape(category)
	subCategory = url.QueryEscape(subCategory)

	// får id fra ids map
	id := storeIds[store]

	// constructer url
	url := transformUrl(store, BASE_URL+id+OPTIONS+category+SEPERATOR+subCategory)

	return url
}

// TODO: vil også finne en måte å automatisere dette på
// noen kategorier er annerledes på joker, spar og meny (men ikke så mange og underkategorier er for det meste det samme)
// denne mappen brukes i getCategoryName for å endre kategorinavn om det er nødvendig
var categoryMappings = map[string]map[string]string{
	"joker": {
		"Meieri & egg":       "Meieriprodukter",
		"Bakeri":             "Bakerivarer",
		"Drikke":             "Drikkevarer",
		"Pålegg & frokost":   "Frokost/Pålegg",
		"Dyr":                "Dyreprodukter",
		"Hus & hjem":         "Hus/Hjem Artikler",
		"Bakevarer og kjeks": "Kaker/Bakevarer",
	},
	"spar": {
		"Bakevarer og kjeks": "Bakeriartikler og kjeks",
		"Dessert og iskrem":  "Dessert",
		"Barneprodukter":     "Barn",
		"Bakeri":             "Brød og bakevarer",
	},
}

func getCategoryName(store, category string) string {
	if mappings, ok := categoryMappings[store]; ok {
		if mappedCategory, exists := mappings[category]; exists {
			return mappedCategory
		}
	}
	return category
}

func transformUrl(store, url string) string {
	// joker bruker "/" i stedet for "&"
	if store == "joker" {
		url = strings.ReplaceAll(url, "+%26+", "%2F")
	}
	// spar bruker "og" i stedet for "&", meget hyggelig
	if store == "spar" {
		url = strings.ReplaceAll(url, "+%26+", "%20og%20")
	}
	return url
}

// henter data for produkter med url-en som blir generert over
func fetchProducts(url string) (ApiResponse, error) {
	// gjør request til url-en
	res, err := http.Get(url)
	if err != nil {
		return ApiResponse{}, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return ApiResponse{}, err
	}

	var produkter ApiResponse
	err = json.Unmarshal(body, &produkter)
	if err != nil {
		return ApiResponse{}, err
	}

	return produkter, nil
}
