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

// func main() {
// 	go database.DBSet()

// 	r := routes.Router()

// 	cfg := &tls.Config{
// 		MinVersion:               tls.VersionTLS12,
// 		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
// 		PreferServerCipherSuites: true,
// 		CipherSuites: []uint16{
// 			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
// 			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
// 			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
// 			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
// 		},
// 	}
// 	srv := &http.Server{
// 		Addr:         ":8080",
// 		Handler:      r,
// 		TLSConfig:    cfg,
// 		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
// 	}
// 	log.Fatal(srv.ListenAndServeTLS("tls.crt", "tls.key"))
// }
