package models

import "gorm.io/gorm"

type Organization struct {
	ID           uint        `gorm:"primary key;autoIncrement" json:"id"`
	Name         string      `json:"name"`
	Bio          string      `json:"bio"`
	AddressLine  string      `json:"addressLine"`
	Pincode      string      `json:"pinCode"`
	District     string      `json:"district"`
	State        string      `json:"state"`
	Country      string      `json:"country"`
	CategoriesIn StringArray `json:"categoriesIn"`
	AvailableIn  string      `json:"availableIn"`
	PageName     string      `json:"pageName"`
	PhoneNo      string      `json:"phoneNo"`
	Email        string      `json:"email"`
	ImageUrl1    string      `json:"imageUrl1"`
	ImageUrl2    string      `json:"imageUrl2"`
	ImageUrl3    string      `json:"imageUrl3"`
	About        string      `json:"about"`
	Experience   string      `json:"experience"`
	Specialities string      `json:"specialities"`
}

type OrganizationResponse struct {
	Overview Organization `json:"owerview"`
}

type Photo struct {
	ID             uint   `gorm:"primary key;autoIncrement" json:"id"`
	OrganizationID string `json:"organizationId"`
	Tag            string `json:"tag"`
	ImageUrl       string `json:"imageUrl"`
}

type Plan struct {
	ID             uint        `gorm:"primary key;autoIncrement" json:"id"`
	OrganizationID string      `json:"organizationId"`
	Type           string      `json:"type"`
	StartingFrom   int         `json:"startingFrom"`
	Services       StringArray `json:"services"`
}

func MigrateOrganization(db *gorm.DB) error {
	return db.AutoMigrate(&Organization{}, &Photo{}, &Plan{})
}
