package controllers

import (
	"fmt"

	"github.com/shayamvlmna/cab-booking-app/app/models"
	"github.com/shayamvlmna/cab-booking-app/app/service/user"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"net/http"
)

func UserAuth(w http.ResponseWriter, r *http.Request) {
	phonenumber := r.FormValue("usrphonenumber")
	if user.IsUserExists(phonenumber) {
		// UserLogin(w, r)
		http.Redirect(w, r, "/user/login", 200)
	} else {
		// UserSignUp(w, r)
		http.Redirect(w, r, "/user/signup", 200)

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
	// password := r.FormValue("usrpassword")
	// phonenumber := r.FormValue("usrphonenumber")

	// user, err := GetUser(phonenumber)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
	// 	fmt.Println(err)
	// }

	// err := user.GetUser()
}
