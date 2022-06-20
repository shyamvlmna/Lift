package routes

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/shayamvlmna/cab-booking-app/pkg/controllers"
)

func UserRoutes(r *mux.Router) {

	userRouter := r.PathPrefix("/user").Subrouter()

	userRouter.HandleFunc("/auth", controllers.UserAuth).Methods("POST")
	userRouter.HandleFunc("/otp", controllers.ValidateOtp).Methods("POST")
	userRouter.HandleFunc("/signup", controllers.UserSignUp).Methods("POST")
	userRouter.HandleFunc("/login", controllers.UserLogin).Methods("POST")
	userRouter.HandleFunc("/logout", controllers.UserLogout).Methods("GET")

	userRouter.HandleFunc("/userhome", controllers.UserHome).Methods("GET")
	userRouter.HandleFunc("/update/{id}", controllers.EditUserProfile).Methods("GET")
	userRouter.HandleFunc("/update/{id}", controllers.UpdateUserProfile).Methods("POST")
	userRouter.HandleFunc("/getusers", controllers.GetUsers).Methods("GET")

	//render enter otp page
	userRouter.HandleFunc("/enterotp", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("submit the otp"))
	}).Methods("GET")

	//render signup page
	userRouter.HandleFunc("/signup", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("user signup page\n\nfirstname\nlastname\nemail\npassword"))
	}).Methods("GET")

	//enter login page to enter password
	userRouter.HandleFunc("/loginpage", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("login page\nOnly submit the password"))
	}).Methods("GET")

}
