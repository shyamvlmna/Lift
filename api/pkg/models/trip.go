package models

import (
	"github.com/shayamvlmna/cab-booking-app/pkg/database"
	"gorm.io/gorm"
)

type Trip struct {
	gorm.Model
	Id            uint64 `gorm:"primaryKey;" json:"tripid"`
	Source        string `json:"source"`
	Destination   string `json:"destination"`
	Distance      uint   `gorm:"not null;"`
	Fare          uint   `gorm:"not null;"`
	ETA           string `json:"timeduration"`
	PaymentMethod string `json:"paymentmethod"`
	Rating        uint8  `json:"triprating"`
	UserId        uint64
	DriverId      uint64
}

func (t *Trip) Add() error {
	db := database.Db
	db.AutoMigrate(&Trip{})

	result := db.Create(&t)
	return result.Error
}

// type Location struct {
// 	Id  uint    `gorm:"primaryKey"`
// 	Lat float64 `json:"latitude"`
// 	Lng float64 `json:"longitude"`
// }

// type Ride struct {
// 	RideId        uint64 `gorm:"autoIncrement;unique;primaryKey" json:"rideid"`
// 	Source        string `json:"source"`
// 	Destination   string `json:"destination"`
// 	ETA           string `json:"eta"`
// 	Fare          uint   `json:"fare"`
// 	PaymentMethod string `json:"paymentmethod"`
// }
