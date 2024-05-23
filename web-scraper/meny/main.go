package main

import (
	"fmt"
	"time"

	"github.com/korneliushen/handlel.app/meny/fetch"
)

// TODO: ting å snakke med henrik om imrgn: næringsinnhold meny og så priser fra andre?
// legge til priser for alle butikker i database skisse

func main() {
	start := time.Now()

	categories := getCategories()

	for i := range categories.Kategorier {
		category := categories.Kategorier[i]
		for j := range category.Underkategorier {
			subCategory := category.Underkategorier[j]
			fetch.GetProducts("meny", category.Navn, subCategory.Navn)
			break
		}
		break
	}

	end := time.Now()
	elapsed := end.Sub(start)
	fmt.Println("Elapsed: ", elapsed)

}
