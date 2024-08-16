package middlewares

import (
	"minireipaz/pkg/domain/models"
	"minireipaz/pkg/domain/services"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func VerifyServiceUserToken(authService *services.AuthService, token string) (bool, error) {
	isValid, err := authService.VerifyServiceUserToken(token)
	if err != nil {
		return false, err
	}
	return isValid, nil
}

func AuthMiddleware(authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.ContentType() != "application/json" {
			c.JSON(http.StatusUnsupportedMediaType, NewUnsupportedMediaTypeError("Only application/json is supported"))
			c.Abort()
			return
		}

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, NewUnauthorizedError(models.AuthInvalid))
			c.Abort()
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			c.JSON(http.StatusUnauthorized, NewUnauthorizedError(models.AuthInvalid))
			c.Abort()
			return
		}

		valid, err := VerifyServiceUserToken(authService, token)
		if err != nil || !valid {
			c.JSON(http.StatusUnauthorized, NewUnauthorizedError(models.AuthInvalid))
			c.Abort()
			return
		}

		c.Next()
	}
}
