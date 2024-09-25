package utils

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

var jwtSecret = []byte("secret_key")

func GenerateJWT(email, role string, id int) (string, error) {

	claims := &jwt.MapClaims{
		"ExpiresAt": time.Now().Add(time.Hour * 72).Unix(),
		"Email":     email,
		"Role":      role,
		"ID":        id,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtSecret)
}

type JWTDecode struct {
	ExpiresAt int64
	Email     string
	Role      string
	ID        int
}

func DecodeJWT(tokenString string) (*JWTDecode, error) {

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected singing method: %v", t.Header["alg"])
		}

		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	decodedJWT := &JWTDecode{}
	if !ok {
		return decodedJWT, fmt.Errorf("error while taking claims in jwt")
	}
	decodedJWT.Email = claims["Email"].(string)
	decodedJWT.ExpiresAt = int64(claims["ExpiresAt"].(float64))
	decodedJWT.Role = claims["Role"].(string)
	decodedId, ok := claims["ID"].(float64)
	if !ok {
		return decodedJWT, fmt.Errorf("id could not be exracted")
	}
	decodedJWT.ID = int(decodedId)

	return decodedJWT, nil

}
