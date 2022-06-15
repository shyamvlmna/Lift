package controllers

import (
	"net/http"

	"github.com/shayamvlmna/cab-booking-app/service/driver"
)

func DriverSignUp(w http.ResponseWriter, r *http.Request) {

	err := driver.AddDriver()
}
func DriverLogin(w http.ResponseWriter, r *http.Request) {

	err := driver.GetDriver()
}
