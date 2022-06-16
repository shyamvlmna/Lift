package database

import (
	"fmt"

	"github.com/shayamvlmna/cab-booking-app/app/models"
	"gorm.io/gorm"
)

func OpenDriverDb() (*gorm.DB, error) {
	Db, err := openDB()
	if err != nil {
		return nil, err
	}
	driver := &models.Driver{}
	err = Db.AutoMigrate(&driver)
	if err != nil {
		return nil, err
	}
	fmt.Println("driver db opened")
	return Db, nil
}

func closeDriverdb(db *gorm.DB) {

	sqlDb, err := db.DB()
	if err != nil {
		fmt.Println(err)
	}
	sqlDb.Close()
	fmt.Println("driver db closed")
}

func CheckDriver(key string) bool {

	db, err := OpenDriverDb()
	if err != nil {
		fmt.Println(err)
	}
	defer closeDriverdb(db)
	driver := &models.Driver{}
	result := db.Where("phone_number=?", key).First(&driver)

	if result.Error == gorm.ErrRecordNotFound {
		return false
	} else {
		return true
	}
}
func InsertDriver(driver *models.Driver) error {
	db, err := OpenDriverDb()
	if err != nil {
		return err
	}
	defer closeDriverdb(db)
	result := db.Create(driver)

	return result.Error
}

func FindDriver(key string) (models.Driver, error) {
	db, err := OpenDriverDb()
	if err != nil {
		fmt.Println(err)
	}
	defer closeDriverdb(db)
	driver := models.Driver{}
	return driver, nil
}

func UpdateDriver() {

}
