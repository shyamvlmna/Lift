package routes

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/shayamvlmna/lift/pkg/controllers"
	"github.com/shayamvlmna/lift/pkg/middleware"
)

func DriverRoutes(r *mux.Router) {
	//subrouter for driver routes
	driverRouter := r.PathPrefix("/driver").Subrouter()

	//check whether phone-number already registerd or is a new entry
	driverRouter.HandleFunc("/auth", controllers.DriverAuth).Methods(http.MethodPost)

	//register new driver
	driverRouter.HandleFunc("/signup", controllers.DriverSignUp).Methods(http.MethodPost)

	//render enter otp page
	driverRouter.HandleFunc("/enterotp", controllers.EnterOTPDriver).Methods(http.MethodGet)

	//validate submitted otp
	driverRouter.Handle("/otp", middleware.ValidateOtp(controllers.DriverSignUpPage)).Methods(http.MethodPost)

	//render login page to enter password since phone-number already exist
	driverRouter.HandleFunc("/loginpage", controllers.DriverLoginPage).Methods(http.MethodGet)

	//validate entered password with phone-number then redirect to home page
	driverRouter.HandleFunc("/login", controllers.DriverLogin).Methods(http.MethodPost)

	//remove stored cookie and remove data from redis
	driverRouter.HandleFunc("/logout", controllers.DriverLogout).Methods(http.MethodGet)

	//render homepage only if authorized with JWT
	driverRouter.Handle("/driverhome", middleware.IsAuthorized(controllers.DriverHome)).Methods(http.MethodGet)

	driverRouter.Handle("/regtodrive", middleware.IsAuthorized(controllers.RegisterToDrive)).Methods(http.MethodGet)

	//add cab to the driver profile
	driverRouter.Handle("/addcab", middleware.IsAuthorized(controllers.AddCab)).Methods(http.MethodPost)

	//add bank details to checkout wallet balance get add //details page
	driverRouter.Handle("/addbank", middleware.IsAuthorized(controllers.AddBankPage)).Methods(http.MethodGet)

	//add bank details to checkout wallet balance //insert details
	driverRouter.Handle("/addbank", middleware.IsAuthorized(controllers.AddBankAccount)).Methods(http.MethodPost)

	driverRouter.Handle("/editbank", middleware.IsAuthorized(controllers.EditBankDetails)).Methods(http.MethodGet)

	driverRouter.Handle("/updatebank", middleware.IsAuthorized(controllers.UpdateBankDetails)).Methods(http.MethodPut)

	//get current driver details to update
	driverRouter.Handle("/editprofile", middleware.IsAuthorized(controllers.EditDriverProfile)).Methods(http.MethodGet)

	//update driver details in to the database
	driverRouter.Handle("/updateprofile", middleware.IsAuthorized(controllers.UpdateDriverProfile)).Methods(http.MethodPost)

	//get current cab details to update
	driverRouter.Handle("/editcab", middleware.IsAuthorized(controllers.EditCab)).Methods(http.MethodGet)

	//update the cab details in to the database
	driverRouter.Handle("/updatecab", middleware.IsAuthorized(controllers.UpdateCab)).Methods(http.MethodPut)

	//get ride from the channel
	driverRouter.Handle("/getride", middleware.IsAuthorized(controllers.GetTrip)).Methods(http.MethodGet)

	//accept the trip and register it
	driverRouter.Handle("/acceptrip", middleware.IsAuthorized(controllers.AcceptTrip)).Methods(http.MethodPost)

	//match trip code with user trip code
	driverRouter.Handle("/matchtripcode", middleware.IsAuthorized(controllers.MatchTripCode)).Methods(http.MethodPost)

	//start trip after matching trip code
	driverRouter.Handle("/startrip", middleware.IsAuthorized(controllers.StartTrip)).Methods(http.MethodGet)

	//end trip after trip completion and get payment
	driverRouter.Handle("/endtrip", middleware.IsAuthorized(controllers.EndTrip))

	//list driver trip history
	driverRouter.Handle("/triphistory", middleware.IsAuthorized(controllers.DriverTripHistory)).Methods(http.MethodGet)

	//get driver wallet data
	driverRouter.Handle("/wallet", middleware.IsAuthorized(controllers.DriverWallet)).Methods(http.MethodGet)

	//payouth the money in the wallet
	driverRouter.Handle("/payout", middleware.IsAuthorized(controllers.PayoutWallet)).Methods(http.MethodPost)

	//get status of submitted payout requests
	driverRouter.Handle("/payoutstatus", middleware.IsAuthorized(controllers.PayoutStatus)).Methods(http.MethodGet)

	driverRouter.Handle("/payouthistory", middleware.IsAuthorized(controllers.PayoutHistory)).Methods(http.MethodGet)

	driverRouter.HandleFunc("/upload", controllers.UploadDocuments).Methods(http.MethodPost)
}
