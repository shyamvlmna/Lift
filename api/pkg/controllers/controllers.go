package controllers

import (
	"html/template"
	"log"
	"net/http"

	"github.com/shayamvlmna/lift/pkg/service/auth"
)

var temp, _ = template.ParseGlob("*.html")

func Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache,no-store,must-revalidate")
	// w.Header().Set("Content-Type", "application/json")

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
		// response := models.Response{
		// 	ResponseStatus:  "success",
		// 	ResponseMessage: "app index",
		// 	ResponseData:    nil,
		// }

		// err := json.NewEncoder(w).Encode(&response)
		// tmp, err := template.ParseFiles("index.html")
		err := temp.ExecuteTemplate(w, "index.html", nil)
		if err != nil {
			log.Fatal(err)
		}
	}
}
