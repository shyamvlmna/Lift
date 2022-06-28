package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	database "github.com/shayamvlmna/cab-booking-app/pkg/database/postgresql"
	redis "github.com/shayamvlmna/cab-booking-app/pkg/database/redis"
	models "github.com/shayamvlmna/cab-booking-app/pkg/models"
	auth "github.com/shayamvlmna/cab-booking-app/pkg/service/auth"
	driver "github.com/shayamvlmna/cab-booking-app/pkg/service/driver"
	trip "github.com/shayamvlmna/cab-booking-app/pkg/service/trip"
)

//Check if the driver already exist in the system.
//Redirect to the driver login page if driver exists.
//Redirect to the driver signup page if driver is new.
func DriverAuth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache,no-store,must-revalidate")
	w.Header().Set("Content-Type", "application/json")

	newDriver := &models.Driver{}

	json.NewDecoder(r.Body).Decode(&newDriver)

	phonenumber := newDriver.PhoneNumber

	if phonenumber != "" {
		auth.StorePhone(phonenumber)
		if driver.IsDriverExists("phone_number", phonenumber) {
			http.Redirect(w, r, "/driver/loginpage", http.StatusSeeOther)
			return
		} else {
			if err := auth.SetOtp(phonenumber); err != nil {
				fmt.Println(err)
				return
			}
			http.Redirect(w, r, "/driver/enterotp", http.StatusSeeOther)
			return
		}
	} else {
		response := models.Response{
			ResponseStatus:  "fail",
			ResponseMessage: "phonenumber required",
			ResponseData:    nil,
		}
		json.NewEncoder(w).Encode(&response)
		return
	}

}

func DriverSignUpPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := &models.Response{
		ResponseStatus:  "success",
		ResponseMessage: "submit driver data",
		ResponseData:    nil,
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
	defer r.Body.Close()

	newDriver.PhoneNumber = auth.GetPhone()
	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(newDriver.Password), bcrypt.DefaultCost)
	newDriver.Password = string(hashPassword)
	//create a driver model with values from the fronted

	//pass the newly created driver model to driver services
	//to insert the new driver to the database
	//after successful signup login the driver and open driver home
	if err := driver.AddDriver(&newDriver); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		response := &models.Response{
			ResponseStatus:  "failed",
			ResponseMessage: "signup failed",
			ResponseData:    nil,
		}
		json.NewEncoder(w).Encode(&response)
		return
	}
	response := &models.Response{
		ResponseStatus:  "success",
		ResponseMessage: "signup sucess",
		ResponseData:    nil,
	}
	json.NewEncoder(w).Encode(&response)
}

//get the existing driver by phone number from the database.
//Validate the entered password with stored hash password.
//Generate a JWT token for the driver after successful login.
//Store the JWT token in the cookie
func DriverLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache,no-store,must-revalidate")
	w.Header().Set("Content-Type", "application/json")

	newDriver := &models.Driver{}

	json.NewDecoder(r.Body).Decode(&newDriver)
	defer r.Body.Close()

	password := newDriver.Password

	// phonenumber := newDriver.PhoneNumber
	phonenumber := auth.GetPhone()

	//get the existing driver by phone number from the database
	Driver := driver.GetDriver("phone_number", phonenumber)

	if !Driver.Active {
		response := &models.Response{
			ResponseStatus:  "failed",
			ResponseMessage: "driver not active",
			ResponseData:    nil,
		}
		json.NewEncoder(w).Encode(&response)
		return
	}

	//validate the entered password with stored hash password
	if err := validPassword(password, Driver.Password); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		response := &models.Response{
			ResponseStatus:  "failed",
			ResponseMessage: "password authentication failed",
			ResponseData:    nil,
		}
		json.NewEncoder(w).Encode(&response)
		return
	}

	token, err := auth.GenerateJWT("driver", phonenumber)
	if err != nil {
		fmt.Println("jwt failed", err)
	}
	redis.StoreData("data", Driver)
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt-token",
		Value:    token,
		Path:     "/",
		MaxAge:   0,
		HttpOnly: true,
	})
	http.Redirect(w, r, "/driver/driverhome", http.StatusSeeOther)
}

func DriverHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache,no-store,must-revalidate")
	w.Header().Set("Content-Type", "application/json")

	c, _ := r.Cookie("jwt-token")
	tokenString := c.Value

	role, phone := auth.ParseJWT(tokenString)

	fmt.Println(role, phone)
	driver := driver.GetDriver("phone_number", phone)

	response := models.Response{
		ResponseStatus:  "success",
		ResponseMessage: "Driver data fetched",
		ResponseData:    driver,
		Token:           tokenString,
	}
	json.NewEncoder(w).Encode(&response)
}

func DriverLogout(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Cache-Control", "no-cache,no-store,must-revalidate")
	w.Header().Set("Content-Type", "application/json")

	c, err := r.Cookie("jwt-token")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusUnauthorized)
		return
	}

	redis.DeleteData("data")

	c.Value = ""
	c.Path = "/"
	c.MaxAge = -1
	http.SetCookie(w, c)

	response := &models.Response{
		ResponseStatus:  "success",
		ResponseMessage: "logout success",
		ResponseData:    nil,
	}
	// http.Redirect(w, r, "/", http.StatusSeeOther)
	json.NewEncoder(w).Encode(&response)
}

func EditDriverProfile(w http.ResponseWriter, r *http.Request) {
}

func UpdateDriverProfile(w http.ResponseWriter, r *http.Request) {
}

func GetDrivers(w http.ResponseWriter, r *http.Request) {
}

func RegisterDriver(w http.ResponseWriter, r *http.Request) {
	// city := r.FormValue("city")
	// dlNumber := r.FormValue("driving_licence")
}

func AddCab(w http.ResponseWriter, r *http.Request) {
	vehicle := models.Vehicle{}
	json.NewDecoder(r.Body).Decode(&vehicle)

	database.Insert(&vehicle)
}

func GetTrip(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	c, _ := r.Cookie("jwt-token")
	tokenString := c.Value
	_, phone := auth.ParseJWT(tokenString)

	driver := driver.GetDriver("phone_number", phone)

	if !driver.Approved {
		json.NewEncoder(w).Encode(&models.Response{
			ResponseStatus:  "failed",
			ResponseMessage: "not an approved driver",
			ResponseData:    nil,
		})
		return
	}
	ride := trip.GetRide()
	json.NewEncoder(w).Encode(&ride)
}

func EndTrip(w http.ResponseWriter, r *http.Request) {

}
