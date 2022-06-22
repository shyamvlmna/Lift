package routes

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/shayamvlmna/cab-booking-app/pkg/controllers"
	"github.com/shayamvlmna/cab-booking-app/pkg/handlers"
	"github.com/shayamvlmna/cab-booking-app/pkg/models"
	"github.com/shayamvlmna/cab-booking-app/pkg/service/auth"
)

func UserRoutes(r *mux.Router) {

	r.HandleFunc("/", controllers.Index)

	userRouter := r.PathPrefix("/user").Subrouter()

	userRouter.HandleFunc("/auth", controllers.UserAuth).Methods("POST")
	userRouter.HandleFunc("/otp", controllers.ValidateOtp).Methods("POST")
	userRouter.HandleFunc("/signup", controllers.UserSignUp).Methods("POST")
	userRouter.HandleFunc("/login", controllers.UserLogin).Methods("POST")
	userRouter.HandleFunc("/googlelogin", auth.GoogleLogin)
	userRouter.HandleFunc("/googleCallback", auth.GoogleCallback)

	userRouter.Handle("/jwt", handlers.IsAuthorized(controllers.Jwt))

	userRouter.HandleFunc("/login", controllers.UserLogin).Methods("POST")

	userRouter.HandleFunc("/logout", controllers.UserLogout).Methods("DELETE")

	userRouter.HandleFunc("/userhome", controllers.UserHome).Methods("GET")
	userRouter.HandleFunc("/update/{id}", controllers.EditUserProfile).Methods("GET")
	userRouter.HandleFunc("/update/{id}", controllers.UpdateUserProfile).Methods("POST")
	userRouter.HandleFunc("/getusers", controllers.GetUsers).Methods("GET")

	//render enter otp page
	userRouter.HandleFunc("/enterotp", func(w http.ResponseWriter, r *http.Request) {
		response := models.Response{
			ResponseStatus:  "success",
			ResponseMessage: "validate otp",
			ResponseData:    nil,
		}
		json.NewEncoder(w).Encode(&response)
		w.Header().Set("Content-Type", "application/json")
	}).Methods("GET")

	//render signup page
	userRouter.HandleFunc("/signup", func(w http.ResponseWriter, r *http.Request) {
		response := models.Response{
			ResponseStatus:  "success",
			ResponseMessage: "new user",
			ResponseData:    nil,
		}
		json.NewEncoder(w).Encode(&response)
		w.Header().Set("Content-Type", "application/json")
	}).Methods("GET")

	//enter login page to enter password
	userRouter.HandleFunc("/loginpage", func(w http.ResponseWriter, r *http.Request) {
		response := models.Response{
			ResponseStatus:  "success",
			ResponseMessage: "existing user",
			ResponseData:    nil,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(&response)
	}).Methods("GET")

}
