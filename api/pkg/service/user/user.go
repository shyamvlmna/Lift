package user

import (
	database "github.com/shayamvlmna/cab-booking-app/pkg/database/postgresql"
	"github.com/shayamvlmna/cab-booking-app/pkg/models"
)

//accepts user models and pass to the
//user database to insert
//retun error if any
func AddUser(newUser *models.User) error {
	return database.InsertUser(newUser)
}

//returns a user model by accepting a key and a value
//eg:if searching using id, key is "id" and value is the id of the user to search
func GetUser(key, value string) models.User {
	user, _ := database.FindUser(key, value)
	return user

}

//return all users in the database
func GetUsers() []models.User {

	return *database.GetUsers()
}

//update a user by accepting the updated user fields
//only update fields with null values
func UpdateUser(user *models.User) error {
	return database.UpdateUser(user)
}

//delete user from the database by the id
func DeleteUser(id string) {

}

//return boolean to check if the user exist or not
func IsUserExists(key, value string) bool {
	_, err := database.FindUser(key, value)
	return err
}
