package models

type UserWallet struct {
	WalletId uint64 `gorm:"primaryKey;"`
	UserId   uint64 `gorm:"foreignKey"`
	Balance  uint
}

type DriverWallet struct {
	WalletId uint64 `gorm:"primaryKey;"`
	DriverId uint64 `json:"driverid"`
	Balance  uint   `json:"balance"`
}

//incase of array field
// var users []User
//   err := db.Model(&User{}).Preload("CreditCard").Find(&users).Error
