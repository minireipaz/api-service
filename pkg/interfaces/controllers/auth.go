package controllers

import (
	"minireipaz/pkg/auth"
	"minireipaz/pkg/config"
	"minireipaz/pkg/domain/models"
	"minireipaz/pkg/domain/repos"
	"minireipaz/pkg/domain/services"
	"minireipaz/pkg/infra/httpclient"
	"minireipaz/pkg/infra/redisclient"
	"minireipaz/pkg/infra/tokenrepo"
	"minireipaz/pkg/interfaces/middlewares"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService repos.AuthService
	once        sync.Once
	config      config.ZitadelConfig
}

func NewAuthContext(cfg config.ZitadelConfig) *AuthController {
	return &AuthController{
		config: cfg,
	}
}

func (ac *AuthController) GetAuthController() *AuthController {
	ac.once.Do(func() {
		zitadelClient := httpclient.NewZitadelClient(
			ac.config.GetZitadelURI(),
			ac.config.GetZitadelServiceUserID(),
			ac.config.GetZitadelServiceUserKeyPrivate(),
			ac.config.GetZitadelServiceUserKeyID(),
			ac.config.GetZitadelProjectID(),
			ac.config.GetZitadelKeyClientID(),
		)

		jwtGenerator := auth.NewJWTGenerator(auth.JWTGeneratorConfig{
			ServiceUser: auth.ServiceUserConfig{
				UserID:     ac.config.GetZitadelServiceUserID(),
				PrivateKey: []byte(ac.config.GetZitadelServiceUserKeyPrivate()),
				KeyID:      ac.config.GetZitadelServiceUserKeyID(),
				ClientID:   ac.config.GetZitadelServiceUserClientID(),
			},
			BackendApp: auth.BackendAppConfig{
				AppID:      ac.config.GetZitadelBackendID(),
				PrivateKey: []byte(ac.config.GetZitadelBackendKeyPrivate()),
				KeyID:      ac.config.GetZitadelBackendKeyID(),
				ClientID:   ac.config.GetZitadelBackendClientID(),
			},
			APIURL:    ac.config.GetZitadelURI(),
			ProjectID: ac.config.GetZitadelProjectID(),
			ClientID:  ac.config.GetZitadelKeyClientID(),
		})
		redisClient := redisclient.NewRedisClient()
		tokenRepo := tokenrepo.NewTokenRepository(redisClient)
		authService := services.NewAuthService(jwtGenerator, zitadelClient, tokenRepo)

		// get cached accestoken for service user
		cachedToken := authService.GetCachedServiceUserAccessToken()
		// in dev state, not rotating service user acces token in servesless functions
		if ac.config.GetEnv("ROTATE_SERVICE_USER_TOKEN", "n") == "y" {
			if cachedToken == nil {
				// Rotate token if it's expired or not found
				_, err := authService.GenerateAccessToken()
				if err != nil { // error saving retry read
					_ = authService.GetCachedServiceUserAccessToken()
				}
			}
		}
		ac.authService = authService
	})
	return ac
}

func (ac *AuthController) GetAuthService() repos.AuthService {
	return ac.GetAuthController().authService
}

func (ac *AuthController) VerifyUserToken(ctx *gin.Context) {
	userToken := ctx.Param("usertoken")
	isValid, isExpired := ac.authService.VerifyUserToken(userToken) // right not controlled to rotate/expire user token

	if !isValid {
		ctx.JSON(http.StatusUnauthorized, middlewares.NewUnauthorizedError(models.AuthInvalid))
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"valid":   isValid,
		"expired": isExpired,
		"error":   "",
	})
}

func (ac *AuthController) VerifyUserTokenForMiddleware(ctx *gin.Context) {
	userToken := ctx.Param("usertoken")
	isValid, isExpired := ac.authService.VerifyUserToken(userToken) // can be rotated

	if !isValid {
		ctx.JSON(http.StatusUnauthorized, middlewares.NewUnauthorizedError(models.AuthInvalid))
		ctx.Abort()
		return
	}

	if isExpired {
		ctx.JSON(http.StatusUnauthorized, middlewares.NewUnauthorizedError(models.UserTokenExpired))
		ctx.Abort()
		return
	}

	ctx.Next()
}
