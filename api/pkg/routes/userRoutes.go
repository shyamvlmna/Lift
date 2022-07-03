package routes

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/shayamvlmna/cab-booking-app/pkg/controllers"
	"github.com/shayamvlmna/cab-booking-app/pkg/middleware"
	"github.com/shayamvlmna/cab-booking-app/pkg/models"
	googleauth "github.com/shayamvlmna/cab-booking-app/pkg/service/googleAuth"
)

func UserRoutes(r *mux.Router) {
	//subrouter for user routes
	userRouter := r.PathPrefix("/user").Subrouter()

	//check wheather phonenumber already registerd or is a new entry
	userRouter.HandleFunc("/auth", controllers.UserAuth).Methods("POST")

	//register new user insert user data to database
	userRouter.HandleFunc("/signup", controllers.UserSignUp).Methods("POST")

	//match entered password with the phone-number then redirect to the home page
	userRouter.HandleFunc("/login", controllers.UserLogin).Methods("POST")

	//login with Google handler
	userRouter.HandleFunc("/googlelogin", googleauth.GoogleLogin)

	//login with Google callback handler
	userRouter.HandleFunc("/googleCallback", googleauth.GoogleCallback)

	//logout session handler
	userRouter.HandleFunc("/logout", controllers.UserLogout)

	//render enter otp page
	userRouter.HandleFunc("/enterotp", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		response := &models.Response{
			ResponseStatus:  "success",
			ResponseMessage: "new user",
			ResponseData:    nil,
		}
		err := json.NewEncoder(w).Encode(&response)
		if err != nil {
			return
		}

	}).Methods("GET")

	//validate submited otp
	userRouter.Handle("/otp", middleware.ValidateOtp(controllers.UserSignupPage)).Methods("POST")

	//render login page to enter password
	userRouter.HandleFunc("/loginpage", controllers.UserLoginPage).Methods("GET")

	//render the homepage only if authorized with JWT
	userRouter.Handle("/userhome", middleware.IsAuthorized(controllers.UserHome)).Methods("GET")

	//get current user profile details to update
	userRouter.HandleFunc("/update/{id}", controllers.EditUserProfile).Methods("GET")

	//update user profile details
	userRouter.HandleFunc("/update/{id}", controllers.UpdateUserProfile).Methods("POST")

	//book new trip with location latitude and longitude from the frontend
	userRouter.Handle("/booktrip", middleware.IsAuthorized(controllers.BookTrip)).Methods("POST")

	//confirm the trip and select payment method
	userRouter.Handle("/confirmtrip", middleware.IsAuthorized(controllers.ConfirmTrip)).Methods("POST")

	//get user trip history
	userRouter.Handle("/triphistory", middleware.IsAuthorized(controllers.UserTripHistory)).Methods("GET")

	//get user wallet details
	userRouter.Handle("/wallet", middleware.IsAuthorized(controllers.UserWallet)).Methods("GET")

	//add money to user wallet
	userRouter.Handle("/addmoney", middleware.IsAuthorized(controllers.AddMoneyToWallet)).Methods("GET")

	userRouter.HandleFunc("/razorpaycallback", controllers.RazorpayCallback)

	// userRouter.HandleFunc("/test", controllers.Test)

}
