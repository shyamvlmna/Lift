package models

import (
	"strconv"

	"github.com/shayamvlmna/cab-booking-app/pkg/database"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Id          uint64     `gorm:"primaryKey;" json:"id"`
	Picture     string     `json:"picture"`
	Phonenumber string     `gorm:"not null;unique;" json:"phonenumber"`
	Firstname   string     `gorm:"not null;" json:"firstname"`
	Lastname    string     `json:"lastname"`
	Email       string     `gorm:"not null;unique;" json:"email"`
	Password    string     `gorm:"not null;" json:"password"`
	Active      bool       `gorm:"default:true;" json:"status"`
	Wallet      UserWallet `json:"userwallet"`
	TripHistory []Trip     `json:"trip_history" gorm:"foreignKey:UserId"`
	Rating      int        `gorm:"default:0" json:"user_rating"`
}

//add new user to database
func (u *User) Add() error {
	db := database.Db
	err := db.AutoMigrate(&User{})
	if err != nil {
		return err
	}

	result := db.Create(&u)
	return result.Error
}

//get a user by key
func (u *User) Get(key, value string) (User, bool) {
	db := database.Db
	err := db.AutoMigrate(&User{})
	if err != nil {
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
func (u *User) GetAll() *[]User {
	db := database.Db

	users := &[]User{}
	db.Find(&users)

	return users
}

// Update existing user by id
func (u *User) Update() error {
	db := database.Db

	user := &User{}

	id := strconv.Itoa(int(u.Id))

	db.Where("user_id=?", id).First(&user)
	user.TripHistory = append(user.TripHistory, u.TripHistory...)
	result := db.Model(&user).Updates(&User{Phonenumber: "",
		Firstname: "",
		Lastname:  "",
		Email:     "",
		Password:  "",
		Active:    false,
	})

	return result.Error
}

// Delete user by id
func (u *User) Delete(id uint64) error {
	db := database.Db

	result := db.Delete(&User{}, id)

	return result.Error
}

// BlockUnblock user by changing user active field
func (u *User) BlockUnblock(id uint64) error {
	db := database.Db
	err := db.AutoMigrate(&User{})
	if err != nil {
		return err
	}

	user := &User{}

	db.Where("id=?", id).First(&user)

	if !user.Active {
		user.Active = true
		result := db.Save(&user)
		return result.Error
	}
	user.Active = false
	result := db.Save(&user)
	return result.Error
}
