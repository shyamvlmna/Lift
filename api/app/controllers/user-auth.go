package controllers

import (
	"fmt"
	"strconv"
	"time"

	// "github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
	"github.com/shayamvlmna/cab-booking-app/app/models"
	"github.com/shayamvlmna/cab-booking-app/app/service/auth"
	"github.com/shayamvlmna/cab-booking-app/app/service/user"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"net/http"
)

func UserAuth(w http.ResponseWriter, r *http.Request) {
	phonenumber := r.FormValue("usrphonenumber")

	//check if the user already exist in the system
	//if user exists redirect to the user login page
	//if user is new redirect to the user signup page
	data := map[string]string{
		"phone": phonenumber,
	}
	if user.IsUserExists("usrphonenumber", phonenumber) {
		UserTemp.ExecuteTemplate(w, "userLoginForm.html", data)
		return
	}
	UserTemp.ExecuteTemplate(w, "userSignupForm.html", data)
}

func UserHome(w http.ResponseWriter, r *http.Request) {

}

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
		UserTemp.ExecuteTemplate(w, "userSignupForm.html", data)
	} else {
		fmt.Println("user added")
		UserLogin(w, r)
	}
}

func UserLogin(w http.ResponseWriter, r *http.Request) {

	password := r.FormValue("usrpassword")
	phonenumber := r.FormValue("usrphonenumber")

	//get the existing user by phone number from the database
	user := user.GetUser("usrphonenumber", phonenumber)

	//validate the entered password with stored hash password
	if err := validPassword(password, user.Password); err != nil {
		fmt.Println(err)
		data := map[any]any{
			"err": "invalid password",
		}
		UserTemp.ExecuteTemplate(w, "userLoginForm.html", data)
		return
	}

	//after successful login, generate a JWT token for the user
	//save the generated token in the cookie
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
	data := map[string]any{
		"userid":    user.ID,
		"firstname": user.FirstName,
		"lastname":  user.LastName,
		"email":     user.Email,
	}
	UserTemp.ExecuteTemplate(w, "userhome.html", data)
}
func EditUserProfile(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	user := user.GetUser("id", id)
	data := map[any]any{
		"userid":    id,
		"firstname": user.FirstName,
		"lastname":  user.LastName,
		"email":     user.Email,
	}
	fmt.Println(user.Email)
	w.Header().Add("id", id)
	UserTemp.ExecuteTemplate(w, "editUserProfile.html", data)
}
func UpdateUserProfile(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])
	fmt.Println("id from head", id)
	firstname := r.FormValue("usrfirstname")
	lastname := r.FormValue("usrlastname")
	email := r.FormValue("usremail")

	newuser := models.User{
		Model: gorm.Model{
			ID: uint(id),
		},
		FirstName:   firstname,
		LastName:    lastname,
		PhoneNumber: "",
		Email:       email,
		Password:    "",
	}
	user.UpdateUser(&newuser)
	// c, err := r.Cookie("jwt-token")
	// if err != nil {
	// 	if err == http.ErrNoCookie {
	// 		fmt.Println(err)
	// 		return
	// 	}
	// 	fmt.Println(err)
	// 	return
	// }
	// validateUser(c.Value)
	// // usrphone, err := auth.ValidateJWT(c.Value)
	// if err != nil {
	// 	if err == jwt.ErrSignatureInvalid {
	// 		fmt.Println(err)
	// 	}
	// 	fmt.Println(err)
	// 	return
	// }
	// user := user.GetUser(usrphone)

	// if authorized {
	// 	fmt.Println("valid user")
	// }

}

func validateUser(tknstr string) {
	// user, err := auth.ValidateJWT(tknstr)
	// if err != nil {

	// }

}

func validPassword(password, hashPassword string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password)); err != nil {
		return err
	}
	return nil
}
