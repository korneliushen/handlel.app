package main

import (
	"fmt"
	"time"

	"github.com/korneliushen/handlel.app/meny/fetch"
	"github.com/korneliushen/handlel.app/meny/scraper"
)

func main() {
	start := time.Now()

	categories := scraper.GetCategories()

	for i := range categories.Kategorier {
		category := categories.Kategorier[i]
		for j := range category.Underkategorier {
			subCategory := category.Underkategorier[j]
			fetch.GetProducts(category.Navn, subCategory.Navn)
			break
		}
		break
	}

	end := time.Now()
	elapsed := end.Sub(start)
	fmt.Println("Elapsed: ", elapsed)

}
