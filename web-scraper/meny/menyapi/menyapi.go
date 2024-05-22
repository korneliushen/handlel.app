package menyapi

import "fmt"

// separerer kategori og sub-kategori
const SEPERATOR = "%3BShoppingListGroups%3A"
const (
	MENY_BASE_URL  = "https://platform-rest-prod.ngdata.no/api/products/1300/7080001150488?page=1&page_size=1000&full_response=true&fieldset=maximal&facets=Category%2CAllergen&facet=Categories%3A"
	JOKER_BASE_URL = "https://platform-rest-prod.ngdata.no/api/products/1210/7080001266110?page=1&page_size=1000&full_response=true&fieldset=maximal&facets="
	SPAR_BASE_URL  = "https://platform-rest-prod.ngdata.no/api/products/1220/7080001420567?page=1&page_size=1000&full_response=true&fieldset=maximal&facets="
)

// enum
var baseUrls = map[string]string{
	"meny":  MENY_BASE_URL,
	"joker": MENY_BASE_URL,
	"spar":  MENY_BASE_URL,
}

func getUrl(shop string, category string, subCategory string) string {
	baseUrl := baseUrls[shop]
	url := baseUrl + category + SEPERATOR + subCategory
	return url
}

func GetProducts(category string, subCategory string) {
	url := getUrl("meny", category, subCategory)
	fmt.Println(url)
}
