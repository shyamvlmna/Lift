package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/shayamvlmna/cab-booking-app/app/routes"
)

func main() {
	fmt.Println("server is up and running")

	r := routes.Router()

	log.Fatal(http.ListenAndServe(":8080", r))
}
