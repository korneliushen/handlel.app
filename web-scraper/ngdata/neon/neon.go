package neon

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/korneliushen/handlel.app/meny/ngdata"
	_ "github.com/lib/pq"
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

// legger til data i neon databasen
func InsertData(products *[]ngdata.Product) {
	// neon
	db := db()
	defer db.Close()

	// lager en waitgroup, som venter på goroutines for å bli ferdig før den
	// starter en ny
	var wg sync.WaitGroup
	// limiter hvor mange go routines som kan kjøre om om gangen
	sem := make(chan struct{}, 4)

	for i := range *products {
		// legger til et item i wait groupen
		wg.Add(1)
		sem <- struct{}{}

		go func(product ngdata.Product) {
			fmt.Println("Legger inn data for:", product.Title)
			// når funksjonen er ferdig, blir waitgroup instansen ferdig + sem
			// (hvor mange ting som kan kjøre om gangen) blir oppdatert
			defer wg.Done()
			defer func() { <-sem }()

			// legger til data i databasen
			if err := query(product, db); err != nil {
				fmt.Printf("Error inserting data into neon db for %s: %v",
					product.Title, err)
			}
		}((*products)[i])
	}

	wg.Wait()
}

func query(product ngdata.Product, db *sql.DB) error {
	// gjør om næringsinnhold (type Næringsinnhold struct) til
	// nutritionalContentJson (basically bare gjør om til json)
	nutritionalContentJson, err := json.Marshal(product.NutritionalContent)
	if err != nil {
		return err
	}

	// gjør om priser til json
	pricesJson, err := json.Marshal(product.Prices)
	if err != nil {
		return err
	}

	// legger til en rad i Products table i databasen. om en rad med samme id
	// (gtin) allerede eksisterer, blir den replaced
	// her gjører bare queryen klart, uten dette blir goroutinene helt fked
	// up og overlapper
	productsStmt, err := db.Prepare(`
		INSERT INTO products (
			id, title, subtitle, imagelink, category, 
			subcategory, onsale, description, weight, origincountry, ingredients, 
			vendor, brand, size, unit, unittype, allergens, mayContainTracesOf, 
			nutritionalcontent, prices
		)
		VALUES (
			$1, $2, $3, $4, $5, $6 , $7, $8, $9, $10,
			$11, $12, $13, $14, $15, $16, $17, $18, $19, $20
		)
		ON CONFLICT (id)
		DO UPDATE SET
			title = EXCLUDED.title,
			subtitle = EXCLUDED.subtitle,
			imagelink = EXCLUDED.imagelink,
			category = EXCLUDED.category,
			subcategory = EXCLUDED.subcategory,
			onsale = EXCLUDED.onsale,
			description = EXCLUDED.description,
			weight = EXCLUDED.weight,
			origincountry = EXCLUDED.origincountry,
			ingredients = EXCLUDED.ingredients,
			vendor = EXCLUDED.vendor,
			brand = EXCLUDED.brand,
			size = EXCLUDED.size,
			unit = EXCLUDED.unit,
			unittype = EXCLUDED.unittype,
			allergens = EXCLUDED.allergens,
			mayContainTracesOf = EXCLUDED.mayContainTracesOf,
			nutritionalcontent = EXCLUDED.nutritionalcontent,
			prices = EXCLUDED.prices
		`)
	if err != nil {
		return err
	}
	defer productsStmt.Close()

	// queryen executes med verdiene fra product
	_, err = productsStmt.Exec(product.Id, product.Title, product.SubTitle,
		product.ImageLink, product.Category, product.SubCategory, product.OnSale,
		product.Description, product.Weight, product.OriginCountry,
		product.Ingredients, product.Vendor, product.Brand, product.Size,
		product.Unit, product.UnitType, product.Allergens,
		product.MayContainTracesOf, nutritionalContentJson, pricesJson)
	if err != nil {
		return err
	}

	// ingen error så returnerer nil
	return nil
}
