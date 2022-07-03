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
	//subrouter for driver routes
	driverRouter := r.PathPrefix("/driver").Subrouter()

	//check whether phone-number already registerd or is a new entry
	driverRouter.HandleFunc("/auth", controllers.DriverAuth).Methods("POST")

	//register new driver
	driverRouter.HandleFunc("/signup", controllers.DriverSignUp).Methods("POST")

	//render enter otp page
	driverRouter.HandleFunc("/enterotp", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		response := models.Response{
			ResponseStatus:  "success",
			ResponseMessage: "new driver",
			ResponseData:    nil,
		}
		err := json.NewEncoder(w).Encode(&response)
		if err != nil {
			return
		}
	}).Methods("GET")

	//validate submitted otp
	driverRouter.Handle("/otp", middleware.ValidateOtp(controllers.DriverSignUpPage)).Methods("POST")

	//render login page to enter password since phone-number already exist
	driverRouter.HandleFunc("/loginpage", controllers.DriverLoginPage).Methods("GET")

	//validate entered password with phone-number then redirect to home page
	driverRouter.HandleFunc("/login", controllers.DriverLogin).Methods("POST")

	//remove stored cookie and remove data from redis
	driverRouter.HandleFunc("/logout", controllers.DriverLogout).Methods("GET")

	//render homepage only if authorized with JWT
	driverRouter.Handle("/driverhome", middleware.IsAuthorized(controllers.DriverHome)).Methods("GET")

	//add cab to the driver profile
	driverRouter.Handle("/addcab", middleware.IsAuthorized(controllers.AddCab)).Methods("POST")

	//get current driver details to update
	driverRouter.Handle("/editprofile", middleware.IsAuthorized(controllers.EditDriverProfile)).Methods("GET")

	//update driver details in to the database
	driverRouter.Handle("/updateprofile", middleware.IsAuthorized(controllers.UpdateDriverProfile)).Methods("POST")

	//get current cab details to update
	driverRouter.Handle("/editcab", middleware.IsAuthorized(controllers.EditCab)).Methods("GET")

	//update the cab details in to the database
	driverRouter.Handle("/updatecab", middleware.IsAuthorized(controllers.UpdateCab)).Methods("POST")

	//get ride from the channel
	driverRouter.Handle("/getride", middleware.IsAuthorized(controllers.GetTrip)).Methods("GET")

	//accept the trip and register it
	driverRouter.Handle("/acceptrip", middleware.IsAuthorized(controllers.AcceptTrip)).Methods("POST")

	//match trip code with user trip code
	driverRouter.Handle("/matchtripcode", middleware.IsAuthorized(controllers.MatchTripCode)).Methods("POST")

	//start trip after matching trip code
	driverRouter.Handle("/startrip", middleware.IsAuthorized(controllers.StartTrip)).Methods("GET")

	//end trip after trip completion and get payment
	driverRouter.Handle("/endtrip", middleware.IsAuthorized(controllers.EndTrip))

	//list driver trip history
	driverRouter.Handle("/triphistory", middleware.IsAuthorized(controllers.DriverTripHistory)).Methods("GET")

}
