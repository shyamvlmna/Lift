package models

import (
	"github.com/shayamvlmna/cab-booking-app/pkg/database"
	"gorm.io/gorm"
)

type Admin struct {
	gorm.Model
	Username string `gorm:"not null;unique;" json:"username"`
	Password string `gorm:"not null;" json:"password"`
	IsAdmin  bool   `gorm:"default:true;" json:"isadmin"`
}

func (a *Admin) Add() error {
	db := database.Db
	db.AutoMigrate(&Admin{})

	result := db.Create(&a)
	return result.Error
}
