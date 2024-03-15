package main

import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"os"
	"strings"
)

const (
	// PostgreSQL connection string
	psqlInfo             = "host=mioxxo-products.ce989v8mdple.us-east-1.rds.amazonaws.com port=5432 user=postgres dbname=products password=Yk9deWhjbYUUR5LEF8SnMf7w9jghPf sslmode=disable"
	recordTypeBrandOwner = "012Hp000001mPmGIAU"
	recordTypeBrand      = "012Hp000001mPmEIAU"
)

var brandsMap = make(map[string]Brand)

type Brand struct {
	ID    string         `db:"_id"`
	Name  sql.NullString `db:"name"`
	Image sql.NullString `db:"image"`
}

type SalesForceEntity struct {
	ID              string
	Name            string
	ReferenceTypeId string
}

type SalesForceCategory struct {
	ID                 string
	Name               string
	CategoryLevel1Name string
	CategoryLevel2Name string
}

type CategoryTree struct {
	ID           string         `db:"id"`
	Name         sql.NullString `db:"name"`
	Image        sql.NullString `db:"image"`
	CategoryId   sql.NullString `db:"subcategory_id"`
	CategoryName sql.NullString `db:"subcategory_name"`
}

type Tax struct {
	Rate        int    `json:"rate"`
	Type        string `json:"type"`
	Withholding bool   `json:"withholding"`
}

type Product struct {
	Sku              string         `db:"sku"`
	Barcode          sql.NullString `db:"barcode"`
	Name             string         `db:"name"`
	ShortDescription string         `db:"short_description"`
	Variant          string         `db:"variant"`
	Brand            sql.NullString `db:"brand_name"`
	Categories       sql.NullString `db:"categories"`
	CategoryName     sql.NullString `db:"category_name"`
	SubcategoryName  sql.NullString `db:"subcategory_name"`
	Image            sql.NullString `db:"image"`
	Tags             sql.NullString `db:"tags"`
	Taxes            sql.NullString `db:"taxes"`
	AgeRestriction   bool           `db:"age_restriction"`
	Restricted       bool           `db:"restricted"`
	Enabled          bool           `db:"is_enabled"`
}

func main() {
	//var brands = readBrands()
	//saveBrandOwnersToCsvFile("brandowners.csv", brands)

	//// TODO import data to the Salesforce here
	//// TODO export Brands csv file from Salesforce
	//
	//brandOwners := readSalesforceEntities("brandExported.csv", recordTypeBrandOwner)
	//for key, value := range brandOwners {
	//	fmt.Printf("Key '%s' Brand owner %+v\n", key, value)
	//}
	//
	//saveBrandsToCsvFile("brands.csv", brands, brandOwners)

	//var categoryTrees = readCategoryTrees()
	//saveCategoriesToCsvFile("categories.csv", categoryTrees)
	//
	////// TODO import data to the Salesforce here
	////// TODO export Categories csv file from Salesforce
	//
	db := createDatabaseConnection()
	var products = readProducts(*db)
	defer db.Close()
	for _, value := range products {
		fmt.Printf("Product %+v \n", value)
	}

	hubs := readSalesforceEntities("migration/production/hubExported.csv", "")
	for key, value := range hubs {
		fmt.Printf("Key '%s' Hub %+v\n", key, value)
	}

	categories := readSalesforceCategories("migration/production/categoryExported.csv")
	for key, value := range categories {
		if strings.Contains(value.Name, "Oxxo: ") {
			fmt.Printf("Key '%s' CategoryTree %+v\n", key, value)
		}
	}

	brands := readSalesforceEntities("migration/production/brandExported.csv", recordTypeBrand)
	for key, value := range brands {
		fmt.Printf("Key '%s' Brand %+v\n", key, value)
	}

	taxRates := readSalesforceEntities("migration/production/taxRateExported.csv", "")
	for key, value := range taxRates {
		fmt.Printf("Key '%s' Tax Rate %+v\n", key, value)
	}

	saveProductsToCsvFile("products.csv", products, taxRates, brands, categories, hubs)
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

func readBrands() []Brand {
	db := createDatabaseConnection()
	var brands []Brand
	err := db.Select(&brands, "SELECT _id, name, image FROM mgo.brands")
	if err != nil {
		log.Fatalf("Could not read brands: %v", err)
	}
	defer db.Close()

	for _, brand := range brands {
		fmt.Printf("Brand %+v\n", brand)
		brandsMap[brand.ID] = brand
	}
	return brands
}

func readCategoryTrees() []CategoryTree {
	db := createDatabaseConnection()
	var categories []CategoryTree
	err := db.Select(&categories,
		`select distinct c._id as id, c.name as name, s._id as subcategory_id, s.name as subcategory_name, c.image as image
					from mgo.categories c
         			join mgo.categories_subcategories cs on c._id = cs.categories_Id
         			join mgo.subcategories s on cs.sub_categories_id = s._id
				union all
				select c._id as id, c.name as name, null as subcategory_id, null as subcategory_name, c.image as image
					from mgo.categories c
					order by id, subcategory_id asc;`)
	defer db.Close()

	if err != nil {
		log.Fatalf("Could not read category trees: %v", err)
	}

	for _, category := range categories {
		fmt.Printf("CategoryTree %+v\n", category)
	}
	return categories
}

func readProducts(db sqlx.DB) []Product {
	var products []Product
	err := db.Select(&products, `SELECT DISTINCT ON (p.sku) p.sku,
                           p.barcode,
                           p.name,
                           p.short_description,
                           p.variant,
                           b.name as brand_name,
                           p.categories,
                           c.name  as category_name,
                           sc.name as subcategory_name,
                           p.image,
                           p.tags,
                           p.taxes,
                           p.age_restriction,
                           p.restricted,
                           p.is_enabled
FROM mgo.products p
         LEFT JOIN mgo.brands b ON p.brand = b._id
         LEFT JOIN LATERAL unnest(p.categories) AS cat(category_id) on true
         JOIN mgo.categories c ON c._id = cat.category_id
         LEFT JOIN LATERAL unnest(p.categories) AS subcat(category_id) on true
         LEFT JOIN mgo.subcategories sc ON sc._id = subcat.category_id 
order by p.sku, category_name, subcategory_name limit 10`)
	if err != nil {
		log.Fatalf("Could not read products: %v", err)
	}
	return products
}

func readBrand(db sqlx.DB, brandId string) Brand {
	var brand []Brand
	err := db.Select(&brand, "SELECT _id, name, image FROM mgo.brands where _id = $1", brandId)
	if err != nil {
		log.Fatalf("Could not read brand by id %s: %v", brandId, err)
	}
	return brand[0]
}

func readSalesforceEntities(filename string, recordTypeId string) map[string]SalesForceEntity {
	data := getDataFromFile(filename)

	var nameColumn int
	var idColumn int
	var referenceTypeIColumn = -1

	result := make(map[string]SalesForceEntity)

	// Print the CSV data
	for i, row := range data {

		if i == 0 {
			for u, header := range row {
				if header == "NAME" {
					nameColumn = u
				}
				if header == "ID" {
					idColumn = u
				}
				if header == "RECORDTYPEID" {
					referenceTypeIColumn = u
				}
			}

			fmt.Printf("Name column index %v\n", nameColumn)
			fmt.Printf("ID column index %v\n", idColumn)
			fmt.Printf("Reference type id column index %v\n", referenceTypeIColumn)

		} else {
			if recordTypeId == "" {
				result[row[nameColumn]] = SalesForceEntity{
					ID:   row[idColumn],
					Name: row[nameColumn],
				}
			} else if row[referenceTypeIColumn] == recordTypeId {
				result[row[nameColumn]] = SalesForceEntity{
					ID:              row[idColumn],
					Name:            row[nameColumn],
					ReferenceTypeId: row[referenceTypeIColumn],
				}
			}
		}
	}
	return result
}

func readSalesforceCategories(filename string) map[string]SalesForceCategory {
	data := getDataFromFile(filename)

	nameColumn := 2
	var idColumn int
	var categoryLevel1Column int
	var categoryLevel2Column int

	result := make(map[string]SalesForceCategory)

	for i, row := range data {

		if i == 0 {
			for u, header := range row {
				if header == "CATEGORY_GROUP__C" {
					nameColumn = u
				}
				if header == "ID" {
					idColumn = u
				}
				if header == "CATEGORY_LEVEL_1__C" {
					categoryLevel1Column = u
				}
				if header == "CATEGORY_LEVEL_2__C" {
					categoryLevel2Column = u
				}
			}

		} else {
			result[row[categoryLevel1Column]+":"+row[categoryLevel2Column]] = SalesForceCategory{
				ID:                 row[idColumn],
				Name:               row[nameColumn],
				CategoryLevel1Name: row[categoryLevel1Column],
				CategoryLevel2Name: row[categoryLevel2Column],
			}
		}
	}
	return result
}

func getDataFromFile(filename string) [][]string {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Could not open file: %v", err)
	}
	defer file.Close()
	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1 // Allow variable number of fields
	data, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}
	return data
}

func saveBrandOwnersToCsvFile(filename string, brands []Brand) {
	file, err := os.Create(filename)
	if err != nil {
		log.Fatalf("Could not create file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Set the headers
	headers := []string{"Name", "RecordTypeId"}
	writer.Write(headers)

	for _, brand := range brands {
		dataRow := make([]string, len(headers))
		dataRow[0] = brand.Name.String
		dataRow[1] = recordTypeBrandOwner
		//dataRow[2] = "Brand Owner"
		writer.Write(dataRow)
	}
}

func saveBrandsToCsvFile(filename string, brands []Brand, brandOwners map[string]SalesForceEntity) {

	file, err := os.Create(filename)
	if err != nil {
		log.Fatalf("Could not create file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Set the headers
	headers := []string{"Name", "RecordTypeId", "Brand_Owner__c"}
	writer.Write(headers)

	for _, brand := range brands {
		dataRow := make([]string, len(headers))
		dataRow[0] = brand.Name.String
		dataRow[1] = recordTypeBrand
		dataRow[2] = brandOwners[brand.Name.String].ID
		writer.Write(dataRow)
	}
}

func saveCategoriesToCsvFile(filename string, categories []CategoryTree) {
	file, err := os.Create(filename)
	if err != nil {
		log.Fatalf("Could not create file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Set the headers
	headers := []string{"Name", "Category_Group__c", "Category_Level_1__c", "Category_Level_2__c"}
	writer.Write(headers)

	for _, category := range categories {
		dataRow := make([]string, len(headers))
		dataRow[0] = "Oxxo: " + clearNameForExport(category.Name.String)
		dataRow[1] = "Oxxo: " + clearNameForExport(category.Name.String)
		dataRow[2] = clearNameForExport(category.Name.String)
		dataRow[3] = clearNameForExport(category.CategoryName.String)

		fmt.Printf("Writing category %+v\n from %+v\n", dataRow, category)

		writer.Write(dataRow)
	}
}

func clearNameForExport(s string) string {
	return strings.ReplaceAll(s, ",", " &")
}

func saveProductsToCsvFile(filename string, products []Product, taxRates map[string]SalesForceEntity, brands map[string]SalesForceEntity, categories map[string]SalesForceCategory, hubs map[string]SalesForceEntity) {

	fmt.Printf("There are %d products\n", len(products))

	file, err := os.Create(filename)
	if err != nil {
		log.Fatalf("Could not create file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Set the headers
	headers := []string{"Barcode__c", "SKU__c", "Name", "Brand__c", "Default_Tax_Rate__c", "Additional_Tax_Rate__c", "Global_Category_Tree__c", "Age_Verification_Required__c",
		"Public_Image_URL__c", "Country__c", "Avalara_Tax_ID__c", "Country_Default_Price__c"}
	writer.Write(headers)

	for _, product := range products {
		dataRow := make([]string, len(headers))
		dataRow[0] = product.Barcode.String
		dataRow[1] = product.Sku
		dataRow[2] = product.Name
		dataRow[3] = brands[product.Brand.String].ID                            // Brand__c
		dataRow[4] = getDefaultTaxRateSalesForceEntity(product, taxRates).ID    // Default_Tax_Rate__c
		dataRow[5] = getAdditionalTaxRateSalesForceEntity(product, taxRates).ID // Additional_Tax_Rate__c
		dataRow[6] = getGlobalCategoryTree(product, categories).ID              // Global_Category_Tree__c
		if product.AgeRestriction {
			dataRow[7] = "Yes - 18+"
		} else {
			dataRow[7] = "No"
		}
		dataRow[8] = product.Image.String
		dataRow[9] = "Mexico" // Country__c
		// dataRow[6] = "" // Avalara_Tax_ID__c
		//dataRow[7] = "" // Country_Default_Price__c

		writer.Write(dataRow)
	}
}

func getFirstElement(input string) string {
	input = strings.Trim(input, "{}")

	if input == "" {
		return ""
	}

	elements := strings.Split(input, ",")

	for _, element := range elements {
		trimmedElement := strings.TrimSpace(element)
		if trimmedElement != "" {
			return trimmedElement // Return the first non-empty element
		}
	}

	return ""
}

func getDefaultTaxRateSalesForceEntity(product Product, taxRates map[string]SalesForceEntity) SalesForceEntity {
	taxRate := getTaxSalesForceEntity(product, taxRates, 0)
	fmt.Printf("Default taxRate %+v set for product %+v \n", taxRate, product)
	return taxRate
}

func getAdditionalTaxRateSalesForceEntity(product Product, taxRates map[string]SalesForceEntity) SalesForceEntity {
	taxRate := getTaxSalesForceEntity(product, taxRates, 1)
	fmt.Printf("Additional taxRate %+v set for product %+v \n", taxRate, product)
	return taxRate
}

func getTaxSalesForceEntity(product Product, taxRates map[string]SalesForceEntity, index int) SalesForceEntity {
	var taxes []Tax
	err := json.Unmarshal([]byte(product.Taxes.String), &taxes)
	if err != nil {
		fmt.Sprintf("Could not unmarshal taxes for product %+v: %+v", product.Taxes, product)
		return SalesForceEntity{}
	}
	if (taxes == nil) || (len(taxes) <= index) {
		return SalesForceEntity{}
	}
	key := fmt.Sprintf("%s%d", taxes[index].Type, taxes[index].Rate)

	return taxRates[key]
}

func getGlobalCategoryTree(product Product, categories map[string]SalesForceCategory) SalesForceCategory {
	cetegoryTree := categories[clearNameForExport(product.CategoryName.String)+":"+clearNameForExport(product.SubcategoryName.String)]
	return cetegoryTree
}
