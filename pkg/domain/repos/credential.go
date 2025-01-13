package repos

import (
	"minireipaz/pkg/domain/models"
	"time"
)

type CredentialService interface {
	CreateCredential(currentCredential *models.RequestCreateCredential) (*models.RequestCreateCredential, error)
	ExchangeGoogleCredential(currentCredential *models.RequestExchangeCredential) (token, refresh *string, expire *time.Time, stateInfo *models.RequestExchangeCredential, err error)
	GetAllCredentials(userID *string) (*models.ResponseGetCredential, bool)
	TransformWorkflow(currenteCredential *models.RequestExchangeCredential, workflow *models.Workflow) *models.Workflow
}

type CredentialHTTPRepository interface {
	GetAllCredentials(userID *string, limitCount uint64) (*[]models.RequestExchangeCredential, error)
}

type CredentialGoogleHTTPRepository interface {
	GenerateAuthURL(credential *models.RequestCreateCredential, credentialCreatedNew *bool) *string
	ExchangeGoogleCredential(currentCredential *models.RequestExchangeCredential) (token, refresh *string, expire *time.Time, stateInfo *models.RequestExchangeCredential, err error)
}

type CredentialFacebookHTTPRepository interface {
}

type CredentialRedisRepository interface {
	SaveTemporalAuthURLData(currentCredential *models.RequestCreateCredential) (inserted bool, err error)
	AddLock(sub *string) (locked bool, err error)
}

type CredentialBrokerRepository interface {
	CreateCredential(token, refresh *string, expire *time.Time, stateInfo *models.RequestExchangeCredential) (sended bool)
}
