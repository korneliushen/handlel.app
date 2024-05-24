package main

import (
	"fmt"
	"log"
	"time"

	"github.com/joho/godotenv"
)

func run() {
	products := &Produkter{}

	categories := getCategories()

	for i := range categories.Kategorier {
		// får kategori
		category := categories.Kategorier[i]
		for j := range category.Underkategorier {
			// får underkategori
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
				gtin := menyData.Hits.Products[k].Data.Ean

				menyProduct := menyData.Hits.Products[k]
				jokerProduct, sparProduct := getPrices(gtin, jokerData, sparData)

				formatData(menyProduct, jokerProduct, sparProduct, products)
			}
			break
		}
		break
	}

	for i := range products.Produkter {
		err := insertData(products.Produkter[i])
		if err != nil {
			fmt.Printf("Error inserting data into db: %v\n", err)
			return
		}
	}
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
