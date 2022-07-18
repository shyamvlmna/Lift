package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-redis/redis/v8"

	"github.com/shayamvlmna/cab-booking-app/pkg/models"
	"github.com/shayamvlmna/cab-booking-app/pkg/service/auth"
)

type Otp struct {
	Otp string `json:"otp"`
}

func ValidateOtp(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		otp := &Otp{}

		err := json.NewDecoder(r.Body).Decode(&otp)
		if err != nil {
			return
		}

		OTP := otp.Otp

		phone := auth.GetPhone()

		auth.StorePhone(phone)

		if err := auth.ValidateOTP(phone, OTP); err != nil {
			if err == redis.Nil {
				fmt.Println("otp expired")
				w.WriteHeader(http.StatusUnauthorized)
				response := models.Response{
					ResponseStatus:  "fail",
					ResponseMessage: "otp expired",
					ResponseData:    nil,
				}
				err := json.NewEncoder(w).Encode(&response)
				if err != nil {
					return
				}
				return
			} else {
				fmt.Println("invalid otp")
				w.WriteHeader(http.StatusUnauthorized)
				response := models.Response{
					ResponseStatus:  "fail",
					ResponseMessage: "invalid otp",
					ResponseData:    nil,
				}
				err := json.NewEncoder(w).Encode(&response)
				if err != nil {
					return
				}
				return
			}

		} else {
			fmt.Println("otp validation success")
			endpoint(w, r)
		}
	})
}
