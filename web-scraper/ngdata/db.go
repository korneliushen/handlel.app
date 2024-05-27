package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func db() *sql.DB {
	connStr := os.Getenv("NEON_SECRET")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	// limiter open og idle connections for å ikke med et uhell overloade databasen
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)

	return db
}
