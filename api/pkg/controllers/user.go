package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

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

	newUser := &models.User{}

	json.NewDecoder(r.Body).Decode(&newUser)

	phonenumber := newUser.PhoneNumber

	fmt.Println(phonenumber)

	if phonenumber != "" {
		auth.StoreUser(phonenumber)
		if user.IsUserExists("phone_number", phonenumber) {
			http.Redirect(w, r, "/user/loginpage", http.StatusSeeOther)
			return
		} else {
			go auth.SetOtp(string(phonenumber))
			http.Redirect(w, r, "/user/enterotp", http.StatusSeeOther)
		}
	} else {
		response := models.Response{
			ResponseStatus:  "fail",
			ResponseMessage: "phonenumber required",
			ResponseData:    nil,
		}
		json.NewEncoder(w).Encode(response)
		return
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

func UserSignupPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := models.Response{
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

	// token, err := auth.GenerateJWT("user", newUser.PhoneNumber)
	// if err != nil {
	// 	fmt.Println("jwt failed", err)
	// }

	response := models.Response{
		ResponseStatus:  "succes",
		ResponseMessage: "signup success",
		ResponseData:    nil,
	}
	json.NewEncoder(w).Encode(&response)
}

//get the existing user by phone number from the database.
//Validate the entered password with stored hash password.
//Generate a JWT token for the user after successful login.
//Store the JWT token in the cookie
func UserLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache,no-store,must-revalidate")
	w.Header().Set("Content-Type", "application/json")

	newUser := &models.User{}

	json.NewDecoder(r.Body).Decode(&newUser)

	password := newUser.Password

	// password := r.FormValue("password")

	// phoneNumber := auth.GetUser()

	phoneNumber := newUser.PhoneNumber

	//get the existing user by phone number from the database
	User := user.GetUser("phone_number", phoneNumber)

	//validate the entered password with stored hash password
	if err := validPassword(password, User.Password); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		response := models.Response{
			ResponseStatus:  "succes",
			ResponseMessage: "password authentication failed",
			ResponseData:    nil,
			Token:           "",
		}
		json.NewEncoder(w).Encode(&response)
		return
	}
	token, err := auth.GenerateJWT("user", phoneNumber)
	if err != nil {
		fmt.Println("jwt failed", err)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "jwt-token",
		Value:    token,
		Path:     "/",
		MaxAge:   0,
		HttpOnly: true,
	})
	// response := models.Response{
	// 	ResponseStatus:  "succes",
	// 	ResponseMessage: "user login success",
	// 	ResponseData:    nil,
	// 	Token:           token,
	// }
	// json.NewEncoder(w).Encode(&response)
	http.Redirect(w, r, "/user/userhome", http.StatusSeeOther)
}

func UserHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache,no-store,must-revalidate")
	w.Header().Set("Content-Type", "application/json")

	// token := r.Header["Token"][0]
	c, _ := r.Cookie("jwt-token")
	tokenString := c.Value
	role, phone := auth.ParseJWT(tokenString)

	fmt.Println(role, phone)

	user := user.GetUser("phone_number", phone)

	response := models.Response{
		ResponseStatus:  "success",
		ResponseMessage: "user data fetched",
		ResponseData:    user,
		Token:           tokenString,
	}
	json.NewEncoder(w).Encode(&response)
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
	w.Header().Set("Token", "")

	response := models.Response{
		ResponseStatus:  "succes",
		ResponseMessage: "logout success",
		ResponseData:    nil,
		Token:           "",
	}
	json.NewEncoder(w).Encode(&response)
	// http.Redirect(w, r, "/", http.StatusSeeOther)
}

//return true if entered password is matching with
//the hash password stored in the database
func validPassword(password, hashPassword string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password)); err != nil {
		return err
	}
	return nil
}
