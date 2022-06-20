package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	database "github.com/shayamvlmna/cab-booking-app/pkg/database/postgresql"
	"github.com/shayamvlmna/cab-booking-app/pkg/models"
	"github.com/shayamvlmna/cab-booking-app/pkg/service/driver"
	user "github.com/shayamvlmna/cab-booking-app/pkg/service/user"
	"golang.org/x/crypto/bcrypt"
)

func CreateAdmin(w http.ResponseWriter, r *http.Request) {
	admin := models.Admin{}
	json.NewDecoder(r.Body).Decode(&admin)
	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(admin.Password), bcrypt.DefaultCost)
	admin.Password = string(hashPassword)
	fmt.Println(admin.Username)
	database.AddAdmin(&admin)
}

func AdminIndex(w http.ResponseWriter, r *http.Request) {
	adminTemp.ExecuteTemplate(w, "adminLoginForm.html", nil)
}

func AdminLogin(w http.ResponseWriter, r *http.Request) {
	admin := models.Admin{}
	json.NewDecoder(r.Body).Decode(&admin)

	Admin, _ := database.GetAdmin(admin.Username)

	if err := validPassword(admin.Password, Admin.Password); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
	} else {
		json.NewEncoder(w).Encode(Admin)
	}
}

func Managedrivers(w http.ResponseWriter, r *http.Request) {
	drivers := []models.Driver{}
	drivers = driver.GetDrivers()
	json.NewEncoder(w).Encode(&drivers)

}
func ManageUsers(w http.ResponseWriter, r *http.Request) {

	users := []models.User{}
	users = user.GetUsers()
	json.NewEncoder(w).Encode(&users)
}
func DriveRequest(w http.ResponseWriter, r *http.Request) {
	adminTemp.ExecuteTemplate(w, "driverRequests.html", nil)
}
func ApproveDriver(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	iid, _ := strconv.Atoi(id)
	database.ApproveDriver(iid)
}
func BlockDriver(w http.ResponseWriter, r *http.Request) {

}
func BlockUser(w http.ResponseWriter, r *http.Request) {

}
