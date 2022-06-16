package user

import (
	"github.com/shayamvlmna/cab-booking-app/app/database"
	"github.com/shayamvlmna/cab-booking-app/app/models"
)

func AddUser(newUser *models.User) error {
	return database.InsertUser(newUser)
}

func GetUser() error {
	database.FindUser()
	return nil
}

func IsUserExists(key string) bool {
	return database.CheckUser(key)
}
