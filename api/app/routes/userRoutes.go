package routes

import (
	"github.com/gorilla/mux"
	"github.com/shayamvlmna/cab-booking-app/app/controllers"
)

func UserRoutes(r *mux.Router) {
	userRouter := r.PathPrefix("/user").Subrouter()

	// userRouter.HandleFunc("/signup", controllers.EnterNumber).Methods("GET")
	// userRouter.HandleFunc("/login", controllers.SearchNumber).Methods("GET")

	userRouter.HandleFunc("/auth", controllers.UserAuth).Methods("POST")

	userRouter.HandleFunc("/signup", controllers.UserSignUp).Methods("POST")
	userRouter.HandleFunc("/login", controllers.UserLogin).Methods("POST")

}
