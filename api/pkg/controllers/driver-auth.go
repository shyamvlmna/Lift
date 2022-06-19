package controllers

import (
	"fmt"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	"gorm.io/gorm"

	models "github.com/shayamvlmna/cab-booking-app/pkg/models"
	auth "github.com/shayamvlmna/cab-booking-app/pkg/service/auth"
	driver "github.com/shayamvlmna/cab-booking-app/pkg/service/driver"
)

//Check if the user already exist in the system.
//Redirect to the user login page if user exists.
//Redirect to the user signup page if user is new.
func DriverAuth(w http.ResponseWriter, r *http.Request) {
	phonenumber := r.FormValue("drvrphonenumber")

	data := map[string]string{
		"phone": phonenumber,
	}
	if driver.IsDriverExists("phone_number", phonenumber) {
		driverTemp.ExecuteTemplate(w, "driverLoginForm", data)
		return
	} else {
		driverTemp.ExecuteTemplate(w, "driverSignupForm", data)
	}
}

func DriverSignUpPage(w http.ResponseWriter, r *http.Request) {
	driverTemp.ExecuteTemplate(w, "driverSignupForm.html", nil)
}
func DriverLoginPage(w http.ResponseWriter, r *http.Request) {
	driverTemp.ExecuteTemplate(w, "driverLoginForm.html", nil)
}

//Create a user model with values from the fronted.
//Pass the newly created user model to user services
//to insert the new user to the database.
//Login the user and open user home after successful signup.
func DriverSignUp(w http.ResponseWriter, r *http.Request) {

	firstname := r.FormValue("drvrfirstname")
	lastname := r.FormValue("drvrlastname")
	phonenumber := r.FormValue("drvrphonenumber")
	email := r.FormValue("drvremail")
	password := r.FormValue("drvrpassword")

	hashpass, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	//create a user model with values from the fronted
	newDriver := models.Driver{
		Model:       gorm.Model{},
		FirstName:   firstname,
		LastName:    lastname,
		PhoneNumber: phonenumber,
		Email:       email,
		Password:    string(hashpass),
	}

	//pass the newly created user model to user services
	//to insert the new user to the database
	//after successful signup login the user and open user home
	if err := driver.AddDriver(&newDriver); err != nil {
		fmt.Println(err)
		data := map[any]any{
			"err": err,
		}
		userTemp.ExecuteTemplate(w, "driverSignupForm.html", data)
	} else {
		fmt.Println("driver added")
		DriverLogin(w, r)
	}
}

//get the existing user by phone number from the database.
//Validate the entered password with stored hash password.
//Generate a JWT token for the user after successful login.
//Store the JWT token in the cookie
func DriverLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache,no-store,must-revalidate")
	password := r.FormValue("drvrpassword")
	phonenumber := r.FormValue("drvrphonenumber")

	//get the existing user by phone number from the database
	driver := driver.GetDriver("phone_number", phonenumber)

	//validate the entered password with stored hash password
	if err := validPassword(password, driver.Password); err != nil {
		fmt.Println(err)
		data := map[any]any{
			"err": "invalid password",
		}
		driverTemp.ExecuteTemplate(w, "driverLoginForm.html", data)
		return
	}

	//after successful login, generate a JWT token for the user
	//save the generated token in the cookie
	jwt, err := auth.GenerateJWT(driver.PhoneNumber)
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
		"userid":    driver.ID,
		"firstname": driver.FirstName,
		"lastname":  driver.LastName,
		"email":     driver.Email,
	}
	driverTemp.ExecuteTemplate(w, "driverhome.html", data)

}
func EditDriverProfile(w http.ResponseWriter, r *http.Request) {

}

func UpdateDriverProfile(w http.ResponseWriter, r *http.Request) {

}

func GetDrivers(w http.ResponseWriter, r *http.Request) {

}
