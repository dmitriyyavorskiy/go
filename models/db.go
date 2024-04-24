package models

import "database/sql"

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

type CategoryTree struct {
	ID           string         `db:"id"`
	Name         sql.NullString `db:"name"`
	Image        sql.NullString `db:"image"`
	CategoryId   sql.NullString `db:"subcategory_id"`
	CategoryName sql.NullString `db:"subcategory_name"`
}

type Inventory struct {
	Sku       string         `db:"sku"`
	Store     string         `db:"store"`
	Zone      sql.NullString `db:"zone"`
	StoreCode string         `db:"cr"`
	Price     float32        `db:"price_list"`
}

type Tax struct {
	Rate        int    `json:"rate"`
	Type        string `json:"type"`
	Withholding bool   `json:"withholding"`
}
