package ngdata

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/gocolly/colly"
	"github.com/korneliushen/handlel.app/scraper/lib"
	"github.com/korneliushen/handlel.app/scraper/model"
)

var stores = [3]string{"meny", "joker", "spar"}

// data som trengs rundt-om-kring i applikasjonen, brukes for å kunne reuse funksjoner
var storeInfo = map[string]struct {
	targetClass   string
	firstCategory string
	url           string
	imageUrl      string
	id            string
}{
	"meny":  {targetClass: "li.cw-categories__item", firstCategory: "Frukt & grønt", url: "https://meny.no/varer/", imageUrl: "https://bilder.ngdata.no", id: "/1300/7080001150488"},
	"joker": {targetClass: "li.product-categories__item", firstCategory: "Bakerivarer", url: "https://joker.no/nettbutikk/varer/", imageUrl: "https://bilder.ngdata.no", id: "/1220/7080001395933"},
	"spar":  {targetClass: "li.product-categories__item", firstCategory: "Bakeartikler og kjeks", url: "https://spar.no/nettbutikk/varer/", imageUrl: "https://bilder.ngdata.no", id: "/1210/7080001097950"},
}

func GetCategories() model.Categories {
	categories := model.Categories{}
	for _, store := range stores {
		c := colly.NewCollector()

		categoriesStarted := false

		c.OnHTML(storeInfo[store].targetClass, func(e *colly.HTMLElement) {
			categoryName := e.ChildText("a span")

			if categoryName == storeInfo[store].firstCategory {
				categoriesStarted = true
			}

			if categoriesStarted {
				// lager instans av kategori med alle verdier jeg har til nå
				categories.Categories = append(categories.Categories, model.Category{
					Name: categoryName, Store: store,
				})
			}
		})

		c.Visit(storeInfo[store].url)
	}

	return categories
}

func GetProductsFromApi(category model.Category, store string) ([]model.BaseProduct, error) {
	// bare meny funker helt for nå
	url := getUrl(store, category.Name)

	// får data om produkter fra api-en
	data, err := fetchProducts(url)
	if err != nil {
		return []model.BaseProduct{}, err
	}

	return data.Hits.Products, nil
}

// genererer en url for norgesgruppen api med butikk id, og kategori
func getUrl(store, category string) string {
	// bruker QueryEscape for å gjøre at man kan putte kategorien i url-en
	// for fetch requesten
	queryCategory := url.QueryEscape(category)

	// får id fra ids map
	id := storeInfo[store].id

	// constructer url, base url + id-en til butikken + options
	// (denne kan endres på for annen data) + kategorien (dette kan også endres
	// på avhengig av options)
	apiUrl := "https://platform-rest-prod.ngdata.no/api/products" +
		id +
		"?page=1&page_size=10000&full_response=true&fieldset=maximal&facets=Category&facet=Categories:" +
		queryCategory

	return apiUrl
}

// henter data for produkter med url-en som blir generert over
func fetchProducts(url string) (model.ApiResponse, error) {
	// gjør request til url-en
	res, err := http.Get(url)
	if err != nil {
		return model.ApiResponse{}, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return model.ApiResponse{}, err
	}

	var produkter model.ApiResponse
	err = json.Unmarshal(body, &produkter)
	if err != nil {
		return model.ApiResponse{}, err
	}

	return produkter, nil
}

func GetProducts(products *model.BaseProducts, categories model.Categories) {
	for _, category := range categories.Categories {
		for _, store := range stores {
			// om kategorien sin butikk og butikken ikke er den samme, er det ikke
			// vits å kjøre request fordi den vil ikke få noe data
			// (og om den får det vil det være duplicate)
			if category.Store != store {
				continue
			}

			// får data om alle produkter i kategorien
			res, err := GetProductsFromApi(category, store)
			if err != nil {
				fmt.Printf("Error getting products from %s in category %s: %v\n",
					store, category, err)
				continue
			}

			// legger til produktet i baseProduct array som mappes over senere,
			// legger også til Store (for senere bruk)
			for _, product := range res {
				// legger til underkategorier, om underkategorien ikke er lagt til
				// underkategorier er jeg ganske sikker på at er basically helt likt
				// på alle sidene, så det vil ikke være duplicates med forskjellig
				// navn, om det er annerledes må jeg bytte til id approach
				if !lib.IsIn(product.Data.SubCategory, category.SubCategories) {
					category.SubCategories = append(
						category.SubCategories, product.Data.SubCategory,
					)
				}

				baseImgUrl := "https://bilder.ngdata.no/"
				// Lager small, medium og large versjoner av image
				product.Data.ImageLinkSmall = fmt.Sprintf("%s%s%s",
					baseImgUrl, product.Data.ImageLink, "/small.jpg")
				product.Data.ImageLinkMedium = fmt.Sprintf("%s%s%s",
					baseImgUrl, product.Data.ImageLink, "/medium.jpg")
				product.Data.ImageLinkLarge = fmt.Sprintf("%s%s%s",
					baseImgUrl, product.Data.ImageLink, "/large.jpg")

				*products = append(*products, product.Extend(store, storeInfo[store].url))
			}
		}
	}
}
