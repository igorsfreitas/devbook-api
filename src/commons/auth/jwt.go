package auth

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/igorsfreitas/devbook-api/src/config"
)

// GenerateJWT generates a JWT token
func GenerateJWT(userID uint64) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["exp"] = time.Now().Add(time.Hour * 6).Unix() // 6 hours
	claims["userID"] = userID

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.SecretKey))
}
