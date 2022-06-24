package middleware

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

		fmt.Println("token from resp", r.Header.Get("Token"))
		// r.Header["Token"]
		c,err:=r.Cookie("jwt-token")

		if err == nil {

			godotenv.Load()
			key := []byte(os.Getenv("SECRET_KEY"))

			claims := &auth.Claims{}

			token, err := jwt.ParseWithClaims(c.Value, claims, func(token *jwt.Token) (interface{}, error) {

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
