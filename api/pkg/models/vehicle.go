package models

import (
	"github.com/shayamvlmna/cab-booking-app/pkg/database"
	"gorm.io/gorm"
)

type Vehicle struct {
	gorm.Model
	VehicleId    uint64 `gorm:"primaryKey;autoIncrement;" json:"vehicleid"`
	Registration string `json:"registration"`
	Brand        string `json:"brand"`
	Category     string `json:"type"`
	VehicleModel string `json:"model"`
	Colour       string `json:"colour"`
	DriverId     uint
}

// Add new vehicle into database
func (v *Vehicle) Add() error {
	db := database.Db
	err := db.AutoMigrate(&Vehicle{})
	if err != nil {
		return err
	}

	result := db.Create(&v)
	return result.Error
}
