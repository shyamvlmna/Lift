package controllers

import (
	"fmt"

	"net/http"

	"github.com/shayamvlmna/cab-booking-app/app/database"
	"github.com/shayamvlmna/cab-booking-app/app/models"
)

func Index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Index page"))

	db, err := database.OpenUserDb()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(db)

	Db, err := database.OpenDriverDb()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(Db)
}
func Cab(w http.ResponseWriter, r *http.Request) {
	Db, err := database.OpenDriverDb()
	if err != nil {
		fmt.Println(err)
	}
	cab := &models.Vehicle{}
	Db.AutoMigrate(&cab)
	fmt.Println(Db)
}
