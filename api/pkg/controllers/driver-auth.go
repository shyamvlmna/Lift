package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/go-redis/redis/v9"
	models "github.com/shayamvlmna/cab-booking-app/pkg/models"
	auth "github.com/shayamvlmna/cab-booking-app/pkg/service/auth"
	driver "github.com/shayamvlmna/cab-booking-app/pkg/service/driver"
)

//Check if the user already exist in the system.
//Redirect to the user login page if user exists.
//Redirect to the user signup page if user is new.
func DriverAuth(w http.ResponseWriter, r *http.Request) {
	phonenumber := r.FormValue("phonenumber")
	auth.StoreDriver(phonenumber)
	if driver.IsDriverExists("phone_number", phonenumber) {
		http.Redirect(w, r, "/driver/loginpage", http.StatusSeeOther)
		return
	} else {
		go auth.SetOtp(phonenumber)
		http.Redirect(w, r, "/driver/enterotp", http.StatusSeeOther)
	}
}
func ValidateDriverOtp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache,no-store,must-revalidate")
	otp := r.FormValue("otp")

	phone := auth.GetDriver()

	if err := auth.ValidateOTP(phone, otp); err != nil {
		if err == redis.Nil {
			w.WriteHeader(http.StatusRequestTimeout)
			w.Write([]byte("otp expired"))
			return
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("invalid otp"))
			return
		}

	}

	http.Redirect(w, r, "/driver/signup", http.StatusSeeOther)
}

func DriverSignUpPage(w http.ResponseWriter, r *http.Request) {
	driverTemp.ExecuteTemplate(w, "driverSignupForm.html", nil)
}
func DriverLoginPage(w http.ResponseWriter, r *http.Request) {
	driverTemp.ExecuteTemplate(w, "driverLoginForm.html", nil)
}

//Create a user model with values from the fronted.
//Pass the newly created user model to user services
//to insert the new user to the database.
//Login the user and open user home after successful signup.
func DriverSignUp(w http.ResponseWriter, r *http.Request) {

	newDriver := models.Driver{}
	json.NewDecoder(r.Body).Decode(&newDriver)
	newDriver.PhoneNumber = auth.GetDriver()
	hashpass, _ := bcrypt.GenerateFromPassword([]byte(newDriver.Password), bcrypt.DefaultCost)
	newDriver.Password = string(hashpass)
	//create a user model with values from the fronted

	//pass the newly created user model to user services
	//to insert the new user to the database
	//after successful signup login the user and open user home
	if err := driver.AddDriver(&newDriver); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Println("driver added")
	setCookie(w, newDriver.PhoneNumber)
	http.Redirect(w, r, "/driver/driverhome", http.StatusSeeOther)

}

//get the existing user by phone number from the database.
//Validate the entered password with stored hash password.
//Generate a JWT token for the user after successful login.
//Store the JWT token in the cookie
func DriverLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache,no-store,must-revalidate")
	password := r.FormValue("drvrpassword")
	phonenumber := auth.GetDriver()

	newDriver := models.Driver{}
	json.NewDecoder(r.Body).Decode(&newDriver)
	//get the existing user by phone number from the database
	driver := driver.GetDriver("phone_number", phonenumber)

	//validate the entered password with stored hash password
	if err := validPassword(password, driver.Password); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if err := setCookie(w, driver.PhoneNumber); err != nil {
		fmt.Println(err)
	}
	//after successful login, generate a JWT token for the user
	//save the generated token in the cookie
	http.Redirect(w, r, "/driver/driverhome", http.StatusOK)
}
func EditDriverProfile(w http.ResponseWriter, r *http.Request) {

}

func UpdateDriverProfile(w http.ResponseWriter, r *http.Request) {

}

func GetDrivers(w http.ResponseWriter, r *http.Request) {

}
