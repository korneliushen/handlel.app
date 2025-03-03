package algolia

import "github.com/korneliushen/handlel.app/scraper/model"

// Currently ikke i bruk, bare om vi må kutte ned på størrelsen til records
// Dette er minumun records trenger for å fungere

type Products []Product

type Product struct {
  Id string `json:"id"`
  Title string `json:"title"`
  Subtitle string `json:"subtitle"`
  Store string `json:"store"`
  Prices []model.Price `json:"prices"`
  Images model.Images `json:"images"`
  OnSale bool `json:"onsale"`
  Unit string `json:"unit"`
  ObjectID string `json:"objectID"`
}
