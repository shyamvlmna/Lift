package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserId      uint64     `gorm:"primaryKey;unique"`
	FirstName   string     `gorm:"not null" json:"first_name"`
	LastName    string     `json:"last_name"`
	PhoneNumber string     `gorm:"not null;unique" json:"phone_number"`
	Email       string     `gorm:"not null;unique" json:"email"`
	Password    string     `gorm:"not null" json:"password"`
	Token       string     `json:"token"`
	Active      bool       `json:"active" gorm:"default:true"`
	Wallet      UserWallet `gorm:"ForeignKey:UserId;references:WalletId;embedded" json:"user_wallet"`
	TripHistory []Trip     `gorm:"ForeignKey:UserId;references:UserId" json:"trip_history" `
}

// BookedTrip  Trip        `gorm:"ForeignKey:UserId;references:TripId;embedded" json:"trip"`
type UserWallet struct {
	gorm.Model
	WalletId uint64 `gorm:"primaryKey;autoIncrement;unique"`
	Balance  string
}
