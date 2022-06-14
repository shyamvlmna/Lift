package controllers

import (
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func UserSignUp(w http.ResponseWriter, r *http.Request) {

	firstname := r.FormValue("firstname")
	lastname := r.FormValue("lastname")
	
	email := r.FormValue("email")
	password := r.FormValue("password")

	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

}
func UserLogin(w http.ResponseWriter, r *http.Request) {
	phonenumber := r.FormValue("phonenumber")
	password := r.FormValue("password")

	user := GetUser(phonenumber)

	err := bcrypt.CompareHashAndPassword([]byte(user.password), []byte(password));err!=nil{
		fmt.Println("Err Mismatched Hash And Password")
	}
}
