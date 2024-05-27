package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

func run() {
	products := &Produkter{}

	categories := getCategories()

	for i := range categories.Kategorier {
		// får kategori
		category := categories.Kategorier[i]
		fmt.Println("Henter data for kategori:", category)
		for j := range category.Underkategorier {
			// får underkategori
			subCategory := category.Underkategorier[j]

			fmt.Println("Henter data for underkategori:", subCategory)

			menyData, err := getProducts("meny", category.Navn, subCategory.Navn)
			if err != nil {
				fmt.Printf("Error getting products: %v\n", err)
			}

			jokerData, err := getProducts("joker", category.Navn, subCategory.Navn)
			if err != nil {
				fmt.Printf("Error getting products: %v\n", err)
			}

			sparData, err := getProducts("spar", category.Navn, subCategory.Navn)
			if err != nil {
				fmt.Printf("Error getting products: %v\n", err)
			}

			// for hvert produkt som er returnert fra api-en, legg til dataen i products structen vi fikk som param
			for k := range menyData.Hits.Products {
				gtin := menyData.Hits.Products[k].Data.Ean

				menyProduct := menyData.Hits.Products[k]
				jokerProduct, sparProduct := getPrices(gtin, jokerData, sparData)

				fmt.Println("Formaterer data for:", gtin, menyProduct.Data.Title)

				formatData(menyProduct, jokerProduct, sparProduct, products)
			}
		}
	}

	db := db()
	defer db.Close()

	// lager en waitgroup, som venter på goroutines for å bli ferdig før den starter en ny
	var wg sync.WaitGroup
	// limiter hvor mange go routines som kan kjøre om om gangen
	sem := make(chan struct{}, 4)

	for i := range products.Produkter {
		// legger til et item i wait groupen
		wg.Add(1)
		sem <- struct{}{}

		go func(product Produkt) {
			// når funksjonen er ferdig, blir waitgroup instansen ferdig + sem (det som keeper track av hvor mange ting som kan kjøre om gangen) blir oppdatert
			defer wg.Done()
			defer func() { <-sem }()

			// legger til data i databasen
			if err := insertData(product, db); err != nil {
				fmt.Printf("Error inserting data for %s: %v", products.Produkter[i].Tittel, err)
			}
		}(products.Produkter[i])
	}

	wg.Wait()
}

func main() {
	start := time.Now()

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Couldn't load env variables: %v\n", err)
	}

	fmt.Println("Running")
	run()

	end := time.Now()
	elapsed := end.Sub(start)
	fmt.Println("Elapsed: ", elapsed)

	fmt.Println("Everything ran successfully :O")
}
