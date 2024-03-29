package routes

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/shayamvlmna/lift/pkg/controllers"
	"github.com/shayamvlmna/lift/pkg/middleware"
)

func AdminRoutes(r *mux.Router) {

	r.HandleFunc("/admin", controllers.AdminIndex).Methods(http.MethodGet)

	adminRouter := r.PathPrefix("/admin").Subrouter()

	adminRouter.HandleFunc("/create", controllers.CreateAdmin).Methods(http.MethodPost)

	adminRouter.HandleFunc("/login", controllers.AdminLogin).Methods(http.MethodPost)

	adminRouter.Handle("/adminhome", middleware.IsAuthorized(controllers.AdminHome)).Methods(http.MethodGet)

	// Manage Drivers

	adminRouter.Handle("/managedrivers", middleware.IsAuthorized(controllers.ManageDrivers)).Methods(http.MethodGet)

	adminRouter.Handle("/driverequst", middleware.IsAuthorized(controllers.DriveRequest)).Methods(http.MethodGet)

	adminRouter.Handle("/approvedriver", middleware.IsAuthorized(controllers.ApproveDriver)).Methods(http.MethodPost)

	adminRouter.Handle("/blockdriver", middleware.IsAuthorized(controllers.BlockDriver)).Methods(http.MethodPost)

	adminRouter.Handle("/unblockdriver", middleware.IsAuthorized(controllers.UnBlockDriver)).Methods(http.MethodPost)

	adminRouter.Handle("/payouts", middleware.IsAuthorized(controllers.PayoutRequests)).Methods(http.MethodGet)

	adminRouter.Handle("/updatepayout", middleware.IsAuthorized(controllers.UpdatePayout)).Methods(http.MethodPost)

	// Manage Users

	adminRouter.Handle("/manageusers", middleware.IsAuthorized(controllers.ManageUsers)).Methods(http.MethodGet)

	adminRouter.Handle("/blockuser", middleware.IsAuthorized(controllers.BlockUser)).Methods(http.MethodPost)

	adminRouter.Handle("/unblockuser", middleware.IsAuthorized(controllers.UnBlockUser)).Methods(http.MethodPost)

	//

	adminRouter.Handle("/addcoupon", middleware.IsAuthorized(controllers.CreateCoupon)).Methods(http.MethodPost)
}
