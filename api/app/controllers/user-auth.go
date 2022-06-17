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
		data := map[string]string{
			"phone": phonenumber,
		}
		// UserLogin(w, r) //get the enter user password page
		UserTemp.ExecuteTemplate(w, "loginform.html", data)

		// http.Redirect(w, r, "/user/login", http.StatusSeeOther)

	} else {
		data := map[string]string{
			"phone": phonenumber,
		}
		UserTemp.ExecuteTemplate(w, "signupform.html", data)
		// UserSignUp(w, r) //get the user signup page

		// http.Redirect(w, r, "/user/signup", http.StatusSeeOther)

	}
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
		data := map[any]any{
			"err": err,
		}
		UserTemp.ExecuteTemplate(w, "signupform.html", data)
	} else {
		fmt.Println("user added")
		UserLogin(w, r)
	}

}

func UserLogin(w http.ResponseWriter, r *http.Request) {

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
	data := map[string]any{
		"userid":    user.ID,
		"firstname": user.FirstName,
		"lastname":  user.LastName,
		"email":     user.Email,
	}
	UserTemp.ExecuteTemplate(w, "userhome.html", data)
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
	validateUser(c.Value)
	// usrphone, err := auth.ValidateJWT(c.Value)
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			fmt.Println(err)
		}
		fmt.Println(err)
		return
	}
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
