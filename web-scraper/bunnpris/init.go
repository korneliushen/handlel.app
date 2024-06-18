package bunnpris

import (
	"context"
	"fmt"

	"github.com/korneliushen/handlel.app/scraper/bunnpris/bunnpris"
	"github.com/korneliushen/handlel.app/scraper/model"
)

func Bunnpris(apiProducts *model.ApiProducts) {
	ctx := context.Background()

	token := "wl1irk2zaqbnbafvvufncju5"

	fmt.Println("Getting categories")

	var categories bunnpris.Categories
	if err := categories.Get(ctx, token); err != nil {
		fmt.Printf("Error getting categories: %v\n", err)
		return
	}

	fmt.Println("Getting products")

	var products bunnpris.BunnprisProducts
	for _, category := range categories {
		err := products.Get(apiProducts, ctx, token, category.Id)
		if err != nil {
			fmt.Printf("Error getting products from %s, %s: %s", category.Name, category.Link, err.Error())
		}
	}

	fmt.Println("Getting product data")

	products.FetchProductPages(ctx, token, apiProducts)
}
