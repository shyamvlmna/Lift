package routes

import (
	"github.com/gorilla/mux"
	"github.com/shayamvlmna/cab-booking-app/pkg/controllers"
)

func Router() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/", controllers.Index)
	go func() {
		UserRoutes(r)
		DriverRoutes(r)
		AdminRoutes(r)
	}()
	return r
}
