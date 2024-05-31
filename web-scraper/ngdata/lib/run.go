package lib

import (
	"fmt"
	"sync"
)

func Run() {
	products := &Products{}

	categories := getCategories()

	for i := range categories.Category {
		// får kategori
		category := categories.Category[i]
		fmt.Println("Henter data for kategori:", category)
		for j := range category.SubCategories {
			// får underkategori
			subCategory := category.SubCategories[j]
			fmt.Println("Henter data for underkategori:", subCategory)

			for _, value := range stores {
				fmt.Println(value)
				// kanskje kjør get products her: getProducts(value, category.Name, subCategory.Name)
			}

			// TODO FOR MEG: gjøre de tre neste getProducts funksjonene uavhengig av navn

			menyData, err := getProducts("meny", category.Name, subCategory.Name)
			if err != nil {
				fmt.Printf("Error getting products: %v\n", err)
			}

			jokerData, err := getProducts("joker", category.Name, subCategory.Name)
			if err != nil {
				fmt.Printf("Error getting products: %v\n", err)
			}

			sparData, err := getProducts("spar", category.Name, subCategory.Name)
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
