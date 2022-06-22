package controllers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/shayamvlmna/cab-booking-app/pkg/models"
	"github.com/shayamvlmna/cab-booking-app/pkg/service/auth"
)

var (
	// userTemp, _   = template.ParseGlob("/home/shyamjith/cab-booking-app/ui/user/*.html")
	// driverTemp, _ = template.ParseGlob("/home/shyamjith/cab-booking-app/ui/driver/*.html")
	adminTemp, _ = template.ParseGlob("/home/shyamjith/cab-booking-app/ui/admin/*.html")
)

func Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache,no-store,must-revalidate")
	w.Header().Set("Content-Type", "application/json")
	// fmt.Println(r.Cookie("jwt-token"))

	c, err := r.Cookie("jwt-token")

	if err != nil {
		response := models.Response{
			ResponseStatus:  "success",
			ResponseMessage: "app index",
			ResponseData:    nil,
		}
		json.NewEncoder(w).Encode(&response)
		return
	} else {
		fmt.Println(c.Value)
		role, _, err := auth.ValidateJWT(r)
		if err == nil {
			if role == "driver" {
				http.Redirect(w, r, "/driver/driverhome", http.StatusSeeOther)
				return
			} else if role == "user" {
				http.Redirect(w, r, "/user/userhome", http.StatusSeeOther)
				response := models.Response{
					ResponseStatus:  "success",
					ResponseMessage: "valid jwt-token",
					ResponseData:    nil,
				}
				json.NewEncoder(w).Encode(&response)
				return
			}
		}

	}

}
