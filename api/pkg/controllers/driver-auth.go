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

//Check if the driver already exist in the system.
//Redirect to the driver login page if driver exists.
//Redirect to the driver signup page if driver is new.
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
		}
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("invalid otp"))
		return
	}

	http.Redirect(w, r, "/driver/signup", http.StatusSeeOther)
}

func DriverHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache,no-store,must-revalidate")
	w.Header().Set("Content-Type", "application/json")

	c, _ := r.Cookie("jwt-token")
	tokenString := c.Value

	role, phone := auth.ParseJWT(tokenString)
	fmt.Println(role)
	// if err != nil {
	// 	if err == errors.New("invalidToken") {
	// 		http.Redirect(w, r, "/", http.StatusUnauthorized)
	// 		return
	// 	}
	// 	http.Redirect(w, r, "/", http.StatusSeeOther)
	// 	return
	// }
	driver := driver.GetDriver("phone_number", phone)

	response := models.Response{
		ResponseStatus:  "success",
		ResponseMessage: "Driver data fetched",
		ResponseData:    driver,
		Token:           "",
	}
	json.NewEncoder(w).Encode(&response)

}

//Create a driver model with values from the fronted.
//Pass the newly created driver model to driver services
//to insert the new driver to the database.
//Login the driver and open driver home after successful signup.
func DriverSignUp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache,no-store,must-revalidate")
	w.Header().Set("Content-Type", "application/json")

	newDriver := models.Driver{}
	json.NewDecoder(r.Body).Decode(&newDriver)
	newDriver.PhoneNumber = auth.GetDriver()
	hashpass, _ := bcrypt.GenerateFromPassword([]byte(newDriver.Password), bcrypt.DefaultCost)
	newDriver.Password = string(hashpass)
	//create a driver model with values from the fronted

	//pass the newly created driver model to driver services
	//to insert the new driver to the database
	//after successful signup login the driver and open driver home
	if err := driver.AddDriver(&newDriver); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Println("driver added")
	// setCookie(w, "driver", newDriver.PhoneNumber)

	http.Redirect(w, r, "/driver/driverhome", http.StatusSeeOther)
}

//get the existing driver by phone number from the database.
//Validate the entered password with stored hash password.
//Generate a JWT token for the driver after successful login.
//Store the JWT token in the cookie
func DriverLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache,no-store,must-revalidate")

	password := r.FormValue("password")

	phonenumber := auth.GetDriver()

	//get the existing driver by phone number from the database
	Driver := driver.GetDriver("phone_number", phonenumber)

	//validate the entered password with stored hash password
	if err := validPassword(password, Driver.Password); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	token, err := auth.GenerateJWT("driver", phonenumber)
	if err != nil {
		fmt.Println("jwt failed", err)
	}
	w.Header().Set("Token", token)
	// if err := setCookie(w, "driver", Driver.PhoneNumber); err != nil {
	// 	fmt.Println(err)
	// }
	//after successful login, generate a JWT token for the driver
	//save the generated token in the cookie
	http.Redirect(w, r, "/driver/driverhome", http.StatusOK)
}

func DriverLogout(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Cache-Control", "no-cache,no-store,must-revalidate")
	w.Header().Set("Token", "")

	// c, err := r.Cookie("jwt-token")
	// if err != nil {
	// 	http.Redirect(w, r, "/", http.StatusForbidden)
	// }
	// c.MaxAge = -1
	// http.SetCookie(w, &http.Cookie{
	// 	Name:   "jwt-token",
	// 	Value:  "",
	// 	Path:   "/",
	// 	Domain: "localhost:8080",
	// 	MaxAge: -1,
	// })
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func EditDriverProfile(w http.ResponseWriter, r *http.Request) {

}

func UpdateDriverProfile(w http.ResponseWriter, r *http.Request) {

}

func GetDrivers(w http.ResponseWriter, r *http.Request) {

}
