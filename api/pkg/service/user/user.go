package user

import (
	"errors"

	"github.com/shayamvlmna/cab-booking-app/pkg/service/coupon"

	"golang.org/x/crypto/bcrypt"

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

// IsUserExists return boolean to check if the user exist or not
func IsUserExists(key, value string) bool {
	_, err := u.Get(key, value)
	return err
}

// AddUser accepts user models and pass to the
//user database to insert
//retun error if any
func AddUser(newUser *models.User) error {
	return newUser.Add()
}

// GetUser returns a user model by accepting a key and a value
//eg:if searching using id, key is "id" and value is the id of the user to search
func GetUser(key, value string) *models.User {
	//p, err := redis.GetData("data")
	//if err != nil {
	user, _ := u.Get(key, value)
	return &user
	//}

	//user := &models.User{}
	//
	//err = json.Unmarshal([]byte(p), &user)
	//if err != nil {
	//	return nil
	//}
	//
	//return user
}

// GetUsers return all users in the database
func GetUsers() (*[]models.User, error) {

	return u.GetAll()
}

// UpdateUser update a user by accepting the updated user fields
//only update fields with null values
func UpdateUser(user *models.User) error {
	return user.Update()
}

func BlockUser(id uint) error {
	return u.BlockUnblock(id)
}

func UnBlockUser(id uint) error {
	return u.BlockUnblock(id)
}

// DeleteUser delete user from the database by the id
func DeleteUser(id uint64) {
	err := u.Delete(id)
	if err != nil {
		return
	}
}

func ApplyCoupon(code string, fare float64) (float64, error) {

	c := coupon.GetCoupon(code)

	if !c.IsApplicable(fare) {
		err := errors.New("coupon not applicable")
		return fare, err
	}
	return fare - c.Amount, nil
}
