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
)

type AuthController struct {
	authService *services.AuthService
}

type AuthContext struct {
	authController *AuthController
	once           sync.Once
	config         config.Config
}

func NewAuthContext(cfg config.Config) *AuthContext {
	return &AuthContext{
		config: cfg,
	}
}

func (ac *AuthContext) GetAuthController() *AuthController {
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

		ac.authController = &AuthController{authService: authService}
	})
	return ac.authController
}

func (ac *AuthContext) GetAuthService() *services.AuthService {
	return ac.GetAuthController().authService
}
