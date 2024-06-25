package algolia

import (
	"fmt"
	"os"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/korneliushen/handlel.app/scraper/model"
)

// Legger til data i algolia index
func InsertRecords(products model.Products) error {
	fmt.Println("Legger data inn i algolia index")

  // Formaterer dataen til et mindre format
  formattedProducts := format(products)

	// Instanse av index brukt til å interacte med algolia indexen
	index := index()

	// Sender alle objekter i products arrayet til algolia indexen
	// Sendes i batches av 1000 om gangen (default).
	// Kan endres til å sende flere, men dette går fort nok så trengs i hver
	// fall ikke ennå
	res, err := index.SaveObjects(formattedProducts)
	if err != nil {
		return err
	}
	res.Wait()

	return nil
}

// Henter indexen som dataen skal sendes til i algolia og lager en instanse av
// *search.Index
func index() *search.Index {
	client := search.NewClient("AA8FDXU3JW", os.Getenv("ALGOLIA_SECRET"))
	return client.InitIndex("test")
}

// Mapper over products array og legger bare til nødvendige fields for å spare
// plass (har bare 1000 bytes per record)
func format(products model.Products) Products {
  formattedProducts := Products{}

  for _, product := range products {
    formattedProducts = append(formattedProducts, 
      Product{
        Id: product.Id,
        Title: product.Title,
        Subtitle: product.Title,
        Store: product.Store,
        Prices: product.Prices,
        Images: product.Images,
        OnSale: product.OnSale,
        Unit: product.Unit,
        ObjectID: product.Id,
      },
    )
  }

  return formattedProducts
}
