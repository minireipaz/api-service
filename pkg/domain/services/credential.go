package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"minireipaz/pkg/common"
	"minireipaz/pkg/domain/models"
	"minireipaz/pkg/domain/repos"
	"net/http"
	"strings"
	"time"
)

type CredentialServiceImpl struct {
	googleOAuthRepo      repos.CredentialGoogleHTTPRepository
	facebookOAuthRepo    repos.CredentialFacebookHTTPRepository
	redisRepo            repos.CredentialRedisRepository
	credentialBrokerRepo repos.CredentialBrokerRepository
	credentialHTTP       repos.CredentialHTTPRepository
}

func NewCredentialService(googleRepo repos.CredentialGoogleHTTPRepository,
	facebookRepo repos.CredentialFacebookHTTPRepository,
	redisCli repos.CredentialRedisRepository,
	brokerRepo repos.CredentialBrokerRepository,
	credentialRepo repos.CredentialHTTPRepository) repos.CredentialService {
	return &CredentialServiceImpl{
		googleOAuthRepo:      googleRepo,
		facebookOAuthRepo:    facebookRepo,
		redisRepo:            redisCli,
		credentialBrokerRepo: brokerRepo,
		credentialHTTP:       credentialRepo,
	}
}

func (c *CredentialServiceImpl) CreateCredential(currentCredential *models.RequestCreateCredential) (*models.RequestCreateCredential, error) {
  // not necesary to generate new ID
  // datasource credentials from clickhouse is using ReplacingMergeTree
	if currentCredential.ID == "none" || !strings.HasPrefix(currentCredential.ID, "credential_") {
		currentCredential.ID = c.generateNewIDCredential(currentCredential)
	}
	switch currentCredential.Type { // TODO: refactor
	case "googlesheets":
		authURL := c.googleOAuthRepo.GenerateAuthURL(currentCredential)
		currentCredential.Data.RedirectURL = *authURL
	default:
		return nil, errors.New("ERROR | currentCredential Type not supported")
	}

	return currentCredential, nil
}

// this function DONT GENERATE GLOBAL UNIQUE ID in theory can be collissions in removed nodes, workflows,...
// TODO: maybe can make general function to create GLOBAL IDs
 func (c *CredentialServiceImpl) generateNewIDCredential(currentCredential *models.RequestCreateCredential) string {
	return fmt.Sprintf("credential_%s_%s_%s_%s", currentCredential.Sub, currentCredential.WorkflowID, currentCredential.NodeID, currentCredential.Type)
}

// func (c *CredentialServiceImpl) generateNewIDCredential() string {
// 	newCredentialID := uuid.New().String()
// 	if newCredentialID == "" { // in case fail
// 		newCredentialID = uuid.New().String()
// 	}
// 	now := time.Now().UTC().Unix()
// 	return fmt.Sprintf("credential_%s_%d", newCredentialID, now)
// }

func (c *CredentialServiceImpl) ExchangeGoogleCredential(currentCredential *models.RequestExchangeCredential) (token, refresh *string, expire *time.Time, stateInfo *models.RequestExchangeCredential, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), models.MaxTimeoutContext)
	defer cancel()
	// check if it's inserted locked
	locked, err := c.insertLocker(currentCredential)
	if err != nil {
		return token, refresh, nil, nil, err
	}

	if !locked {
		return token, refresh, nil, nil, fmt.Errorf("ERROR | Wait 5 seconds")
	}

	err = c.retryTemplateWithError(ctx, models.MaxAttempts, func() error {
		var lastError error
		token, refresh, expire, stateInfo, lastError = c.googleOAuthRepo.ExchangeGoogleCredential(currentCredential)
		return lastError
	})

	if err != nil {
		return token, refresh, nil, nil, err
	}
	sended := c.saveCredentialExchange(token, refresh, expire, stateInfo)
	if !sended {
		log.Printf("ERROR | Needs to Added to dead letter %s", currentCredential.Sub)
	}
	return token, refresh, expire, stateInfo, err
}

func (c *CredentialServiceImpl) insertLocker(currentCredential *models.RequestExchangeCredential) (locked bool, err error) {
	locked, err = c.redisRepo.AddLock(&currentCredential.Sub)
	if err != nil {
		log.Printf("ERROR | Needs to Added to dead letter %s", currentCredential.Sub)
	}

	if !locked {
		log.Printf("WARN | Cannot created lock for user %s", currentCredential.Sub)
	}

	return locked, err
}

func (c *CredentialServiceImpl) saveCredentialExchange(token, refresh *string, expire *time.Time, stateInfo *models.RequestExchangeCredential) (sended bool) {
	sended = c.credentialBrokerRepo.CreateCredential(token, refresh, expire, stateInfo)
	return sended
}

func (c *CredentialServiceImpl) retryTemplateWithError(ctx context.Context, maxAttempts int, operation func() error) error {
	for i := 1; i < maxAttempts; i++ {
		if err := operation(); err == nil {
			return nil
		}
		// not necesaary last wait time
		waitTime := common.RandomDuration(models.MaxRangeSleepDuration, models.MinRangeSleepDuration, i)
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(waitTime):
		}
	}
	return fmt.Errorf("ERROR | operation failed after %d attempts", maxAttempts)
}

func (c *CredentialServiceImpl) GetAllCredentials(userID *string) (response *models.ResponseGetCredential, isok bool) {
	credentials, err := c.credentialHTTP.GetAllCredentials(userID, 1000) // 1000 credential as max
	if err != nil {
		return &models.ResponseGetCredential{
			Status: http.StatusBadRequest,
			Error:  err.Error(),
		}, false
	}

	return &models.ResponseGetCredential{
		Status:      http.StatusOK,
		Credentials: credentials,
	}, true
}
