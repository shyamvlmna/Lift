package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/gorm"

	models "github.com/shayamvlmna/cab-booking-app/pkg/models"
	auth "github.com/shayamvlmna/cab-booking-app/pkg/service/auth"
	user "github.com/shayamvlmna/cab-booking-app/pkg/service/user"
)

func UserHome(w http.ResponseWriter, r *http.Request) {
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

func GetUsers(w http.ResponseWriter, r *http.Request) {
	users := user.GetUsers()
	for _, user := range users {
		fmt.Println(user.FirstName)
	}
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
	userTemp.ExecuteTemplate(w, "editUserProfile.html", data)
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
	err := user.UpdateUser(&newuser)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("user updated successfully")
	w.Write([]byte("user updated successfully"))

}
