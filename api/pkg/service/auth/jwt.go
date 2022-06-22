package auth

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"github.com/joho/godotenv"
)

type Claims struct {
	Usrphone string `json:"phonenumber"`
	Role     string
	jwt.RegisteredClaims
}

//create jwt token with claims: role,phonenumber
func GenerateJWT(role, usrphone string) (string, error) {

	godotenv.Load()
	key := []byte(os.Getenv("SECRET_KEY"))

	claims := Claims{
		Usrphone: usrphone,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(time.Minute * 25)},
		},
	}
	
	token := jwt.NewWithClaims(&jwt.SigningMethodHMAC{}, claims)

	tokenstring, err := token.SignedString(key)
	if err != nil {
		return "", err
	}
	return tokenstring, nil
}

//returns role,phonenumber and error if any
func ValidateJWT(r *http.Request) (string, string, error) {
	godotenv.Load()
	key := []byte(os.Getenv("jwtSecretKey"))

	c, err := r.Cookie("jwt-token")
	if err == http.ErrNoCookie {
		fmt.Println("no cookie")
		return "", "", err
	}
	tokenstring := c.Value

	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(tokenstring, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})

	if err != nil {
		return "", "", err
	}

	if !tkn.Valid {
		fmt.Println("invalid token")
		return "", "", errors.New("invalidToken")
	}
	return claims.Role, claims.Usrphone, nil
}
