package middlewares

import (
	"github.com/gin-gonic/gin"
	"minireipaz/pkg/domain/services"
	"net/http"
	"strings"
)

func VerifyServiceUserToken(authService *services.AuthService, token string) (bool, error) {
	isValid, err := authService.VerifyToken(token)
	if err != nil {
		return false, err
	}
	return isValid, nil
}

func AuthMiddleware(authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
			c.Abort()
			return
		}

		token := parts[1]

		valid, err := VerifyServiceUserToken(authService, token)
		if err != nil || !valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

    // TODO: Client Access Token

		c.Next()
	}
}
