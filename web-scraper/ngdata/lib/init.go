package lib

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

func run() {
	products := &Products{}

	categories := getCategories()

	for _, category := range categories.Category {
		for _, subCategory := range category.SubCategories {
			storeData := []StoreData{}

			for _, store := range stores {
				res, err := getProducts(store, category.Name, subCategory.Name)
				if err != nil {
					fmt.Printf("Error getting products from %s in sub-category %s: %v\n", store, category.Name, err)
				}
				storeData = append(storeData, StoreData{Store: store, ApiRes: res, Category: category.Name, SubCategory: subCategory.Name})
			}

			// for hvert produkt som er returnert fra api-en, legg til dataen i products structen vi fikk som param
			for i := range stores {
				// bruker no goofy kode fra chatgpt for å kunne rotere hvilken rolle hver butikk har
				// For eks: første iterasjon vil meny være hoved butikken som itereres over, andre gang blir det spar osv.
				firstIndex := i
				secondIndex := (i + 1) % len(stores)
				thirdIndex := (i + 2) % len(stores)

				for _, product := range storeData[firstIndex].ApiRes.Hits.Products {
					gtin := product.Data.Ean

					secondProduct, thirdProduct := getPrices(gtin, storeData[secondIndex].ApiRes, storeData[thirdIndex].ApiRes)

					formatData(product, secondProduct, thirdProduct, storeData[firstIndex].Store, storeData[secondIndex].Store, storeData[thirdIndex].Store, products)
				}
			}
		}
		break
	}

	db := db()
	defer db.Close()

	// lager en waitgroup, som venter på goroutines for å bli ferdig før den starter en ny
	var wg sync.WaitGroup
	// limiter hvor mange go routines som kan kjøre om om gangen
	sem := make(chan struct{}, 4)

	for i := range products.Products {
		// legger til et item i wait groupen
		wg.Add(1)
		sem <- struct{}{}

		go func(product Product) {
			// når funksjonen er ferdig, blir waitgroup instansen ferdig + sem (det som keeper track av hvor mange ting som kan kjøre om gangen) blir oppdatert
			defer wg.Done()
			defer func() { <-sem }()

			// legger til data i databasen
			if err := insertData(product, db); err != nil {
				fmt.Printf("Error inserting data for %s: %v", products.Products[i].Title, err)
			}
		}(products.Products[i])
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
