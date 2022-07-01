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
	Distance      string `gorm:"not null;" json:"distance"`
	Fare          uint   `json:"fare"`
	ETA           string `json:"timeduration"`
	PaymentMethod string `json:"paymentmethod"`
	UserId        uint64 `json:"userid"`
	DriverId      uint64 `json:"driverid"`
	Rating        uint8  `json:"triprating"`
}

func (t *Trip) Add(trip *Trip) error {
	db := database.Db
	err := db.AutoMigrate(&Trip{})
	if err != nil {
		return err
	}

	result := db.Create(&trip)
	return result.Error
}

func (t *Trip) Update() error {
	db := database.Db
	err := db.AutoMigrate(&Trip{})
	if err != nil {
		return err
	}

	trip := &Trip{}

	db.Where("id=?", t.Id).First(&trip)

	result := db.Model(&trip).Updates(&Trip{})

	return result.Error
}

// type Location struct {
// 	Id  uint    `gorm:"primaryKey"`
// 	Lat float64 `json:"latitude"`
// 	Lng float64 `json:"longitude"`
// }

type Ride struct {
	Source        string `json:"source"`
	Destination   string `json:"destination"`
	Distance      string `json:"distance"`
	ETA           string `json:"eta"`
	Fare          uint   `json:"fare"`
	PaymentMethod string `json:"paymentmethod"`
	UserId        uint64 `json:"userid"`
	DriverId      uint64 `json:"driverid"`
	Rating        int    `json:"rating"`
}
