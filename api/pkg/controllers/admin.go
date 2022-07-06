package controllers

import (
	"encoding/json"
	"net/http"

	database "github.com/shayamvlmna/cab-booking-app/pkg/database/postgresql"
	"github.com/shayamvlmna/cab-booking-app/pkg/models"
	"github.com/shayamvlmna/cab-booking-app/pkg/service/driver"
	"github.com/shayamvlmna/cab-booking-app/pkg/service/user"
	"golang.org/x/crypto/bcrypt"
)

func CreateAdmin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	admin := models.Admin{}

	err := json.NewDecoder(r.Body).Decode(&admin)
	if err != nil {
		return
	}

	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(admin.Password), bcrypt.DefaultCost)

	admin.Password = string(hashPassword)
	err = database.AddAdmin(&admin)
	if err != nil {
		return
	}

	err = json.NewEncoder(w).Encode(&models.Response{
		ResponseStatus:  "success",
		ResponseMessage: "created admin",
		ResponseData:    nil,
	})
	if err != nil {
		return
	}
}

func AdminIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := &models.Response{
		ResponseStatus:  "success",
		ResponseMessage: "admin index",
		ResponseData:    nil,
	}

	err := json.NewEncoder(w).Encode(&response)
	if err != nil {
		return
	}
}

func AdminLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	admin := models.Admin{}

	err := json.NewDecoder(r.Body).Decode(&admin)
	if err != nil {
		return
	}

	Admin, _ := database.GetAdmin(admin.Username)

	if err := validPassword(admin.Password, Admin.Password); err != nil {
		w.WriteHeader(http.StatusUnauthorized)

		err := json.NewEncoder(w).Encode(&models.Response{
			ResponseStatus:  "fail",
			ResponseMessage: "password authentication failed",
			ResponseData:    nil,
		})
		if err != nil {
			return
		}
		return
	}

	err = json.NewEncoder(w).Encode(&models.Response{
		ResponseStatus:  "success",
		ResponseMessage: "admin login success",
		ResponseData:    Admin,
	})
	if err != nil {
		return
	}
}

type Data struct {
	Id uint64 `json:"driver_id"`
}

func ManageDrivers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	drivers := []models.Driver{}

	drivers = driver.GetAllDrivers()

	err := json.NewEncoder(w).Encode(&models.Response{
		ResponseStatus:  "success",
		ResponseMessage: "fetched drivers data",
		ResponseData:    &drivers,
	})
	if err != nil {
		return
	}
}

func ManageUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	users := user.GetUsers()

	err := json.NewEncoder(w).Encode(&models.Response{
		ResponseStatus:  "success",
		ResponseMessage: "fetched users data",
		ResponseData:    users,
	})
	if err != nil {
		return
	}
}

func DriveRequest(w http.ResponseWriter, r *http.Request) {

}

func ApproveDriver(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	data := &Data{}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return
	}

	id := data.Id

	if err := driver.ApproveDriver(id); err != nil {
		response := &models.Response{
			ResponseStatus:  "failed",
			ResponseMessage: "error approving driver",
			ResponseData:    nil,
		}

		if err := json.NewEncoder(w).Encode(&response); err != nil {
			return
		}
		return
	}

	response := &models.Response{
		ResponseStatus:  "success",
		ResponseMessage: "driver approval success",
		ResponseData:    nil,
	}

	if err := json.NewEncoder(w).Encode(&response); err != nil {
		return
	}
}

func BlockDriver(w http.ResponseWriter, r *http.Request) {
	data := &Data{}
	id := data.Id
	if err := driver.BlockDriver(id); err != nil {
		return
	}

	//RESP
}

func UnBlockDriver(w http.ResponseWriter, r *http.Request) {
	data := &Data{}
	id := data.Id
	if err := driver.UnBlockDriver(id); err != nil {
		return
	}

	//RESP

}

func BlockUser(w http.ResponseWriter, r *http.Request) {

}
