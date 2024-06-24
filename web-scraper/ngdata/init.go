package ngdata

import (
	"github.com/korneliushen/handlel.app/scraper/model"
)

func Ngdata(products *model.BaseProducts) {
	// henter alle kategorier
	categories := GetCategories()

	// henter alle produkter og appender til baseProducts arrayet
	GetProducts(products, categories)
}
