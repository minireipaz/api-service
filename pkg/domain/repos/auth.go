package repos

import (
	"minireipaz/pkg/infra/tokenrepo"
	"time"
)

type AuthRepository interface {
	GenerateIntrospectJWT(duration time.Duration) string
	GenerateNewToken() (string, error)
	GetServiceUserAccessToken() (string, error)
	VerifyServiceUserToken(token string) (bool, error)
	verifyWithIDProvider(token *tokenrepo.Token) (bool, error)
	VerifyUserToken(userToken string) bool
}

type JWTGenerator interface {
	GenerateInstrospectJWT(duration time.Duration) (string, error)
	GenerateServiceUserJWT(duration time.Duration) (string, error)
}

type ZitadelClient interface {
	GetServiceUserAccessToken(jwt string) (string, time.Duration, error)
	ValidateUserToken(userToken, introspectJWT string) (bool, error)
}

type TokenRepository interface {
	SaveToken(token *tokenrepo.Token) error
	GetToken() (*tokenrepo.Token, error)
}
