package services

import (
	"fmt"
	"log"
	"minireipaz/pkg/config"
	"minireipaz/pkg/domain/models"
	"minireipaz/pkg/domain/repos"
	"time"
)

type AuthServiceImpl struct {
	jwtGenerator  repos.JWTGenerator
	zitadelClient repos.ZitadelClient
	tokenRepo     repos.TokenRepository
}

func NewAuthService(jwtGenerator repos.JWTGenerator, zitadelClient repos.ZitadelClient, tokenRepo repos.TokenRepository) repos.AuthService {
	return &AuthServiceImpl{
		jwtGenerator:  jwtGenerator,
		zitadelClient: zitadelClient,
		tokenRepo:     tokenRepo,
	}
}

func (s *AuthServiceImpl) GenerateAccessToken() (*string, error) {
	assertionJWT, err := s.jwtGenerator.GenerateServiceUserAssertionJWT(time.Hour)
	if err != nil {
		log.Panicf("ERROR | Cannot generate JWT %v", err)
	}

	accessToken, expiresIn, err := s.zitadelClient.GenerateServiceUserAccessToken(assertionJWT)
	if err != nil {
		log.Printf("ERROR | Cannot acces to ACCESS token %v", err)
		return nil, fmt.Errorf("ERROR | Cannot acces to ACCESS token %v", err)
	}

	err = s.tokenRepo.SaveToken(accessToken, &expiresIn)
	if err != nil {
		log.Printf("ERROR | Failed to save token, %v", err)
		return nil, fmt.Errorf("ERROR | Failed to save token, %v", err)
	}

	return accessToken, nil
}

func (s *AuthServiceImpl) GetCachedServiceUserAccessToken() *string {
	existingToken, err := s.tokenRepo.GetToken()
	if err != nil && (err.Error() == "token expired" || err.Error() == "no token found in redis") {
		return nil
	}

	if existingToken == nil {
		return nil
	}
	if config.GetEnv("ROTATE_SERVICE_USER_TOKEN", "n") == "y" {
		// to verify
		isValid, err := s.verifyOnlineServiceUserToken(existingToken.AccessToken)
		if !isValid || err != nil {
			token, _ := s.GenerateAccessToken()
			return token
		}
	}
	return existingToken.AccessToken
}

func (s *AuthServiceImpl) verifyCachedServiceUserToken(token *string) (isOk bool, err error) {
	cachedAccesToken := s.GetCachedServiceUserAccessToken()
	if config.GetEnv("ROTATE_SERVICE_USER_TOKEN", "n") == "y" {
		if cachedAccesToken == nil {
			cachedAccesToken, err = s.GenerateAccessToken()
		}
	}

	if cachedAccesToken == nil || err != nil {
		return false, fmt.Errorf("ERROR | AccessToken cannot be empty")
	}

	if *cachedAccesToken == *token {
		return true, nil
	}
	return false, fmt.Errorf("ERROR | AccessToken cannot be empty")
}

func (s *AuthServiceImpl) verifyOnlineServiceUserToken(token *string) (isValid bool, err error) {
	assertionJWT, err := s.jwtGenerator.GenerateAppInstrospectJWT(time.Hour)
	if err != nil {
		log.Panicf("ERROR | Cannot generate JWT %v", err)
	} // not validate needs to generate
	isValid, err = s.zitadelClient.ValidateServiceUserAccessToken(token, &assertionJWT)
	if err != nil {
		log.Printf("ERROR | Cannot get UserToken %s error: %v", *token, err)
		return false, err
	}
	return isValid, err
}

func (s *AuthServiceImpl) VerifyServiceUserToken(token string) (isOk bool, err error) {
	if token == "" {
		return false, fmt.Errorf("ERROR | AccessToken cannot be empty")
	}

	isOk, err = s.verifyCachedServiceUserToken(&token)
	if err == nil && isOk {
		return isOk, err
	}

	isOk, err = s.verifyOnlineServiceUserToken(&token)
	return isOk, err
}

func (s *AuthServiceImpl) VerifyUserToken(userToken string) (bool, bool) {
	if userToken == "" {
		return false, true
	}
	assertionJWT, err := s.jwtGenerator.GenerateAppInstrospectJWT(time.Hour)
	if err != nil {
		log.Panicf("ERROR | Cannot generate JWT %v", err)
		return false, true
	}

	isValid, expire, err := s.zitadelClient.ValidateUserToken(userToken, assertionJWT)
	if err != nil {
		log.Printf("ERROR | Cannot get UserToken %s error: %v", userToken, err)
		return false, true
	}
	// drift for jwt expire early for 10 minutes
	isExpired := (time.Now().UTC().Unix() - models.TimeDriftForExpire) > expire
	return isValid, isExpired
}
