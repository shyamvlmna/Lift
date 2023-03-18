package main

import (
	"fmt"
	"github.com/shayamvlmna/lift/internal/database"
	"github.com/shayamvlmna/lift/pkg/routes"
	"log"
	"net/http"
)

func main() {
	go database.DBSet()

	r := routes.Router()

	fmt.Println("server is up and running")
	log.Fatal(http.ListenAndServe(":8080", r))
}
