package routes

import (
	"github.com/gorilla/mux"
	"github.com/shayamvlmna/cab-booking-app/pkg/controllers"
)

func UserRoutes(r *mux.Router) {
	userRouter := r.PathPrefix("/user").Subrouter()

	// userRouter.HandleFunc("/signup", controllers.EnterNumber).Methods("GET")
	// userRouter.HandleFunc("/login", controllers.SearchNumber).Methods("GET")

	userRouter.HandleFunc("/signup", controllers.UserSignUp).Methods("POST")
	userRouter.HandleFunc("/login", controllers.UserLogin).Methods("POST")

}
