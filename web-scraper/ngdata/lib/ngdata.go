package lib

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

// separerer kategori og sub-kategori
const BASE_URL = "https://platform-rest-prod.ngdata.no/api/products"
const OPTIONS = "?page=1&page_size=10000&full_response=true&fieldset=maximal&facets=Category&facet=Categories:"

var stores = [3]string{"meny", "joker", "spar"}

// data som trengs rundt-om-kring i applikasjonen, brukes for å kunne reuse funksjoner
var storeData = map[string]struct {
	targetClass   string
	firstCategory string
	url           string
	id            string
}{
	"meny":  {targetClass: "li.cw-categories__item", firstCategory: "Frukt & grønt", url: "https://meny.no/varer/", id: "/1300/7080001150488"},
	"joker": {targetClass: "li.product-categories__item", firstCategory: "Bakerivarer", url: "https://joker.no/nettbutikk/varer/", id: "/1220/7080001395933"},
	"spar":  {targetClass: "li.product-categories__item", firstCategory: "Bakeartikler og kjeks", url: "https://spar.no/nettbutikk/varer/", id: "/1210/7080001097950"},
}

func getProducts(store, category string) ([]ApiProduct, error) {
	// bare meny funker helt for nå
	url := getUrl(store, category)

	// får data om produkter fra api-en
	data, err := fetchProducts(url)
	if err != nil {
		return []ApiProduct{}, err
	}

	return data.Hits.Products, nil
}

// genererer en url for norgesgruppen api med butikk id, og kategori
func getUrl(store, category string) string {
	// bruker QueryEscape for å gjøre at man kan putte kategorien i url-en for fetch requesten
	queryCategory := url.QueryEscape(category)

	// får id fra ids map
	id := storeData[store].id

	// constructer url
	apiUrl := BASE_URL + id + OPTIONS + queryCategory

	return apiUrl
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
