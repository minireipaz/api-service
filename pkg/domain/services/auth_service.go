package services

import (
	"fmt"
	"log"
	"minireipaz/pkg/auth"
	"minireipaz/pkg/infra/httpclient"
	"minireipaz/pkg/infra/tokenrepo"
	"time"
)

const (
	twoDays = 172_800 * time.Second
)

type AuthService struct {
	jwtGenerator  *auth.JWTGenerator
	zitadelClient *httpclient.ZitadelClient
	tokenRepo     *tokenrepo.TokenRepository
}

func NewAuthService(jwtGenerator *auth.JWTGenerator, zitadelClient *httpclient.ZitadelClient, tokenRepo *tokenrepo.TokenRepository) *AuthService {
	return &AuthService{
		jwtGenerator:  jwtGenerator,
		zitadelClient: zitadelClient,
		tokenRepo:     tokenRepo,
	}
}

func (s *AuthService) GenerateNewToken() (string, error) {
	jwt, err := s.jwtGenerator.GenerateJWT(twoDays)
	if err != nil {
		log.Panicf("ERROR | Cannot generate JWT %v", err)
	}

	accessToken, expiresIn, err := s.zitadelClient.GetAccessToken(jwt)
	if err != nil {
		log.Panicf("ERROR | Cannot acces to ACCESS token %v", err)
	}

	token := &tokenrepo.Token{
		AccessToken: accessToken,
		ExpiresIn:   expiresIn,
		ObtainedAt:  time.Now(),
	}

	err = s.tokenRepo.SaveToken(token)
	if err != nil {
		log.Panicf("ERROR | Failed to save token, %v", err)
	}

	return accessToken, nil
}

func (s *AuthService) GetAccessToken() (string, error) {
	existingToken, err := s.tokenRepo.GetToken()
	if err != nil && (err.Error() == "token expired" || err.Error() == "no token found in redis") {
		return s.GenerateNewToken()
	}

	if existingToken == nil || time.Now().After(existingToken.ObtainedAt.Add(existingToken.ExpiresIn*time.Second)) {
		// Rotate token if it's expired or not found
		return s.GenerateNewToken()
	}
	// TODO: Verify Service USER access token with ID Provider
	isValid, err := s.verifyWithIDProvider(existingToken)
	if !isValid || err != nil {
		return s.GenerateNewToken()
	}
	return existingToken.AccessToken, nil
}

func (s *AuthService) VerifyServiceUserToken(token string) (bool, error) {
	masterToken, err := s.GetAccessToken()
	if err != nil {
		return false, err
	}
	return masterToken == token, err
}

func (s *AuthService) verifyWithIDProvider(token *tokenrepo.Token) (bool, error) {
	// TODO: verify with IDProvider
	if token.AccessToken == "" { /// dummy check
		return false, fmt.Errorf("ERROR | AccessToken cannot be empty")
	}
	return true, nil
}

func (s *AuthService) VerifyUserToken(token *tokenrepo.Token) (bool, error) {
  return false, nil
}

// func (s *AuthService) verifyUserAccessTokenWithIDProvider(token *tokenrepo.Token) (bool, error) {
//   // TODO: verify with IDProvider
//   if token.AccessToken == "" { /// dummy check
//     return false, nil
//   }
// 	return true, nil
// }
