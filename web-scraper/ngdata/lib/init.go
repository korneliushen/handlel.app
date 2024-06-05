package lib

import (
	"fmt"
	"log"
	"time"

	"github.com/joho/godotenv"
)

func run() {
	products := &[]Product{}

	categories := &[]Category{}

	for _, store := range stores {
		getCategories(store, categories)
	}

	// har produkter som allerede er sjekket i et array, så det ikke blir duplicates av samme produkt (sparer også tid fordi den exiter tidlig)
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
			break
		}
		break
	}

	// mapper over alle produkter vi har fått fra databasen og formatterer dataen i egne structs
	for _, firstProduct := range apiProducts {
		gtin := firstProduct.Data.Ean

		// om produktet allerede er sjekket, continue
		if isIn(gtin, checkedGtins) {
			continue
		}
		checkedGtins = append(checkedGtins, gtin)

		// finner andre produkter med samme gtin og legger til i et array
		sameProduct := []ApiProduct{firstProduct}
		for _, secondProduct := range apiProducts {
			if gtin == secondProduct.Data.Ean && firstProduct.Store != secondProduct.Store {
				sameProduct = append(sameProduct, secondProduct)
			}
		}

		// formaterer dataen til alle produkter med samme gtin
		formatData(sameProduct, products)
	}

	// legger data inn i neon database og legger til records i algolia
	insertData(products)

	if err := insertRecords(*products); err != nil {
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
