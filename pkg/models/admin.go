package models

import "gorm.io/gorm"

type Admin struct {
	gorm.Model
	Username string
	Password string
	IsAdmin  bool
}
