package controllers

import (
	"net/http"

	"github.com/shayamvlmna/cab-booking-app/app/models"
	"github.com/shayamvlmna/cab-booking-app/app/service/driver-service"

	"gorm.io/gorm"
)

func DriverSignUp(w http.ResponseWriter, r *http.Request) {

	newDriver := models.Driver{
		Model:       gorm.Model{},
		FirstName:   "",
		LastName:    "",
		PhoneNumber: 0,
		Email:       "",
		Password:    "",
	}

	err := driver.InsertDriver(&newDriver)
}
func DriverLogin(w http.ResponseWriter, r *http.Request) {

	err := driver.GetDriver()
}
