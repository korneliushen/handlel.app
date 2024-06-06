package algolia

import (
	"fmt"
	"os"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/korneliushen/handlel.app/meny/ngdata"
)

// TODO: bare legge til data om det er ny data

// henter indexen som dataen skal sendes til i algolia
func index() *search.Index {
	client := search.NewClient("AA8FDXU3JW", os.Getenv("ALGOLIA_SECRET"))
	return client.InitIndex("test")
}

// legger til data i algolia index
func InsertRecords(products []ngdata.Product) error {
	fmt.Println("Legger data inn i algolia index")
	// instanse av index brukt til å interacte med algolia indexen
	index := index()

	// sender alle objekter i products arrayet til algolia indexen
	res, err := index.SaveObjects(products)
	if err != nil {
		return err
	}
	res.Wait()

	return nil
}
