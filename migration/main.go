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
	"strconv"
	"strings"
)

const (
	// PostgreSQL connection string
	psqlInfo             = "host=mioxxo-products.ce989v8mdple.us-east-1.rds.amazonaws.com port=5432 user=postgres dbname=products password=Yk9deWhjbYUUR5LEF8SnMf7w9jghPf sslmode=disable"
	recordTypeBrandOwner = "012Hp000001mPmGIAU"
	recordTypeBrand      = "012Hp000001mPmEIAU"
	recordTypeVendor     = "012Hp000001mPmDIAU"
	recordTypeHub        = "012Hp000001mPmBIAU"
	jokrEntity           = "001Hp00002kvRgtIAE"
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

func main() {
	//var brands = readBrands()
	//saveBrandOwnersToCsvFile("brandowners.csv", brands)

	//// TODO import data to the Salesforce here
	//// TODO export Brands csv file from Salesforce
	//
	//brandOwners := readSalesforceEntities("migration/sandbox/brandExported.csv", recordTypeBrandOwner)
	//for key, value := range brandOwners {
	//	fmt.Printf("Key '%s' Brand owner %+v\n", key, value)
	//}
	//
	//saveBrandsToCsvFile("brands.csv", brands, brandOwners)

	//// TODO import data to the Salesforce here
	//// TODO export Brands csv file from Salesforce

	//var categoryTrees = readCategoryTrees()
	//saveCategoriesToCsvFile("categories.csv", categoryTrees)

	////// TODO import data to the Salesforce here
	////// TODO export Categories csv file from Salesforce

	//productSuppliers, suppliers := readSuppliers()
	//
	//for _, value := range productSuppliers {
	//	fmt.Printf("Product suppliers %+v\n", value)
	//}
	//
	//for key, value := range suppliers {
	//	fmt.Printf("Supplier Key %+v Value %+v\n", key, value)
	//}

	db := createDatabaseConnection()
	var products = readProducts(*db)
	defer db.Close()
	for _, value := range products {
		fmt.Printf("Product %+v \n", value)
	}

	hubs := readSalesforceEntities("migration/sandbox/hubExported.csv", "")
	for key, value := range hubs {
		fmt.Printf("Key '%s' Hub %+v\n", key, value)
	}

	categories := readSalesforceCategories("migration/sandbox/categoryExported.csv")
	//for key, value := range categories {
	//	if strings.Contains(value.Name, "Oxxo: ") {
	//		fmt.Printf("Key '%s' CategoryTree %+v\n", key, value)
	//	}
	//}

	brands := readSalesforceEntities("migration/sandbox/brandExported.csv", recordTypeBrand)
	//for key, value := range brands {
	//	fmt.Printf("Key '%s' Brand %+v\n", key, value)
	//}

	taxRates := readSalesforceEntities("migration/sandbox/taxRateExported.csv", "")
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

func readProducts(db sqlx.DB) []Product {
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
		variantList := strings.Split(product.Variant, " ")

		dataRow[13] = strings.ReplaceAll(variantList[0], ",", ".")
		if len(variantList) > 1 {
			if contains([]string{"tabs", "Tabs"}, variantList[1]) {
				dataRow[14] = "Tabs"
			} else if contains([]string{"h"}, variantList[1]) {
				dataRow[14] = "h"
			} else if contains([]string{"rollo", "rollos", "Rollo", "Rollos"}, variantList[1]) {
				dataRow[14] = "Rollos"
			} else {
				dataRow[14] = variantList[1]
			}
		}

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

	fmt.Sprintf("There are %v rows\n", len(rows))

	return suppliers, suppliersMap
}
