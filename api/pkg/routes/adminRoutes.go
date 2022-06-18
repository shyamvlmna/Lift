package routes

import (
	"github.com/gorilla/mux"
	"github.com/shayamvlmna/cab-booking-app/pkg/controllers"
)

func AdminRoutes(r *mux.Router) {

	r.HandleFunc("/", controllers.Index)

	r.HandleFunc("/admin", controllers.AdminIndex)
	adminRouter := r.PathPrefix("/admin").Subrouter()
	// adminRouter.HandleFunc("/", controllers.AdminIndex)

	adminRouter.HandleFunc("/login", controllers.AdminLogin)
	adminRouter.HandleFunc("/managedrivers", controllers.Managedrivers)
	adminRouter.HandleFunc("/manageusers", controllers.ManageUsers)
	adminRouter.HandleFunc("/driverequst", controllers.DriveRequest)
	adminRouter.HandleFunc("/blockdriver", controllers.BlockDriver)
	adminRouter.HandleFunc("/blockuser", controllers.BlockUser)
}
