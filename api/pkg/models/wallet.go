package models

type UserWallet struct {
	Id      uint64 `gorm:"primaryKey;"`
	Balance uint
	UserId  uint64
}

type DriverWallet struct {
	Id       uint64 `gorm:"primaryKey;"`
	Balance  uint
	DriverId uint64
}

//incase of array field
// var users []User
//   err := db.Model(&User{}).Preload("CreditCard").Find(&users).Error
