package models

type ProductChild struct {
	Sku      string
	Quantity string
}

type ProductSupplier struct {
	Item         string
	SupplierId   int
	SupplierName string
}

type ProductSupplierKey struct {
	Item       string
	SupplierId int
}

type TaxRate struct {
	Name      string
	Rate      string
	Type      string
	Secondary bool
}
