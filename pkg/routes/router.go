package routes

import (
	"github.com/gorilla/mux"
	"github.com/shayamvlmna/cab-booking-app/pkg/controllers"
)

func Router() *mux.Router {
	r := mux.NewRouter()
	go func() {
		UserRoutes(r)
		DriverRoutes(r)
		AdminRoutes(r)
	}()

	r.HandleFunc("/", controllers.Index).Methods("GET")
	r.HandleFunc("/cab", controllers.Cab).Methods("GET")

	return r
}
