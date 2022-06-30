package models

type Response struct {
	ResponseStatus  string `json:"responseStatus"`
	ResponseMessage string `json:"responseMessage"`
	ResponseData    any    `json:"responseData"`
}

type UserData struct {
	Id          uint64 `json:"id"`
	Phonenumber string `json:"phonenumber"`
	Firstname   string `json:"firstname"`
	Lastname    string `json:"lastname"`
	Email       string `json:"email"`
}
