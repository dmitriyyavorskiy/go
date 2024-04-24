package models

type CreateProduct struct {
	Name               string  `json:"name"`
	Barcode            string  `json:"barcode"`
	Sku                string  `json:"sku"`
	Brand              string  `json:"brand"`
	DefaultTaxRate     string  `json:"default_tax_rate"`
	AdditionalTaxRate  string  `json:"additional_tax_rate"`
	GlobalCategoryTree string  `json:"global_category_tree"`
	AgeRestriction     bool    `json:"age_restriction"`
	Restricted         bool    `json:"restricted"`
	Image              string  `json:"image"`
	Price              float32 `json:"price"`
	Amount             float32 `json:"amount"`
	Measurement        string  `json:"measurement"`
	AssortmentType     string  `json:"assortment_type"`
	PackSize           string  `json:"pack_size"`
}
