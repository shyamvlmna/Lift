package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/shayamvlmna/lift/pkg/database/redis"
	"github.com/shayamvlmna/lift/pkg/models"
	"github.com/shayamvlmna/lift/pkg/service/auth"
	"github.com/shayamvlmna/lift/pkg/service/driver"
	"github.com/shayamvlmna/lift/pkg/service/trip"
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
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(r.Body)

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
			ResponseStatus:  "failed",
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

func EnterOTPDriver(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := models.Response{
		ResponseStatus:  "success",
		ResponseMessage: "new driver",
		ResponseData:    "verify with otp",
	}
	err := json.NewEncoder(w).Encode(&response)
	if err != nil {
		return
	}
}

func DriverSignUpPage(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := &models.Response{
		ResponseStatus:  "OTP validation success",
		ResponseMessage: "submit driver data",
		ResponseData:    "render signup page",
	}

	err := json.NewEncoder(w).Encode(&response)
	if err != nil {
		return
	}
}

func DriverLoginPage(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	phonenumber := auth.GetPhone()

	response := models.Response{
		ResponseStatus:  "success",
		ResponseMessage: "existing driver",
		ResponseData:    driver.GetDriver("phone_number", phonenumber).FirstName,
	}

	auth.StorePhone(phonenumber)

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
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(r.Body)

	if err := driver.RegisterDriver(&newDriver); err != nil {
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
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(r.Body)

	password := newDriver.Password

	// phonenumber := newDriver.PhoneNumber

	phonenumber := auth.GetPhone()

	//get the existing driver by phone number from the database
	Driver := driver.GetDriver("phone_number", phonenumber)

	if !Driver.Active {
		w.WriteHeader(http.StatusForbidden)
		response := &models.Response{
			ResponseStatus:  "failed",
			ResponseMessage: "driver not active",
			ResponseData:    nil,
		}
		err := json.NewEncoder(w).Encode(&response)
		if err != nil {
			return
		}
		return
	}

	//validate the entered password with stored hash password
	if !ValidPassword(password, Driver.Password) {
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

// DriverHome get the logged-in driver data from redis
//if err get from the primary database
//render the driver home page
func DriverHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache,no-store,must-revalidate")
	w.Header().Set("Content-Type", "application/json")

	driver, err := GetDriverFromCookie(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(&ErrorParsingCookie)
		return
	}

	driverData := &models.DriverData{
		Id:          driver.DriverId,
		Phonenumber: driver.PhoneNumber,
		Firstname:   driver.FirstName,
		Lastname:    driver.LastName,
		Email:       driver.Email,
		City:        driver.City,
		LicenceNum:  driver.LicenceNum,
	}

	response := models.Response{
		ResponseStatus:  "success",
		ResponseMessage: "Driver data fetched",
		ResponseData:    &driverData,
	}

	err = json.NewEncoder(w).Encode(&response)
	if err != nil {
		return
	}
}

// DriverLogout delete the stored driver data from redis
//also expire the cookie stored
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

	err = json.NewEncoder(w).Encode(&response)
	if err != nil {
		return
	}
	// http.Redirect(w, r, "/", http.StatusSeeOther)
}

func RegisterToDrive(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	//d, err := GetDriverFromCookie(r)
	//
	//driver.RegisterToDrive(d)

	//TODO : complete driver request to drive
}

func UploadDocuments(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(20 << 30)

	licence, licencehandler, err := r.FormFile("licence")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(licencehandler.Filename)
	fmt.Println(licencehandler.Size)
	fmt.Println(licencehandler.Header)

	tmpFile, err := ioutil.TempFile("licence-doc", "licence-*.png")
	if err != nil {
		fmt.Println(err)
	}
	defer tmpFile.Close()

	fileBytes, err := ioutil.ReadAll(licence)
	if err != nil {
		fmt.Println(err)
	}

	tmpFile.Write(fileBytes)
}

func EditDriverProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	driver, err := GetDriverFromCookie(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(&ErrorParsingCookie)
		return
	}

	cab := &models.CabData{
		VehicleId:    driver.Cab.VehicleId,
		Registration: driver.Cab.Registration,
		Brand:        driver.Cab.Brand,
		Category:     driver.Cab.Category,
		VehicleModel: driver.Cab.VehicleModel,
		Colour:       driver.Cab.Colour,
	}

	driverData := &models.DriverData{
		Id:          driver.DriverId,
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

	err = json.NewEncoder(w).Encode(&response)
	if err != nil {
		return
	}
}

func UpdateDriverProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	d := models.Driver{}

	err := json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		return
	}

	err = d.Update(d)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := &models.Response{
			ResponseStatus:  "failed",
			ResponseMessage: "driver profile update failed",
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
		ResponseMessage: "driver profile update success",
		ResponseData:    nil,
	}
	err = json.NewEncoder(w).Encode(&response)
	if err != nil {
		return
	}
}

func AddCab(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vehicle := &models.Vehicle{}

	err := json.NewDecoder(r.Body).Decode(&vehicle)
	if err != nil {
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(r.Body)

	driver, err := GetDriverFromCookie(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(&ErrorParsingCookie)
		return
	}

	vehicle.DriverId = driver.DriverId
	driver.Cab = vehicle

	err = driver.Update(*driver)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := &models.Response{
			ResponseStatus:  "failed",
			ResponseMessage: "error adding cab",
			ResponseData:    nil,
		}
		if err := json.NewEncoder(w).Encode(&response); err != nil {
			return
		}
		return
	}

	response := &models.Response{
		ResponseStatus:  "success",
		ResponseMessage: "cab added",
		ResponseData:    nil,
	}

	if err := json.NewEncoder(w).Encode(&response); err != nil {
		return
	}
}

func EditCab(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	driver, err := GetDriverFromCookie(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(&ErrorParsingCookie)
		return
	}

	cab := &models.CabData{
		VehicleId:    driver.Cab.VehicleId,
		Registration: driver.Cab.Registration,
		Brand:        driver.Cab.Brand,
		Category:     driver.Cab.Category,
		VehicleModel: driver.Cab.VehicleModel,
		Colour:       driver.Cab.Colour,
	}

	response := &models.Response{
		ResponseStatus:  "success",
		ResponseMessage: "edit cab data",
		ResponseData:    cab,
	}

	json.NewEncoder(w).Encode(&response)
}

func UpdateCab(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	cab := &models.CabData{}

	err := json.NewDecoder(r.Body).Decode(&cab)
	if err != nil {
		return
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(r.Body)

	driver, err := GetDriverFromCookie(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(&ErrorParsingCookie)
		return
	}

	driver.Cab.Registration = cab.Registration
	driver.Cab.Brand = cab.Brand
	driver.Cab.Category = cab.Category
	driver.Cab.VehicleModel = cab.VehicleModel
	driver.Cab.Colour = cab.Colour

	err = driver.Update(*driver)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := &models.Response{
			ResponseStatus:  "failed",
			ResponseMessage: "error updating cab",
			ResponseData:    nil,
		}
		if err := json.NewEncoder(w).Encode(&response); err != nil {
			return
		}
		return
	}

	response := &models.Response{
		ResponseStatus:  "success",
		ResponseMessage: "cab updated",
		ResponseData:    cab,
	}

	if err := json.NewEncoder(w).Encode(&response); err != nil {
		return
	}

}

func AddBankPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := &models.Response{
		ResponseStatus:  "success",
		ResponseMessage: "add bank account details",
		ResponseData:    nil,
	}

	json.NewEncoder(w).Encode(&response)
}

func AddBankAccount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	bank := &models.Bank{}
	json.NewDecoder(r.Body).Decode(&bank)
	r.Body.Close()

	d, err := GetDriverFromCookie(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(&ErrorParsingCookie)
		return
	}

	d.BankAccount = bank

	err = d.Update(*d)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := &models.Response{
			ResponseStatus:  "failed",
			ResponseMessage: "bank account not added",
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
		ResponseMessage: "bank account successfully added",
		ResponseData:    nil,
	}
	err = json.NewEncoder(w).Encode(&response)
	if err != nil {
		return
	}
}

func EditBankDetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	d, err := GetDriverFromCookie(r)
	if err != nil {
		json.NewEncoder(w).Encode(&ErrorParsingCookie)
		return
	}

	bankAccount, err := driver.GetBankDetails(d.DriverId)
	if err != nil {
		resp := &models.Response{
			ResponseStatus:  "failed",
			ResponseMessage: "fetching bank account details failed",
			ResponseData:    err,
		}
		json.NewEncoder(w).Encode(resp)
		return
	}
	resp := &models.Response{
		ResponseStatus:  "success",
		ResponseMessage: "bank account details fetched",
		ResponseData:    bankAccount,
	}
	json.NewEncoder(w).Encode(resp)
}

func UpdateBankDetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	d, err := GetDriverFromCookie(r)
	if err != nil {
		json.NewEncoder(w).Encode(&ErrorParsingCookie)
		return
	}

	bankDetails := &models.Bank{}
	json.NewDecoder(r.Body).Decode(&bankDetails)

	if err := driver.UpdateBankDetails(d.DriverId, bankDetails); err != nil {
		resp := &models.Response{
			ResponseStatus:  "failed",
			ResponseMessage: "updating bank account failed",
			ResponseData:    err,
		}
		json.NewEncoder(w).Encode(resp)
		return
	}
	resp := &models.Response{
		ResponseStatus:  "success",
		ResponseMessage: "bank account successfully updated",
		ResponseData:    bankDetails,
	}
	json.NewEncoder(w).Encode(resp)

}

//GetTrip returns the available trips booked by the users
func GetTrip(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	driver, err := GetDriverFromCookie(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(&ErrorParsingCookie)
		return
	}

	if !driver.Approved {
		w.WriteHeader(http.StatusUnauthorized)
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

	if driver.Cab.Registration == "" {
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(&models.Response{
			ResponseStatus:  "failed",
			ResponseMessage: "cab not added",
			ResponseData:    nil,
		})
		if err != nil {
			return
		}
		return
	}

	ride := trip.GetRide()
	err = json.NewEncoder(w).Encode(&ride)
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
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(r.Body)

	curDriver, err := GetDriverFromCookie(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(&ErrorParsingCookie)
		return
	}

	ride.DriverId = curDriver.DriverId

	if err := trip.RegisterTrip(ride); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := &models.Response{
			ResponseStatus:  "failed",
			ResponseMessage: "error registering trip",
			ResponseData:    nil,
		}
		if err := json.NewEncoder(w).Encode(&response); err != nil {
			return
		}
		return
	}

	response := &models.Response{
		ResponseStatus:  "success",
		ResponseMessage: "trip successfully registered",
		ResponseData:    ride,
	}

	otp, _ := auth.TripCode()

	if err = redis.Set("tripcode-"+strconv.Itoa(int(ride.DriverId))+strconv.Itoa(int(ride.UserId)), otp); err != nil {
		return
	}

	if err := redis.StoreTrip("trip-"+strconv.Itoa(int(curDriver.DriverId)), ride); err != nil {
		return
	}

	if err = json.NewEncoder(w).Encode(&response); err != nil {
		return
	}
}

type otp struct {
	TripCode int `json:"tripcode"`
}

//MatchTripCode match the trip code from user
func MatchTripCode(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	code := &otp{}

	err := json.NewDecoder(r.Body).Decode(&code)
	if err != nil {
		return
	}
	driver, err := GetDriverFromCookie(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(&ErrorParsingCookie)
		return
	}

	ride, err := redis.GetTrip("trip-" + strconv.Itoa(int(driver.DriverId)))
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		response := &models.Response{
			ResponseStatus:  "failed",
			ResponseMessage: "couldn't find a trip",
			ResponseData:    err,
		}
		json.NewEncoder(w).Encode(&response)
		return
	}

	tripCode, _ := redis.Get("tripcode-" + strconv.Itoa(int(ride.DriverId)) + strconv.Itoa(int(ride.UserId)))
	fmt.Println(tripCode)

	matchCode, _ := strconv.Atoi(tripCode)
	fmt.Println(matchCode)
	fmt.Println(tripCode)
	if code.TripCode == matchCode {
		redis.DeleteData("tripcode-" + strconv.Itoa(int(ride.DriverId)) + strconv.Itoa(int(ride.UserId)))
		http.Redirect(w, r, "/driver/startrip", http.StatusSeeOther)
		return
	}
	w.WriteHeader(http.StatusBadRequest)
	response := &models.Response{
		ResponseStatus:  "failed",
		ResponseMessage: "trip code doesn't match",
		ResponseData:    nil,
	}
	err = json.NewEncoder(w).Encode(&response)
	if err != nil {
		return
	}
}

func StartTrip(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	curDriver, err := GetDriverFromCookie(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(&ErrorParsingCookie)
		return
	}

	driverID := strconv.Itoa(int(curDriver.DriverId))

	ride, err := redis.GetTrip("trip-" + driverID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := &models.Response{
			ResponseStatus:  "failed",
			ResponseMessage: "error geting trip data",
			ResponseData:    err,
		}
		json.NewEncoder(w).Encode(&response)
		return
	}
	response := &models.Response{
		ResponseStatus:  "success",
		ResponseMessage: "start trip",
		ResponseData:    ride,
	}
	if err = json.NewEncoder(w).Encode(&response); err != nil {
		return
	}
}

func EndTrip(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	curDriver, err := GetDriverFromCookie(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(&ErrorParsingCookie)
		return
	}

	trip, err := redis.GetTrip("trip-" + strconv.Itoa(int(curDriver.DriverId)))
	if err != nil {
		return
	}

	if trip.PaymentMethod == "Wallet" {
		err = models.WalletTransactions(trip.UserId, trip.DriverId, trip.Fare)
		if err != nil {
			return
		}
	} else {
		err = models.CashTransactions(trip.DriverId, trip.Fare)
		if err != nil {
			return
		}
	}

	response := &models.Response{
		ResponseStatus:  "success",
		ResponseMessage: "trip completed",
		ResponseData:    nil,
	}
	//go user.TripEnd()
	if err = json.NewEncoder(w).Encode(&response); err != nil {
		return
	}
}

func DriverTripHistory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	d, err := GetDriverFromCookie(r)
	if err != nil {
		json.NewEncoder(w).Encode(&ErrorParsingCookie)
		return
	}

	tripHistory := trip.GetTripHistory("driver_id", d.DriverId)

	response := &models.Response{
		ResponseStatus:  "success",
		ResponseMessage: "fetched trip history",
		ResponseData:    tripHistory,
	}
	err = json.NewEncoder(w).Encode(&response)
	if err != nil {
		return
	}
}

func DriverWallet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	driver, err := GetDriverFromCookie(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(&ErrorParsingCookie)
		return
	}

	response := &models.Response{
		ResponseStatus:  "success",
		ResponseMessage: "driver wallet balance",
		ResponseData:    driver.WalletBalance,
	}
	json.NewEncoder(w).Encode(&response)
}

type payout struct {
	Amount string `json:"amount"`
}

func PayoutWallet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	d, err := GetDriverFromCookie(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(&ErrorParsingCookie)
		return
	}

	if d.BankAccount.AccountNumber == "" {
		w.WriteHeader(http.StatusBadRequest)
		resp := &models.Response{
			ResponseStatus:  "failed",
			ResponseMessage: "bank account not added",
			ResponseData:    nil,
		}
		json.NewEncoder(w).Encode(&resp)
		return
	}

	amount := &payout{}
	json.NewDecoder(r.Body).Decode(&amount)
	r.Body.Close()

	payoutAmount, _ := strconv.Atoi(amount.Amount)

	if uint(payoutAmount) > d.WalletBalance {
		w.WriteHeader(http.StatusBadRequest)
		resp := &models.Response{
			ResponseStatus:  "failed",
			ResponseMessage: "invalid payout amount",
			ResponseData:    fmt.Sprintf("availiable balance:%v", d.WalletBalance),
		}
		json.NewEncoder(w).Encode(&resp)
		return
	}
	if err = driver.Payout(amount.Amount, d.DriverId); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp := &models.Response{
			ResponseStatus:  "failed",
			ResponseMessage: "checkout request failed",
			ResponseData:    err.Error(),
		}
		json.NewEncoder(w).Encode(resp)
		return
	}
	resp := &models.Response{
		ResponseStatus:  "success",
		ResponseMessage: "checkout request added",
		ResponseData:    nil,
	}
	json.NewEncoder(w).Encode(&resp)
}

//PayoutStatus returns the latest payout request
func PayoutStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	d, err := GetDriverFromCookie(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(&ErrorParsingCookie)
		return
	}

	payouts, err := driver.PayoutRequests(d.DriverId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp := &models.Response{
			ResponseStatus:  "failed",
			ResponseMessage: "payout request fetching failed",
			ResponseData:    err,
		}
		json.NewEncoder(w).Encode(&resp)
		return
	}
	if len(payouts) == 0 {
		resp := &models.Response{
			ResponseStatus:  "success",
			ResponseMessage: "payout request status",
			ResponseData:    "no pending payout requests",
		}
		json.NewEncoder(w).Encode(&resp)
		return
	}
	resp := &models.Response{
		ResponseStatus:  "success",
		ResponseMessage: "payout request status",
		ResponseData:    &payouts,
	}
	json.NewEncoder(w).Encode(&resp)
}

//PayoutHistory returns the history of payouts requested by the driver
func PayoutHistory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	d, err := GetDriverFromCookie(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(&ErrorParsingCookie)
		return
	}
	payouts := driver.PayoutHistory(d.DriverId)
	resp := &models.Response{
		ResponseStatus:  "success",
		ResponseMessage: "payout request status",
		ResponseData:    payouts,
	}
	json.NewEncoder(w).Encode(&resp)
}

// GetDriverFromCookie returns the logged-in user from the stored cookie in session
func GetDriverFromCookie(r *http.Request) (*models.Driver, error) {
	c, err := r.Cookie("jwt-token")

	if err != nil {
		return &models.Driver{}, err
	}

	tokenString := c.Value

	_, phone := auth.ParseJWT(tokenString)

	d := driver.GetDriver("phone_number", phone)

	return d, nil
}

var ErrorParsingCookie = &models.Response{
	ResponseStatus:  "failed",
	ResponseMessage: "error parsing cookie",
	ResponseData:    "login to continue",
}
