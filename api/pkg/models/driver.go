package models

import "gorm.io/gorm"

type Driver struct {
	gorm.Model
	FirstName   string  `gorm:"not null" json:"drvrfirstname"`
	LastName    string  `json:"drvrlastname"`
	PhoneNumber string     `gorm:"not null;unique" json:"drvrphonenumber"`
	Email       string  `gorm:"not null;unique" json:"drvremail"`
	Password    string  `gorm:"not null" json:"drvrpassword"`
}
