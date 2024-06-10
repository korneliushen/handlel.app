package ngdata

import (
	"github.com/korneliushen/handlel.app/scraper/model"
)

func Ngdata(apiProducts *model.ApiProducts) {
	// henter alle kategorier
	categories := GetCategories()

	// henter alle produkter og appender til apiProducts arrayet
	GetProducts(apiProducts, categories)
}
