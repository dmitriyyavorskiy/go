package main

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
)

const (
	// PostgreSQL connection string
	psqlInfo = "host=mioxxo-products.ce989v8mdple.us-east-1.rds.amazonaws.com port=5432 user=postgres dbname=products password=Yk9deWhjbYUUR5LEF8SnMf7w9jghPf sslmode=disable"
)

type Brand struct {
	ID    string         `db:"_id"`
	Name  sql.NullString `db:"name"`
	Image sql.NullString `db:"image"`
}

type Category struct {
	ID    string         `db:"_id"`
	Name  sql.NullString `db:"name"`
	Image sql.NullString `db:"image"`
}

type Product struct {
	Sku              string         `db:"sku"`
	Barcode          sql.NullString `db:"barcode"`
	Name             string         `db:"name"`
	ShortDescription string         `db:"short_description"`
	Variant          string         `db:"variant"`
	Brand            sql.NullString `db:"brand"`
	Categories       sql.NullString `db:"categories"`
	Image            sql.NullString `db:"image"`
}

func main() {
	db := createDatabaseConnection()
	readBrands(*db)
	readCategories(*db)

	defer db.Close()

	db = createDatabaseConnection()
	readProducts(*db)
	defer db.Close()
}

func createDatabaseConnection() *sqlx.DB {
	// Open the connection
	db, err := sqlx.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func readBrands(db sqlx.DB) []Brand {
	var brands []Brand
	err := db.Select(&brands, "SELECT _id, name, image FROM mgo.brands")
	if err != nil {
		log.Fatalf("Could not read brands: %v", err)
	}

	for _, brand := range brands {
		fmt.Printf("Brand %+v\n", brand)
	}
	return brands
}

func readCategories(db sqlx.DB) []Category {
	var categories []Category
	err := db.Select(&categories, "SELECT _id, name, image FROM mgo.categories")
	if err != nil {
		log.Fatalf("Could not read categories: %v", err)
	}

	for _, category := range categories {
		fmt.Printf("Category %+v\n", category)
	}
	return categories
}

func readProducts(db sqlx.DB) []Product {
	var products []Product
	err := db.Select(&products, "SELECT sku, barcode, name, short_description, variant, brand, categories, image FROM mgo.products")
	if err != nil {
		log.Fatalf("Could not read products: %v", err)
	}

	for _, product := range products {
		fmt.Printf("Product %+v\n", product)
	}
	return products

}
