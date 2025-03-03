package algolia

import (
	"fmt"
	"os"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/korneliushen/handlel.app/scraper/model"
)

// Henter indexen som dataen skal sendes til i algolia og lager en instanse av
// *search.Index
func index() *search.Index {
	client := search.NewClient("AA8FDXU3JW", os.Getenv("ALGOLIA_SECRET"))
	return client.InitIndex("test")
}

// Legger til data i algolia index
func InsertRecords(products model.Products) error {
	fmt.Println("Legger data inn i algolia index")
	// Instanse av index brukt til å interacte med algolia indexen
	index := index()

	// Sender alle objekter i products arrayet til algolia indexen
	// Sendes i batches av 1000 om gangen (default).
	// Kan endres til å sende flere, men dette går fort nok så trengs i hver
	// fall ikke ennå
	res, err := index.SaveObjects(products)
	if err != nil {
		return err
	}
	res.Wait()

	return nil
}
