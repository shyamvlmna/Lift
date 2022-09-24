package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/shayamvlmna/lift/pkg/database"
	"github.com/shayamvlmna/lift/pkg/routes"
)

func main() {
	go database.DBSet()

	r := routes.Router()

	fmt.Println("server is up and running")
	log.Fatal(http.ListenAndServe(":8080", r))
}
