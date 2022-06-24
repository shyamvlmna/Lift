package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FirstName   string `gorm:"not null" json:"first_name" validate:"required,min=2,max=30"`
	LastName    string `json:"last_name"`
	PhoneNumber string `gorm:"not null;unique" json:"phone_number" validate:"required,min=10"`
	Email       string `gorm:"not null;unique" json:"email" validate:"required,min=20"`
	Password    string `gorm:"not null" json:"password" validate:"required,min=4"`
	Token       string `json:"token"`
	// Wallet      UserWallet `json:"user_wallet" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	// TripHistory []Trip     `json:"trip_history" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type UserWallet struct {
	gorm.Model
	Balance uint
}

// func (u User) CreateUser() {

// }

// func (u User) UpdateUser() {

// }
// func (u User) DeleteUser() {

// }
// func (u User) GetUser(key, value string) User {
// 	return user.GetUser(key, value)
// }
