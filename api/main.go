package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/shayamvlmna/cab-booking-app/pkg/routes"
)

func main() {
	fmt.Println("server is up and running")

	r := routes.Router()
	// database.OpenTripDb()
	// database.OpenWalletDb()
	// database.OpenLocationDb()
	log.Fatal(http.ListenAndServe(":8080", r))
}
