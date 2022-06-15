package controllers

import (
	// "fmt"
	"net/http"

	"github.com/shayamvlmna/cab-booking-app/pkg/models"
	"github.com/shayamvlmna/cab-booking-app/service/user"
	"gorm.io/gorm"
	// "golang.org/x/crypto/bcrypt"
)

func UserSignUp(w http.ResponseWriter, r *http.Request) {

	firstname := r.FormValue("usrfirstname")
	lastname := r.FormValue("usrlastname")
	phonenumber := r.FormValue("usrphonenumber")
	email := r.FormValue("usremail")
	password := r.FormValue("usrpassword")

	// hashPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	newUser := models.User{
		Model:       gorm.Model{},
		FirstName:   firstname,
		LastName:    lastname,
		PhoneNumber: phonenumber,
		Email:       email,
		Password:    password,
	}

	err :=user.AddUser(&newUser)
}
func UserLogin(w http.ResponseWriter, r *http.Request) {
	// phonenumber := r.FormValue("phonenumber")
	// password := r.FormValue("password")

	// user, err := GetUser(phonenumber)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
	// 	fmt.Println(err)
	// }

	err := user.GetUser()
}
