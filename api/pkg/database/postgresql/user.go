package database

import (
	"strconv"

	"gorm.io/gorm"

	"github.com/shayamvlmna/cab-booking-app/pkg/database"
	"github.com/shayamvlmna/cab-booking-app/pkg/models"
)

// func closeUserdb(db *gorm.DB) {

// 	sqlDb, err := db.DB()
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	sqlDb.Close()
// 	fmt.Println("user db closed")
// }

//receive a user model and insert it into the user database
func InsertUser(user *models.User) error {

	db := database.UserData(database.Db)

	result := db.Create(&user)

	return result.Error
}

func FindUser(key, value string) (models.User, bool) {

	db := database.UserData(database.Db)

	user := &models.User{}
	result := db.Where(key+"=?", value).First(&user)

	if result.Error == gorm.ErrRecordNotFound {
		return *user, false
	}
	return *user, true

}

//get and return all users from the driver database
func GetUsers() *[]models.User {
	db := database.DriverData(database.Db)
	// defer closeUserdb(db)

	users := &[]models.User{}
	db.Find(&users)

	return users
}

//update a user by getting updated user fields
//only update the not null user fields
func UpdateUser(updatedUser *models.User) error {
	db := database.UserData(database.Db)

	user := &models.User{}

	id := strconv.Itoa(int(updatedUser.UserId))

	db.Where("user_id=?", id).First(&user)
	user.TripHistory = append(user.TripHistory, updatedUser.TripHistory...)
	result := db.Model(&user).Updates(&models.User{
		FirstName:   updatedUser.FirstName,
		LastName:    updatedUser.LastName,
		PhoneNumber: updatedUser.PhoneNumber,
		Email:       updatedUser.Email,
		Password:    updatedUser.Password,
		Active:      updatedUser.Active,
		Wallet:      updatedUser.Wallet,
		TripHistory: updatedUser.TripHistory,
	})

	return result.Error
}

func AppendTrip(updatedUser *models.User, trip *models.Ride) error {
	db := database.UserData(database.Db)
	user := &models.User{}
	db.AutoMigrate(&user)

	id := strconv.Itoa(int(updatedUser.UserId))

	db.Where("user_id=?", id).First(&user)
	user.TripHistory = append(user.TripHistory, *trip)
	result := db.Model(&user).Updates(&models.User{
		TripHistory: user.TripHistory,
	})
	// db.Save(&user)
	return result.Error
}

//delete user by id
//returns err if any
func DeleteUser(id uint64) error {
	db := database.UserData(database.Db)

	user := &models.User{}

	result := db.Delete(&user, id)

	return result.Error
}
