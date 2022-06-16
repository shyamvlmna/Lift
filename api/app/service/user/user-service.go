package user

import (
	"github.com/shayamvlmna/cab-booking-app/app/database"
	"github.com/shayamvlmna/cab-booking-app/app/models"
)

func AddUser(newUser *models.User) error {
	return database.InsertUser(newUser)
}

func GetUser(key string) models.User {
	user, _ := database.FindUser(key)
	return user

}

func IsUserExists(key string) bool {
	_, err := database.FindUser(key)
	return err
}
