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
		tokenString := extractTokenFromHeader(c.Request)
		if tokenString == "" {
			fmt.Println("1")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte("SuperSecretKey"), nil
		})

		if err != nil || !token.Valid {
			fmt.Println("2")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			fmt.Println("3")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		userID, ok := claims["user_Id"].(float64)
		if !ok {
			fmt.Println("4")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		c.Set("User_id", int(userID))

		userRole, ok := claims["role"].(string)
		if !ok {
			fmt.Println("5")
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
func (a *Application) Guest(allowedRoles ...pkg.Roles) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := extractTokenFromHeader(c.Request)
		if tokenString == "" {
			c.Set("User_id", int(0))
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte("SuperSecretKey"), nil
		})

		if err != nil || !token.Valid {
			c.Set("User_id", int(0))
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.Set("User_id", int(0))
			return
		}

		userID, ok := claims["user_Id"].(float64)
		if !ok {
			c.Set("User_id", int(0))
			return
		}

		c.Set("User_id", int(userID))

		userRole, ok := claims["role"].(string)
		if !ok {
			c.Set("User_id", int(0))
			return
		}

		if !isRoleAllowed(string(userRole), allowedRoles) {
			c.Set("User_id", int(0))
			return
		}

		c.Next()
	}
}

// func extractTokenFromHeader(request *http.Request) {
// 	panic("unimplemented")
// }

func extractTokenFromHeader(req *http.Request) string {
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
