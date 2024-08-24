package controllers

import (
	"log"
	"minireipaz/pkg/auth"
	"minireipaz/pkg/config"
	"minireipaz/pkg/domain/models"
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
	authService *services.AuthService
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
			},
			BackendApp: auth.BackendAppConfig{
				AppID:      ac.config.GetZitadelBackendID(),
				PrivateKey: []byte(ac.config.GetZitadelBackendKeyPrivate()),
				KeyID:      ac.config.GetZitadelBackendKeyID(),
			},
			APIURL:    ac.config.GetZitadelURI(),
			ProjectID: ac.config.GetZitadelProjectID(),
			ClientID:  ac.config.GetZitadelKeyClientID(),
		})
		redisClient := redisclient.NewRedisClient()
		tokenRepo := tokenrepo.NewTokenRepository(redisClient)
		authService := services.NewAuthService(jwtGenerator, zitadelClient, tokenRepo)

		_, err := authService.GetServiceUserAccessToken()
		if err != nil {
			log.Panicf("ERROR | %v", err)
		}
		ac.authService = authService
		// ac.authController = &AuthController{authService: authService}
	})
	return ac
}

func (ac *AuthController) GetAuthService() *services.AuthService {
	return ac.GetAuthController().authService
}

func (ac *AuthController) VerifyUserToken(ctx *gin.Context) {
	userToken := ctx.Param("id")
	isValid := ac.authService.VerifyUserToken(userToken)

	if !isValid {
		ctx.JSON(http.StatusUnauthorized, middlewares.NewUnauthorizedError(models.AuthInvalid))
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"valid": isValid,
		"error": "",
	})
}
