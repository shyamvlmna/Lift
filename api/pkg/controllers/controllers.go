package controllers

import (
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
	if err == nil {
		http.Redirect(w, r, "/user/userhome", http.StatusSeeOther)
		return
	}
	// tokenstring := c.Value
	// phone,err!=auth.ValidateJWT(tokenstring)

	IndexTemp.ExecuteTemplate(w, "appIndex.html", nil)
}

// func validateCookie(w http.ResponseWriter, r *http.Request) {

// }
