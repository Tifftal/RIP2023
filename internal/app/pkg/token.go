package pkg

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

func GenerateToken(userID uint, role string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	currentTime := time.Now()
	claims := token.Claims.(jwt.MapClaims)
	claims["user_Id"] = userID
	claims["role"] = role

	// Используйте Unix(), чтобы преобразовать время в метку Unix
	expirationTime := currentTime.Add(time.Hour * 4).Unix()
	claims["exp"] = expirationTime

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_KEY")))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
