package models

import "gorm.io/gorm"

type User struct {

	gorm.Model
	// Id          int64
	FirstName   string `gorm:"not null" json:"usrfirstname"`
	LastName    string `json:"usrlastname"`
	PhoneNumber string    `gorm:"not null;unique" json:"usrphonenumber"`
	Email       string `gorm:"not null;unique" json:"usremail"`
	Password    string `gorm:"not null" json:"usrpassword"`
}
