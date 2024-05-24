package main

import (
	"fmt"
	"time"
)

func run() {
	products := &Produkter{}

	categories := getCategories()

	for i := range categories.Kategorier {
		category := categories.Kategorier[i]
		for j := range category.Underkategorier {
			subCategory := category.Underkategorier[j]

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
				insertData(menyData.Hits.Products[k], jokerData.Hits.Products[k], sparData.Hits.Products[k], products)
			}
			break
		}
		break
	}

	err := WriteData(products, "./data.json")
	if err != nil {
		fmt.Printf("Error writing to json: %v\n", err)
		return
	}
}

func main() {
	start := time.Now()

	run()

	end := time.Now()
	elapsed := end.Sub(start)
	fmt.Println("Elapsed: ", elapsed)

}
