package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/gorilla/mux"

	redis "github.com/shayamvlmna/cab-booking-app/pkg/database/redis"
	models "github.com/shayamvlmna/cab-booking-app/pkg/models"
	auth "github.com/shayamvlmna/cab-booking-app/pkg/service/auth"
	trip "github.com/shayamvlmna/cab-booking-app/pkg/service/trip"
	user "github.com/shayamvlmna/cab-booking-app/pkg/service/user"
)

//Check if the user already exist in the system.
//Redirect to the user login page if user exists.
//Redirect to the user signup page if user is new.
func UserAuth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache,no-store,must-revalidate")
	w.Header().Set("Content-Type", "application/json")

	newUser := &models.User{}

	json.NewDecoder(r.Body).Decode(&newUser)

	phonenumber := newUser.Phonenumber

	fmt.Println(phonenumber)

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
		json.NewEncoder(w).Encode(&response)
		return
	}
}

//render the signup page to submit the details of the new user
func UserSignupPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := &models.Response{
		ResponseStatus:  "success",
		ResponseMessage: "submit user data",
		ResponseData:    nil,
	}
	json.NewEncoder(w).Encode(&response)
}

//Create a user model with values from the fronted.
//Pass the newly created user model to user services
//to insert the new user to the database.
//Login the user and open user home after successful signup.
func UserSignUp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache,no-store,must-revalidate")
	w.Header().Set("Content-Type", "application/json")

	newUser := models.User{}
	json.NewDecoder(r.Body).Decode(&newUser)
	defer r.Body.Close()

	newUser.Phonenumber = auth.GetPhone()
	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	newUser.Password = string(hashPassword)
	//create a user model with values from the fronted

	//pass the newly created user model to user services
	//to insert the new user to the database
	//after successful signup login the user and open user home
	if err := user.AddUser(&newUser); err != nil {
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
		ResponseStatus:  "succes",
		ResponseMessage: "signup success",
		ResponseData:    nil,
	}
	json.NewEncoder(w).Encode(&response)
}

//get the existing user by phone number from the database.
//Validate the entered password with stored hash password.
//Generate a JWT token for the user after successful login.
//Store the JWT token in the http only cookie
func UserLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache,no-store,must-revalidate")
	w.Header().Set("Content-Type", "application/json")

	newUser := &models.User{}

	json.NewDecoder(r.Body).Decode(&newUser)
	defer r.Body.Close()

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
		json.NewEncoder(w).Encode(&response)
		return
	}

	//validate the entered password with stored hash password
	if err := validPassword(password, User.Password); err != nil {
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
	token, err := auth.GenerateJWT("user", phoneNumber)
	if err != nil {
		fmt.Println("jwt failed", err)
	}

	redis.StoreData("data", User)

	http.SetCookie(w, &http.Cookie{
		Name:     "jwt-token",
		Value:    token,
		Path:     "/",
		MaxAge:   0,
		HttpOnly: true,
	})
	http.Redirect(w, r, "/user/userhome", http.StatusSeeOther)
}

//get the logged in user data from redis
//if err get from the primary database
//render the user home page
func UserHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache,no-store,must-revalidate")
	w.Header().Set("Content-Type", "application/json")

	c, _ := r.Cookie("jwt-token")
	tokenString := c.Value

	_, phone := auth.ParseJWT(tokenString)

	user := user.GetUser("phonenumber", phone)
	response := models.Response{
		ResponseStatus:  "success",
		ResponseMessage: "user data fetched",
		ResponseData:    user,
	}
	json.NewEncoder(w).Encode(&response)
}

//delete the stored user data from redis
//also expire the cookie stored
func UserLogout(w http.ResponseWriter, r *http.Request) {
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
		ResponseStatus:  "succes",
		ResponseMessage: "logout success",
		ResponseData:    nil,
	}
	json.NewEncoder(w).Encode(&response)
	// http.Redirect(w, r, "/", http.StatusSeeOther)
}

func EditUserProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache,no-store,must-revalidate")
	params := mux.Vars(r)
	id := params["id"]
	user := user.GetUser("id", id)

	json.NewEncoder(w).Encode(&user)
}

func UpdateUserProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache,no-store,must-revalidate")

	newUser := models.User{}
	json.NewDecoder(r.Body).Decode(&newUser)

	err := user.UpdateUser(&newUser)
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

//get the pickup point and destination from the booktrip call from the user
func BookTrip(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	newRide := &trip.Ride{}

	json.NewDecoder(r.Body).Decode(&newRide)

	trip := trip.CreateTrip(newRide)

	ride := &models.Trip{
		Source:      trip.Source,
		Destination: trip.Destination,
		Fare:        trip.Fare,
		ETA:         trip.ETA,
	}
	response := &models.Response{
		ResponseStatus:  "success",
		ResponseMessage: "trip created successfully",
		ResponseData:    ride,
	}
	json.NewEncoder(w).Encode(&response)
}

func TripHistory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	c, _ := r.Cookie("jwt-token")
	tokenString := c.Value

	_, phone := auth.ParseJWT(tokenString)

	user := user.GetUser("phonenumber", phone)

	tripHistory := trip.GetTripHistory(user.Id)

	response := &models.Response{
		ResponseStatus:  "success",
		ResponseMessage: "fetched trip history",
		ResponseData:    tripHistory,
	}
	json.NewEncoder(w).Encode(&response)
}

// func Test(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	result := mapservice.DistanceAPI()
// 	distance := result.Rows[0].Element[0].Distance.Text
// 	fmt.Println(distance)
// 	json.NewEncoder(w).Encode(&result)
// }

func ConfirmTrip(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// cnftrip := &models.Trip{}

	// c, _ := r.Cookie("jwt-token")
	// tokenString := c.Value

	// _, phone := auth.ParseJWT(tokenString)

	// curUser := user.GetUser("phone_number", phone)

	// // if err := user.AppendTrip(&curUser, cnftrip); err != nil {
	// // 	fmt.Println(err)
	// // 	return
	// // }

	// trip.FindCab(&cnftrip)

	// json.NewDecoder(r.Body).Decode(&cnftrip)

}
