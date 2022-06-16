package controllers

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/shayamvlmna/cab-booking-app/app/models"
	"github.com/shayamvlmna/cab-booking-app/app/service/auth"
	"github.com/shayamvlmna/cab-booking-app/app/service/user"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"net/http"
)

func UserAuth(w http.ResponseWriter, r *http.Request) {
	phonenumber := r.FormValue("usrphonenumber")
	if user.IsUserExists(phonenumber) {

		UserLogin(w, r) //get the enter user password page

		// http.Redirect(w, r, "/user/login", http.StatusSeeOther)

	} else {

		UserSignUp(w, r) //get the user signup page

		// http.Redirect(w, r, "/user/signup", http.StatusSeeOther)

	}
}
func UserSignUp(w http.ResponseWriter, r *http.Request) {

	firstname := r.FormValue("usrfirstname")
	lastname := r.FormValue("usrlastname")
	phonenumber := r.FormValue("usrphonenumber")
	email := r.FormValue("usremail")
	password := r.FormValue("usrpassword")

	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	newUser := models.User{
		Model:       gorm.Model{},
		FirstName:   firstname,
		LastName:    lastname,
		PhoneNumber: phonenumber,
		Email:       email,
		Password:    string(hashPassword),
	}

	if err := user.AddUser(&newUser); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("user added")
	}

}
func UserLogin(w http.ResponseWriter, r *http.Request) {

	// email := r.FormValue("usremail")
	password := r.FormValue("usrpassword")
	phonenumber := r.FormValue("usrphonenumber")

	user := user.GetUser(phonenumber)

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		fmt.Println(err)
	} else {
		fmt.Print("user login success")
	}

	jwt, err := auth.GenerateJWT(user.Email)
	if err != nil {
		fmt.Println(err)
	}
	http.SetCookie(w, &http.Cookie{
		Name:    "jwt-token",
		Value:   jwt,
		Path:    "/user",
		Expires: time.Now().Add(time.Minute * 30),
	})

}

func UpdateUserProfile(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("jwt-token")
	if err != nil {
		if err == http.ErrNoCookie {
			fmt.Println(err)
			return
		}
		fmt.Println(err)
		return
	}

	isAuthorized, err := auth.IsAuthorized(c.Value)
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			fmt.Println(err)
		}
		fmt.Println(err)
		return
	}
	if isAuthorized {
		fmt.Println("valid user")
	}

}
