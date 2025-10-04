package ngdata

import "github.com/korneliushen/handlel.app/scraper/model"

func Fetch(products *model.Products) {
	// Henter alle kategorier
	categories := GetCategories()

	// Henter alle produkter fra ngdata api og returnere et array med formaterte
  // produkter
  ngdataProducts := GetProducts(categories)

  // Appender alle formaterte produkter vi får fra ngdata api-et til products
  // arrayet
  for _, product := range ngdataProducts {
    *products = append(*products, product)
  }
}
