package routes

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/shayamvlmna/cab-booking-app/pkg/controllers"
	"github.com/shayamvlmna/cab-booking-app/pkg/middleware"
	"github.com/shayamvlmna/cab-booking-app/pkg/models"
	"github.com/shayamvlmna/cab-booking-app/pkg/service/auth"
)

func UserRoutes(r *mux.Router) {

	r.HandleFunc("/", controllers.Index)

	userRouter := r.PathPrefix("/user").Subrouter()

	userRouter.HandleFunc("/auth", controllers.UserAuth).Methods("POST")
	userRouter.HandleFunc("/signup", controllers.UserSignUp).Methods("POST")
	userRouter.HandleFunc("/login", controllers.UserLogin).Methods("POST")
	userRouter.HandleFunc("/googlelogin", auth.GoogleLogin)
	userRouter.HandleFunc("/googleCallback", auth.GoogleCallback)
	userRouter.HandleFunc("/logout", controllers.UserLogout)

	userRouter.Handle("/otp", middleware.ValidateOtp(controllers.UserSignupPage)).Methods("POST")
	userRouter.Handle("/userhome", middleware.IsAuthorized(controllers.UserHome)).Methods("GET")

	userRouter.HandleFunc("/update/{id}", controllers.EditUserProfile).Methods("GET")
	userRouter.HandleFunc("/update/{id}", controllers.UpdateUserProfile).Methods("POST")

	//render enter otp page
	userRouter.HandleFunc("/enterotp", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		response := models.Response{
			ResponseStatus:  "success",
			ResponseMessage: "new user",
			ResponseData:    nil,
		}
		json.NewEncoder(w).Encode(&response)

	}).Methods("GET")

	//enter login page to enter password
	userRouter.HandleFunc("/loginpage", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		response := models.Response{
			ResponseStatus:  "success",
			ResponseMessage: "existing user",
			ResponseData:    nil,
		}
		json.NewEncoder(w).Encode(&response)
	}).Methods("GET")

	userRouter.Handle("/booktrip", middleware.IsAuthorized(controllers.BookTrip))

}
