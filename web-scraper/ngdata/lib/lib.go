package lib

import (
	"database/sql"
	"log"
	"os"
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

func isIn[T comparable](e T, arr []T) bool {
	for i := range arr {
		if e == arr[i] {
			return true
		}
	}
	return false
}
