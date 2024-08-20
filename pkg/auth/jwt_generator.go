package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTGenerator struct {
	ServiceUser ServiceUserConfig
	BackendApp  BackendAppConfig
	apiURL      string
	clientID    string
	projectID   string
}

type ServiceUserConfig struct {
	UserID     string
	PrivateKey []byte
	KeyID      string
}

type BackendAppConfig struct {
	KeyID      string
	PrivateKey []byte
	AppID      string
}

type JWTGeneratorConfig struct {
	ServiceUser ServiceUserConfig
	BackendApp  BackendAppConfig
	APIURL      string
	ProjectID   string
	ClientID    string
}

func NewJWTGenerator(config JWTGeneratorConfig) *JWTGenerator {
	return &JWTGenerator{
		ServiceUser: config.ServiceUser,
		BackendApp:  config.BackendApp,
		apiURL:      config.APIURL,
		projectID:   config.ProjectID,
		clientID:    config.ClientID,
	}
}

func (j *JWTGenerator) GenerateInstrospectJWT(timeExpire time.Duration) (string, error) {
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iss": j.clientID,
		"sub": j.clientID,
		"aud": j.apiURL,
		"exp": now.Add(timeExpire).Unix(),
		"iat": now.Unix(),
	})
	token.Header["kid"] = j.BackendApp.KeyID

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(j.BackendApp.PrivateKey)
	if err != nil {
		return "", err
	}

	return token.SignedString(privateKey)
}

func (j *JWTGenerator) GenerateServiceUserJWT(timeExpire time.Duration) (string, error) {
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iss": j.ServiceUser.UserID,
		"sub": j.ServiceUser.UserID,
		"aud": j.apiURL,
		"exp": now.Add(timeExpire).Unix(),
		"iat": now.Unix(),
	})
	token.Header["kid"] = j.ServiceUser.KeyID

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(j.ServiceUser.PrivateKey)
	if err != nil {
		return "", err
	}

	return token.SignedString(privateKey)
}
