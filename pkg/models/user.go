package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	// Id          int64
	FirstName   string `gorm:"not null" json:"firstname"`
	LastName    string `json:"lastname"`
	PhoneNumber int    `gorm:"not null;unique" json:"phonenumber"`
	Email       string `gorm:"not null;unique" json:"email"`
	Password    string `gorm:"not null" json:"password"`
}
