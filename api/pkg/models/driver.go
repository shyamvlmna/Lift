package models

import (
	"errors"
	"strconv"

	"gorm.io/gorm"

	"github.com/shayamvlmna/cab-booking-app/pkg/database"
)

type Driver struct {
	gorm.Model
	DriverId      uint     `gorm:"primaryKey;unique;autoIncrement;" json:"driverid"`
	PhoneNumber   string   `gorm:"not null;unique" json:"phonenumber"`
	FirstName     string   `gorm:"not null" json:"firstname"`
	LastName      string   `json:"lastname"`
	Email         string   `gorm:"not null;unique" json:"email"`
	Password      string   `gorm:"not null" json:"password"`
	City          string   `json:"city"`
	LicenceNum    string   `json:"licence"`
	Rating        int      `gorm:"default:0" json:"driver_rating"`
	Approved      bool     `gorm:"default:false" json:"approved"`
	Active        bool     `gorm:"default:true" json:"status"`
	Cab           *Vehicle `json:"cab" gorm:"embedded"`
	WalletBalance uint     `json:"driverwallet"  gorm:"default:0;"`
	BankAccount   *Bank    `json:"bank_account" gorm:"embedded"`
}

type Bank struct {
	AccountHolderName string `json:"account_holder_name"`
	BankName          string `json:"bank_name"`
	AccountNumber     string `json:"account_number"`
	IFSC              string `gorm:"ifsc" json:"ifsc"`
}

type DriverTransaction struct {
	gorm.Model
	DriverId uint   `gorm:"index:driver_id" json:"driver_id"`
	Type     string `gorm:"" json:"type"`
	Amount   string `json:"amount"`
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

// Get driver by key
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

// GetAll drivers in the database
func (*Driver) GetAll() (*[]Driver, error) {
	db := database.Db

	drivers := &[]Driver{}
	result := db.Find(&drivers)

	return drivers, result.Error
}

func DriverRequests() (*[]Driver, error) {
	db := database.Db

	db.AutoMigrate(&Driver{})

	drivers := &[]Driver{}

	result := db.Find(&drivers, "approved=?", false)

	return drivers, result.Error
}

// Update a driver by getting updated driver fields
//only update the not null driver fields
func (*Driver) Update(d Driver) error {
	db := database.Db

	driver := &Driver{}

	id := strconv.Itoa(int(d.DriverId))

	db.Where("driver_id=?", id).First(&driver)

	result := db.Model(&driver).Updates(Driver{
		PhoneNumber:   d.PhoneNumber,
		FirstName:     d.FirstName,
		LastName:      d.LastName,
		Email:         d.Email,
		Password:      d.Password,
		City:          d.City,
		LicenceNum:    d.LicenceNum,
		Rating:        d.Rating,
		Cab:           d.Cab,
		WalletBalance: 0,
		BankAccount:   d.BankAccount,
	})

	db.Save(&driver)

	return result.Error
}

// Delete driver by id
func (d *Driver) Delete(id uint64) error {
	db := database.Db
	err := db.AutoMigrate(&Driver{})
	if err != nil {
		return err
	}

	result := db.Delete(&Driver{}, id)
	return result.Error
}

// BlockUnblock driver by toggling driver approved field
func (*Driver) BlockUnblock(id uint) error {
	db := database.Db

	driver := &Driver{}

	db.Where("driver_id=?", id).First(&driver)

	if driver.Active {
		result := db.Model(&driver).Update("active", false)
		return result.Error
	}

	result := db.Model(&driver).Update("active", true)
	return result.Error
}

func (*Driver) ApproveToDrive(id uint) error {

	db := database.Db

	driver := &Driver{}

	db.Where("driver_id=?", id).First(&driver)

	if driver.Approved {
		result := db.Model(&driver).Update("approved", false)
		return result.Error
	}

	result := db.Model(&driver).Update("approved", true)
	return result.Error
}

func GetBankDetails(id uint) (*Bank, error) {
	db := database.Db

	db.AutoMigrate(&Driver{})

	driver := &Driver{}

	db.Where("driver_id=?", id).First(&driver)

	bank := driver.BankAccount

	if bank.AccountNumber == "" {
		return nil, errors.New("bank account not added")
	}
	return bank, nil
}

func (b *Bank) UpdateBank(id uint, bank *Bank) error {

	db := database.Db
	db.AutoMigrate(&Driver{})

	driver := &Driver{}

	db.Where("driver_id=?", id).First(&driver)
	result := db.Model(&driver).Updates(Driver{
		BankAccount: bank,
	})
	db.Save(&driver)
	return result.Error
}
