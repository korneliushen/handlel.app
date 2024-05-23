package main

import (
	"fmt"
	"time"
)

func main() {
	start := time.Now()

	categories := getCategories()

	for i := range categories.Kategorier {
		category := categories.Kategorier[i]
		for j := range category.Underkategorier {
			subCategory := category.Underkategorier[j]
			err := getProducts("meny", category.Navn, subCategory.Navn)
			if err != nil {
				fmt.Printf("Error getting products: %v\n", err)
			}
			break
		}
		break
	}

	end := time.Now()
	elapsed := end.Sub(start)
	fmt.Println("Elapsed: ", elapsed)

}
