package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

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
		userTemp.ExecuteTemplate(w, "validateOtp.html", data)
		go auth.SetOtp(phonenumber)
		// userTemp.ExecuteTemplate(w, "userSignupForm.html", data)
	}
	// if err := ; err != nil {
	// 	fmt.Println(err)
	// }
}

func ValidateOtp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache,no-store,must-revalidate")
	otp := r.FormValue("otp")

	phone := r.FormValue("phone")
	phonedata := map[string]string{
		"phone": phone,
	}
	if err := auth.ValidateOTP(phone, otp); err != nil {
		if err == redis.Nil {
			data := map[string]string{
				"phone": phone,
				"err":   "otp expired",
			}
			userTemp.ExecuteTemplate(w, "validateOtp.html", data)
			return
		} else {
			data := map[string]string{
				"phone": phone,
				"err":   "invalid otp",
			}
			userTemp.ExecuteTemplate(w, "validateOtp.html", data)
			return
		}

	}
	// http.Redirect(w,r,"/user/signup",)
	userTemp.ExecuteTemplate(w, "userSignupForm.html", phonedata)

}

func UserHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache,no-store,must-revalidate")
	c, err := r.Cookie("jwt-token")
	if err == http.ErrNoCookie {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	tokenstring := c.Value
	phone, errr := auth.ValidateJWT(tokenstring)
	fmt.Println("phone from jwt", phone)
	if errr != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	user := user.GetUser("phone_number", phone)
	data := map[string]any{
		"userid":    user.ID,
		"firstname": user.FirstName,
		"lastname":  user.LastName,
		"email":     user.Email,
	}
	userTemp.ExecuteTemplate(w, "userhome.html", data)

}

//Create a user model with values from the fronted.
//Pass the newly created user model to user services
//to insert the new user to the database.
//Login the user and open user home after successful signup.
func UserSignUp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache,no-store,must-revalidate")

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

func EditUserProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache,no-store,must-revalidate")
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
	userTemp.ExecuteTemplate(w, "editUserProfile.html", data)
}
func UpdateUserProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache,no-store,must-revalidate")
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
	err := user.UpdateUser(&newuser)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("user updated successfully")
	w.Write([]byte("user updated successfully"))

}

func UserLogout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache,no-store,must-revalidate")
	http.SetCookie(w, &http.Cookie{
		Name:    "jwt-token",
		Value:   "",
		Path:    "/",
		Domain:  "localhost:8080",
		Expires: time.Time{},
		MaxAge:  -1,
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	users := user.GetUsers()
	for _, user := range users {
		fmt.Println(user.FirstName)
	}
}

//return true if entered password is matching with
//the hash password stored in the database
func validPassword(password, hashPassword string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password)); err != nil {
		return err
	}
	return nil
}
