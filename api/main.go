package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/shayamvlmna/cab-booking-app/pkg/database"
	"github.com/shayamvlmna/cab-booking-app/pkg/routes"
)

func main() {
	go database.DBSet()

	r := routes.Router()

	fmt.Println("server is up and running")
	log.Fatal(http.ListenAndServe(":8080", r))
}
