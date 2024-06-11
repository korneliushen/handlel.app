package bunnpris

import (
	"context"
	"fmt"

	"github.com/korneliushen/handlel.app/scraper/bunnpris/bunnpris"
	"github.com/korneliushen/handlel.app/scraper/model"
)

func Bunnpris(apiProducts *model.ApiProducts) {
	ctx := context.Background()

	token := "wlli2e4n52sjbydqekk2nnki"

	var categories bunnpris.Categories
	if err := categories.Get(ctx, token); err != nil {
		fmt.Printf("Error getting categories: %v\n", err)
		return
	}

	for _, category := range categories {
		err := bunnpris.GetProducts(apiProducts, ctx, token, category.Link)
		if err != nil {
			fmt.Printf("Error getting products from %s, %s: %s", category.Name, category.Link, err.Error())
		}
	}
}
