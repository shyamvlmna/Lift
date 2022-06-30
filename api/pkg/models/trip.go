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
	Distance      int   `gorm:"not null;" json:"distance"`
	Fare          int
	ETA           string `json:"timeduration"`
	PaymentMethod string `json:"paymentmethod"`
	Rating        uint8  `json:"triprating"`
	UserId        uint64 `json:"userid"`
	DriverId      uint64 `json:"driverid"`
}

func (t *Trip) Add() error {
	db := database.Db
	db.AutoMigrate(&Trip{})

	result := db.Create(&t)
	return result.Error
}

func (t *Trip) Update() error {
	db := database.Db
	db.AutoMigrate(&Trip{})

	trip := &Trip{}

	db.Where("id=?", t.Id).First(&trip)

	result := db.Model(&trip).Updates(&Trip{
		Model:         gorm.Model{},
		Id:            0,
		Source:        "",
		Destination:   "",
		Distance:      0,
		Fare:          0,
		ETA:           "",
		PaymentMethod: "",
		Rating:        0,
		UserId:        0,
		DriverId:      0,
	})

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
