package main

import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/xuri/excelize/v2"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const (
	environment = "test"
	// PostgreSQL connection string
	psqlInfo              = "host=mioxxo-products.ce989v8mdple.us-east-1.rds.amazonaws.com port=5432 user=postgres dbname=products password=Yk9deWhjbYUUR5LEF8SnMf7w9jghPf sslmode=disable"
	recordTypeEntity      = "012Hp000001mPmCIAU"
	recordTypeBrandOwner  = "012Hp000001mPmGIAU"
	recordTypeBrand       = "012Hp000001mPmEIAU"
	recordTypeVendor      = "012Hp000001mPmDIAU"
	recordTypeHub         = "012Hp000001mPmBIAU"
	jokrEntity            = "001D400000p9IxFIAU" // "001Hp00002kvRgtIAE"
	customerAccessAccount = "001D400000p9LdVIAU" // "001Hp00002kuyTBIAY"
	areaAccessAccount     = "001D400000p9LdWIAU"
	integrationOfBtlr     = "005Hp00000igLP3IAM"
)

var brandsMap = make(map[string]Brand)

type ProductChild struct {
	Sku      string
	Quantity string
}

type ProductSupplier struct {
	Item         string
	SupplierId   int
	SupplierName string
}

type Store struct {
	Name       string
	Code       string
	Id         string
	Plaza      string
	Address    string
	PostalCode string
	City       string
	State      string
	Latitude   string
	Longitude  string
	Surface    string
	Area       string
}

type ProductSupplierKey struct {
	Item       string
	SupplierId int
}

type Brand struct {
	ID    string         `db:"_id"`
	Name  sql.NullString `db:"name"`
	Image sql.NullString `db:"image"`
}

type Hub struct {
	ID     string `db:"_id"`
	Name   string `db:"name"`
	Status string `db:"status"`
	ZoneId string `db:"zone_id"`
	Type   string `db:"type"`
	Code   string `db:"cr"`
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

type SalesForceAccount struct {
	ID              string
	Name            string
	BusinessName    string
	InternalId      string
	ReferenceTypeId string
	Status          string
}

type SalesForceProduct struct {
	ID   string
	Name string
	Sku  string
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

type TaxRate struct {
	Name      string
	Rate      string
	Type      string
	Secondary bool
}

type Product struct {
	Sku              string         `db:"sku"`
	Barcode          sql.NullString `db:"barcode"`
	Name             string         `db:"name"`
	ShortDescription string         `db:"short_description"`
	Variant          string         `db:"variant"`
	MinPrice         float32        `db:"min_price"`
	MaxPrice         float32        `db:"max_price"`
	MaxQuantity      int            `db:"max_quantity"`
	Brand            sql.NullString `db:"brand_name"`
	Categories       sql.NullString `db:"categories"`
	CategoryName     sql.NullString `db:"category_name"`
	SubcategoryName  sql.NullString `db:"subcategory_name"`
	Image            sql.NullString `db:"image"`
	Tags             sql.NullString `db:"tags"`
	Taxes            sql.NullString `db:"taxes"`
	Children         sql.NullString `db:"children"`
	AgeRestriction   bool           `db:"age_restriction"`
	Restricted       bool           `db:"restricted"`
	Enabled          bool           `db:"is_enabled"`
	Discount         sql.NullString `db:"discount"`
	DiscountType     sql.NullString `db:"discount_type"`
	DiscountTitle    sql.NullString `db:"discount_title"`
}

type Inventory struct {
	Sku       string         `db:"sku"`
	Store     string         `db:"store"`
	Zone      sql.NullString `db:"zone"`
	StoreCode string         `db:"cr"`
	Price     float32        `db:"price_list"`
}

func main() {
	//var brands = readBrands()
	//saveBrandOwnersToCsvFile("brandowners.csv", brands)

	//// TODO import data to the Salesforce here
	//// TODO export Brands csv file from Salesforce

	//brandOwners := readSalesforceEntities(fmt.Sprintf("migration/%s/brandExported.csv", environment), recordTypeBrandOwner)
	//for key, value := range brandOwners {
	//	fmt.Printf("Key '%s' Brand owner %+v\n", key, value)
	//}
	//
	//saveBrandsToCsvFile("brands.csv", brands, brandOwners)

	//// TODO import data to the Salesforce here
	//// TODO export Brands csv file from Salesforce

	//var categoryTrees = readCategoryTrees()
	//saveCategoriesToCsvFile("categories.csv", categoryTrees)

	//saveUtilityEntitiesToCsvFile("utilityEntities.csv")

	////// TODO import data to the Salesforce here
	////// TODO export Categories csv file from Salesforce

	//taxRates := readTaxRates()
	//saveTaxRatesToCsvFile("taxrates.csv", taxRates)

	// TODO import data to the Salesforce here
	// TODO export Tax Rates csv file from Salesforce

	//var products = readProducts()
	//
	//for _, value := range products {
	//	fmt.Printf("Product %+v \n", value)
	//}
	//
	//categories := readSalesforceCategories(fmt.Sprintf("migration/%s/categoryExported.csv", environment))
	//for key, value := range categories {
	//	if strings.Contains(value.Name, "Oxxo: ") {
	//		fmt.Printf("Key '%s' CategoryTree %+v\n", key, value)
	//	}
	//}
	//
	//brands := readSalesforceEntities(fmt.Sprintf("migration/%s/brandExported.csv", environment), recordTypeBrand)
	//for key, value := range brands {
	//	fmt.Printf("Key '%s' Brand %+v\n", key, value)
	//}
	//
	//taxRates := readSalesforceEntities(fmt.Sprintf("migration/%s/taxRateExported.csv", environment), "")
	//for key, value := range taxRates {
	//	fmt.Printf("Key '%s' Tax Rate %+v\n", key, value)
	//}
	//
	//saveProductsToCsvFile("products.csv", products, taxRates, brands, categories)

	//productSuppliers, suppliers := readSuppliers()
	//
	//productSuppliers = filterUniqueProductSuppliers(productSuppliers)
	//
	//for _, value := range productSuppliers {
	//	fmt.Printf("Product suppliers %+v\n", value)
	//}
	//
	//for key, value := range suppliers {
	//	fmt.Printf("Supplier Key %+v Value %+v\n", key, value)
	//}
	//
	//saveVendorsToCsvFile("vendors.csv", suppliers)

	//_, suppliers := readSuppliers()

	// TODO import data to the Salesforce here
	// TODO export Account and Product csv file from Salesforce

	//salesForceVendors := readSalesforceAccounts(fmt.Sprintf("migration/%s/accountExported.csv", environment), recordTypeVendor, "PARTNER_CODE__C")
	//
	//salesForceProducts := readSalesforceProducts(fmt.Sprintf("migration/%s/productExported.csv", environment))
	//
	//for key, value := range products {
	//	fmt.Printf("Key %s  Product %+v\n", key, value)
	//}
	//
	productSuppliers, _ := readSuppliers()
	productSuppliers = filterUniqueProductSuppliers(productSuppliers)
	//saveVendorProductsToCsvFile("vendorproducts.csv", productSuppliers, products, salesForceVendors, salesForceProducts)

	// TODO import data to the Salesforce here
	// TODO export Accounts csv file from Salesforce

	hubs := readHubs()

	for _, hub := range hubs {
		fmt.Printf("Hub %+v \n", hub)
	}

	_, storesMap := readStores()
	for key, store := range storesMap {
		fmt.Printf("Key %s Store %+v \n", key, store)
	}

	//saveHubAccountsToCsvFile("hubaccounts.csv", hubs, storesMap)

	// TODO import data to the Salesforce here
	// TODO export Accounts csv file from Salesforce

	salesForceHubAccounts := readSalesforceAccounts(fmt.Sprintf("migration/%s/accountExported.csv", environment), recordTypeHub, "NAME")

	//
	//fmt.Printf("There are %d hub accounts \n", len(salesForceHubAccounts))
	//for key, value := range salesForceHubAccounts {
	//	fmt.Printf("Key %s  Hub %+v\n", key, value)
	//}
	//inventory := readInventory()

	//for _, value := range inventory {
	//	fmt.Printf("Inventory %+v\n", value)
	//}

	//saveHubsToCsvFile("hubs.csv", salesForceHubAccounts, storesMap)

	//salesForceVendorProducts := readSalesforceEntities(fmt.Sprintf("migration/%s/vendorProductExported.csv", environment), "")
	salesForceHubs := readSalesforceEntities(fmt.Sprintf("migration/%s/hubExported.csv", environment), "")
	//saveHubProductsToCsvFile("hubproducts.csv", salesForceHubs, productSuppliers, inventory, salesForceProducts, salesForceVendorProducts)

	updateHubAccountsToCsvFile("hubaccounts-update.csv", salesForceHubAccounts, storesMap)
	updateHubToCsvFile("hubs-update.csv", salesForceHubs, storesMap)

}

func readTaxRates() []TaxRate {
	result := make([]TaxRate, 0)
	result = append(result, TaxRate{
		Name:      "IVA0",
		Rate:      "0.0",
		Type:      "IVA",
		Secondary: false,
	})
	result = append(result, TaxRate{
		Name:      "IVA16",
		Rate:      "16.0",
		Type:      "IVA",
		Secondary: false,
	})
	result = append(result, TaxRate{
		Name:      "IEPS0",
		Rate:      "0.0",
		Type:      "IEPS",
		Secondary: true,
	})
	result = append(result, TaxRate{
		Name:      "IEPS8",
		Rate:      "8.0",
		Type:      "IEPS",
		Secondary: true,
	})
	result = append(result, TaxRate{
		Name:      "IEPS25",
		Rate:      "25.0",
		Type:      "IEPS",
		Secondary: true,
	})
	result = append(result, TaxRate{
		Name:      "IEPS53",
		Rate:      "53.0",
		Type:      "IEPS",
		Secondary: true,
	})
	return result
}

func filterUniqueProductSuppliers(productSuppliers []ProductSupplier) []ProductSupplier {
	uniqueMap := make(map[ProductSupplierKey]bool)
	for _, productSupplier := range productSuppliers {
		key := ProductSupplierKey{
			Item:       productSupplier.Item,
			SupplierId: productSupplier.SupplierId,
		}
		uniqueMap[key] = true
	}

	result := make([]ProductSupplier, 0, len(uniqueMap))
	for key := range uniqueMap {
		result = append(result, ProductSupplier{
			Item:         key.Item,
			SupplierId:   key.SupplierId,
			SupplierName: "", // SupplierName is not used for uniqueness
		})
	}
	return result
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
	order by id, subcategory_id asc; `)
	defer db.Close()

	if err != nil {
		log.Fatalf("Could not read category trees: %v", err)
	}

	for _, category := range categories {
		fmt.Printf("CategoryTree %+v\n", category)
	}
	return categories
}

func readProducts() []Product {
	db := createDatabaseConnection()
	var products []Product
	err := db.Select(&products, `SELECT DISTINCT ON (p.sku) p.sku,
                           p.barcode,
                           p.name,
                           p.short_description,
                           p.variant,
                           (select min(price_list) as min_price from mgo.products_inventory where sku = p.sku and zone is not null),
                           (select max(price_list) as max_price from mgo.products_inventory where sku = p.sku and zone is not null),
                           (select max(max_per_order) as max_quantity from mgo.products_inventory where sku = p.sku and zone is not null),
                           b.name  as brand_name,
                           p.categories,
                           c.name  as category_name,
                           sc.name as subcategory_name,
                           p.image,
                           p.tags,
                           p.taxes,
                           p.children,
                           p.age_restriction,
                           p.restricted,
                           p.is_enabled
							FROM mgo.products p
         					LEFT JOIN mgo.brands b ON p.brand = b._id
         					LEFT JOIN LATERAL unnest(p.categories) AS cat(category_id) on true
         					JOIN mgo.categories c ON c._id = cat.category_id
         					LEFT JOIN LATERAL unnest(p.categories) AS subcat(category_id) on true
         					LEFT JOIN mgo.subcategories sc ON sc._id = subcat.category_id
         					WHERE (c.name, coalesce(sc.name, '')) in (select c.name, s.name
         					from mgo.categories c
         					join mgo.categories_subcategories cs on c._id = cs.categories_Id
         					join mgo.subcategories s on cs.sub_categories_id = s._id
         					union all
         					select c.name as name, '' as subcategory_name from mgo.categories c where c.name != 'Promociones')
							order by p.sku, category_name, subcategory_name desc;`)
	defer db.Close()
	if err != nil {
		log.Fatalf("Could not read products: %v", err)
	}
	return products
}

func readInventory() []Inventory {
	db := createDatabaseConnection()
	var inventories []Inventory
	err := db.Select(&inventories, `select sku, store, zone, cr, price_list from mgo.products_inventory pi
												LEFT JOIN mgo.products_inventory_stores pis on pi.store = pis._id 
												WHERE cr is not null`)
	if err != nil {
		log.Fatalf("Could not read inventories: %v", err)
	}
	defer db.Close()
	return inventories
}

func readHubs() []Hub {
	db := createDatabaseConnection()
	var hubs []Hub
	err := db.Select(&hubs, `select _id,name,status,zone_id,cr,type from mgo.products_inventory_stores where cr is not null`)
	if err != nil {
		log.Fatalf("Could not read hubs: %v", err)
	}
	defer db.Close()
	return hubs
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
				if strings.EqualFold(header, "Id") {
					idColumn = u
				}
				if header == "RECORDTYPEID" {
					referenceTypeIColumn = u
				}
			}

			fmt.Printf("Name column index %v\n", nameColumn)
			fmt.Printf("Id column index %v\n", idColumn)
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
				if header == "Id" {
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

func readSalesforceAccounts(filename string, recordTypeId string, keyHeader string) map[string]SalesForceAccount {
	data := getDataFromFile(filename)

	var nameColumn int
	var idColumn int
	var externalIdColumn int
	var referenceTypeColumn int
	var businessNameColumn int
	var statusColumn int

	result := make(map[string]SalesForceAccount)

	for i, row := range data {

		if i == 0 {
			for u, header := range row {
				if header == "NAME" {
					nameColumn = u
				}
				if strings.EqualFold(header, "Id") {
					idColumn = u
				}
				if header == keyHeader {
					externalIdColumn = u
				}
				if header == "RECORDTYPEID" {
					referenceTypeColumn = u
				}
				if header == "BUSINESS_NAME__C" {
					businessNameColumn = u
				}
				if header == "STATUS__C" {
					statusColumn = u
				}
			}

		} else {
			if row[referenceTypeColumn] == recordTypeId {
				result[row[externalIdColumn]] = SalesForceAccount{
					ID:              row[idColumn],
					Name:            row[nameColumn],
					InternalId:      row[externalIdColumn],
					ReferenceTypeId: row[referenceTypeColumn],
					BusinessName:    row[businessNameColumn],
					Status:          row[statusColumn],
				}
			}
		}
	}
	return result
}

func readSalesforceProducts(filename string) map[string]SalesForceProduct {
	data := getDataFromFile(filename)

	var nameColumn int
	var idColumn int
	var skuColumn int

	result := make(map[string]SalesForceProduct)

	for i, row := range data {
		if i == 0 {
			for u, header := range row {
				if header == "NAME" {
					nameColumn = u
				}
				if header == "Id" {
					idColumn = u
				}
				if header == "SKU__C" {
					skuColumn = u
				}
			}
		} else {
			result[row[skuColumn]] = SalesForceProduct{
				ID:   row[idColumn],
				Name: row[nameColumn],
				Sku:  row[skuColumn],
			}
		}
	}

	fmt.Printf("There are %d products \n", len(result))

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

func saveUtilityEntitiesToCsvFile(filename string) {
	file, err := os.Create(filename)
	if err != nil {
		log.Fatalf("Could not create file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Set the headers
	headers := []string{"Name", "BILLINGCOUNTRY", "SHIPPINGCOUNTRY", "RecordTypeId"}
	writer.Write(headers)

	var dataRow []string

	dataRow = make([]string, len(headers))
	dataRow[0] = "MiOXXO Mexico"
	dataRow[1] = "Mexico"
	dataRow[2] = "Mexico"
	dataRow[3] = recordTypeEntity
	writer.Write(dataRow)

	dataRow = make([]string, len(headers))
	dataRow[0] = "CustomerSuccessAccount"
	dataRow[1] = "Mexico"
	dataRow[2] = "Mexico"
	dataRow[3] = recordTypeEntity
	writer.Write(dataRow)

	dataRow = make([]string, len(headers))
	dataRow[0] = "Area Account"
	dataRow[1] = "Mexico"
	dataRow[2] = "Mexico"
	dataRow[3] = recordTypeEntity
	writer.Write(dataRow)
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

func saveHubAccountsToCsvFile(filename string, hubs []Hub, stores map[string]Store) {
	file, err := os.Create(filename)
	if err != nil {
		log.Fatalf("Could not create file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Set the headers
	headers := []string{"Name", "RecordTypeId", "JOKR_ENTITY__C", "DEFAULT_ORDER_CYCLE__C", "STATUS__C", "BUSINESS_NAME__C", "PARENTID", "COUNTRY__C",
		"BILLINGCITY", "BILLINGCOUNTRY", "BILLINGPOSTALCODE", "BILLINGSTATE", "BILLINGSTREET", "SHIPPINGCITY", "SHIPPINGCOUNTRY", "SHIPPINGPOSTALCODE", "SHIPPINGSTATE", "SHIPPINGSTREET"}
	writer.Write(headers)

	for _, hub := range hubs {
		dataRow := make([]string, len(headers))
		dataRow[0] = hub.Code
		dataRow[1] = recordTypeHub
		dataRow[2] = jokrEntity
		dataRow[3] = "Every Week"
		if hub.Status == "active" {
			dataRow[4] = "Active"
		} else {
			dataRow[4] = "Inactive"
		}
		dataRow[5] = hub.Name
		dataRow[6] = jokrEntity // MiOxxo Mexico Account
		dataRow[7] = "Mexico"
		dataRow[8] = stores[hub.Code].City
		dataRow[9] = "Mexico"
		dataRow[10] = stores[hub.Code].PostalCode
		dataRow[11] = stores[hub.Code].State
		dataRow[12] = stores[hub.Code].Address
		dataRow[13] = stores[hub.Code].City
		dataRow[14] = "Mexico"
		dataRow[15] = stores[hub.Code].PostalCode
		dataRow[16] = stores[hub.Code].State
		dataRow[17] = stores[hub.Code].Address
		writer.Write(dataRow)
	}
}

func updateHubAccountsToCsvFile(filename string, hubs map[string]SalesForceAccount, stores map[string]Store) {
	file, err := os.Create(filename)
	if err != nil {
		log.Fatalf("Could not create file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Set the headers
	headers := []string{"ID", "NAME", "COUNTRY__C", "BILLINGCITY", "BILLINGCOUNTRY", "BILLINGPOSTALCODE", "BILLINGSTATE", "BILLINGSTREET", "SHIPPINGCITY", "SHIPPINGCOUNTRY", "SHIPPINGPOSTALCODE", "SHIPPINGSTATE", "SHIPPINGSTREET", "STATE_REGISTRATION__C", "CITY_REGISTRATION__C"}
	writer.Write(headers)

	for key, hub := range hubs {
		dataRow := make([]string, len(headers))
		dataRow[0] = hub.ID
		dataRow[1] = hub.Name
		//dataRow[2] = "Mexico"
		dataRow[3] = stores[key].City
		dataRow[4] = "Mexico"
		dataRow[5] = stores[key].PostalCode
		dataRow[6] = stores[key].State
		dataRow[7] = stores[key].Address
		dataRow[8] = stores[key].City
		dataRow[9] = "Mexico"
		dataRow[10] = stores[key].PostalCode
		dataRow[11] = stores[key].State
		dataRow[12] = stores[key].Address
		dataRow[13] = stores[key].State
		dataRow[14] = stores[key].City
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
		//dataRow[0] = "Oxxo: " + clearNameForExport(category.Name.String)
		//dataRow[1] = "Oxxo: " + clearNameForExport(category.Name.String)
		dataRow[0] = clearNameForExport(category.Name.String)
		dataRow[1] = clearNameForExport(category.Name.String)
		dataRow[2] = clearNameForExport(category.Name.String)
		dataRow[3] = clearNameForExport(category.CategoryName.String)

		fmt.Printf("Writing category %+v\n from %+v\n", dataRow, category)

		writer.Write(dataRow)
	}
}

func saveTaxRatesToCsvFile(filename string, taxRates []TaxRate) {
	file, err := os.Create(filename)
	if err != nil {
		log.Fatalf("Could not create file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Set the headers
	headers := []string{"NAME", "NAME__C", "Country__c", "PERCENTAGE__C", "SECONDARY_TAX__C", "JOKR_ENTITY__C"}
	writer.Write(headers)

	for _, taxRate := range taxRates {
		dataRow := make([]string, len(headers))
		dataRow[0] = taxRate.Name
		dataRow[1] = taxRate.Name
		dataRow[2] = "Mexico"
		dataRow[3] = taxRate.Rate
		if taxRate.Secondary {
			dataRow[4] = "TRUE"
		} else {
			dataRow[4] = "FALSE"
		}
		dataRow[5] = jokrEntity

		fmt.Printf("Writing Tax Rate %+v\n from %+v\n", dataRow, taxRate)

		writer.Write(dataRow)
	}
}

func clearNameForExport(s string) string {
	return strings.ReplaceAll(s, ",", " &")
}

func saveProductsToCsvFile(filename string, products []Product, taxRates map[string]SalesForceEntity, brands map[string]SalesForceEntity, categories map[string]SalesForceCategory) {

	fmt.Printf("There are %d products\n", len(products))

	file, err := os.Create(filename)
	if err != nil {
		log.Fatalf("Could not create file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Set the headers
	headers := []string{"Barcode__c", "SKU__c", "Name", "Product_Title__c", "Brand__c", "Default_Tax_Rate__c", "Additional_Tax_Rate__c", "Global_Category_Tree__c", "Age_Verification_Required__c",
		"Public_Image_URL__c", "JOKR_Entity__c", "Product_Ownership__c", "Country_Default_Price__c", "UI_Content_1__c", "UI_Content_1_UOM__c", "Assortment_Type__c", "Maximum_quantity__C", "Minimum_quantity__C",
		"Hub_Temperature_Storage__C", "Expiry_Date_Tracking_Required__c", "Reception_Temperature_Tracking_required__c", "UOM2__c", "Customer_Fulfillment_Type__c", "CRWN_Content_Ready__c", "Alcoholic_Beverage__c",
		"Alcohol_by_Volume__c"}

	writer.Write(headers)

	for _, product := range products {
		dataRow := make([]string, len(headers))
		dataRow[0] = product.Barcode.String
		dataRow[1] = product.Sku
		dataRow[2] = product.Name
		dataRow[3] = product.Name
		dataRow[4] = brands[product.Brand.String].ID
		defaultTaxRate, additionalTaxRate := getDefaultTaxRateSalesForceEntity(product, taxRates)
		dataRow[5] = defaultTaxRate.ID                             // Default_Tax_Rate__c
		dataRow[6] = additionalTaxRate.ID                          // Additional_Tax_Rate__c
		dataRow[7] = getGlobalCategoryTree(product, categories).ID // Global_Category_Tree__c
		if product.AgeRestriction {
			dataRow[8] = "Yes - 18+"
		} else {
			dataRow[8] = "No"
		}
		dataRow[9] = product.Image.String
		dataRow[10] = jokrEntity // JOKR_Entity__с
		dataRow[11] = "JOKR owned"
		dataRow[12] = fmt.Sprintf("%.2f", product.MaxPrice/100) // Country_Default_Price__c

		amount, measurement := getVariantFields(product.Variant)
		dataRow[13] = amount
		dataRow[14] = measurement

		_, assortmentType := getProductChildrenAdAssortmentType(product)
		dataRow[15] = assortmentType // Assortment_Type__c
		dataRow[16] = fmt.Sprintf("%d", product.MaxQuantity)
		dataRow[17] = "1"                // Minimum_quantity__C
		dataRow[18] = "Room Temperature" // Hub_Temperature_Storage__C  Room Temperature, Refrigerated, Frozen
		dataRow[19] = "No"
		dataRow[20] = "No"
		dataRow[21] = getPackSize(product.Name) //Underlying barcode quantity e.g. ea, pk-3, pk-4
		dataRow[22] = "Core"                    // Customer_Fulfillment_Type__c
		if product.Restricted {
			dataRow[24] = "Yes"
		}

		writer.Write(dataRow)
	}
}

func saveVendorsToCsvFile(filename string, vendors map[int]ProductSupplier) {
	file, err := os.Create(filename)
	if err != nil {
		log.Fatalf("Could not create file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Set the headers
	headers := []string{"Name", "Business_name__C", "RecordTypeId", "JOKR_Entity__c", "Portal_url__c", "Status__c", "Spot_Buy_Vendor__c", "Commercial_Owner__c", "PO_Reception_Method__c", "Finance_contact_email__c",
		"Finance_contact_full_name__c", "Finance_contact_phone__c", "Default_delivery_time_setting__c", "Default_lead_time__c", "Country__c", "Partner_code__c", "Location_scope__c", "BILLINGCOUNTRY", "SHIPPINGCOUNTRY"}
	writer.Write(headers)

	uniqueVendorNames := make(map[string]bool)

	i := 0
	for _, vendor := range vendors {
		//if strings.Contains(vendor.SupplierName, "N/A") {
		//	continue
		//}
		//if i < 10 {
		name := strings.ReplaceAll(vendor.SupplierName, ",", ".")

		if uniqueVendorNames[name] {
			newName := fmt.Sprintf("%s - %d", name, vendor.SupplierId)
			fmt.Printf("Duplicate vendor name %s Name will be %s\n", name, newName)
			name = newName
		} else {
			uniqueVendorNames[name] = true
		}

		dataRow := make([]string, len(headers))
		dataRow[0] = name
		dataRow[1] = name
		dataRow[2] = recordTypeVendor
		dataRow[3] = jokrEntity                        // JOKR_Entity__с
		dataRow[5] = "Active"                          // status__c
		dataRow[6] = "No"                              // Spot_Buy_Vendor__c
		dataRow[7] = integrationOfBtlr                 // Commercial_Owner__c
		dataRow[8] = "N/A"                             // PO_Reception_Method__c
		dataRow[9] = "dmitriy.yavorskiy@icemobile.com" // Finance_contact_email__c
		dataRow[10] = "Dmitriy"                        // Finance_contact_full_name__c
		dataRow[11] = "1234567890"                     // Finance_contact_phone__c
		dataRow[12] = "Calendar Days"                  // Default_delivery_time__setting_c
		dataRow[13] = "1"                              // Default_lead_time__c
		//dataRow[14] = "Mexico"                             // Country__c
		dataRow[15] = fmt.Sprintf("%d", vendor.SupplierId) // Partner_code__c
		dataRow[16] = "Country-wide"                       // Location_scope__c
		dataRow[17] = "Mexico"
		dataRow[18] = "Mexico"

		writer.Write(dataRow)
		i++

		//}
	}
}

func saveVendorProductsToCsvFile(filename string, productSuppliers []ProductSupplier, products []Product, salesForceVendors map[string]SalesForceAccount, salesForceProducts map[string]SalesForceProduct) {
	file, err := os.Create(filename)
	if err != nil {
		log.Fatalf("Could not create file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Set the headers
	headers := []string{"Name", "Account__c", "Product__c", "PURCHASING_STATUS__C", "LAST_COST__C", "MINIMUM_QUANTITY__C"}
	writer.Write(headers)

	productsMap := convertProductsToMap(products)

	uniqueNames := make(map[string]bool)

	i := 0
	for _, vendorProduct := range productSuppliers {
		if (salesForceProducts[vendorProduct.Item] == SalesForceProduct{}) {
			fmt.Printf("Product not found %+v\n", vendorProduct)
			continue
		}
		//if i < 10 {
		name := fmt.Sprintf("%s - %d", vendorProduct.Item, vendorProduct.SupplierId)

		dataRow := make([]string, len(headers))
		dataRow[0] = name
		dataRow[1] = salesForceVendors[strconv.Itoa(vendorProduct.SupplierId)].ID
		dataRow[2] = salesForceProducts[vendorProduct.Item].ID
		dataRow[3] = "Available to Purchase"
		dataRow[4] = fmt.Sprintf("%.2f", productsMap[vendorProduct.Item].MinPrice/100)
		dataRow[5] = "1"
		writer.Write(dataRow)
		i++
		//}

		if uniqueNames[name] {
			fmt.Printf("Duplicate vendor product %s\n", name)
		} else {
			uniqueNames[name] = true
		}
	}

	fmt.Printf("There are %d/%d product salesForceVendors\n", i, len(productSuppliers))
}

func saveHubProductsToCsvFile(filename string, salesforceHubs map[string]SalesForceEntity, productSuppliers []ProductSupplier, inventory []Inventory, salesForceProducts map[string]SalesForceProduct, salesForceVendorProducts map[string]SalesForceEntity) {
	file, err := os.Create(filename)
	if err != nil {
		log.Fatalf("Could not create file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Set the headers
	headers := []string{"Name", "HUB__C", "MAIN_VENDOR_PRODUCT__C", "RETAIL_PRICE__C", "HUB_SKU_TEXT__C", "PRODUCT__C"}
	writer.Write(headers)

	uniqueVendorProductMap := make(map[string]ProductSupplier)
	for _, vendorProduct := range productSuppliers {
		uniqueVendorProductMap[vendorProduct.Item] = vendorProduct
	}

	hubSkuMap := make(map[string]Inventory)

	i := 0
	for _, inventory := range inventory {

		vendorProductName := fmt.Sprintf("%s - %d", inventory.Sku, uniqueVendorProductMap[inventory.Sku].SupplierId)
		hubSku := fmt.Sprintf("%s - %s", inventory.StoreCode, inventory.Sku)

		if (salesForceVendorProducts[vendorProductName] == SalesForceEntity{}) {
			fmt.Printf("Vendor Product %s not found %+v\n", vendorProductName, inventory)
			continue
		}

		if (hubSkuMap[hubSku] != Inventory{}) {
			fmt.Printf("Duplicated Hub Product %s not found %+v\n", hubSku, inventory)
			continue
		}
		hubSkuMap[hubSku] = inventory

		//if i >= 20 {
		dataRow := make([]string, len(headers))
		dataRow[0] = fmt.Sprintf("%s - %d - %s", inventory.Sku, uniqueVendorProductMap[inventory.Sku].SupplierId, inventory.StoreCode)
		dataRow[1] = salesforceHubs[inventory.StoreCode].ID
		dataRow[2] = salesForceVendorProducts[vendorProductName].ID
		dataRow[3] = fmt.Sprintf("%.2f", inventory.Price/100)
		dataRow[4] = hubSku
		dataRow[5] = salesForceProducts[inventory.Sku].ID
		writer.Write(dataRow)
		//}
		i++
	}

	fmt.Printf("There are %d/%d hub products\n", i, len(inventory))
}

func saveHubsToCsvFile(filename string, hubs map[string]SalesForceAccount, stores map[string]Store) {
	file, err := os.Create(filename)
	if err != nil {
		log.Fatalf("Could not create file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Set the headers
	headers := []string{"Name", "JOKR_ENTITY__C", "HUB_TYPE__C", "HUB_ACCOUNT__C", "HUB_NAME__C", "REGIONAL_ACCOUNT__C", "AREA_ACCOUNT__C", "HUB_STATUS__C", "HUB_COUNTRY__C", "HUB_CITY__C", "HUB_GEOLOCATION__LATITUDE__S", "HUB_GEOLOCATION__LONGITUDE__S",
		"HUB_POSTAL_CODE__C", "HUB_STATE__C", "HUB_STREET__C"}
	writer.Write(headers)

	for _, hub := range hubs {
		dataRow := make([]string, len(headers))
		dataRow[0] = hub.Name
		dataRow[1] = jokrEntity
		dataRow[2] = "Hub"
		dataRow[3] = hub.ID
		dataRow[4] = hub.BusinessName
		dataRow[5] = customerAccessAccount // Customer access account
		dataRow[6] = areaAccessAccount     // Area access account
		if strings.EqualFold(hub.Status, "Active") {
			dataRow[7] = "Active"
		} else {
		}
		dataRow[8] = "Mexico"
		dataRow[9] = stores[hub.Name].City
		dataRow[10] = stores[hub.Name].Latitude
		dataRow[11] = stores[hub.Name].Longitude
		dataRow[12] = stores[hub.Name].PostalCode
		dataRow[13] = stores[hub.Name].State
		dataRow[14] = stores[hub.Name].Address
		writer.Write(dataRow)
	}
}

func updateHubToCsvFile(filename string, hubs map[string]SalesForceEntity, stores map[string]Store) {
	file, err := os.Create(filename)
	if err != nil {
		log.Fatalf("Could not create file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Set the headers
	headers := []string{"ID", "Name", "HUB_COUNTRY__C", "HUB_CITY__C", "HUB_GEOLOCATION__LATITUDE__S", "HUB_GEOLOCATION__LONGITUDE__S", "HUB_POSTAL_CODE__C", "HUB_STATE__C", "HUB_STREET__C"}
	writer.Write(headers)

	for _, hub := range hubs {
		dataRow := make([]string, len(headers))
		dataRow[0] = hub.ID
		dataRow[1] = hub.Name
		dataRow[2] = "Mexico"
		dataRow[3] = stores[hub.Name].City
		dataRow[4] = stores[hub.Name].Latitude
		dataRow[5] = stores[hub.Name].Longitude
		dataRow[6] = stores[hub.Name].PostalCode
		dataRow[7] = stores[hub.Name].State
		dataRow[8] = stores[hub.Name].Address
		writer.Write(dataRow)
	}
}

func getVariantFields(variant string) (string, string) {
	variant = strings.ReplaceAll(variant, ",", ".")
	variant = strings.ReplaceAll(variant, "..", ".")
	variant = strings.ReplaceAll(variant, "17*3", "51")
	variant = insertSpace(variant)

	variantList := strings.Split(variant, " ")

	amount := variantList[0]

	var measurement string
	if len(variantList) > 1 {
		if contains([]string{"tabs", "Tabs"}, variantList[1]) {
			measurement = "Tabs"
		} else if contains([]string{"rollo", "rollos", "Rollo", "Rollos"}, variantList[1]) {
			measurement = "Rollos"
		} else {
			measurement = variantList[1]
		}
	}

	return amount, measurement
}

func insertSpace(s string) string {
	re := regexp.MustCompile(`([0-9]+)([a-zA-Z]+)`)
	return re.ReplaceAllString(s, `$1 $2`)
}

func notContains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return false
		}
	}
	return true
}

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

func getPackSize(productName string) string {
	var quantity int
	_, err := fmt.Sscanf(productName, "%d", &quantity)
	if err != nil || !strings.Contains(productName, " Pack ") {
		return "ea"
	}
	result := fmt.Sprintf("pk-%d", quantity)
	//fmt.Printf("Product is not ea %+v. Pack code is %+v\n", productName, result)
	return result
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

func getDefaultTaxRateSalesForceEntity(product Product, taxRates map[string]SalesForceEntity) (SalesForceEntity, SalesForceEntity) {
	firstTaxRate := getTaxSalesForceEntity(product, taxRates, 0)
	if strings.Contains(firstTaxRate.Name, "IEPS") {
		defaultTaxRate, err := getTaxSalesForceEntityByName(taxRates, "IVA0")
		if err != nil {
			log.Fatal(err)
		}
		//fmt.Printf("Default Tax rate set to  %+v and additional to %+v for product %+v \n", defaultTaxRate.Name, firstTaxRate.Name, product)
		return defaultTaxRate, firstTaxRate
	}
	additionalTaxRate := getTaxSalesForceEntity(product, taxRates, 1)
	//fmt.Printf("Default Tax rate set to  %+v and additional to %+v for product %+v \n", firstTaxRate.Name, additionalTaxRate.Name, product)

	if firstTaxRate.Name == additionalTaxRate.Name {
		additionalTaxRate = SalesForceEntity{}
	}

	if strings.HasPrefix(firstTaxRate.Name, "IVA0") && strings.HasPrefix(additionalTaxRate.Name, "IVA") {
		return additionalTaxRate, SalesForceEntity{}
	}

	return firstTaxRate, additionalTaxRate
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

	if len(taxes) > 2 {
		fmt.Sprintf("There are more then 2 taxes for the product: %+v", product)
	}

	key := fmt.Sprintf("%s%d", taxes[index].Type, taxes[index].Rate)

	return taxRates[key]
}

func getProductChildrenAdAssortmentType(product Product) ([]ProductChild, string) {
	var children []ProductChild
	err := json.Unmarshal([]byte(product.Children.String), &children)
	if err != nil {
		fmt.Sprintf("Could not unmarshal children for product %+v: %+v", product, product)
		return []ProductChild{}, "Standard"
	}

	if len(children) > 0 {
		fmt.Printf("Product %+v has children %+v\n", product.Name, children)
		return children, "Bundle"
	}

	return []ProductChild{}, "Standard"
}

func getTaxSalesForceEntityByName(taxRates map[string]SalesForceEntity, name string) (SalesForceEntity, error) {
	for _, taxRate := range taxRates {
		if taxRate.Name == name {
			return taxRate, nil
		}
	}
	return SalesForceEntity{}, errors.New("ax rate with the given name not found")
}

func getGlobalCategoryTree(product Product, categories map[string]SalesForceCategory) SalesForceCategory {
	cetegoryTree := categories[clearNameForExport(product.CategoryName.String)+":"+clearNameForExport(product.SubcategoryName.String)]
	return cetegoryTree
}

func convertProductsToMap(products []Product) map[string]Product {
	result := make(map[string]Product)
	for _, product := range products {
		result[product.Sku] = product
	}
	return result
}

func readSuppliers() ([]ProductSupplier, map[int]ProductSupplier) {
	f, err := excelize.OpenFile("migration/Suppliers_store.xlsx")
	if err != nil {
		log.Fatal(err)
	}

	rows, err := f.GetRows("ITEM_SUPPLIER")
	if err != nil {
		log.Fatal(err)
	}

	var suppliers []ProductSupplier
	suppliersMap := make(map[int]ProductSupplier)
	for i, row := range rows {
		if i == 0 {
			continue
		}
		supplierId, err := strconv.Atoi(row[1])
		if err != nil {
			supplierId = 0
		}
		supplier := ProductSupplier{
			Item:         row[0],
			SupplierId:   supplierId,
			SupplierName: row[2],
		}
		suppliers = append(suppliers, supplier)
		suppliersMap[supplierId] = supplier
	}

	fmt.Printf("There are %d rows\n", len(rows))

	fmt.Printf("There are %d unique suppliers\n", len(suppliersMap))

	return suppliers, suppliersMap
}

func readStores() ([]Store, map[string]Store) {
	f, err := excelize.OpenFile("migration/Stores.xlsx")
	if err != nil {
		log.Fatal(err)
	}

	rows, err := f.GetRows("MFCs")
	if err != nil {
		log.Fatal(err)
	}

	var stores []Store
	storesMap := make(map[string]Store)
	for i, row := range rows {
		if i <= 1 {
			continue
		}
		storeId := row[1]
		address := strings.Split(row[5], ",")

		fmt.Printf("Address is %s Size %d \n", address, len(address))

		var store Store

		if len(address) > 4 {
			store = Store{
				Name:       row[0],
				Code:       storeId,
				Id:         row[2],
				Plaza:      row[3],
				Address:    strings.TrimSpace(address[0] + " " + address[1]),
				PostalCode: strings.TrimSpace(strings.ReplaceAll(address[2], "C.P. ", "")),
				City:       strings.TrimSpace(address[3]),
				State:      strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(address[4], "N.L.", "Nuevo León"), "N.L", "Nuevo León")),
				Latitude:   row[7],
				Longitude:  row[8],
				Surface:    row[9],
				Area:       row[10],
			}
		} else {
			store = Store{
				Name:  row[0],
				Code:  storeId,
				Id:    row[2],
				Plaza: row[3],
			}
		}
		stores = append(stores, store)
		storesMap[storeId] = store
	}

	fmt.Printf("There are %d rows\n", len(rows))
	fmt.Printf("There are %d unique stortes\n", len(storesMap))

	return stores, storesMap
}
