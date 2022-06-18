package models

import "gorm.io/gorm"

type Vehicle struct {
	gorm.Model
	Owner        Driver `gorm:"foreignKey:Registration"`
	Registration string `gorm:"not null;unique" json:"registration"`
	Brand        string `gorm:"not null" json:"cabrand"`
	Category     string `gorm:"not null" json:"cabtype"`
	VehicleModel string `gorm:"not null" json:"cabmodel"`
	Colour       string `gorm:"not null" json:"cabcolour"`
}
