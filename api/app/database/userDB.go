package database

import (
	"fmt"
	"strconv"

	"github.com/shayamvlmna/cab-booking-app/app/models"
	"gorm.io/gorm"
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
	} else {
		return *user, true
	}

}

func UpdateUser(updatedUser *models.User) {
	db, err := OpenUserDb()
	if err != nil {
		fmt.Println(err)
	}
	user := &models.User{}
	id := strconv.Itoa(int(updatedUser.ID))
	db.Where("id=?", id).First(&user)
	db.Model(&user).Updates(models.User{
		FirstName: updatedUser.FirstName,
		LastName:  updatedUser.LastName,
		Email:     updatedUser.Email,
	})
}
