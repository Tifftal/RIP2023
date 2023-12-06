package app

import (
	"MSRM/internal/app/pkg"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

const prefix = "Bearer"

func (a *Application) RoleMiddleware(aloowedRoles ...pkg.Roles) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := extractTokenFromHaeder(c.Request)
		fmt.Println(c.Request)
		if tokenString == "" {
			fmt.Println("HERE1")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte("SuperSecretKey"), nil
		})

		if err != nil || !token.Valid {
			fmt.Println("HERE2")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			fmt.Println("HERE3")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		userID, ok := claims["user_Id"].(float64)
		if !ok {
			fmt.Println("HERE4")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		c.Set("User_id", int(userID))

		userRole, ok := claims["role"].(string)
		if !ok {
			fmt.Println("HERE5")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		if !isRoleAllowed(string(userRole), aloowedRoles) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Доступ отказан"})
			return
		}

		c.Next()
	}
}

func extractTokenFromHaeder(req *http.Request) string {
	bearerToken := req.Header.Get("Authorization")
	if bearerToken == "" {
		return ""
	}

	if strings.Split(bearerToken, " ")[0] != prefix {
		return ""
	}

	return strings.Split(bearerToken, " ")[1]
}

func isRoleAllowed(userRole string, allowedRoles []pkg.Roles) bool {
	for _, allowedRoles := range allowedRoles {
		if userRole == string(allowedRoles) {
			return true
		}
	}
	return false
}
