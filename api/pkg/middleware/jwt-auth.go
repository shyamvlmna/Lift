package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/shayamvlmna/cab-booking-app/pkg/models"

	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
	"github.com/shayamvlmna/cab-booking-app/pkg/service/auth"
)

func IsAuthorized(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		reqToken := r.Header.Get("Authorization")

		// c, err := r.Cookie("jwt-token")

		if reqToken != "" {
			splitToken := strings.Split(reqToken, "Bearer ")
			tokenString := splitToken[1]
			err := godotenv.Load()
			if err != nil {
				return
			}
			key := []byte(os.Getenv("JWT-SECRET_KEY"))

			claims := &auth.Claims{}

			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {

				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return "", fmt.Errorf(("invalid signing method"))
				}

				if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
					return "", fmt.Errorf(("expired token"))
				}

				return key, nil
			})

			if err != nil {
				response := &models.Response{
					ResponseStatus:  "failed",
					ResponseMessage: "No Authorization Token provided",
					ResponseData:    err.Error(),
				}
				err := json.NewEncoder(w).Encode(&response)
				if err != nil {
					return
				}
				return
			}

			if token.Valid {
				endpoint(w, r)
			}
		} else {
			response := &models.Response{
				ResponseStatus:  "failed",
				ResponseMessage: "No Authorization Token provided",
				ResponseData:    nil,
			}
			err := json.NewEncoder(w).Encode(&response)
			if err != nil {
				return
			}
			return
		}
	})
}
