package routes

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/shayamvlmna/cab-booking-app/pkg/controllers"
	"github.com/shayamvlmna/cab-booking-app/pkg/middleware"
	"github.com/shayamvlmna/cab-booking-app/pkg/models"
)

func DriverRoutes(r *mux.Router) {
	driverRouter := r.PathPrefix("/driver").Subrouter()

	driverRouter.HandleFunc("/auth", controllers.DriverAuth).Methods("POST")
	driverRouter.HandleFunc("/signup", controllers.DriverSignUp).Methods("POST")

	//render enter otp page
	driverRouter.HandleFunc("/enterotp", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		response := models.Response{
			ResponseStatus:  "success",
			ResponseMessage: "new driver",
			ResponseData:    nil,
		}
		json.NewEncoder(w).Encode(&response)
	}).Methods("GET")

	//enter login page to enter password
	driverRouter.HandleFunc("/loginpage", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		response := models.Response{
			ResponseStatus:  "success",
			ResponseMessage: "existing driver",
			ResponseData:    nil,
		}
		json.NewEncoder(w).Encode(&response)
	}).Methods("GET")

	driverRouter.HandleFunc("/login", controllers.DriverLogin).Methods("POST")
	driverRouter.HandleFunc("/logout", controllers.DriverLogout).Methods("GET")

	driverRouter.Handle("/otp", middleware.ValidateOtp(controllers.DriverSignUpPage)).Methods("POST")
	driverRouter.Handle("/driverhome", middleware.IsAuthorized(controllers.DriverHome)).Methods("GET")

	driverRouter.Handle("/regtodrive", middleware.IsAuthorized(controllers.RegisterDriver)).Methods("POST")
	driverRouter.Handle("/addcab", middleware.IsAuthorized(controllers.AddCab)).Methods("POST")

	driverRouter.Handle("/getride", middleware.IsAuthorized(controllers.GetTrip)).Methods("GET")

}
