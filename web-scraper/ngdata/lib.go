package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func writeData(data any, path string) error {
	jsonData, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return err
	}

	err = os.WriteFile(path, jsonData, 0666)
	if err != nil {
		return err
	}

	return nil
}

func db() *sql.DB {
	connStr := os.Getenv("NEON_SECRET")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	return db
}
