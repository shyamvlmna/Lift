package database

import (
	"github.com/shayamvlmna/cab-booking-app/pkg/models"
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
	return Db, nil
}

func AddUser(newUser *models.User) {

}

func FindUser() {

}

func UpdateUser() {

}
