package models

import (
	"strconv"

	"github.com/shayamvlmna/cab-booking-app/pkg/database"
	"gorm.io/gorm"
)

type Driver struct {
	gorm.Model
	Id          uint64       `gorm:"primaryKey;" json:"id"`
	PhoneNumber string       `gorm:"not null;unique" json:"phonenumber"`
	FirstName   string       `gorm:"not null" json:"firstname"`
	LastName    string       `json:"lastname"`
	Email       string       `gorm:"not null;unique" json:"email"`
	Password    string       `gorm:"not null" json:"password"`
	City        string       `json:"city"`
	LicenceNum  string       `json:"licence"`
	Approved    bool         `gorm:"default:false" json:"approved"`
	Active      bool         `gorm:"default:true" json:"status"`
	Cab         *Vehicle     `json:"cab" gorm:"embedded"`
	Wallet      DriverWallet `json:"driverwallet"`
	TripHistory []Trip       `json:"trip_history" gorm:"foreignKey:DriverId;"`
	Rating      int          `gorm:"default:0" json:"driver_rating"`
}

// Add new driver to database
func (d *Driver) Add() error {
	db := database.Db
	err := db.AutoMigrate(&Driver{})
	if err != nil {
		return err
	}

	result := db.Create(&d)
	return result.Error
}

//get driver by key
func (d *Driver) Get(key, value string) (Driver, bool) {
	db := database.Db

	err := db.AutoMigrate(&Driver{})
	if err != nil {
		return Driver{}, false
	}

	driver := &Driver{}
	result := db.Where(key+"=?", value).First(&driver)

	if result.Error == gorm.ErrRecordNotFound {
		return *driver, false
	} else {
		return *driver, true
	}
}

//get all drivers in the database
func (d *Driver) GetAll() *[]Driver {
	db := database.Db

	drivers := &[]Driver{}
	db.Find(&drivers)

	return drivers
}

//update a driver by getting updated driver fields
//only update the not null driver fields
func (*Driver) Update(d Driver) error {
	db := database.Db

	driver := &Driver{}

	id := strconv.Itoa(int(d.Id))

	db.Where("id=?", id).First(&driver)

	result := db.Model(&driver).Updates(Driver{
		Cab: d.Cab,
	})

	db.Save(&driver)

	return result.Error
}

//delete driver by id
func (d *Driver) Delete(id uint64) error {
	db := database.Db
	db.AutoMigrate(&Driver{})

	result := db.Delete(&Driver{}, id)
	return result.Error
}

//block/unblock driver by toggling driver approved field
func (d *Driver) BlockUnblock(id uint64) error {
	db := database.Db

	driver := &Driver{}

	db.Where("id=?", id).First(&driver)

	// if !driver.Approved {
	// 	return errors.New("accesDenied")
	// }

	if driver.Approved {
		driver.Approved = false
		result := db.Save(&driver)
		return result.Error
	}

	driver.Approved = true
	result := db.Save(&driver)
	return result.Error
}
