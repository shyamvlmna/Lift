package models

import (
	"errors"
	"log"
	"strconv"
	"time"

	"gorm.io/gorm"

	"github.com/shayamvlmna/cab-booking-app/pkg/database"
)

type Payout struct {
	gorm.Model
	DriverId uint   `gorm:"index:driver_id" json:"driver_id"`
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

func AddPayoutRequest(amount string, driverId uint) error {
	db := database.Db

	d := &Driver{}

	driver, exist := d.Get("driver_id", strconv.Itoa(int(driverId)))
	if !exist {
		return errors.New("driver not found")
	}

	if err := db.AutoMigrate(&Payout{}); err != nil {
		return err
	}

	payouts := []Payout{}
	result := db.Where("driver_id = ? AND status = ?", driverId, "pending").Find(&payouts)
	//db.Where("driver_id <> ?", driverId).Find(&payouts)

	//check whether if there is an existing payout request to proceed
	//if exists no new payout request to be allowed
	if len(payouts) == 0 {

		payout := &Payout{
			DriverId: driverId,
			Amount:   amount,
			Bank:     driver.BankAccount,
			Status:   "pending",
		}

		result = db.Create(payout)
		if result.Error != nil {
			if result.Error != nil {
				log.Println(result.Error)
				return errors.New("something went wrong")
			}
		}

		if err := db.AutoMigrate(&Driver{}); err != nil {
			log.Println(err)
			return errors.New("something went wrong")
		}

		if err := db.Model(&Driver{}).Where("driver_id=?", driverId).UpdateColumn("wallet_balance", gorm.Expr("wallet_balance - ?", amount)).Error; err != nil {
			log.Println(err)
			return errors.New("something went wrong")
		}

		return nil
	}
	return errors.New("pending request exist")

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

type PayoutResponse struct {
	RequestDate time.Time `json:"requestdate"`
	Amount      string    `json:"amount"`
	Status      string    `json:"status"`
	Bank        *Bank     `json:"bank"`
}

func GetPayoutStatus(id uint) ([]PayoutResponse, error) {
	db := database.Db

	err := db.AutoMigrate(&Payout{})
	if err != nil {
		log.Println(err)
		return nil, errors.New("something went wrong")
	}

	payouts := []Payout{}

	db.Find(&payouts)

	pendingPayouts := []PayoutResponse{}

	for _, p := range payouts {
		if p.DriverId != id || p.Status == "paid" {
			continue
		}
		payoutresp := PayoutResponse{}
		payoutresp.RequestDate = p.CreatedAt
		payoutresp.Amount = p.Amount
		payoutresp.Status = p.Status
		payoutresp.Bank = p.Bank
		pendingPayouts = append(pendingPayouts, payoutresp)
	}

	return pendingPayouts, nil
}

func PayoutHistory(id uint) []PayoutResponse {
	db := database.Db
	db.AutoMigrate(&Payout{})

	payouts := []Payout{}

	payoutHistory := []PayoutResponse{}

	db.Find(&payouts)

	for _, p := range payouts {
		if p.DriverId != id {
			continue
		}

		payoutresp := PayoutResponse{
			RequestDate: p.CreatedAt,
			Amount:      p.Amount,
			Status:      p.Status,
			Bank:        p.Bank,
		}

		payoutHistory = append(payoutHistory, payoutresp)
	}

	return payoutHistory
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
