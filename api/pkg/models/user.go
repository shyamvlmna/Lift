package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserId      uint64      `gorm:"primaryKey;unique"`
	FirstName   string      `gorm:"not null" json:"first_name" validate:"required,min=2,max=30"`
	LastName    string      `json:"last_name"`
	PhoneNumber string      `gorm:"not null;unique" json:"phone_number" validate:"required,min=10"`
	Email       string      `gorm:"not null;unique" json:"email" validate:"required,min=20"`
	Password    string      `gorm:"not null" json:"password" validate:"required,min=4"`
	Token       string      `json:"token"`
	Active      bool        `json:"active" gorm:"default:true"`
	Wallet      UserWallet  `gorm:"ForeignKey:UserId;references:WalletId;embedded" json:"user_wallet"`
	BookedTrip  Trip        `gorm:"ForeignKey:UserId;references:TripId;embedded" json:"trip"`
	// TripHistory TripHistory `gorm:"ForeignKey:UserId;references:TripHistoryId;embedded" json:"trip_history" `
}

type UserWallet struct {
	gorm.Model
	WalletId uint64 `gorm:"primaryKey;autoIncrement;unique"`
	Balance  string
}

type TripHistory struct {
	gorm.Model
	TripHistoryId uint64 `gorm:"autoIncrement;unique;primaryKey" json:"triphistory"`
	Trips         []Trip
}
type Trip struct {
	TripId        uint64 `gorm:"autoIncrement;unique;primaryKey" json:"tripid"`
	Distance      uint   `gorm:"not null"`
	Fare          uint   `gorm:"not null"`
	PaymentMethod string
	Rating        uint
}
