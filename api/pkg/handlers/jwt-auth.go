package handlers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
	"github.com/shayamvlmna/cab-booking-app/pkg/service/auth"
)

func IsAuthorized(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Header["Token"] != nil {

			godotenv.Load()
			key := []byte(os.Getenv("SECRET_KEY"))

			claims := &auth.Claims{}

			token, err := jwt.ParseWithClaims(r.Header["Token"][0], claims, func(token *jwt.Token) (interface{}, error) {

				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return "", fmt.Errorf(("invalid signing method"))
				}

				if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
					return "", fmt.Errorf(("expired token"))
				}

				return key, nil
			})

			if err != nil {
				fmt.Fprint(w, err.Error())
			}

			if token.Valid {
				endpoint(w, r)
			}
		} else {
			fmt.Fprintf(w, "No Authorization Token provided")
		}
	})
}
