package controllers

import (
	"encoding/json"
	"net/http"

	database "github.com/shayamvlmna/cab-booking-app/pkg/database/postgresql"
	"github.com/shayamvlmna/cab-booking-app/pkg/models"
)

func RegisterDriver(w http.ResponseWriter, r *http.Request) {
	// city := r.FormValue("city")
	// dlNumber := r.FormValue("driving_licence")

}
func AddCab(w http.ResponseWriter, r *http.Request) {
	vehicle := models.Vehicle{}
	json.NewDecoder(r.Body).Decode(&vehicle)
	database.Insert(&vehicle)
}
