package models

// Source        Location
// Trip          Location `gorm:"ForeignKey:TripId"`

type Location struct {
	Id  uint `gorm:"primaryKey"`
	Lon string
	Lat string
}

type Payment struct {
	Wallet bool
	Cash   bool
}
