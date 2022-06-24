package models

import "gorm.io/gorm"

type Trip struct {
	gorm.Model
	Source        Location 
	Destination   Location 
	Distance      uint
	Fare          uint
	PaymentMethod string
	Rating        uint8
}

type Location struct {
	Lon string
	Lat string
}

type Payment struct {
	Wallet bool
	Cash   bool
}
