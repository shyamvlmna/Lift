package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/shayamvlmna/cab-booking-app/app/controllers"
)

func DriverRoutes(r *mux.Router) {
	driverRouter := r.PathPrefix("/driver").Subrouter()

	driverRouter.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		controllers.DriverTemp.ExecuteTemplate(w, "driverSignupForm.html", nil)
	})
	driverRouter.HandleFunc("/auth", controllers.DriverAuth).Methods("POST")

	driverRouter.HandleFunc("/signup", controllers.DriverSignUp).Methods("POST")
	driverRouter.HandleFunc("/login", controllers.DriverLogin).Methods("POST")
	driverRouter.HandleFunc("/addcab", controllers.AddCab)
	driverRouter.HandleFunc("/regtodrive", controllers.RegisterDriver)

}
