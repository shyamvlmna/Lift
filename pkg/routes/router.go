package routes

import (
	"github.com/gorilla/mux"
	"github.com/shayamvlmna/lift/pkg/controllers"
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
