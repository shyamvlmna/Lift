package routes

import (
	"github.com/shayamvlmna/lift/pkg/controllers"
	"github.com/shayamvlmna/lift/pkg/middleware"
	"github.com/shayamvlmna/lift/pkg/service/googleAuth"
	"net/http"

	"github.com/gorilla/mux"
)

func UserRoutes(r *mux.Router) {
	//subrouter for user routes
	userRouter := r.PathPrefix("/user").Subrouter()

	//check wheather phonenumber already registerd or is a new entry
	userRouter.HandleFunc("/auth", controllers.UserAuth).Methods(http.MethodPost)

	//register new user insert user data to database
	userRouter.HandleFunc("/signup", controllers.UserSignUp).Methods(http.MethodPost)

	//match entered password with the phone-number then redirect to the home page
	userRouter.HandleFunc("/login", controllers.UserLogin).Methods(http.MethodPost)

	//login with Google handler
	userRouter.HandleFunc("/googlelogin", googleauth.googleauth.GoogleLogin)

	//login with Google callback handler
	userRouter.HandleFunc("/googleCallback", googleauth.GoogleCallback)

	//logout session handler
	userRouter.HandleFunc("/logout", controllers.UserLogout)

	//render enter otp page
	userRouter.HandleFunc("/enterotp", controllers.EnterOTPUser).Methods(http.MethodGet)

	//validate submited otp
	userRouter.Handle("/otp", middleware.ValidateOtp(controllers.UserSignupPage)).Methods(http.MethodPost)

	//render login page to enter password
	userRouter.HandleFunc("/loginpage", controllers.UserLoginPage).Methods(http.MethodGet)

	//render the homepage only if authorized with JWT
	userRouter.Handle("/userhome", middleware.IsAuthorized(controllers.UserHome)).Methods(http.MethodGet)

	//get current user profile details to update
	userRouter.HandleFunc("/editprofile", controllers.EditUserProfile).Methods(http.MethodGet)

	//update user profile details
	userRouter.HandleFunc("/updateprofile", controllers.UpdateUserProfile).Methods(http.MethodPut)

	//book new trip with location latitude and longitude from the frontend
	userRouter.Handle("/booktrip", middleware.IsAuthorized(controllers.BookTrip)).Methods(http.MethodPost)

	//confirm the trip and select payment method
	userRouter.Handle("/confirmtrip", middleware.IsAuthorized(controllers.ConfirmTrip)).Methods(http.MethodPost)

	//get user trip history
	userRouter.Handle("/triphistory", middleware.IsAuthorized(controllers.UserTripHistory)).Methods(http.MethodGet)

	//get user wallet details
	userRouter.Handle("/wallet", middleware.IsAuthorized(controllers.UserWallet)).Methods(http.MethodGet)

	//add money to user wallet
	userRouter.Handle("/addmoney", middleware.IsAuthorized(controllers.AddMoneyToWallet)).Methods(http.MethodGet)

	userRouter.HandleFunc("/razorpaycallback", controllers.RazorpayCallback)

	userRouter.HandleFunc("/razorpay/webhook", controllers.RazorpayWebhook).Methods(http.MethodPost)

	//returns all available coupons
	userRouter.Handle("/coupons", middleware.IsAuthorized(controllers.GetCoupons)).Methods(http.MethodGet)

	userRouter.Handle("/applycoupon", middleware.IsAuthorized(controllers.ApplyCoupon)).Methods(http.MethodPost)
}
