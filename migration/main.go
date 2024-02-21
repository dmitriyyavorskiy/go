package main

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/xuri/excelize/v2"
	"log"
)

const (
	// PostgreSQL connection string
	psqlInfo         = "host=mioxxo-products.ce989v8mdple.us-east-1.rds.amazonaws.com port=5432 user=postgres dbname=products password=Yk9deWhjbYUUR5LEF8SnMf7w9jghPf sslmode=disable"
	lastModifiedById = "005Hp00000igHBGIA2"
	recordTypeId     = "012Hp000001mPmEIAU"
	createdDate      = "2021-08-25T20:00:00.000Z"
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
	Tags             sql.NullString `db:"tags"`
}

func main() {
	db := createDatabaseConnection()
	var brands = readBrands(*db)
	saveToExcelFile("brands.xlsx", brands)

	//readCategories(*db)

	defer db.Close()

	//db = createDatabaseConnection()
	//readProducts(*db)
	//defer db.Close()
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
	err := db.Select(&products, "SELECT sku, barcode, name, short_description, variant, brand, categories, image, tags FROM mgo.products")
	if err != nil {
		log.Fatalf("Could not read products: %v", err)
	}

	for _, product := range products {
		fmt.Printf("Product %+v\n", product)
	}
	return products

}

func saveToExcelFile(filename string, brands []Brand) {
	sheetName := "Sheet1"
	file := excelize.NewFile()
	index, err := file.NewSheet(sheetName)
	if err != nil {
		log.Fatalf("Could not create sheet: %v", err)
	}
	file.SetActiveSheet(index)
	// Set the headers
	headers := []string{"_", "Id", "OwnerId", "IsDeleted", "Name", "CurrencyIsoCode", "RecordTypeId", "CreatedDate", "CreatedById", "LastModifiedDate", "LastModifiedById", "SystemModstamp", "LastActivityDate", "LastViewedDate", "LastReferencedDate", "Brand_Owner__c", "Type__c"}
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		file.SetCellValue(sheetName, cell, header)
	}

	for i, brand := range brands {
		row := i + 2
		setCellValue(file, sheetName, row, 1, "[Brand__c]")
		setCellValue(file, sheetName, row, 2, brand.ID)
		setCellValue(file, sheetName, row, 3, "005Hp00000igHBGIA2")
		setCellValue(file, sheetName, row, 4, "FALSE")
		setCellValue(file, sheetName, row, 5, brand.Name.String)
		setCellValue(file, sheetName, row, 6, "MXN")
		setCellValue(file, sheetName, row, 7, recordTypeId)
		setCellValue(file, sheetName, row, 8, createdDate)
		setCellValue(file, sheetName, row, 9, lastModifiedById)
		setCellValue(file, sheetName, row, 10, createdDate)
		setCellValue(file, sheetName, row, 11, lastModifiedById)
		setCellValue(file, sheetName, row, 12, createdDate)
		setCellValue(file, sheetName, row, 16, "Brand")
	}
	if err := file.SaveAs(filename); err != nil {
		fmt.Println(err)
	}
}

func setCellValue(file *excelize.File, sheetName string, row int, column int, value string) {
	cell, _ := excelize.CoordinatesToCellName(column, row)
	file.SetCellValue(sheetName, cell, value)
}
