package user

import (
	"github.com/shayamvlmna/cab-booking-app/app/database"
	"github.com/shayamvlmna/cab-booking-app/app/models"
)

func AddUser(newUser *models.User) error {

	database.InsertUser(newUser)
	return nil
}

func GetUser() error {
	database.FindUser()
	return nil
}
