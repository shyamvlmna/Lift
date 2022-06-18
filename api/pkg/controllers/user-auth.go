package controllers

import (
	"fmt"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	models "github.com/shayamvlmna/cab-booking-app/pkg/models"
	auth "github.com/shayamvlmna/cab-booking-app/pkg/service/auth"
	user "github.com/shayamvlmna/cab-booking-app/pkg/service/user"
)

//Check if the user already exist in the system.
//Redirect to the user login page if user exists.
//Redirect to the user signup page if user is new.
func UserAuth(w http.ResponseWriter, r *http.Request) {
	phonenumber := r.FormValue("usrphonenumber")

	//check if the user already exist in the system
	//if user exists redirect to the user login page
	//if user is new redirect to the user signup page
	data := map[string]string{
		"phone": phonenumber,
	}
	if user.IsUserExists("phone_number", phonenumber) {
		userTemp.ExecuteTemplate(w, "userLoginForm.html", data)
		return
	} else {
		userTemp.ExecuteTemplate(w, "userSignupForm.html", data)
	}
}

//Create a user model with values from the fronted.
//Pass the newly created user model to user services
//to insert the new user to the database.
//Login the user and open user home after successful signup.
func UserSignUp(w http.ResponseWriter, r *http.Request) {

	firstname := r.FormValue("usrfirstname")
	lastname := r.FormValue("usrlastname")
	phonenumber := r.FormValue("usrphonenumber")
	email := r.FormValue("usremail")
	password := r.FormValue("usrpassword")

	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	//create a user model with values from the fronted
	newUser := models.User{
		Model:       gorm.Model{},
		FirstName:   firstname,
		LastName:    lastname,
		PhoneNumber: phonenumber,
		Email:       email,
		Password:    string(hashPassword),
	}

	//pass the newly created user model to user services
	//to insert the new user to the database
	//after successful signup login the user and open user home
	if err := user.AddUser(&newUser); err != nil {
		fmt.Println(err)
		data := map[any]any{
			"err": err,
		}
		userTemp.ExecuteTemplate(w, "userSignupForm.html", data)
	} else {
		fmt.Println("user added")
		UserLogin(w, r)
	}
}

//get the existing user by phone number from the database.
//Validate the entered password with stored hash password.
//Generate a JWT token for the user after successful login.
//Store the JWT token in the cookie
func UserLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache,no-store,must-revalidate")
	phonenumber := r.FormValue("usrphonenumber")
	password := r.FormValue("usrpassword")

	//get the existing user by phone number from the database
	user := user.GetUser("phone_number", phonenumber)

	//validate the entered password with stored hash password
	if err := validPassword(password, user.Password); err != nil {
		fmt.Println(err)
		data := map[any]any{
			"err": "invalid password",
		}
		userTemp.ExecuteTemplate(w, "userLoginForm.html", data)
		return
	}

	//after successful login, generate a JWT token for the user
	//save the generated token in the cookie
	jwt, err := auth.GenerateJWT(user.PhoneNumber)
	if err != nil {
		fmt.Println(err)
	}
	http.SetCookie(w, &http.Cookie{
		Name:    "jwt-token",
		Value:   jwt,
		Path:    "/",
		Expires: time.Now().Add(time.Minute * 30),
	})
	data := map[string]any{
		"userid":    user.ID,
		"firstname": user.FirstName,
		"lastname":  user.LastName,
		"email":     user.Email,
	}
	userTemp.ExecuteTemplate(w, "userhome.html", data)
}

//return true if entered password is matching with
//the hash password stored in the database
func validPassword(password, hashPassword string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password)); err != nil {
		return err
	}
	return nil
}
