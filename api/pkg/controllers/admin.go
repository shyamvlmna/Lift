package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"

	database "github.com/shayamvlmna/cab-booking-app/pkg/database/postgresql"
	"github.com/shayamvlmna/cab-booking-app/pkg/models"
	"github.com/shayamvlmna/cab-booking-app/pkg/service/auth"
	"github.com/shayamvlmna/cab-booking-app/pkg/service/coupon"
	"github.com/shayamvlmna/cab-booking-app/pkg/service/driver"
	"github.com/shayamvlmna/cab-booking-app/pkg/service/user"
)

// AdminIndex render index page for admins to login
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

// CreateAdmin create a new admin by the super admin
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

func AdminLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	admin := models.Admin{}

	err := json.NewDecoder(r.Body).Decode(&admin)
	if err != nil {
		return
	}

	Admin, _ := database.GetAdmin(admin.Username)

	if !ValidPassword(admin.Password, Admin.Password) {
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

	token, err := auth.GenerateJWT("admin", Admin.Username)
	if err != nil {
		fmt.Println(err)
		return
	}
	if err = json.NewEncoder(w).Encode(&models.Response{
		ResponseStatus:  "success",
		ResponseMessage: "admin login success",
		ResponseData:    token,
	}); err != nil {
		return
	}
}

// AdminHome admin home page to manage users and drivers
func AdminHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := &models.Response{
		ResponseStatus:  "success",
		ResponseMessage: "admin home page",
		ResponseData:    nil,
	}
	json.NewEncoder(w).Encode(&response)
}

type DrvrData struct {
	Id uint `json:"driver_id"`
}

func CreateCoupon(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	c := &coupon.AmountCoupon{}
	json.NewDecoder(r.Body).Decode(&c)
	r.Body.Close()
	c.FinishDate = time.Now().AddDate(0, 0, 20)
	if err := c.CreateCoupon(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp := &models.Response{
			ResponseStatus:  "failed",
			ResponseMessage: "creating coupon failed",
			ResponseData:    c,
		}
		json.NewEncoder(w).Encode(resp)
		return
	}

	resp := &models.Response{
		ResponseStatus:  "success",
		ResponseMessage: "successfully created new coupon",
		ResponseData:    c,
	}
	json.NewEncoder(w).Encode(resp)
}

//ManageDrivers fetch the drivers details for the admin
func ManageDrivers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	drivers := []models.Driver{}

	drivers, err := driver.GetAllDrivers()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		json.NewEncoder(w).Encode(&models.Response{
			ResponseStatus:  "failed",
			ResponseMessage: "fetching drivers data failed",
			ResponseData:    nil,
		})
		return
	}
	err = json.NewEncoder(w).Encode(&models.Response{
		ResponseStatus:  "success",
		ResponseMessage: "fetched drivers data",
		ResponseData:    &drivers,
	})
	if err != nil {
		return
	}
}

//DriveRequest fetch the not approved drivers details
func DriveRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	requests, err := driver.DriverRequests()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp := &models.Response{
			ResponseStatus:  "failed",
			ResponseMessage: "fetching driver requests failed",
			ResponseData:    nil,
		}
		json.NewEncoder(w).Encode(resp)
		return
	}
	resp := &models.Response{
		ResponseStatus:  "success",
		ResponseMessage: "fetched driver requests",
		ResponseData:    requests,
	}
	json.NewEncoder(w).Encode(resp)

}

//ApproveDriver approves a driver to accept trips
func ApproveDriver(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	data := &DrvrData{}

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
	w.Header().Set("Content-Type", "application/json")
	data := &DrvrData{}
	id := data.Id
	if err := driver.BlockDriver(id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp := &models.Response{
			ResponseStatus:  "failed",
			ResponseMessage: "blocking driver failed",
			ResponseData:    nil,
		}
		json.NewEncoder(w).Encode(resp)
		return
	}
	resp := &models.Response{
		ResponseStatus:  "success",
		ResponseMessage: "successfully blocked driver",
		ResponseData:    nil,
	}
	json.NewEncoder(w).Encode(resp)

}

func UnBlockDriver(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	data := &DrvrData{}
	id := data.Id
	if err := driver.UnBlockDriver(id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp := &models.Response{
			ResponseStatus:  "failed",
			ResponseMessage: "unblock driver failed",
			ResponseData:    nil,
		}
		json.NewEncoder(w).Encode(resp)
		return
	}
	resp := &models.Response{
		ResponseStatus:  "success",
		ResponseMessage: "successfully unblocked driver",
		ResponseData:    nil,
	}
	json.NewEncoder(w).Encode(resp)
}

func PayoutRequests(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	payouts := models.GetPayouts()

	err := json.NewEncoder(w).Encode(&payouts)
	if err != nil {
		return
	}
}

type PayoutData struct {
	PayoutId string `json:"payoutId"`
	Status   string `json:"status"`
}

func UpdatePayout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	data := PayoutData{}
	json.NewDecoder(r.Body).Decode(&data)
	r.Body.Close()

	id, _ := strconv.Atoi(data.PayoutId)
	status := data.Status
	if err := models.UpdateCompletedPayoutRequest(uint(id), status); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp := &models.Response{
			ResponseStatus:  "failed",
			ResponseMessage: "payout request update failed",
			ResponseData:    err,
		}
		json.NewEncoder(w).Encode(resp)
		return
	}

	resp := &models.Response{
		ResponseStatus:  "success",
		ResponseMessage: "updated payout request",
		ResponseData:    nil,
	}
	json.NewEncoder(w).Encode(resp)
}

type UsrData struct {
	Id uint `json:"user_id"`
}

//ManageUsers fetch the users details for the admin
func ManageUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	users := &[]models.User{}

	users, err := user.GetUsers()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(&models.Response{
			ResponseStatus:  "failed",
			ResponseMessage: "fetching users data faiiled",
			ResponseData:    nil,
		})
		return
	}

	err = json.NewEncoder(w).Encode(&models.Response{
		ResponseStatus:  "success",
		ResponseMessage: "fetched users data",
		ResponseData:    users,
	})
	if err != nil {
		return
	}
}

func BlockUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	data := &UsrData{}
	json.NewDecoder(r.Body).Decode(&data)

	id := data.Id

	if err := user.BlockUser(id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp := &models.Response{
			ResponseStatus:  "failed",
			ResponseMessage: "blocking user failed",
			ResponseData:    nil,
		}
		json.NewEncoder(w).Encode(resp)
		return
	}
	resp := &models.Response{
		ResponseStatus:  "success",
		ResponseMessage: "successfully blocked user",
		ResponseData:    nil,
	}
	json.NewEncoder(w).Encode(resp)
}

func UnBlockUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	data := &UsrData{}
	json.NewDecoder(r.Body).Decode(&data)

	id := data.Id

	if err := user.UnBlockUser(id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp := &models.Response{
			ResponseStatus:  "failed",
			ResponseMessage: "unblocking user failed",
			ResponseData:    nil,
		}
		json.NewEncoder(w).Encode(resp)
		return
	}
	resp := &models.Response{
		ResponseStatus:  "success",
		ResponseMessage: "successfully unblocked user",
		ResponseData:    nil,
	}
	json.NewEncoder(w).Encode(resp)
}
