package controllers

import (
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/shayamvlmna/cab-booking-app/pkg/database/redis"
	"github.com/shayamvlmna/cab-booking-app/pkg/models"
	"github.com/shayamvlmna/cab-booking-app/pkg/service/auth"
	"github.com/shayamvlmna/cab-booking-app/pkg/service/trip"
	"github.com/shayamvlmna/cab-booking-app/pkg/service/user"
)

// UserAuth Check if the user already exist in the system.
//Redirect to the user login page if user exists.
//Redirect to the user signup page if user is new.
func UserAuth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache,no-store,must-revalidate")
	w.Header().Set("Content-Type", "application/json")

	newUser := &models.User{}

	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		return
	}

	phonenumber := newUser.Phonenumber

	if phonenumber != "" {
		auth.StorePhone(phonenumber)
		if user.IsUserExists("phonenumber", phonenumber) {
			http.Redirect(w, r, "/user/loginpage", http.StatusSeeOther)
			return
		} else {
			if err := auth.SetOtp(phonenumber); err != nil {
				fmt.Println(err)
				return
			}
			http.Redirect(w, r, "/user/enterotp", http.StatusSeeOther)
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

// UserSignupPage render the signup page to submit the details of the new user
func UserSignupPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := &models.Response{
		ResponseStatus:  "success",
		ResponseMessage: "submit user data",
		ResponseData:    nil,
	}
	err := json.NewEncoder(w).Encode(&response)
	if err != nil {
		return
	}
}

func UserLoginPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	phonenumber := auth.GetPhone()

	response := models.Response{
		ResponseStatus:  "success",
		ResponseMessage: "existing user",
		ResponseData:    user.GetUser("phonenumber", phonenumber).Firstname,
	}

	auth.StorePhone(phonenumber)

	err := json.NewEncoder(w).Encode(&response)
	if err != nil {
		return
	}
}

// UserSignUp Create a user model with values from the fronted.
//Pass the newly created user model to user services
//to insert the new user to the database.
//Login the user and open user home after successful signup.
func UserSignUp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache,no-store,must-revalidate")
	w.Header().Set("Content-Type", "application/json")

	newUser := models.User{}

	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(r.Body)

	if err := user.RegisterUser(&newUser); err != nil {
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
		ResponseMessage: "signup success",
		ResponseData:    nil,
	}
	err = json.NewEncoder(w).Encode(&response)
	if err != nil {
		return
	}
}

// UserLogin get the existing user by phone number from the database.
//Validate the entered password with stored hash password.
//Generate a JWT token for the user after successful login.
//Store the JWT token in the http only cookie
func UserLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache,no-store,must-revalidate")
	w.Header().Set("Content-Type", "application/json")

	newUser := &models.User{}

	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(r.Body)

	password := newUser.Password

	// phoneNumber := newUser.PhoneNumber

	phoneNumber := auth.GetPhone()

	//get the existing user by phone number from the database
	User := user.GetUser("phonenumber", phoneNumber)

	if !User.Active {
		response := &models.Response{
			ResponseStatus:  "failed",
			ResponseMessage: "user not active",
			ResponseData:    nil,
		}
		err := json.NewEncoder(w).Encode(&response)
		if err != nil {
			return
		}
		return
	}

	//validate the entered password with stored hash password
	if err := validPassword(password, User.Password); err != nil {
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

	token, err := auth.GenerateJWT("user", phoneNumber)
	if err != nil {
		fmt.Println("jwt failed", err)
	}

	err = redis.StoreData("data", User)
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

	// http.Redirect(w, r, "/user/userhome", http.StatusSeeOther)
}

// UserHome get the logged-in user data from redis
//if err get from the primary database
//render the user home page
func UserHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache,no-store,must-revalidate")
	w.Header().Set("Content-Type", "application/json")

	c, err := r.Cookie("jwt-token")
	if err != nil {
		response := &models.Response{
			ResponseStatus:  "failed",
			ResponseMessage: "no cookie",
			ResponseData:    nil,
		}
		if err := json.NewEncoder(w).Encode(&response); err != nil {
			return
		}
		return
	}

	tokenString := c.Value

	_, phone := auth.ParseJWT(tokenString)

	user := user.GetUser("phonenumber", phone)

	userData := &models.UserData{
		Id:          user.UserId,
		Phonenumber: user.Phonenumber,
		Firstname:   user.Firstname,
		Lastname:    user.Lastname,
		Email:       user.Email,
	}

	response := models.Response{
		ResponseStatus:  "success",
		ResponseMessage: "user data fetched",
		ResponseData:    userData,
	}

	err = json.NewEncoder(w).Encode(&response)
	if err != nil {
		return
	}
}

// UserLogout delete the stored user data from redis
//also expire the cookie stored
func UserLogout(w http.ResponseWriter, r *http.Request) {
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
		ResponseStatus:  "succes",
		ResponseMessage: "logout success",
		ResponseData:    nil,
	}

	err = json.NewEncoder(w).Encode(&response)
	if err != nil {
		return
	}
	// http.Redirect(w, r, "/", http.StatusSeeOther)
}

func EditUserProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache,no-store,must-revalidate")

	params := mux.Vars(r)
	id := params["id"]

	user := user.GetUser("id", id)

	err := json.NewEncoder(w).Encode(&user)
	if err != nil {
		return
	}
}

func UpdateUserProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache,no-store,must-revalidate")

	newUser := models.User{}

	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		return
	}

	err = user.UpdateUser(&newUser)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Println("user updated successfully")
	// w.Write([]byte("user updated successfully"))
	w.WriteHeader(http.StatusOK)
}

//match the entered password with
//the hash password stored in the database
func validPassword(password, hashPassword string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password)); err != nil {
		return err
	}
	return nil
}

// BookTrip get the pickup point and destination from the booktrip call from the user
func BookTrip(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	c, err := r.Cookie("jwt-token")
	if err != nil {
		response := &models.Response{
			ResponseStatus:  "failed",
			ResponseMessage: "no cookie",
			ResponseData:    nil,
		}
		if err := json.NewEncoder(w).Encode(&response); err != nil {
			return
		}
		return
	}
	_, phone := auth.ParseJWT(c.Value)

	user := user.GetUser("phonenumber", phone)

	newRide := &trip.Ride{}

	err = json.NewDecoder(r.Body).Decode(&newRide)
	if err != nil {
		return
	}

	newTrip := trip.CreateTrip(newRide)

	ride := &models.Ride{
		Source:        newTrip.Source,
		Destination:   newTrip.Destination,
		Distance:      newTrip.Distance,
		ETA:           newTrip.ETA,
		Fare:          newTrip.Fare,
		PaymentMethod: "",
	}

	if user.WalletBalance < newTrip.Fare {
		response := &models.Response{
			ResponseStatus:  "failed",
			ResponseMessage: "not enough balance in wallet",
			ResponseData:    &ride,
		}
		if err = json.NewEncoder(w).Encode(&response); err != nil {
			return
		}
		return
	}

	response := &models.Response{
		ResponseStatus:  "success",
		ResponseMessage: "trip created successfully",
		ResponseData:    &ride,
	}
	if err = json.NewEncoder(w).Encode(&response); err != nil {
		return
	}
}

var OTPchan = make(chan int)

//ConfirmTrip returns the trip code to match with the driver to start the ride
func ConfirmTrip(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	cnftrip := &models.Ride{}

	err := json.NewDecoder(r.Body).Decode(&cnftrip)
	if err != nil {
		return
	}

	c, err := r.Cookie("jwt-token")
	if err != nil {
		response := &models.Response{
			ResponseStatus:  "failed",
			ResponseMessage: "no cookie",
			ResponseData:    nil,
		}
		if err := json.NewEncoder(w).Encode(&response); err != nil {
			return
		}
		return
	}

	tokenString := c.Value

	_, phone := auth.ParseJWT(tokenString)

	curUser := user.GetUser("phone_number", phone)

	cnftrip.UserId = curUser.UserId
	go trip.FindCab(&cnftrip)

	//otp, err := auth.TripCode()
	//Tripcode, err := strconv.Atoi(otp)

	//redis.Set("tripcode"+strconv.Itoa(int(curUser.Id)), otp)

	err = json.NewEncoder(w).Encode(&models.Response{
		ResponseStatus:  "success",
		ResponseMessage: "waiting to accept ride",
		ResponseData:    nil,
	})
	if err != nil {
		return
	}
}

//TripHistory returns saved trips for the user
func UserTripHistory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	c, err := r.Cookie("jwt-token")
	if err != nil {
		response := &models.Response{
			ResponseStatus:  "failed",
			ResponseMessage: "no cookie",
			ResponseData:    nil,
		}
		if err := json.NewEncoder(w).Encode(&response); err != nil {
			return
		}
		return
	}

	tokenString := c.Value

	_, phone := auth.ParseJWT(tokenString)

	user := user.GetUser("phonenumber", phone)

	tripHistory := trip.GetTripHistory("user_id", user.UserId)

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
