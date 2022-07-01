package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/shayamvlmna/cab-booking-app/pkg/database/redis"
	"github.com/shayamvlmna/cab-booking-app/pkg/models"
	"github.com/shayamvlmna/cab-booking-app/pkg/service/auth"
	"github.com/shayamvlmna/cab-booking-app/pkg/service/driver"
	"github.com/shayamvlmna/cab-booking-app/pkg/service/trip"
)

// DriverAuth Check if the driver already exist in the system.
//Redirect to the driver login page if driver exists.
//Redirect to the driver signup page if driver is new.
func DriverAuth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache,no-store,must-revalidate")
	w.Header().Set("Content-Type", "application/json")

	newDriver := &models.Driver{}

	err := json.NewDecoder(r.Body).Decode(&newDriver)
	if err != nil {
		return
	}

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
		err := json.NewEncoder(w).Encode(&response)
		if err != nil {
			return
		}
		return
	}

}

func DriverSignUpPage(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := &models.Response{
		ResponseStatus:  "success",
		ResponseMessage: "submit driver data",
		ResponseData:    nil,
	}
	err := json.NewEncoder(w).Encode(&response)
	if err != nil {
		return
	}
}

func DriverLoginPage(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := models.Response{
		ResponseStatus:  "success",
		ResponseMessage: "existing driver",
		ResponseData:    driver.GetDriver("phone_number", auth.GetPhone()).FirstName,
	}
	err := json.NewEncoder(w).Encode(&response)
	if err != nil {
		return
	}
}

// DriverSignUp Create a driver model with values from the fronted.
//Pass the newly created driver model to driver services
//to insert the new driver to the database.
//Login the driver and open driver home after successful signup.
func DriverSignUp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache,no-store,must-revalidate")
	w.Header().Set("Content-Type", "application/json")

	newDriver := models.Driver{}
	err := json.NewDecoder(r.Body).Decode(&newDriver)
	if err != nil {
		return
	}
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
		err := json.NewEncoder(w).Encode(&response)
		if err != nil {
			return
		}
		return
	}
	response := &models.Response{
		ResponseStatus:  "success",
		ResponseMessage: "signup sucess",
		ResponseData:    nil,
	}
	err = json.NewEncoder(w).Encode(&response)
	if err != nil {
		return
	}
}

// DriverLogin get the existing driver by phone number from the database.
//Validate the entered password with stored hash password.
//Generate a JWT token for the driver after successful login.
//Store the JWT token in the cookie
func DriverLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache,no-store,must-revalidate")
	w.Header().Set("Content-Type", "application/json")

	newDriver := &models.Driver{}

	err := json.NewDecoder(r.Body).Decode(&newDriver)
	if err != nil {
		return
	}
	defer r.Body.Close()

	password := newDriver.Password

	// phonenumber := newDriver.PhoneNumber
	phonenumber := auth.GetPhone()

	//get the existing driver by phone number from the database
	Driver := driver.GetDriver("phone_number", phonenumber)

	//validate the entered password with stored hash password
	if err := validPassword(password, Driver.Password); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		response := &models.Response{
			ResponseStatus:  "failed",
			ResponseMessage: "password authentication failed",
			ResponseData:    nil,
		}
		err := json.NewEncoder(w).Encode(&response)
		if err != nil {
			return
		}
		return
	}

	token, err := auth.GenerateJWT("driver", phonenumber)
	if err != nil {
		fmt.Println("jwt failed", err)
	}
	err = redis.StoreData("data", Driver)
	if err != nil {
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt-token",
		Value:    token,
		Path:     "/",
		MaxAge:   0,
		HttpOnly: true,
	})

	response := &models.Response{
		ResponseStatus:  "success",
		ResponseMessage: "login success",
		ResponseData:    token,
	}

	err = json.NewEncoder(w).Encode(&response)
	if err != nil {
		return
	}
	// http.Redirect(w, r, "/driver/driverhome", http.StatusSeeOther)
}

func DriverHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache,no-store,must-revalidate")
	w.Header().Set("Content-Type", "application/json")

	c, _ := r.Cookie("jwt-token")
	tokenString := c.Value

	role, phone := auth.ParseJWT(tokenString)

	fmt.Println(role, phone)
	driver := driver.GetDriver("phone_number", phone)
	cab := &models.CabData{
		VehicleId:    driver.Cab.VehicleId,
		Registration: driver.Cab.Registration,
		Brand:        driver.Cab.Brand,
		Category:     driver.Cab.Category,
		VehicleModel: driver.Cab.VehicleModel,
		Colour:       driver.Cab.Colour,
	}
	driverData := &models.DriverData{
		Id:          driver.Id,
		Phonenumber: driver.PhoneNumber,
		Firstname:   driver.FirstName,
		Lastname:    driver.LastName,
		Email:       driver.Email,
		City:        driver.City,
		LicenceNum:  driver.LicenceNum,
		Cab:         cab,
	}

	response := models.Response{
		ResponseStatus:  "success",
		ResponseMessage: "Driver data fetched",
		ResponseData:    &driverData,
	}
	err := json.NewEncoder(w).Encode(&response)
	if err != nil {
		return
	}
}

func DriverLogout(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Cache-Control", "no-cache,no-store,must-revalidate")
	w.Header().Set("Content-Type", "application/json")

	c, err := r.Cookie("jwt-token")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusUnauthorized)
		return
	}

	err = redis.DeleteData("data")
	if err != nil {
		return
	}

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
	err = json.NewEncoder(w).Encode(&response)
	if err != nil {
		return
	}
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

func AddCab(_ http.ResponseWriter, r *http.Request) {
	vehicle := &models.Vehicle{}
	err := json.NewDecoder(r.Body).Decode(&vehicle)
	if err != nil {
		return
	}
	c, _ := r.Cookie("jwt-token")
	tokenString := c.Value
	_, phone := auth.ParseJWT(tokenString)

	driver := driver.GetDriver("phone_number", phone)
	vehicle.DriverId = driver.Id
	driver.Cab = vehicle

	err = driver.Update(*driver)
	if err != nil {
		return
	}
}

func GetTrip(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	c, _ := r.Cookie("jwt-token")
	tokenString := c.Value
	_, phone := auth.ParseJWT(tokenString)

	driver := driver.GetDriver("phone_number", phone)

	if !driver.Approved {
		err := json.NewEncoder(w).Encode(&models.Response{
			ResponseStatus:  "failed",
			ResponseMessage: "not an approved driver",
			ResponseData:    nil,
		})
		if err != nil {
			return
		}
		return
	}
	ride := trip.GetRide()
	err := json.NewEncoder(w).Encode(&ride)
	if err != nil {
		return
	}
}

// AcceptTrip register the trip by user id and driver id
func AcceptTrip(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ride := &models.Ride{}
	err := json.NewDecoder(r.Body).Decode(&ride)
	if err != nil {
		return
	}
	c, _ := r.Cookie("jwt-token")
	tokenString := c.Value

	_, phone := auth.ParseJWT(tokenString)

	curDriver := driver.GetDriver("phone_number", phone)

	ride.DriverId = curDriver.Id

	if err := trip.RegisterTrip(ride); err != nil {
		response := &models.Response{
			ResponseStatus:  "failed",
			ResponseMessage: "error registering trip",
			ResponseData:    nil,
		}
		err := json.NewEncoder(w).Encode(&response)
		if err != nil {
			return
		}
		return
	}

	response := &models.Response{
		ResponseStatus:  "success",
		ResponseMessage: "trip successfully registered",
		ResponseData:    ride,
	}

	err = json.NewEncoder(w).Encode(&response)
	if err != nil {
		return
	}
}

func EndTrip(w http.ResponseWriter, r *http.Request) {

}
