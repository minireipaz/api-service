package middlewares

import (
	"minireipaz/pkg/domain/models"
	"minireipaz/pkg/domain/repos"

	// "minireipaz/pkg/domain/services"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(authService *repos.AuthService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.ContentType() != "application/json" {
			ctx.JSON(http.StatusUnsupportedMediaType, NewUnsupportedMediaTypeError("Only application/json is supported"))
			ctx.Abort()
			return
		}

		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, NewUnauthorizedError(models.AuthInvalid))
			ctx.Abort()
			return
		}

		accessToken := strings.TrimPrefix(authHeader, "Bearer ")
		if accessToken == "" {
			ctx.JSON(http.StatusUnauthorized, NewUnauthorizedError(models.AuthInvalid))
			ctx.Abort()
			return
		}

		valid, err := verifyServiceUserToken(*authService, accessToken)
		if err != nil || !valid {
			ctx.JSON(http.StatusUnauthorized, NewUnauthorizedError(models.AuthInvalid))
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

// TODO: can be expired btw rightnow not rotated
func verifyServiceUserToken(authService repos.AuthService, token string) (bool, error) {
	isValid, err := authService.VerifyServiceUserToken(token)
	if err != nil {
		return false, err
	}
	return isValid, nil
}
