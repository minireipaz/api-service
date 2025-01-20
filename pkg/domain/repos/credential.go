package repos

import (
	"minireipaz/pkg/domain/models"
	"time"
)

type CredentialService interface {
	CreateCredential(credentialFrontend *models.RequestCreateCredential) (*models.RequestExchangeCredential, error)
	CreateTokenCredential(credentialFrontend *models.RequestCreateCredential) (saved bool, transformedCredentialID *string, err error)
	ExchangeGoogleCredential(currentCredential *models.RequestExchangeCredential) (token, refresh *string, expire *time.Time, stateInfo *models.RequestExchangeCredential, err error)
	GetAllCredentials(userID *string) (*models.ResponseGetCredential, bool)
	TransformWorkflow(currenteCredential *models.RequestExchangeCredential, workflow *models.Workflow) *models.Workflow
	GetCredentialByID(userID *string, credentialID *string) (response *models.ResponseGetCredential)
}

type CredentialHTTPRepository interface {
	GetAllCredentials(userID *string, limitCount uint64) (*[]models.RequestExchangeCredential, error)
	GetCredentialByID(userID *string, credentialID *string, limitCount uint64) (*[]models.RequestExchangeCredential, error)
}

type CredentialGoogleHTTPRepository interface {
	GenerateAuthURL(credential *models.RequestExchangeCredential, credentialCreatedNew *bool) *string
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
