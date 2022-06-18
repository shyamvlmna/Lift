package database

import (
	"fmt"
	"strconv"

	"gorm.io/gorm"

	"github.com/shayamvlmna/cab-booking-app/pkg/models"
)

func OpenUserDb() (*gorm.DB, error) {
	Db, err := openDB()
	if err != nil {
		return nil, err
	}
	user := &models.User{}
	err = Db.AutoMigrate(&user)
	if err != nil {
		return nil, err
	}
	fmt.Println("user db opened")
	return Db, nil
}

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

	db, err := OpenUserDb()
	if err != nil {
		return err
	}
	// defer closeUserdb(db)
	result := db.Create(user)

	return result.Error
}

func FindUser(key, value string) (models.User, bool) {

	db, err := OpenUserDb()
	if err != nil {
		fmt.Println(err)
	}
	// defer closeUserdb(db)
	user := &models.User{}
	result := db.Where(key+"=?", value).First(&user)

	if result.Error == gorm.ErrRecordNotFound {
		return *user, false
	}
	return *user, true

}

//get and return all users from the driver database
func GetUsers() *[]models.User {
	db, err := OpenUserDb()
	if err != nil {
		fmt.Println(err)
	}

	users := &[]models.User{}
	db.Find(&users)

	return users
}

//update a user by getting updated user fields
//only update the not null user fields
func UpdateUser(updatedUser *models.User) error {
	db, err := OpenUserDb()
	if err != nil {
		return err
	}
	user := &models.User{}
	id := strconv.Itoa(int(updatedUser.ID))
	db.Where("id=?", id).First(&user)
	result := db.Model(&user).Updates(models.User{
		FirstName: updatedUser.FirstName,
		LastName:  updatedUser.LastName,
		Email:     updatedUser.Email,
	})
	return result.Error
}

//delete user by id
//returns err if any
func DeleteUser(id string) error {

	return nil
}
