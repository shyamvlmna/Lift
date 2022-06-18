package controllers

import (
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"gorm.io/gorm"

	models "github.com/shayamvlmna/cab-booking-app/pkg/models"
	driver "github.com/shayamvlmna/cab-booking-app/pkg/service/driver"
)

func DriverAuth(w http.ResponseWriter, r *http.Request) {
	phonenumber := r.FormValue("drvrphonenumber")
	data := map[string]string{
		"phone": phonenumber,
	}
	if driver.IsDriverExists("phone_number", phonenumber) {

		driverTemp.ExecuteTemplate(w, "driverLoginForm", data)
		return
	} else {

		driverTemp.ExecuteTemplate(w, "driverSignupForm", data)

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

	password := r.FormValue("drvrpassword")
	phonenumber := r.FormValue("drvrphonenumber")

	driver := driver.GetDriver("phone_number", phonenumber)

	if err := validPassword(password, driver.Password); err != nil {
		fmt.Println(err)
		data := map[any]any{
			"err": "invalid password",
		}
		driverTemp.ExecuteTemplate(w, "driverLoginForm.html", data)
		return
	}

}
