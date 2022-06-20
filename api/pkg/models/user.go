package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FirstName   string `gorm:"not null" json:"firstname"`
	LastName    string `json:"lastname"`
	PhoneNumber string `gorm:"not null;unique" json:"phonenumber"`
	Email       string `gorm:"not null;unique" json:"email"`
	Password    string `gorm:"not null" json:"password"`
}

// func (u *User) CreateUser() {

// }

// func (u *User) UpdateUser() {

// }
// func (u *User) DeleteUser() {

// }
// func (u *User) GetUser(key, value string) User {
// 	return user.GetUser(key, value)
// }
