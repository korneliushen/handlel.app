package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// separerer kategori og sub-kategori
const BASE_URL = "https://platform-rest-prod.ngdata.no/api/products"
const OPTIONS = "?page=1&page_size=1000&full_response=true&fieldset=maximal&facets=Category%2CAllergen&facet=Categories%3A"
const SEPERATOR = "%3BShoppingListGroups%3A"

// id til butikkene
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

func getProducts(shop string, category string, subCategory string) (ApiResponse, error) {
	// bare meny funker helt for nå
	url := getUrl(shop, category, subCategory)

	// får data om produkter fra api-en
	data, err := fetchProducts(url)
	if err != nil {
		return ApiResponse{}, err
	}

	return data, nil
}

// genererer en url for norgesgruppen api med butikk, kategori og underkategori
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
