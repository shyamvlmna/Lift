package driver

import "github.com/shayamvlmna/cab-booking-app/pkg/database"

func AddDriver() error {

	database.AddDriver()
	return nil
}

func GetDriver() error {
	database.FindDriver()
	return nil
}
