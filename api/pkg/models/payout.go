package models

import (
	"errors"
	"strconv"

	"gorm.io/gorm"

	"github.com/shayamvlmna/cab-booking-app/pkg/database"
)

type Payout struct {
	gorm.Model
	DriverId uint   `gorm:"primaryKey;unique" json:"driver_id"`
	Amount   string `json:"amount" gorm:"not null"`
	Bank     *Bank  `gorm:"embedded" json:"bank"`
	Status   string `json:"payout_status"`
}

type DriverWalletPayout struct {
	RequestId uint        `json:"requestId"`
	Driver    *DriverData `json:"driver"`
	Amount    string      `json:"amount" gorm:"not null"`
	Bank      *Bank       `gorm:"embedded" json:"bank"`
	Status    string      `json:"payout_status"`
}

func AddPayout(amount string, driverId uint) error {
	db := database.Db

	d := &Driver{}
	driver, er := d.Get("driver_id", strconv.Itoa(int(driverId)))
	if er != true {
		return errors.New("driver not found")
	}

	payout := &Payout{
		DriverId: driverId,
		Amount:   amount,
		Bank:     driver.BankAccount,
		Status:   "requested",
	}

	err := db.AutoMigrate(&Payout{})
	if err != nil {
		return err
	}

	result := db.Create(payout)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func GetPayouts() *[]DriverWalletPayout {
	db := database.Db
	db.AutoMigrate(&Payout{})

	payouts := []Payout{}

	db.Find(&payouts)

	db.AutoMigrate(&Driver{})

	payoutRequests := []DriverWalletPayout{}
	for _, val := range payouts {

		if val.Status == "paid" || val.Status == "closed" {
			continue
		}
		driver := &Driver{}

		db.Where("driver_id=?", val.DriverId).First(&driver)

		driverPayout := DriverWalletPayout{}

		cab := &CabData{
			VehicleId:    driver.Cab.VehicleId,
			Registration: driver.Cab.Registration,
			Brand:        driver.Cab.Brand,
			Category:     driver.Cab.Category,
			VehicleModel: driver.Cab.VehicleModel,
			Colour:       driver.Cab.Colour,
		}
		driverData := &DriverData{
			Id:          driver.DriverId,
			Phonenumber: driver.PhoneNumber,
			Firstname:   driver.FirstName,
			Lastname:    driver.LastName,
			Email:       driver.Email,
			City:        driver.City,
			LicenceNum:  driver.LicenceNum,
			Cab:         cab,
		}
		driverPayout.RequestId = val.ID
		driverPayout.Driver = driverData
		driverPayout.Bank = val.Bank
		driverPayout.Amount = val.Amount
		driverPayout.Status = val.Status

		payoutRequests = append(payoutRequests, driverPayout)

	}

	return &payoutRequests
}

func GetPayoutStatus(id uint) *Payout {
	db := database.Db
	db.AutoMigrate(&Payout{})

	payout := &Payout{}

	db.Where("driver_id=?", id).First(&payout)
	return payout
}

func UpdateCompletedPayoutRequest(id uint, status string) error {
	db := database.Db
	db.AutoMigrate(&Payout{})

	payout := &Payout{}

	db.Where("id=?", id).First(&payout)

	payout.Status = status
	db.Save(&payout)

	return nil
}
