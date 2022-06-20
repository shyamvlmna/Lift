package models

import "gorm.io/gorm"

type Admin struct {
	gorm.Model
	Username string `gorm:"not null;unique" json:"username"`
	Password string `gorm:"not null" json:"password"`
	IsAdmin  bool   `gorm:"default:true" json:"isadmin"`
}

// func (a Admin) Add() {
// 	// database.AddAdmin(&a)
// }

// func (a Admin) Get() {

// }
