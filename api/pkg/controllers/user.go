package controllers

import "net/http"

func Jwt(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("secret information"))
}
