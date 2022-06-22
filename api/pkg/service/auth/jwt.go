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

// //returns role,phonenumber and error if any
// func ValidateJWT(r *http.Request) (string, error) {
// 	godotenv.Load()
// 	key := []byte(os.Getenv("jwtSecretKey"))

// 	// c, err := r.Cookie("jwt-token")
// 	// if err == http.ErrNoCookie {
// 	// 	fmt.Println("no cookie")
// 	// 	return "", "", err
// 	// }
// 	// tokenstring := c.Value
// 	// tokenstring
// 	claims := &Claims{}
// 	token, err := jwt.ParseWithClaims(r.Header["Bearer"][0], claims, func(token *jwt.Token) (interface{}, error) {
// 		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 			return "", fmt.Errorf(("Invalid Signing Method"))
// 		}

// 		if _, ok := token.Claims.(jwt.RegisteredClaims); !ok && !token.Valid {
// 			return "", fmt.Errorf(("Expired token"))
// 		}

// 		return key, nil
// 	})

// 	if err != nil {
// 		return "", err
// 	}

// 	if !token.Valid {
// 		return "", fmt.Errorf("invalidToken")
// 	}
// 	return claims.Role, claims.Usrphone, nil
// }
