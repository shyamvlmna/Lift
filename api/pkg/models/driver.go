package models

import "gorm.io/gorm"

type Driver struct {
	gorm.Model
	// DriverId    uint64  `gorm:"primaryKey;autoIncrement"`
	FirstName   string  `gorm:"not null" json:"first_name"`
	LastName    string  `json:"last_name"`
	PhoneNumber string  `gorm:"not null;unique" json:"phone_number"`
	Email       string  `gorm:"not null;unique" json:"email"`
	Password    string  `gorm:"not null" json:"password"`
	City        string  `json:"city"`
	LicenceNum  string  `json:"licence"`
	Approved    bool    `gorm:"default:false" json:"approved"`
	Active      bool    `gorm:" default:true" json:"active"`
	Cab         Vehicle `gorm:"ForeignKey:ID;references:ID"`
}
