package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/shayamvlmna/cab-booking-app/app/routes"
)

func main() {
	fmt.Println("lets Go...")

	r := routes.Router()

	log.Fatal(http.ListenAndServe(":8080", r))
}
