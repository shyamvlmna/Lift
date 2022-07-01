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

	//check wheather phonenumber already registerd or is a new entry
	driverRouter.HandleFunc("/auth", controllers.DriverAuth).Methods("POST")

	//insert data to the database
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

	//validate submited otp
	driverRouter.Handle("/otp", middleware.ValidateOtp(controllers.DriverSignUpPage)).Methods("POST")

	//render login page to enter password since phonenumber alredy exist
	driverRouter.HandleFunc("/loginpage", controllers.DriverLoginPage).Methods("GET")

	//validate entered password with phonenumber and render home page
	driverRouter.HandleFunc("/login", controllers.DriverLogin).Methods("POST")

	//remove stored cookie and remove data from redis
	driverRouter.HandleFunc("/logout", controllers.DriverLogout).Methods("GET")

	//render homepage only if authorized with JWT
	driverRouter.Handle("/driverhome", middleware.IsAuthorized(controllers.DriverHome)).Methods("GET")

	driverRouter.Handle("/regtodrive", middleware.IsAuthorized(controllers.RegisterDriver)).Methods("POST")
	driverRouter.Handle("/addcab", middleware.IsAuthorized(controllers.AddCab)).Methods("POST")

	//get ride from the channel
	driverRouter.Handle("/getride", middleware.IsAuthorized(controllers.GetTrip)).Methods("GET")

	//accept the trip and register it
	driverRouter.Handle("/acceptrip", middleware.IsAuthorized(controllers.AcceptTrip)).Methods("POST")

	driverRouter.Handle("/matchtripcode", middleware.IsAuthorized(controllers.MatchTripCode)).Methods("POST")

	driverRouter.Handle("/startrip", middleware.IsAuthorized(controllers.StartTrip)).Methods("GET")

}
