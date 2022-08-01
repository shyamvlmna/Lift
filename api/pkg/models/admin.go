package models

import (
	"github.com/shayamvlmna/lift/pkg/database"
	"gorm.io/gorm"
)

type Admin struct {
	gorm.Model
	AdminId       uint64 `gorm:"primaryKey;unique;autoIncrement;" json:"adminid"`
	Username      string `gorm:"not null;unique;" json:"username"`
	Password      string `gorm:"not null;" json:"password"`
	IsAdmin       bool   `gorm:"default:true;" json:"isadmin"`
	WalletBalance uint   `gorm:"default:0" json:"adminwallet"`
}

func (a *Admin) Add() error {
	db := database.Db
	err := db.AutoMigrate(&Admin{})
	if err != nil {
		return err
	}

	result := db.Create(&a)
	return result.Error
}
