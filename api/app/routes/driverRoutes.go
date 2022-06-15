package routes

import (
	"github.com/gorilla/mux"
	"github.com/shayamvlmna/cab-booking-app/app/controllers"
)

func DriverRoutes(r *mux.Router) {
	driverRouter := r.PathPrefix("/driver").Subrouter()
	// driverRouter.HandleFunc("/signup", controllers.DriverSignUp).Methods("GET")
	// driverRouter.HandleFunc("/login", controllers.DriverLogin).Methods("GET")

	driverRouter.HandleFunc("/signup", controllers.DriverSignUp).Methods("POST")
	driverRouter.HandleFunc("/login", controllers.DriverLogin).Methods("POST")

}
