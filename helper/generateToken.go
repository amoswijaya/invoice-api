package helper

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)
func GenerateToken(userID uint) (string, error) {
    secret := os.Getenv("JWT_SECRET")
    claims := jwt.MapClaims{
        "userID": userID,
        "exp":    time.Now().Add(24 * time.Hour).Unix(),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(secret))
}