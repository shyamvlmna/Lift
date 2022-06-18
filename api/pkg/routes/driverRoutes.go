package routes

import (
	"github.com/gorilla/mux"

	"github.com/shayamvlmna/cab-booking-app/pkg/controllers"
)

func DriverRoutes(r *mux.Router) {
	driverRouter := r.PathPrefix("/driver").Subrouter()

	driverRouter.HandleFunc("/auth", controllers.DriverAuth).Methods("POST")
	driverRouter.HandleFunc("/signup", controllers.DriverSignUp).Methods("POST")
	driverRouter.HandleFunc("/login", controllers.DriverLogin).Methods("POST")

	driverRouter.HandleFunc("/addcab", controllers.AddCab)
	driverRouter.HandleFunc("/regtodrive", controllers.RegisterDriver)

}