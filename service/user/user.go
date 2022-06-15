package user

import (
	"github.com/shayamvlmna/cab-booking-app/pkg/database"
	"github.com/shayamvlmna/cab-booking-app/pkg/models"
)

func AddUser(newUser *models.User) error {

	database.AddUser(newUser)
	return nil
}

func GetUser() error {
	database.FindUser()
	return nil
}
