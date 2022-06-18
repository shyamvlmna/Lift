package controllers

import "net/http"

func AdminIndex(w http.ResponseWriter, r *http.Request) {
	adminTemp.ExecuteTemplate(w, "adminLoginForm.html", nil)
}
func AdminLogin(w http.ResponseWriter, r *http.Request) {

}
func Managedrivers(w http.ResponseWriter, r *http.Request) {
	adminTemp.ExecuteTemplate(w, "viewDrivers.html", nil)
}
func ManageUsers(w http.ResponseWriter, r *http.Request) {
	adminTemp.ExecuteTemplate(w, "viewUsers.html", nil)
}
func DriveRequest(w http.ResponseWriter, r *http.Request) {
	adminTemp.ExecuteTemplate(w, "driverRequests.html", nil)
}
func ApproveDriver(w http.ResponseWriter, r *http.Request) {

}
func BlockDriver(w http.ResponseWriter, r *http.Request) {

}
func BlockUser(w http.ResponseWriter, r *http.Request) {

}
