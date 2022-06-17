package auth

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
)

type Claims struct {
	Usrphone string `json:"userphonenumber"`
	jwt.RegisteredClaims
}

func GenerateJWT(usrphone string) (string, error) {

	godotenv.Load()
	key := []byte(os.Getenv("jwtSecretKey"))

	claims := Claims{
		Usrphone: usrphone,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(time.Minute * 25)},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenstring, err := token.SignedString(key)
	if err != nil {
		return "", err
	}
	return tokenstring, nil
}

func ValidateJWT(tokenstring string) (string, error) {
	godotenv.Load()
	key := []byte(os.Getenv("jwtSecretKey"))

	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(tokenstring, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})

	if err != nil {
		return "", err
	}

	if !tkn.Valid {
		return "", errors.New("http.StatusUnauthorized")
	}
	return claims.Usrphone, nil
}
