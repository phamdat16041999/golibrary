package jwtp

import (
	"fmt"

	"github.com/golang-jwt/jwt"
)

func Decrypt() {
	var hmacSampleSecret []byte
	tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJyb2xlIjoiQWRtaW4iLCJ1c2VybmFtZSI6IlBoYW1EYXQifQ.aa6QAK0048-REKz_F9rcLuhnR267sejrL6O2GelX8io"
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return hmacSampleSecret, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims["username"], claims["role"])
	} else {
		fmt.Println(err)
	}

}
