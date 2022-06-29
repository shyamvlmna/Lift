package models

import (
	"github.com/shayamvlmna/cab-booking-app/pkg/database"
	"gorm.io/gorm"
)

type Vehicle struct {
	gorm.Model
	VehicleId    uint64 `gorm:"primaryKey;autoIncrement;" json:"vehicleid"`
	Registration string `gorm:"not null;unique;" json:"registration"`
	Brand        string `gorm:"not null" json:"brand"`
	Category     string `gorm:"not null" json:"type"`
	VehicleModel string `gorm:"not null" json:"model"`
	Colour       string `gorm:"not null" json:"colour"`
	DriverId     uint64
}

//add new vehicle into database
func (v *Vehicle) Add() error {
	db := database.Db
	db.AutoMigrate(&Vehicle{})

	result := db.Create(&v)
	return result.Error
}
