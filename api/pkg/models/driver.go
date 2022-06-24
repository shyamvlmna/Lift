package models

import "gorm.io/gorm"

type Driver struct {
	gorm.Model
	FirstName   string `gorm:"not null" json:"firstname"`
	LastName    string `json:"lastname"`
	PhoneNumber string `gorm:"not null;unique" json:"phonenumber"`
	Email       string `gorm:"not null;unique" json:"email"`
	Password    string `gorm:"not null" json:"password"`
	City        string `json:"city"`
	LicenceNum  string `json:"licence"`
	Approved    bool   `gorm:"default:false" json:"approved"`
	Active      bool
	Cab         Vehicle `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
