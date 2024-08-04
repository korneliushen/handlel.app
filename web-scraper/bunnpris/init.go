package bunnpris

import (
	"context"
	"fmt"

	"github.com/korneliushen/handlel.app/scraper/bunnpris/bunnpris"
	"github.com/korneliushen/handlel.app/scraper/model"
)

func Bunnpris(products *model.Products) {
  originalLength := len(*products)

	ctx := context.Background()

	token, err := bunnpris.ReadToken()
	if err != nil {
		fmt.Printf("Error getting token: %v\n", err)
		return
	}

	fmt.Println("Getting categories")

	var categories bunnpris.Categories
	if err := categories.Get(ctx, token.Value); err != nil {
		fmt.Printf("Error getting categories: %v\n", err)
		return
	}

	fmt.Println("Getting products")

	var productLinks bunnpris.BunnprisProducts
	for _, category := range categories {
		err := productLinks.Get(ctx, token.Value, category.Id)
		if err != nil {
			fmt.Printf("Error getting products from %s, %s: %s", category.Name, category.Link, err.Error())
		}
	}

	fmt.Println("Getting product data")

	productLinks.FetchProductPages(ctx, token.Value, products)

  if originalLength == len(*products) {
    fmt.Println("Got no products from bunnpris")
  }
}
