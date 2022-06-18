package routes

import (
	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	r := mux.NewRouter()
	go func() {
		UserRoutes(r)
		DriverRoutes(r)
		AdminRoutes(r)
	}()
	return r
}
