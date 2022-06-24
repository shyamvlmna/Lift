package auth

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"github.com/joho/godotenv"
)

type Claims struct {
	Usrphone   string `json:"phonenumber"`
	Role       string
	Authorized bool
	jwt.RegisteredClaims
}

//create jwt token with claims: role,phonenumber
func GenerateJWT(role, usrphone string) (string, error) {

	godotenv.Load()
	key := []byte(os.Getenv("SECRET_KEY"))

	claims := Claims{
		Usrphone:   usrphone,
		Role:       role,
		Authorized: true,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(time.Minute * 25)},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(key)
	if err != nil {
		return "", err
	}
	fmt.Print(string(tokenString))
	return tokenString, nil
}

//parse the given JWT token string and returns the role and phonenumber
func ParseJWT(tokenString string) (string, string) {
	claims := &Claims{}

	godotenv.Load()
	key := []byte(os.Getenv("SECRET_KEY"))

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
		fmt.Println(err)
	}

	if !token.Valid {
		fmt.Println("invalid token")
		return "", ""
	}

	return claims.Role, claims.Usrphone

}
