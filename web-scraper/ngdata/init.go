package main

import (
	"fmt"
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/korneliushen/handlel.app/meny/algolia"
	"github.com/korneliushen/handlel.app/meny/lib"
	"github.com/korneliushen/handlel.app/meny/neon"
	"github.com/korneliushen/handlel.app/meny/ngdata"
)

func run() {
	products := &[]lib.Product{}

	categories := &[]lib.Category{}

	// henter alle kategorier
	for _, store := range ngdata.Stores {
		ngdata.GetCategories(store, categories)
	}

	// har produkter som allerede er sjekket i et array, så det ikke blir
	// duplicates av samme produkt (sparer også tid fordi den exiter tidlig)
	var checkedGtins []string

	apiProducts := []lib.ApiProduct{}

	for _, category := range *categories {
		for _, store := range ngdata.Stores {
			// om kategorien sin butikk og butikken ikke er den samme, er det ikke
			// vits å kjøre request fordi den vil ikke få noe data
			// (og om den får det vil det være duplicate)
			if category.Store != store {
				continue
			}

			// får data om alle produkter i kategorien
			res, err := ngdata.GetProducts(store, category.Name)
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
				apiProducts = append(apiProducts, lib.ApiProduct{
					Store: store, Data: product.Data, BaseUrl: ngdata.StoreInfo[store].Url,
				})
			}
			break
		}
		break
	}

	// mapper over alle produkter vi har fått fra databasen og formatterer
	// dataen i egne structs
	for _, firstProduct := range apiProducts {
		gtin := firstProduct.Data.Ean

		// om produktet allerede er sjekket, continue
		if lib.IsIn(gtin, checkedGtins) {
			continue
		}
		checkedGtins = append(checkedGtins, gtin)

		// finner andre produkter med samme gtin og legger til i et array
		sameProduct := []lib.ApiProduct{firstProduct}
		for _, secondProduct := range apiProducts {
			if gtin == secondProduct.Data.Ean &&
				firstProduct.Store != secondProduct.Store {
				sameProduct = append(sameProduct, secondProduct)
			}
		}

		// formaterer dataen til alle produkter med samme gtin
		firstProduct.FormatData(sameProduct, products)
	}

	// legger data inn i neon database og legger til records i algolia
	neon.InsertData(products)

	if err := algolia.InsertRecords(*products); err != nil {
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
