package models

type Vehicle struct {
	Id           uint   `gorm:"primaryKey"`
	Registration string `gorm:"not null;unique" json:"registration"`
	// Owner        Driver `gorm:"foreignKey:id" json:"owner"`
	Brand        string `gorm:"not null" json:"cabrand"`
	Category     string `gorm:"not null" json:"cabtype"`
	VehicleModel string `gorm:"not null" json:"cabmodel"`
	Colour       string `gorm:"not null" json:"cabcolour"`
}

//
