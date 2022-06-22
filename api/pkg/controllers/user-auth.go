package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/go-redis/redis/v9"
	"github.com/gorilla/mux"

	models "github.com/shayamvlmna/cab-booking-app/pkg/models"
	auth "github.com/shayamvlmna/cab-booking-app/pkg/service/auth"
	user "github.com/shayamvlmna/cab-booking-app/pkg/service/user"
)

//Check if the user already exist in the system.
//Redirect to the user login page if user exists.
//Redirect to the user signup page if user is new.
func UserAuth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache,no-store,must-revalidate")
	w.Header().Set("Content-Type", "application/json")
	phonenumber := r.FormValue("phonenumber")
	auth.StoreUser(phonenumber)

	if user.IsUserExists("phone_number", phonenumber) {
		http.Redirect(w, r, "/user/loginpage", http.StatusSeeOther)
		return
	} else {
		go auth.SetOtp(phonenumber)
		http.Redirect(w, r, "/user/enterotp", http.StatusSeeOther)
	}
}

func ValidateOtp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache,no-store,must-revalidate")
	w.Header().Set("Content-Type", "application/json")
	otp := r.FormValue("otp")

	phone := auth.GetUser()

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
	http.Redirect(w, r, "/user/signup", http.StatusSeeOther)
}

func UserHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache,no-store,must-revalidate")
	w.Header().Set("Content-Type", "application/json")

	_, phone, err := auth.ValidateJWT(r)

	if err != nil {
		if err == errors.New("invalidToken") {
			http.Redirect(w, r, "/", http.StatusUnauthorized)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	user := user.GetUser("phone_number", phone)

	response := models.Response{
		ResponseStatus:  "success",
		ResponseMessage: "user data fetched",
		ResponseData:    user,
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
	newUser.PhoneNumber = auth.GetUser()
	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	newUser.Password = string(hashPassword)
	//create a user model with values from the fronted

	//pass the newly created user model to user services
	//to insert the new user to the database
	//after successful signup login the user and open user home
	if err := user.AddUser(&newUser); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Println("user added")
	setCookie(w, "user", newUser.PhoneNumber)

	http.Redirect(w, r, "/user/userhome", http.StatusSeeOther)
}

//get the existing user by phone number from the database.
//Validate the entered password with stored hash password.
//Generate a JWT token for the user after successful login.
//Store the JWT token in the cookie
func UserLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache,no-store,must-revalidate")
	w.Header().Set("Content-Type", "application/json")

	password := r.FormValue("password")

	phoneNumber := auth.GetUser()

	//get the existing user by phone number from the database
	User := user.GetUser("phone_number", phoneNumber)

	//validate the entered password with stored hash password
	if err := validPassword(password, User.Password); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if err := setCookie(w, "user", User.PhoneNumber); err != nil {
		fmt.Println(err)
	}
	//after successful login, generate a JWT token for the user
	//save the generated token in the cookie

	http.Redirect(w, r, "/user/userhome", http.StatusSeeOther)
}

func setCookie(w http.ResponseWriter, role, key string) error {
	jwt, err := auth.GenerateJWT(role, key)
	if err != nil {
		return err
	}
	http.SetCookie(w, &http.Cookie{
		Name:   "jwt-token",
		Value:  jwt,
		Path:   "/",
		MaxAge: 0,
	})
	return nil
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

func UserLogout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache,no-store,must-revalidate")

	c, _ := r.Cookie("jwt-token")

	http.SetCookie(w, &http.Cookie{
		Name:    "jwt-token",
		Value:   "",
		Path:    "/",
		MaxAge:  -1,
		Expires: time.Now().Add(-100 * time.Hour),
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	users := user.GetUsers()
	for _, user := range users {
		fmt.Println(user.FirstName)
	}
	strings.Split()
}

//return true if entered password is matching with
//the hash password stored in the database
func validPassword(password, hashPassword string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password)); err != nil {
		return err
	}
	return nil
}
