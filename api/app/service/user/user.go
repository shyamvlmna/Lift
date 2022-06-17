package user

import (
	"github.com/shayamvlmna/cab-booking-app/app/database"
	"github.com/shayamvlmna/cab-booking-app/app/models"
)

//accepts user models and pass to the
//user database to insert
//retun error if any
func AddUser(newUser *models.User) error {
	return database.InsertUser(newUser)
}

//accepts a key and a value to a user from the
//user database. key indicate the field in the database
//value is the value of the field
//eg:if searching using id key is "id" and value is the
//id of the user to search
//returns a user model
func GetUser(key, value string) models.User {
	user, _ := database.FindUser(key, value)
	return user

}

//return boolean to check if the user exist or not
func IsUserExists(key, value string) bool {
	_, err := database.FindUser(key, value)
	return err
}

func UpdateUser(user *models.User) {
	database.UpdateUser(user)
}
