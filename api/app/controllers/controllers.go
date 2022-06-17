package controllers

import (
	"fmt"
	"html/template"
	"net/http"
)

var (
	UserTemp, _   = template.ParseGlob("/home/shyamjith/cab-booking-app/ui/user/*.html")
	DriverTemp, _ = template.ParseGlob("/home/shyamjith/cab-booking-app/ui/driver/*.html")
	IndexTemp, _  = template.ParseGlob("/home/shyamjith/cab-booking-app/ui/index/*.html")
	AdminTemp, _  = template.ParseGlob("/home/shyamjith/cab-booking-app/ui/admin/*.html")
)

func Index(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("jwt-token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			IndexTemp.ExecuteTemplate(w, "index.html", nil)
			fmt.Println("Index Page")
		}
	}
	// tokenstring := c.Value
	// phone,err!=auth.ValidateJWT(tokenstring)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
	}
	IndexTemp.ExecuteTemplate(w, "userhome.html", nil)
}
