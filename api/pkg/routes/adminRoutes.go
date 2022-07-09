package routes

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/shayamvlmna/cab-booking-app/pkg/controllers"
	"github.com/shayamvlmna/cab-booking-app/pkg/middleware"
)

func AdminRoutes(r *mux.Router) {

	r.HandleFunc("/admin", controllers.AdminIndex)

	adminRouter := r.PathPrefix("/admin").Subrouter()

	adminRouter.HandleFunc("/create", controllers.CreateAdmin).Methods(http.MethodPost)

	adminRouter.HandleFunc("/login", controllers.AdminLogin).Methods(http.MethodPost)

	adminRouter.Handle("/adminhome", middleware.IsAuthorized(controllers.AdminHome)).Methods(http.MethodGet)

	//

	adminRouter.Handle("/managedrivers", middleware.IsAuthorized(controllers.ManageDrivers))

	adminRouter.Handle("/driverequst", middleware.IsAuthorized(controllers.DriveRequest)).Methods(http.MethodGet)

	adminRouter.Handle("/approvedriver", middleware.IsAuthorized(controllers.ApproveDriver)).Methods(http.MethodPost)

	adminRouter.Handle("/blockdriver", middleware.IsAuthorized(controllers.BlockDriver)).Methods(http.MethodPost)

	adminRouter.Handle("/payouts", middleware.IsAuthorized(controllers.PayoutRequests)).Methods(http.MethodGet)

	adminRouter.Handle("/updatepayout", middleware.IsAuthorized(controllers.UpdatePayout)).Methods(http.MethodPost)
	//

	adminRouter.Handle("/manageusers", middleware.IsAuthorized(controllers.ManageUsers)).Methods(http.MethodGet)

	adminRouter.Handle("/blockuser", middleware.IsAuthorized(controllers.BlockUser)).Methods(http.MethodPost)

	adminRouter.Handle("/addcoupon", middleware.IsAuthorized(controllers.AddCoupon)).Methods(http.MethodPost)
}
