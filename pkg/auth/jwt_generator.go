package auth

import (
	"time"
	"github.com/golang-jwt/jwt/v5"
)

type JWTGenerator struct {
	serviceUserID string
	privateKey    []byte
	keyID         string
	apiURL        string
}

func NewJWTGenerator(serviceUserID, privateKey, keyID, apiURL string) *JWTGenerator {
	return &JWTGenerator{
		serviceUserID: serviceUserID,
		privateKey:    []byte(privateKey),
		keyID:         keyID,
		apiURL:        apiURL,
	}
}

func (j *JWTGenerator) GenerateJWT(timeExpire time.Duration) (string, error) {
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iss": j.serviceUserID,
		"sub": j.serviceUserID,
		"aud": j.apiURL,
		"exp": now.Add(timeExpire).Unix(),
		"iat": now.Unix(),
	})
	token.Header["kid"] = j.keyID

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(j.privateKey)
	if err != nil {
		return "", err
	}

	return token.SignedString(privateKey)
}
