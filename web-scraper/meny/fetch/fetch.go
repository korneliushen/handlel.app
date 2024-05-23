package fetch

import (
	"fmt"
	"net/url"
	"strings"
)

// separerer kategori og sub-kategori
const SEPERATOR = "%3BShoppingListGroups%3A"
const BASE_URL = "https://platform-rest-prod.ngdata.no/api/products"
const OPTIONS = "?page=1&page_size=1000&full_response=true&fieldset=maximal&facets=Category%2CAllergen&facet=Categories%3A"

const (
	MENY_ID  = "/1300/7080001150488"
	JOKER_ID = "/1220/7080001395933"
	SPAR_ID  = "/1210/7080001097950"
)

// enum
var ids = map[string]string{
	"meny":  MENY_ID,
	"joker": JOKER_ID,
	"spar":  SPAR_ID,
}

func getUrl(shop string, category string, subCategory string) string {
	// joker trenger stor bokstav på starten av alle ord i kategori eller underkategori
	if shop == "joker" {
		category = strings.Title(category)
		subCategory = strings.Title(subCategory)
	}

	// gjør kategori og underkategori mulig å putte inn i query
	category = url.QueryEscape(category)
	subCategory = url.QueryEscape(subCategory)

	// får id fra ids map
	id := ids[shop]

	// constructer url
	url := BASE_URL + id + OPTIONS + category + SEPERATOR + subCategory

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

func GetProducts(shop string, category string, subCategory string) {
	// bare meny funker helt for nå
	url := getUrl(shop, category, subCategory)
	fmt.Println(url)
}
