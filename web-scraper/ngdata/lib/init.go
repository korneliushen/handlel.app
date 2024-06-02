package lib

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

func run() {
	products := &[]Product{}

	categories := &[]Category{}

	for _, store := range stores {
		getCategories(store, categories)
	}

	var checkedGtins []string

	apiProducts := []ApiProduct{}

	for _, category := range *categories {
		for _, store := range stores {
			// om kategorien sin butikk og butikken ikke er den samme, er det ikke vits å kjøre request fordi den vil ikke få noe data (og om den får det vil det være duplicate)
			if category.Store != store {
				continue
			}

			// får data om alle produkter i kategorien
			res, err := getProducts(store, category.Name)
			if err != nil {
				fmt.Printf("Error getting products from %s in category %s: %v\n", store, category, err)
				continue
			}

			// legger til produktet i apiProducts array som mappes over senere, legger også til Store (for senere bruk)
			for _, product := range res {
				// legger til underkategorier, legger bare til om underkategorien ikke allerede er lagt til
				// underkategorier er jeg ganske sikker på at er basically helt likt på alle sidene, så det vil ikke være duplicates med forskjellig navn, om det er annerledes må jeg bytte til id approach
				if !isIn(product.Data.SubCategory, category.SubCategories) {
					category.SubCategories = append(category.SubCategories, product.Data.SubCategory)
				}
				apiProducts = append(apiProducts, ApiProduct{Store: store, Data: product.Data, BaseUrl: storeData[store].url})
			}
		}
	}

	// mapper over alle produkter vi har fått fra databasen og formatterer dataen i egne structs
	for _, firstProduct := range apiProducts {
		gtin := firstProduct.Data.Ean

		if isIn(gtin, checkedGtins) {
			continue
		}
		checkedGtins = append(checkedGtins, gtin)

		sameProduct := []ApiProduct{firstProduct}

		for _, secondProduct := range apiProducts {
			if gtin == secondProduct.Data.Ean {
				sameProduct = append(sameProduct, secondProduct)
			}
		}

		formatData(sameProduct, products)
	}

	db := db()
	defer db.Close()

	// lager en waitgroup, som venter på goroutines for å bli ferdig før den starter en ny
	var wg sync.WaitGroup
	// limiter hvor mange go routines som kan kjøre om om gangen
	sem := make(chan struct{}, 4)

	for i := range *products {
		// legger til et item i wait groupen
		wg.Add(1)
		sem <- struct{}{}

		go func(product Product) {
			// når funksjonen er ferdig, blir waitgroup instansen ferdig + sem (det som keeper track av hvor mange ting som kan kjøre om gangen) blir oppdatert
			defer wg.Done()
			defer func() { <-sem }()

			// legger til data i databasen
			if err := insertData(product, db); err != nil {
				fmt.Printf("Error inserting data for %s: %v", product.Title, err)
			}
		}((*products)[i])
	}

	wg.Wait()
}

func Init() {
	start := time.Now()

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Couldn't load env variables: %v\n", err)
	}

	run()
	elapsed := time.Now().Sub(start)
	fmt.Println("Elapsed: ", elapsed)
}
