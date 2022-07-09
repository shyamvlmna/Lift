package models

import (
	"errors"
	"strconv"

	"github.com/shayamvlmna/cab-booking-app/pkg/database"
	"gorm.io/gorm"
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

type Payouts struct {
	gorm.Model
	DriverId uint   `gorm:"primaryKey;unique" json:"driver_id"`
	Amount   string `json:"amount" gorm:"not null"`
	Bank     *Bank  `gorm:"embedded" json:"bank"`
	Status   string `json:"payout_status"`
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
func (d *Driver) GetAll() *[]Driver {
	db := database.Db

	drivers := &[]Driver{}
	db.Find(&drivers)

	return drivers
}

// Update a driver by getting updated driver fields
//only update the not null driver fields
func (*Driver) Update(d Driver) error {
	db := database.Db

	driver := &Driver{}

	id := strconv.Itoa(int(d.DriverId))

	db.Where("id=?", id).First(&driver)

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
func (d *Driver) BlockUnblock(id uint64) error {
	db := database.Db

	driver := &Driver{}

	db.Where("driver_id=?", id).First(&driver)

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

func AddPayout(amount string, driverId uint) error {
	db := database.Db

	d := &Driver{}
	driver, er := d.Get("driver_id", strconv.Itoa(int(driverId)))
	if er != true {
		return errors.New("driver not found")
	}

	payout := &Payouts{
		DriverId: driverId,
		Amount:   amount,
		Bank:     driver.BankAccount,
		Status:   "requested",
	}

	err := db.AutoMigrate(&Payouts{})
	if err != nil {
		return err
	}

	result := db.Create(payout)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func GetPayouts() *[]Payouts {
	db := database.Db
	db.AutoMigrate(&Payouts{})

	payouts := &[]Payouts{}

	db.Find(&payouts)

	return payouts
}

func GetPayoutStatus(id uint) *Payouts {
	db := database.Db
	db.AutoMigrate(&Payouts{})

	payout := &Payouts{}

	db.Where("driver_id=?", id).First(&payout)
	return payout
}
