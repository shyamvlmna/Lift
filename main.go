package main

import (
	"fmt"
	"log"
	"net/http"

	routes "github.com/shayamvlmna/cab-booking-app/pkg/routes"
)

func main() {
	fmt.Println("lets Go...")

	r := routes.Router()

	log.Fatal(http.ListenAndServe(":8080", r))
}
