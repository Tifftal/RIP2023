package pkg

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

func GenerateToken(userID uint, role string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["user_Id"] = userID
	claims["role"] = role
	claims["exp"] = time.Hour * 1

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_KEY")))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
