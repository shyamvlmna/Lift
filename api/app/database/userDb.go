package database

import (
	"fmt"

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
func closeUserdb(db *gorm.DB) {

	sqlDb, err := db.DB()
	if err != nil {
		fmt.Println(err)
	}
	sqlDb.Close()
	fmt.Println("user db closed")
}
func FindUser(key string) (models.User, bool) {

	db, err := OpenUserDb()
	if err != nil {
		fmt.Println(err)
	}
	defer closeUserdb(db)
	user := &models.User{}
	result := db.Where("phone_number=?", key).First(&user)

	if result.Error == gorm.ErrRecordNotFound {
		return *user, false
	} else {
		return *user, true
	}

}

func InsertUser(user *models.User) error {

	db, err := OpenUserDb()
	if err != nil {
		return err
	}
	defer closeUserdb(db)
	result := db.Create(user)

	return result.Error
}

func UpdateUser() {

}
