package user

import (
	"encoding/json"

	"golang.org/x/crypto/bcrypt"

	"github.com/shayamvlmna/cab-booking-app/pkg/database/redis"
	"github.com/shayamvlmna/cab-booking-app/pkg/models"
	"github.com/shayamvlmna/cab-booking-app/pkg/service/auth"
)

var u = &models.User{}

func RegisterUser(newUser *models.User) error {

	newUser.Phonenumber = auth.GetPhone()
	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	newUser.Password = string(hashPassword)

	if err := AddUser(newUser); err != nil {
		return err
	}
	auth.StorePhone(newUser.Phonenumber)
	return nil
}

func GoogleAuthUser(user *models.User) {

	// usr,ok:=u.Get("email", user.Email)

	// if!ok{
	// 	http.newr
	// }

}

func GoogleSignupUser() {

}

func GoogleLoginUser() {

}

//return boolean to check if the user exist or not
func IsUserExists(key, value string) bool {
	_, err := u.Get(key, value)
	return err
}

//accepts user models and pass to the
//user database to insert
//retun error if any
func AddUser(newUser *models.User) error {
	return newUser.Add()
}

//returns a user model by accepting a key and a value
//eg:if searching using id, key is "id" and value is the id of the user to search
func GetUser(key, value string) *models.User {
	p, err := redis.GetData("data")
	if err != nil {
		user, _ := u.Get(key, value)
		return &user
	}

	user := &models.User{}

	json.Unmarshal([]byte(p), &user)

	return user
}

//return all users in the database
func GetUsers() []models.User {
	return *u.GetAll()
}

//update a user by accepting the updated user fields
//only update fields with null values
func UpdateUser(user *models.User) error {
	return user.Update()
}

//delete user from the database by the id
func DeleteUser(id uint64) {
	u.Delete(id)
}

// func AppendTrip(user *models.User, trip *models.Trip) error {
// 	return database.AppendTrip(user, trip)
// }
