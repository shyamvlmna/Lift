package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	database "github.com/shayamvlmna/cab-booking-app/pkg/database/postgresql"
	"github.com/shayamvlmna/cab-booking-app/pkg/models"
	"github.com/shayamvlmna/cab-booking-app/pkg/service/driver"
	user "github.com/shayamvlmna/cab-booking-app/pkg/service/user"
	"golang.org/x/crypto/bcrypt"
)

func CreateAdmin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	admin := models.Admin{}
	json.NewDecoder(r.Body).Decode(&admin)
	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(admin.Password), bcrypt.DefaultCost)
	admin.Password = string(hashPassword)
	database.AddAdmin(&admin)

	json.NewEncoder(w).Encode(&models.Response{
		ResponseStatus:  "success",
		ResponseMessage: "created admin",
		ResponseData:    nil,
	})

}

func AdminIndex(w http.ResponseWriter, r *http.Request) {
}

func AdminLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	admin := models.Admin{}
	json.NewDecoder(r.Body).Decode(&admin)

	Admin, _ := database.GetAdmin(admin.Username)

	if err := validPassword(admin.Password, Admin.Password); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(&models.Response{
			ResponseStatus:  "fail",
			ResponseMessage: "password authentication failed",
			ResponseData:    nil,
		})
		return
	}

	json.NewEncoder(w).Encode(&models.Response{
		ResponseStatus:  "success",
		ResponseMessage: "admin login success",
		ResponseData:    Admin,
	})
}

type Data struct {
	Id uint64 `json:"id"`
}

func Managedrivers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	drivers := []models.Driver{}
	drivers = driver.GetDrivers()
	json.NewEncoder(w).Encode(&models.Response{
		ResponseStatus:  "success",
		ResponseMessage: "fetched drivers data",
		ResponseData:    &drivers,
	})
}

func ManageUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	users := user.GetUsers()
	json.NewEncoder(w).Encode(&models.Response{
		ResponseStatus:  "success",
		ResponseMessage: "fetched users data",
		ResponseData:    users,
	})
}

func DriveRequest(w http.ResponseWriter, r *http.Request) {

}

func ApproveDriver(w http.ResponseWriter, r *http.Request) {
	data := &Data{}

	json.NewDecoder(r.Body).Decode(&data)

	id := data.Id

	database.ApproveDriver(id)
}

func BlockDriver(w http.ResponseWriter, r *http.Request) {
	data := &Data{}
	id := data.Id
	database.ApproveDriver(id)
}

func BlockUser(w http.ResponseWriter, r *http.Request) {

}
