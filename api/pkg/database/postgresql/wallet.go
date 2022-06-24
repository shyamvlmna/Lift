package database

import (
	"fmt"

	"github.com/shayamvlmna/cab-booking-app/pkg/models"
	"gorm.io/gorm"
)

func OpenWalletDb() (*gorm.DB, error) {
	Db, err := openDB()
	if err != nil {
		return nil, err
	}
	wallet := &models.UserWallet{}
	err = Db.AutoMigrate(&wallet)
	if err != nil {
		return nil, err
	}
	fmt.Println("wallet db opened")
	return Db, nil
}
