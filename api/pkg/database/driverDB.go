package database

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/shayamvlmna/cab-booking-app/pkg/models"
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

// func closeDriverdb(db *gorm.DB) {

// 	sqlDb, err := db.DB()
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	sqlDb.Close()
// 	fmt.Println("driver db closed")
// }

//insert driver model to the driver database
//return error if any
func InsertDriver(driver *models.Driver) error {
	db, err := OpenDriverDb()
	if err != nil {
		return err
	}
	// defer closeDriverdb(db)
	result := db.Create(driver)

	return result.Error
}

func FindDriver(key, value string) (models.Driver, bool) {

	db, err := OpenDriverDb()
	if err != nil {
		fmt.Println(err)
	}
	// defer closeDriverdb(db)
	driver := &models.Driver{}
	result := db.Where(key+"=?", value).First(&driver)

	if result.Error == gorm.ErrRecordNotFound {
		return *driver, false
	} else {
		return *driver, true
	}
}

//get and return all drivers from the driver database
func GetDrivers() *[]models.Driver {
	db, err := OpenDriverDb()
	if err != nil {
		fmt.Println(err)
	}

	drivers := &[]models.Driver{}
	db.Find(&drivers)

	return drivers
}

//update a driver by getting updated driver fields
//only update the not null driver fields
func UpdateDriver(updatedDriver *models.Driver) error {

	return nil
}

//delete driver by id
//returns err if any
func DeleteDriver(id string) error {

	return nil
}
