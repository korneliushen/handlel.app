package ngdata

import (
	"cmp"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"slices"
	"strings"

	"github.com/gocolly/colly"
	"github.com/korneliushen/handlel.app/meny/lib"
)

//////////////////////////
///// Category(ies) /////
////////////////////////

func (categories *Categories) GetCategories() {
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
				categories.Categories = append(categories.Categories, Category{Name: categoryName, Store: store})
			}
		})

		c.Visit(storeInfo[store].Url)
	}
}

func (category Category) GetProducts(store string) ([]ApiProduct, error) {
	// bare meny funker helt for nå
	url := getUrl(store, category.Name)

	// får data om produkter fra api-en
	data, err := fetchProducts(url)
	if err != nil {
		return []ApiProduct{}, err
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

//////////////////////////
///// ApiProduct(s) /////
////////////////////////

func (products *ApiProducts) GetProducts(categories Categories) {
	for _, category := range categories.Categories {
		for _, store := range stores {
			// om kategorien sin butikk og butikken ikke er den samme, er det ikke
			// vits å kjøre request fordi den vil ikke få noe data
			// (og om den får det vil det være duplicate)
			if category.Store != store {
				continue
			}

			// får data om alle produkter i kategorien
			res, err := category.GetProducts(store)
			if err != nil {
				fmt.Printf("Error getting products from %s in category %s: %v\n",
					store, category, err)
				continue
			}

			// legger til produktet i apiProducts array som mappes over senere,
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
				*products = append(*products, product.Extend(store, storeInfo[store].Url))
			}
			break
		}
		break
	}
}

// oppdaterer data i ApiProduct
func (product ApiProduct) Extend(store, baseUrl string) ApiProduct {
	product.Store = store
	product.BaseUrl = baseUrl
	return product
}

// legger en method til i ApiProduct struct så vi kan accesse apiProduct
func (apiProduct *ApiProduct) FormatData(productData []ApiProduct, products *Products) {
	product := Product{}

	// for algolia
	product.ObjectID = apiProduct.Data.Ean

	product.Id = apiProduct.Data.Ean

	// legger til alle fields fra ApiProduct som har samme navn som Product
	vDest := reflect.ValueOf(&product).Elem()
	vSrc := reflect.ValueOf(&apiProduct.Data).Elem()
	for i := range vDest.NumField() {
		fieldDest := vDest.Field(i)
		fieldSrc := vSrc.FieldByName(vDest.Type().Field(i).Name)

		if fieldSrc.IsValid() && fieldSrc.Type() == fieldDest.Type() {
			fieldDest.Set(fieldSrc)
		}
	}

	// lager en string for vekt, med value og unit
	product.Weight = fmt.Sprintf("%v%s",
		apiProduct.Data.Weight, apiProduct.Data.WeightMeasurementType)

	// lager hele url-en for bildelinker for ulike størrelser
	product.ImageLink = fmt.Sprintf("%s%s",
		"https://bilder.ngdata.no/", apiProduct.Data.ImageLink)

	// fikser kategori navn
	// (bruker hard-coda kategori navn for å gjøre ting til samme kategori)
	product.Category = getCorrectCategoryName(apiProduct.Data.Category)

	// lager et array av priser, å gjøre det på denne måten gjør det lettere
	// når dataen skal sendes til database
	var prices []Price
	storeMap := map[string]bool{}
	// sjekker at prisen ikke er 0, om den er det er det ikke vits å sende til
	// databasen
	for _, product := range productData {
		if _, exists := storeMap[product.Store]; exists {
			continue
		}
		storeMap[product.Store] = true
		prices = append(prices, Price{
			Store:         product.Store,
			Price:         product.Data.Price,
			OriginalPrice: product.Data.OriginalPrice,
			UnitPrice:     product.Data.ComparePricePerUnit,
			Url:           fmt.Sprintf("%s%s", product.BaseUrl, product.Data.Slug),
		})
	}

	// sorterer basert på pris, så det første elementet i arrayet vil være det
	// billigste
	priceCmp := func(a, b Price) int {
		return cmp.Compare(a.Price, b.Price)
	}
	slices.SortFunc(prices, priceCmp)
	product.Prices = prices

	// mapper over allergener array som vi fikk fra databasen
	// i databasen så bestemmer koden hva itemet i arrayen betyr for produktet
	// om koden er JA, blir det lagt til i allergens, om det er kan blir det
	// lagt til i mayContainTracesOf
	var allergens []string
	var mayContainTracesOf []string
	for _, allergen := range apiProduct.Data.Allergens {
		if allergen.Code == "JA" {
			allergens = append(allergens, allergen.Name)
		} else if allergen.Code == "KAN" {
			mayContainTracesOf = append(mayContainTracesOf, allergen.Name)
		}
	}
	product.Allergens = strings.Join(allergens, ", ")
	product.MayContainTracesOf = strings.Join(mayContainTracesOf, ", ")

	// næringsinnhold
	nutritionalContentData := apiProduct.Data.NutritionalContent

	// om det ikke er noe næringsinnhold
	if len(nutritionalContentData) == 0 {
		product.NutritionalContent = nil
	} else {
		nutritionalContent := NutritionalContent{}
		v := reflect.ValueOf(&nutritionalContent).Elem()
		// legger til næringsinnhold data i fields med navn som matcher dataen
		// fra api-en (reflect)
		for i := range len(nutritionalContentData) {
			field := v.FieldByName(nutritionalContentData[i].Name)
			if field.CanSet() {
				// value er bygd opp av mengde og enhet (xxxg)
				field.SetString(fmt.Sprintf("%v%s",
					nutritionalContentData[i].Amount, nutritionalContentData[i].Unit))
			}
		}
		product.NutritionalContent = &nutritionalContent
	}

	// legger til produktet i products array
	*products = append(*products, product)
}

// temporary løsning til jeg finner ut hvordan vi kan automatisere det
// lager et array med kategorinavn som er forskjellige, som så blir brukt
// etterpå for å gjøre alt i samme kategori til samme kategori navn
// MENY, JOKER, SPAR
var categoryNames = [...][3]string{
	// ting som basically kan automatiseres
	{"Frukt & grønt", "Frukt/Grønt", "Frukt og grønt"},
	{"Fisk & skalldyr", "Fisk/Skalldyr", "Fisk og skalldyr"},
	{"Personlige artikler", "Personlige Artikler"},

	// ting som er mobbing (må kjøre noe contains greier ig)
	{"Bakeri", "Bakerivarer", "Brød og bakervarer"},
	{"Meieri & egg", "Meieriprodukter", "Meieri og egg"},
	{"Pålegg & frokost", "Frokost/Pålegg", "Pålegg og frokost"},
	{"Drikke", "Drikkevarer"},
	{"Dessert og iskrem", "Dessert"},
	{"Bakevarer og kjeks", "Kaker/Bakevarer", "Bakeartikler og kjeks"},
	{"Barneprodukter", "Barn"},
	{"Hus & hjem", "Hus/Hjem Artikler", "Hus og hjem"},
	{"Snacks & godteri", "Snacks/Godteri", "Snacks og godteri"},
	{"Dyr", "Dyreprodukter"},
}

// bruker slicen av slices over til å gjøre alt til samme kategori navn
func getCorrectCategoryName(category string) string {
	for i := range categoryNames {
		for j := range categoryNames[i] {
			if categoryNames[i][j] == category {
				return categoryNames[i][0]
			}
		}
	}
	return category
}

//////////////////////////
///// Product(s) ////////
////////////////////////

func (products *Products) Format(apiProducts ApiProducts) {
	// lagrer alle produkter som allerede har blitt sjekket i et array
	// da blir det ikke duplicates og vi kan returnere tidlig om produktet
	// allerede er ferdig
	var checkedGtins []string

	// mapper over alle produkter vi har fått fra databasen og formatterer
	// dataen i egne structs
	for _, firstProduct := range apiProducts {
		gtin := firstProduct.Data.Ean

		// om produktet allerede er sjekket, skip dette produktet
		if lib.IsIn(gtin, checkedGtins) {
			continue
		}
		checkedGtins = append(checkedGtins, gtin)

		// finner andre produkter med samme gtin og legger til i et array
		sameProduct := []ApiProduct{firstProduct}

		for _, secondProduct := range apiProducts {
			if gtin == secondProduct.Data.Ean &&
				firstProduct.Store != secondProduct.Store {
				// legger til produktet i sameProduct array, alle produkter i dette
				// arrayet sjekkes nå priser legges inn
				sameProduct = append(sameProduct, secondProduct)
			}
		}

		// formaterer dataen til alle produkter med samme gtin
		firstProduct.FormatData(sameProduct, products)
	}
}
