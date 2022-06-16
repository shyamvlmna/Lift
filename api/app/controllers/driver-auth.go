package controllers

import (
	"fmt"
	"net/http"

	"github.com/shayamvlmna/cab-booking-app/app/models"
	"github.com/shayamvlmna/cab-booking-app/app/service/driver"
	"golang.org/x/crypto/bcrypt"

	"gorm.io/gorm"
)

func DriverAuth(w http.ResponseWriter, r *http.Request) {
	phonenumber := r.FormValue("drvrphonenumber")
	if driver.IsDriverExists(phonenumber) {
		// UserLogin(w, r)
		http.Redirect(w, r, "/driver/login", 200)
	} else {
		DriverSignUp(w, r)
		// http.Redirect(w, r, "/driver/signup", 200)
	}
}
func DriverSignUp(w http.ResponseWriter, r *http.Request) {

	firstname := r.FormValue("drvrfirstname")
	lastname := r.FormValue("drvrlastname")
	phonenumber := r.FormValue("drvrphonenumber")
	email := r.FormValue("drvremail")
	password := r.FormValue("drvrpassword")

	hashpass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
		return
	}

	newDriver := models.Driver{
		Model:       gorm.Model{},
		FirstName:   firstname,
		LastName:    lastname,
		PhoneNumber: phonenumber,
		Email:       email,
		Password:    string(hashpass),
	}

	if err = driver.AddDriver(&newDriver); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("driver added")
	}
}
func DriverLogin(w http.ResponseWriter, r *http.Request) {

	// err := driver.GetDriver()
	// if err != nil {
	// 	fmt.Println(err)
	// }
}
