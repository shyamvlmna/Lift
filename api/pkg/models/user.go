package models

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/shayamvlmna/cab-booking-app/pkg/database"
)

type User struct {
	gorm.Model
	UserId        uint    `gorm:"primaryKey;autoIncrement;unique" json:"userid"`
	Picture       string  `json:"picture"`
	Phonenumber   string  `gorm:"not null;unique;" json:"phonenumber"`
	Firstname     string  `gorm:"not null;" json:"firstname"`
	Lastname      string  `json:"lastname"`
	Email         string  `gorm:"not null;unique;" json:"email"`
	Password      string  `gorm:"not null;" json:"password"`
	Rating        int     `gorm:"default:0" json:"user_rating"`
	Active        bool    `gorm:"default:true;" json:"status"`
	WalletBalance float64 `json:"userwallet" gorm:"default:0;"`
}

// Add new user to database
func (u *User) Add() error {
	db := database.Db
	err := db.AutoMigrate(&User{})
	if err != nil {
		return err
	}

	result := db.Create(&u)
	return result.Error
}

// Get a user by key
func (u *User) Get(key, value string) (User, bool) {
	db := database.Db
	if err := db.AutoMigrate(&User{}); err != nil {
		fmt.Println("error migrating")
		return User{}, false
	}

	user := &User{}
	result := db.Where(key+"=?", value).First(&user)

	if result.Error == gorm.ErrRecordNotFound {
		return *user, false
	}
	return *user, true
}

// GetAll users in the database
func (u *User) GetAll() (*[]User, error) {
	db := database.Db

	users := &[]User{}
	result := db.Find(&users)

	return users, result.Error
}

// Update existing user by id
func (u *User) Update(id uint) error {
	db := database.Db

	user := &User{}

	db.Where("user_id=?", id).First(&user)

	result := db.Model(&user).Updates(&User{Phonenumber: "",
		Firstname: u.Firstname,
		Lastname:  u.Lastname,
		Email:     u.Email,
		Password:  u.Password,
	})

	return result.Error
}

// Delete user by id
func (*User) Delete(id uint64) error {
	db := database.Db

	result := db.Delete(&User{}, id)

	return result.Error
}

// BlockUnblock user by changing user active field
func (*User) BlockUnblock(id uint) error {
	db := database.Db

	user := &User{}

	db.Where("user_id=?", id).First(&user)

	if user.Active {
		result := db.Model(&user).Update("active", false)
		return result.Error
	}

	result := db.Model(&user).Update("active", true)
	return result.Error
}
