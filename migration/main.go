package main

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/xuri/excelize/v2"
	"log"
	"strings"
)

const (
	// PostgreSQL connection string
	psqlInfo         = "host=mioxxo-products.ce989v8mdple.us-east-1.rds.amazonaws.com port=5432 user=postgres dbname=products password=Yk9deWhjbYUUR5LEF8SnMf7w9jghPf sslmode=disable"
	ownerId          = "005Hp00000igpOJIAY"
	lastModifiedById = "005Hp00000igHBGIA2"
	jokrEntity       = "001Hp00002kvRgtIAE"
	taxRateId        = "a0dHp00003rSHNVIA4"
	createdDate      = "2021-08-25T20:00:00.000Z"
	sheetName        = "Sheet1"
)

var brandsMap = make(map[string]Brand)

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
	Enabled          bool           `db:"is_enabled"`
}

func main() {
	db := createDatabaseConnection()
	var brands = readBrands(*db)
	saveBrandsToExcelFile("brands.xlsx", brands)

	//var categories = readCategories(*db)
	//saveCategoriesToExcelFile("categories.xlsx", categories)
	//
	//defer db.Close()
	//
	//db = createDatabaseConnection()
	//var products = readProducts(*db)
	//saveProductsToExcelFile("products.xlsx", products)
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
		brandsMap[brand.ID] = brand
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
	err := db.Select(&products, "SELECT sku, barcode, name, short_description, variant, brand, categories, image, tags, is_enabled FROM mgo.products")
	if err != nil {
		log.Fatalf("Could not read products: %v", err)
	}

	for _, product := range products {
		fmt.Printf("Product %+v\n", product)
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

func saveBrandsToExcelFile(filename string, brands []Brand) {
	file := excelize.NewFile()
	index, err := file.NewSheet(sheetName)
	if err != nil {
		log.Fatalf("Could not create sheet: %v", err)
	}
	file.SetActiveSheet(index)
	// Set the headers
	headers := []string{"_", "Id", "OwnerId", "IsDeleted", "Name", "CurrencyIsoCode", "RecordTypeId", "CreatedDate", "CreatedById", "LastModifiedDate", "LastModifiedById", "SystemModstamp", "LastActivityDate", "LastViewedDate", "LastReferencedDate",
		"Brand_Owner__c", "Type__c"}
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		file.SetCellValue(sheetName, cell, header)
	}

	for i, brand := range brands {
		row := i + 2
		setCellValue(file, sheetName, row, 1, "[Brand__c]")
		setCellValue(file, sheetName, row, 2, brand.ID)
		setCellValue(file, sheetName, row, 3, ownerId)
		setCellValue(file, sheetName, row, 4, "FALSE")
		setCellValue(file, sheetName, row, 5, brand.Name.String)
		setCellValue(file, sheetName, row, 6, "MXN")
		setCellValue(file, sheetName, row, 7, "012Hp000001mPmEIAU")
		//setCellValue(file, sheetName, row, 8, createdDate)
		//setCellValue(file, sheetName, row, 9, lastModifiedById)
		//setCellValue(file, sheetName, row, 10, createdDate)
		//setCellValue(file, sheetName, row, 11, lastModifiedById)
		//setCellValue(file, sheetName, row, 12, createdDate)
		setCellValue(file, sheetName, row, 16, "a0FHp00000uorcKMAQ")
		setCellValue(file, sheetName, row, 17, "Brand")
	}
	if err := file.SaveAs(filename); err != nil {
		fmt.Println(err)
	}
}

func saveCategoriesToExcelFile(filename string, categories []Category) {
	file := excelize.NewFile()
	index, err := file.NewSheet(sheetName)
	if err != nil {
		log.Fatalf("Could not create sheet: %v", err)
	}
	file.SetActiveSheet(index)
	// Set the headers
	headers := []string{"_", "Id", "OwnerId", "IsDeleted", "Name", "CurrencyIsoCode", "CreatedDate", "CreatedById", "LastModifiedDate", "LastModifiedById", "SystemModstamp", "LastActivityDate", "LastViewedDate", "LastReferencedDate",
		"CS_Product_Object_Type__c", "Category_Group__c", "Category_Level_1__c", "Category_Level_2__c", "Category_Level_3__c", "Category_Level_4__c", "Category_Level_5__c", "Category_Level_6__c", "Global_Category_Group__c", "Global_Category_Level_1__c",
		"Is_Supermarket__c", "JOKR_Entity__c"}

	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		file.SetCellValue(sheetName, cell, header)
	}

	for i, category := range categories {
		row := i + 2
		setCellValue(file, sheetName, row, 1, "[Global_Category_Tree__c]")
		setCellValue(file, sheetName, row, 2, category.ID)
		setCellValue(file, sheetName, row, 3, ownerId)
		setCellValue(file, sheetName, row, 4, "FALSE")
		//setCellValue(file, sheetName, row, 5, category.Name.String)
		setCellValue(file, sheetName, row, 6, "MXN")
		//setCellValue(file, sheetName, row, 7, createdDate)
		//setCellValue(file, sheetName, row, 8, lastModifiedById)
		//setCellValue(file, sheetName, row, 9, createdDate)
		//setCellValue(file, sheetName, row, 10, lastModifiedById)
		//setCellValue(file, sheetName, row, 11, createdDate)
		setCellValue(file, sheetName, row, 16, category.Name.String)
		setCellValue(file, sheetName, row, 25, "FALSE")
	}
	if err := file.SaveAs(filename); err != nil {
		fmt.Println(err)
	}
}

func saveProductsToExcelFile(filename string, products []Product) {
	file := excelize.NewFile()
	index, err := file.NewSheet(sheetName)
	if err != nil {
		log.Fatalf("Could not create sheet: %v", err)
	}
	file.SetActiveSheet(index)
	// Set the headers
	headers := []string{"_", "Id", "IsDeleted", "Name", "CurrencyIsoCode", "CreatedDate", "CreatedById", "LastModifiedDate", "LastModifiedById", "SystemModstamp", "LastActivityDate", "LastViewedDate", "LastReferencedDate",
		"JOKR_Entity__c", "Active_Hubs__c", "Additional_Tax_Rate__c", "Age_Verification_Required__c", "Allow_for_Soon_to_Expire__c", "Archived_Hubs__c", "Assortment_Type__c", "Avalara_Tax_ID__c",
		"Barcode_2__c", "Barcode_3__c", "Barcode__c", "Brand_Division__c", "Brand_Line_OLD__c", "Brand_Name_OLD__c", "Brand_OLD__c", "Brand_Owner_OLD__c", "Brand_Story__c", "Brand__c",
		"Buffer__c", "CMS_Product_Link__c", "CRWN_Content_Ready__c", "Chile_Alcohol_Schedule_Available_To_Sale__c", "Clusters_in_Discontinuation__c", "Content_Relevance__c", "Country_Default_Compare_at_Price__c",
		"Country_Default_Price__c", "Country_Recommended_Retail_Price__c", "Country__c", "Customer_Delivery_Option__c", "Customer_Fulfillment_Type__c", "Cycle_count_Frequency__c", "DC_Temperature_Differs_from_Hub_Temp__c",
		"DC_Temperature_Storage_Requirements__c", "Default_Tax_Rate__c", "Draft_Hubs__c", "Exclude_from_Promotions__c", "Expiry_date_tracking_required__c", "Fresh_Product__c", "Global_Category_Tree__c", "Global_Name__c",
		"Gross_Weight_UOM__c", "Gross_Weight__c", "Height_UOM__c", "Hub_Additional_Storage_Requirements__c", "Hub_Temperature_Storage__c", "IEPS__c",
		"Is_Bundle__c", "Is_Generic_Product__c", "JOKR_UPID__c", "Key_Value_Item__c", "Length_UOM__c", "Local_Hero_Product__c", "Long_Description__c", "Maximum_Acceptable_Reception_Temp_C__c", "Maximum_Quantity__c", "Minimum_Acceptable_Reception_Temp_C__c",
		"Minimum_Quantity__c", "Net_Weight_UOM__c", "Net_Weight__c", "NetsuiteId__c", "Notes__c", "Nutritional_Table_Required__c", "Open_Catalogue_URL__c", "Private_Label_Product__c", "Product_Content_Ready__c", "Product_Description__c",
		"Product_Height__c", "Product_Length__c", "Product_Ownership__c", "Product_Registration_Step__c", "Product_Title__c", "Product_Type__c", "Product_Width__c", "Public_Image_URL__c", "Purchasing_Status__c", "Quality_Technical_Card__c",
		"Reception_Temperature_Tracking_Required__c", "SKU__c", "Search_terms__c", "Selling_Quantity_UOM__c", "Selling_Quantity__c", "Selling_UOM__c", "Shelf_Life_Inbound_Minimum__c", "Shelf_Life_Outbound_Minimum__c", "Shelf_Life_Tracking_Required__c",
		"Shelf_Life_Type__c", "Short_Description__c", "Soon_to_Expire_Default_Discount__c", "Special_Packaging_Requirements__c", "Status__c", "Sub_Category_CO__c", "Sub_Category_MX__c",
		"Sub_Category_PE__c", "Sub_Category_US__c", "Substitute_1__c", "Substitute_2__c", "Substitute_3__c", "Taxation_Ready__c", "TestDateValue__c", "To_discontinue_when_Out_of_Stock__c", "Total_Shelf_Life_in_Days__c", "Transport_Temperature__c",
		"UI_Categories__c", "UI_Content_1_UOM__c", "UI_Content_1__c", "UI_Content_2_UOM__c", "UI_Content_2__c", "UOM2__c", "Valid_Clusters__c", "Vendor_Account__c", "Vendor_Test_Item__c", "Vendor_minimum_price__c", "Vendor_revenue_share_percentage__c",
		"Volume_UOM__c", "Volume__c	Warranty_Offered__c", "Width_UOM__c", "X18_Char_Id__c", "Yield_UOM__c", "Yield__c", "competitor_price_benchmark_URL_1__c", "competitor_price_benchmark_URL_2__c", "competitor_price_benchmark_URL_3__c", "Available_to_Purchase_Vendor_Products__c",
		"Alcohol_by_Volume__c", "Alcoholic_Beverage__c", "Brand_Name__c", "Brand_Owner__c", "Category_Group__c", "Category_Level_1__c", "Category_Level_2__c", "Category_Level_3__c", "Category_Level_4__c", "Category_Level_5__c", "Category_Level_6__c", "External_Unique_Identifier__c",
		"Is_Key_Value_Item__c", "Next_Status__c", "Purchasing_Unavailable__c", "Tax_Code_2__c", "Tax_Code__c", "Tax_Percentage_2__c", "Tax_Percentage__c", "TopSort_Vendor__c", "Total_Tax_Rate__c", "Unit_Volume_in_cm3__c", "Vendor_Name__c",
	}

	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		file.SetCellValue(sheetName, cell, header)
	}

	for i, product := range products {
		row := i + 2

		productId := product.Sku + "-" + product.Barcode.String

		setCellValue(file, sheetName, row, 1, "[Product__c]")
		setCellValue(file, sheetName, row, 2, productId)
		setCellValue(file, sheetName, row, 3, "FALSE")
		setCellValue(file, sheetName, row, 4, product.Name)
		setCellValue(file, sheetName, row, 5, "MXN")
		setCellValue(file, sheetName, row, 6, createdDate)
		setCellValue(file, sheetName, row, 7, lastModifiedById)
		setCellValue(file, sheetName, row, 8, createdDate)
		setCellValue(file, sheetName, row, 9, lastModifiedById)
		setCellValue(file, sheetName, row, 10, createdDate)
		setCellValue(file, sheetName, row, 14, jokrEntity)
		setCellValue(file, sheetName, row, 15, 0)
		setCellValue(file, sheetName, row, 16, taxRateId)
		setCellValue(file, sheetName, row, 17, "No")
		setCellValue(file, sheetName, row, 18, "FALSE")
		setCellValue(file, sheetName, row, 19, 0)
		setCellValue(file, sheetName, row, 20, "Standard")
		setCellValue(file, sheetName, row, 21, 50171800)
		setCellValue(file, sheetName, row, 24, product.Barcode.String)
		setCellValue(file, sheetName, row, 31, product.Brand.String)
		setCellValue(file, sheetName, row, 35, "FALSE")
		setCellValue(file, sheetName, row, 38, 100) // TODO Set here the price
		setCellValue(file, sheetName, row, 40, "Mexico")
		setCellValue(file, sheetName, row, 41, "Rapid Delivery")
		setCellValue(file, sheetName, row, 42, "Core") // could be 13
		setCellValue(file, sheetName, row, 44, "FALSE")
		setCellValue(file, sheetName, row, 46, taxRateId)
		setCellValue(file, sheetName, row, 47, 0)
		setCellValue(file, sheetName, row, 48, "FALSE")
		setCellValue(file, sheetName, row, 49, "Yes")                     // expiry_date_tracking_required__c
		setCellValue(file, sheetName, row, 50, "FALSE")                   // fresh_product__c
		setCellValue(file, sheetName, row, 51, product.Categories.String) // global_category_tree__c // TODO Get category tree here
		setCellValue(file, sheetName, row, 57, "Room Temperature")        // hub_temperature_storage__c
		setCellValue(file, sheetName, row, 59, "FALSE")                   // is_bundle__c
		setCellValue(file, sheetName, row, 60, "FALSE")                   // is_generic_product__c
		setCellValue(file, sheetName, row, 62, "FALSE")                   // key_value_item__c
		setCellValue(file, sheetName, row, 64, "FALSE")                   // local_hero_product__c
		setCellValue(file, sheetName, row, 67, 100)                       // maximum_quantity__c // TODO Set here
		setCellValue(file, sheetName, row, 69, 1)                         // minimum_quantity__c // TODO Set here
		setCellValue(file, sheetName, row, 76, "FALSE")                   // private_label_product__c
		setCellValue(file, sheetName, row, 77, "FALSE")                   // product_content_ready__c
		setCellValue(file, sheetName, row, 81, "JOKR owned")              // product_ownership__c
		setCellValue(file, sheetName, row, 82, "Draft")                   // product_registration_step__c
		setCellValue(file, sheetName, row, 83, product.Name)              // product_title__c
		setCellValue(file, sheetName, row, 86, product.Image.String)      // public_image_url__c
		setCellValue(file, sheetName, row, 89, "No")
		setCellValue(file, sheetName, row, 90, product.Sku) // sku__c
		setCellValue(file, sheetName, row, 95, 200)         //shelf_life_inbound_minimum__c // TODO Set here
		setCellValue(file, sheetName, row, 96, 300)         // shelf_life_outbound_maximum__c // TODO Set here
		var status string
		if product.Enabled {
			status = "Active"
		} else {
			status = "Inactive"
		}
		setCellValue(file, sheetName, row, 102, status)        // status__c
		setCellValue(file, sheetName, row, 110, "FALSE")       // taxation_ready__c
		setCellValue(file, sheetName, row, 111, "16:41:59.00") // test_date_value__c
		setCellValue(file, sheetName, row, 112, "FALSE")       // to_discontinue_when_out_of_stock__c
		setCellValue(file, sheetName, row, 113, 360)           // total_shelf_life_in_days__c
		setCellValue(file, sheetName, row, 117, "ml")          // UI_Content_1_UOM__c // TODO Set here
		setCellValue(file, sheetName, row, 119, "ea")          // UOM2__c  // TODO Set here
		setCellValue(file, sheetName, row, 130, productId)
		setCellValue(file, sheetName, row, 136, 1)
		setCellValue(file, sheetName, row, 139, product.Brand.String)
		setCellValue(file, sheetName, row, 140, product.Brand.String)
		setCellValue(file, sheetName, row, 141, getFirstElement(product.Categories.String))
		setCellValue(file, sheetName, row, 148, productId)
		setCellValue(file, sheetName, row, 149, 0)
		setCellValue(file, sheetName, row, 150, "01: Ready to be updated to Inactive")
		setCellValue(file, sheetName, row, 151, "FALSE") // TODO Set here
		setCellValue(file, sheetName, row, 152, "IEPS0") //  TODO Set here
		setCellValue(file, sheetName, row, 153, "IVA0")  // TODO Set here
		setCellValue(file, sheetName, row, 154, 0)
		setCellValue(file, sheetName, row, 155, 0)
		setCellValue(file, sheetName, row, 156, brandsMap[product.Brand.String].Name.String) // TODO Set here
		setCellValue(file, sheetName, row, 157, 0)
		setCellValue(file, sheetName, row, 158, 0)
	}

	if err := file.SaveAs(filename); err != nil {
		fmt.Println(err)
	}
}

func setCellValue(file *excelize.File, sheetName string, row int, column int, value interface{}) {
	cell, _ := excelize.CoordinatesToCellName(column, row)
	file.SetCellValue(sheetName, cell, value)
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
