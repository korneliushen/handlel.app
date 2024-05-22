package fetch

import (
	"fmt"
	"net/url"
	"strings"
)

// separerer kategori og sub-kategori
const SEPERATOR = "%3BShoppingListGroups%3A"
const (
	MENY_BASE_URL  = "https://platform-rest-prod.ngdata.no/api/products/1300/7080001150488?page=1&page_size=1000&full_response=true&fieldset=maximal&facets=Category%2CAllergen&facet=Categories%3A"
	JOKER_BASE_URL = "https://platform-rest-prod.ngdata.no/api/products/1220/7080001395933?page=1&page_size=1000&full_response=true&fieldset=maximal&facets=Category%2CAllergen&facet=Categories%3A"
	SPAR_BASE_URL  = "https://platform-rest-prod.ngdata.no/api/products/1210/7080001097950?page=1&page_size=1000&full_response=true&fieldset=maximal&facets=Category%2CAllergen&facet=Categories%3A"
)

// enum
var baseUrls = map[string]string{
	"meny":  MENY_BASE_URL,
	"joker": JOKER_BASE_URL,
	"spar":  SPAR_BASE_URL,
}

func getUrl(shop string, category string, subCategory string) string {
	baseUrl := baseUrls[shop]

	// joker trenger stor bokstav på starten av alle ord i kategori eller underkategori
	if shop == "joker" {
		category = strings.Title(category)
		subCategory = strings.Title(subCategory)
	}

	// gjør kategori og underkategori mulig å putte inn i query
	category = url.QueryEscape(category)
	subCategory = url.QueryEscape(subCategory)

	url := baseUrl + category + SEPERATOR + subCategory

	// joker aksepterer bare denne versjonen av "&" av en eller annen grunn
	if shop == "joker" {
		url = strings.ReplaceAll(url, "+%26+", "%2F")
	}

	// spar bruker "og" i stedet for "&", meget hyggelig
	if shop == "spar" {
		url = strings.ReplaceAll(url, "+%26+", "%20og%20")
	}

	return url
}

func GetProducts(category string, subCategory string) {
	// bare meny funker helt for nå
	url := getUrl("spar", category, subCategory)
	fmt.Println(url)
}
