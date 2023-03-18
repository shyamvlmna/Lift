package models

type Response struct {
	ResponseStatus  string `json:"responseStatus"`
	ResponseMessage string `json:"responseMessage"`
	ResponseData    any    `json:"responseData"`
}

type UserData struct {
	Id          uint   `json:"id"`
	Phonenumber string `json:"phonenumber"`
	Firstname   string `json:"firstname"`
	Lastname    string `json:"lastname"`
	Email       string `json:"email"`
}

type DriverData struct {
	Id          uint     `json:"id"`
	Phonenumber string   `json:"phonenumber"`
	Firstname   string   `json:"firstname"`
	Lastname    string   `json:"lastname"`
	Email       string   `json:"email"`
	City        string   `json:"city"`
	LicenceNum  string   `json:"licence"`
	Cab         *CabData `json:"cab"`
}

type CabData struct {
	VehicleId    uint64 `json:"vehicleid"`
	Registration string `json:"registration"`
	Brand        string `json:"brand"`
	Category     string `json:"type"`
	VehicleModel string `json:"model"`
	Colour       string `json:"colour"`
}
