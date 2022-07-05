package routes

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/shayamvlmna/cab-booking-app/pkg/controllers"
)

func AdminRoutes(r *mux.Router) {

	r.HandleFunc("/admin", controllers.AdminIndex)

	adminRouter := r.PathPrefix("/admin").Subrouter()

	adminRouter.HandleFunc("/create", controllers.CreateAdmin).Methods(http.MethodPost)

	adminRouter.HandleFunc("/login", controllers.AdminLogin)
	adminRouter.HandleFunc("/managedrivers", controllers.ManageDrivers)
	adminRouter.HandleFunc("/approve", controllers.ApproveDriver).Methods(http.MethodPost)
	
	adminRouter.HandleFunc("/manageusers", controllers.ManageUsers).Methods(http.MethodGet)
	adminRouter.HandleFunc("/driverequst", controllers.DriveRequest)
	adminRouter.HandleFunc("/blockdriver", controllers.BlockDriver)
	adminRouter.HandleFunc("/blockuser", controllers.BlockUser)
}
