package repos

import (
	"minireipaz/pkg/infra/tokenrepo"
	"time"
)

type AuthService interface {
	GenerateAccessToken() (*string, error)
	GetCachedServiceUserAccessToken() *string
	VerifyServiceUserToken(token string) (isOk bool, err error)
	VerifyUserToken(userToken string) (bool, bool)
	GetActionUserAccessToken() (*string, error)
}

type AuthRepository interface {
	GenerateIntrospectJWT(duration time.Duration) string
	GenerateAccessToken() (string, error)
	VerifyServiceUserToken(token string) (bool, error)
	verifyWithIDProvider(token *tokenrepo.Token) (bool, error)
	VerifyUserToken(userToken string) (bool, bool)
}

type JWTGenerator interface {
	GenerateServiceUserAssertionJWT(duration time.Duration) (string, error)
	GenerateAppInstrospectJWT(duration time.Duration) (string, error)
}

type ZitadelClient interface {
	GenerateServiceUserAccessToken(jwt string) (*string, time.Duration, error)
	ValidateUserToken(userToken, introspectJWT string) (bool, int64, error)
	ValidateServiceUserAccessToken(userToken, introspectJWT *string) (bool, error)
}

type TokenRepository interface {
	SaveServiceUserToken(accessToken *string, expiresIn *time.Duration) error
	GetServiceUserToken() (*tokenrepo.Token, error)
	GetActionUserToken() (*tokenrepo.Token, error)
}
