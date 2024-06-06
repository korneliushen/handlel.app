package lib

import (
	"database/sql"
	"log"
	"os"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
)

func db() *sql.DB {
	connStr := os.Getenv("NEON_URL")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	// limiter open og idle connections for å ikke med et uhell overloade databasen
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return db
}

// henter indexen som dataen skal sendes til i algolia
func index() *search.Index {
	client := search.NewClient("AA8FDXU3JW", os.Getenv("ALGOLIA_SECRET"))
	return client.InitIndex("test")
}

func isIn[T comparable](e T, arr []T) bool {
	for i := range arr {
		if e == arr[i] {
			return true
		}
	}
	return false
}
