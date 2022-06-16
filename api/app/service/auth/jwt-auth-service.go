package auth

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
)

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func GenerateJWT(username string) (string, error) {

	godotenv.Load()
	key := []byte(os.Getenv("jwtSecretKey"))

	claims := Claims{
		Username: username,
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

func IsAuthorized(tokenstring string) (bool, error) {
	godotenv.Load()
	key := []byte(os.Getenv("jwtSecretKey"))

	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(tokenstring, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})

	if err != nil {
		return false, err
	}

	if !tkn.Valid {
		return false, errors.New("http.StatusUnauthorized")
	}
	return true, nil
}
