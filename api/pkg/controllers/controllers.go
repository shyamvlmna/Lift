package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/shayamvlmna/cab-booking-app/pkg/models"
	auth "github.com/shayamvlmna/cab-booking-app/pkg/service/auth"
)

func Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache,no-store,must-revalidate")
	w.Header().Set("Content-Type", "application/json")

	// headers := r.Header
	// _, ok := headers["Token"]

	c, err := r.Cookie("jwt-token")

	if err == nil {
		tokenString := c.Value
		role, _ := auth.ParseJWT(tokenString)

		if role == "driver" {
			http.Redirect(w, r, "/driver/driverhome", http.StatusSeeOther)
			return
		} else if role == "user" {
			http.Redirect(w, r, "/user/userhome", http.StatusSeeOther)
			return
		}
	} else {
		response := models.Response{
			ResponseStatus:  "success",
			ResponseMessage: "app index",
			ResponseData:    nil,
		}
		json.NewEncoder(w).Encode(&response)
	}
}
