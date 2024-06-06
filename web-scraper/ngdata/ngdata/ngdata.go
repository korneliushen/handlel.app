package ngdata

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/gocolly/colly"
	"github.com/korneliushen/handlel.app/meny/lib"
)

// separerer kategori og sub-kategori
const BASE_URL = "https://platform-rest-prod.ngdata.no/api/products"
const OPTIONS = "?page=1&page_size=10000&full_response=true&fieldset=maximal&facets=Category&facet=Categories:"

var Stores = [3]string{"meny", "joker", "spar"}

// data som trengs rundt-om-kring i applikasjonen, brukes for å kunne reuse funksjoner
var StoreInfo = map[string]struct {
	targetClass   string
	firstCategory string
	Url           string
	id            string
}{
	"meny":  {targetClass: "li.cw-categories__item", firstCategory: "Frukt & grønt", Url: "https://meny.no/varer/", id: "/1300/7080001150488"},
	"joker": {targetClass: "li.product-categories__item", firstCategory: "Bakerivarer", Url: "https://joker.no/nettbutikk/varer/", id: "/1220/7080001395933"},
	"spar":  {targetClass: "li.product-categories__item", firstCategory: "Bakeartikler og kjeks", Url: "https://spar.no/nettbutikk/varer/", id: "/1210/7080001097950"},
}

func GetCategories(store string, categories *[]lib.Category) {
	c := colly.NewCollector()

	categoriesStarted := false

	c.OnHTML(StoreInfo[store].targetClass, func(e *colly.HTMLElement) {
		categoryName := e.ChildText("a span")

		if categoryName == StoreInfo[store].firstCategory {
			categoriesStarted = true
		}

		if categoriesStarted {
			// lager instans av kategori med alle verdier jeg har til nå
			*categories = append(*categories, lib.Category{Name: categoryName, Store: store})
		}
	})

	c.Visit(StoreInfo[store].Url)
}

func GetProducts(store, category string) ([]lib.ApiProduct, error) {
	// bare meny funker helt for nå
	url := getUrl(store, category)

	// får data om produkter fra api-en
	data, err := fetchProducts(url)
	if err != nil {
		return []lib.ApiProduct{}, err
	}

	return data.Hits.Products, nil
}

// genererer en url for norgesgruppen api med butikk id, og kategori
func getUrl(store, category string) string {
	// bruker QueryEscape for å gjøre at man kan putte kategorien i url-en for fetch requesten
	queryCategory := url.QueryEscape(category)

	// får id fra ids map
	id := StoreInfo[store].id

	// constructer url
	apiUrl := BASE_URL + id + OPTIONS + queryCategory

	return apiUrl
}

// henter data for produkter med url-en som blir generert over
func fetchProducts(url string) (lib.ApiResponse, error) {
	// gjør request til url-en
	res, err := http.Get(url)
	if err != nil {
		return lib.ApiResponse{}, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return lib.ApiResponse{}, err
	}

	var produkter lib.ApiResponse
	err = json.Unmarshal(body, &produkter)
	if err != nil {
		return lib.ApiResponse{}, err
	}

	return produkter, nil
}
