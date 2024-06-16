package main

import (
	"fmt"
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/korneliushen/handlel.app/meny/algolia"
	"github.com/korneliushen/handlel.app/meny/neon"
	"github.com/korneliushen/handlel.app/meny/ngdata"
)

func run() {
	// henter alle kategorier
	categories := ngdata.Categories{}
	categories.GetCategories()

	// henter alle produkter, lagres i et array med ApiProduct objekter
	apiProducts := ngdata.ApiProducts{}
	apiProducts.GetProducts(categories)

	// mapper over alle produkter vi har fått fra databasen og formatterer
	// dataen i egne structs
	products := ngdata.Products{}
	products.Format(apiProducts)

	// legger data inn i neon databasen
	neon.InsertData(products)

	// legger data inn i algolia indexen
	if err := algolia.InsertRecords(products); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
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
