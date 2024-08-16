package controllers

import (
	"log"
	"minireipaz/pkg/auth"
	"minireipaz/pkg/config"
	"minireipaz/pkg/domain/services"
	"minireipaz/pkg/infra/httpclient"
	"minireipaz/pkg/infra/redisclient"
	"minireipaz/pkg/infra/tokenrepo"
	"sync"

	"github.com/gin-gonic/gin"
)

// type AuthController struct {
// }

type AuthController struct {
	authService *services.AuthService
	// authController *AuthController
	once   sync.Once
	config config.ZitadelConfig
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
			ac.config.GetZitadelKeyUserID(),
			ac.config.GetZitadelKeyPrivate(),
			ac.config.GetZitadelKeyID(),
		)

		jwtGenerator := auth.NewJWTGenerator(
			ac.config.GetZitadelKeyUserID(),
			ac.config.GetZitadelKeyPrivate(),
			ac.config.GetZitadelKeyID(),
			ac.config.GetZitadelURI(),
		)
		redisClient := redisclient.NewRedisClient()
		tokenRepo := tokenrepo.NewTokenRepository(redisClient)
		authService := services.NewAuthService(jwtGenerator, zitadelClient, tokenRepo)

		_, err := authService.GetAccessToken()
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
	log.Printf("%v", userToken)
	return
}
