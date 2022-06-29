package models

type Response struct {
	ResponseStatus  string      `json:"responseStatus"`
	ResponseMessage string      `json:"responseMessage"`
	ResponseData    interface{} `json:"responseData"`
}
