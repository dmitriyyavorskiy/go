package models

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
