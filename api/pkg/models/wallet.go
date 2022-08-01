package models

import (
	"github.com/shayamvlmna/lift/pkg/database"
	"gorm.io/gorm"
)

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

type AdminWallet struct {
	WalletId uint64 `gorm:"primaryKey;"`
	AdminId  uint64 `json:"adminid"`
	Balance  uint   `json:"balance"`
}

func WalletTransactions(userId, driverId uint, fare float64) error {
	db := database.Db

	if err := db.Transaction(func(tx *gorm.DB) error {

		if err := tx.AutoMigrate(&User{}); err != nil {
			return err
		}

		if err := tx.Model(&User{}).Where("user_id=?", userId).UpdateColumn("wallet_balance", gorm.Expr("wallet_balance - ?", fare)).Error; err != nil {
			return err
		}

		if err := tx.AutoMigrate(&Driver{}); err != nil {
			return err
		}

		if err := tx.Model(&Driver{}).Where("driver_id=?", driverId).UpdateColumn("wallet_balance", gorm.Expr("wallet_balance + ?", fare*80/100)).Error; err != nil {
			return err
		}

		if err := tx.AutoMigrate(&Admin{}); err != nil {
			return err
		}

		if err := tx.Model(&Admin{}).Where("admin_id=1").UpdateColumn("wallet_balance", gorm.Expr("wallet_balance + ?", fare*20/100)).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}
	return nil
}
func CashTransactions(driverId uint, fare float64) error {
	db := database.Db

	if err := db.Transaction(func(tx *gorm.DB) error {

		if err := tx.AutoMigrate(&Driver{}); err != nil {
			return err
		}

		if err := tx.Model(&Driver{}).Where("driver_id=?", driverId).UpdateColumn("wallet_balance", gorm.Expr("wallet_balance - ?", fare*15/100)).Error; err != nil {
			return err
		}

		if err := tx.AutoMigrate(&Admin{}); err != nil {
			return err
		}

		if err := tx.Model(&Admin{}).Where("admin_id=1").UpdateColumn("wallet_balance", gorm.Expr("wallet_balance + ?", fare*15/100)).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}
	return nil
}
