package routes

import (
	"github.com/gorilla/mux"

	"github.com/shayamvlmna/cab-booking-app/pkg/controllers"
)

func UserRoutes(r *mux.Router) {

	userRouter := r.PathPrefix("/user").Subrouter()

	userRouter.HandleFunc("/auth", controllers.UserAuth).Methods("POST")
	userRouter.HandleFunc("/signup", controllers.UserSignUp).Methods("POST")
	userRouter.HandleFunc("/login", controllers.UserLogin).Methods("POST")

	userRouter.HandleFunc("/userhome", controllers.UserHome).Methods("GET")
	userRouter.HandleFunc("/update/{id}", controllers.EditUserProfile).Methods("GET")
	userRouter.HandleFunc("/update/{id}", controllers.UpdateUserProfile).Methods("POST")
	userRouter.HandleFunc("/logout", controllers.UserLogin).Methods("GET")

	userRouter.HandleFunc("/getusers", controllers.GetUsers).Methods("GET")
}
