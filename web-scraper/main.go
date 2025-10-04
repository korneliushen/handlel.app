package main

import (
	"fmt"
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/korneliushen/handlel.app/scraper/bunnpris"
	"github.com/korneliushen/handlel.app/scraper/model"
	"github.com/korneliushen/handlel.app/scraper/neon"
	"github.com/korneliushen/handlel.app/scraper/ngdata"
)

func run() {
	// Instansierer en variabel med type BaseProduct
	// Denne brukes til å sende data til databasen og algolia, så alle butikker
	// legger bare til data i dette arrayet
  products := &model.Products{}

	// meny, joker, spar
  ngdata.Ngdata(products)

  // bunnpris
  bunnpris.Bunnpris(products)

	// LEGG TIL NYE BUTIKKER UNDER HER:
	//
	//
	//
	//

	// Mapper over alle produkter vi har fått fra databasen og formatterer
	// dataen i egne structs
  formattedProducts := products.Format()

	if len(*products) == 0 {
		fmt.Println("No products")
		return
	}
  fmt.Println("Formated: ", len(formattedProducts))

	// legger data inn i neon databasen
	neon.InsertData(formattedProducts)
}

func main() {
	start := time.Now()

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Couldn't load env variables: %v\n", err)
	}

	run()

	elapsed := time.Since(start)
	fmt.Println("Elapsed: ", elapsed)
}
