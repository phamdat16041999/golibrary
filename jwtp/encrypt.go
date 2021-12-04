package jwtp

import (
	"fmt"

	"github.com/golang-jwt/jwt"
)

func Encrypt() {
	var hmacSampleSecret []byte
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": "PhamDat",
		"role":     "Admin",
	})
	tokenString, err := token.SignedString(hmacSampleSecret)
	fmt.Println(tokenString, err)
}
