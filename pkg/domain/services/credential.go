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

	"github.com/google/uuid"
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
	credentialCreatedNew := false
	// not necessary to generate new ID
	// datasource credentials from clickhouse is using ReplacingMergeTree
	// TODO: requestid is not setted to global uuid and not saved to DB
	if currentCredential.ID == "none" || !strings.HasPrefix(currentCredential.ID, "credential_") {
		credentialCreatedNew = true
		currentCredential.ID = c.generateNewIDCredential(currentCredential)
	}
	currentCredential.RequestID = c.generateNewRequestID()
	switch currentCredential.Type { // TODO: refactor
	case "googlesheets":
		authURL := c.googleOAuthRepo.GenerateAuthURL(currentCredential, &credentialCreatedNew)
		currentCredential.Data.RedirectURL = *authURL
	default:
		return nil, errors.New("ERROR | currentCredential Type not supported")
	}

	return currentCredential, nil
}

//  func (c *CredentialServiceImpl) generateNewIDCredential(currentCredential *models.RequestCreateCredential) string {
// 	now := time.Now().UTC().Unix()
// 	return fmt.Sprintf("credential_%s_%s_%s_%s_%d", currentCredential.Sub, currentCredential.WorkflowID, currentCredential.NodeID, currentCredential.Type, now)
// }

// this function DONT GENERATE GLOBAL UNIQUE ID in theory can be collissions in removed nodes, workflows,...
// TODO: maybe can make general function to create GLOBAL IDs
func (c *CredentialServiceImpl) generateNewIDCredential(currentCredential *models.RequestCreateCredential) string {
	newCredentialID := uuid.New().String()
	if newCredentialID == "" { // in case fail
		newCredentialID = uuid.New().String()
	}
	return fmt.Sprintf("credential_%s_%s", currentCredential.Sub, newCredentialID)
}

// this function DONT GENERATE GLOBAL UNIQUE ID in theory can be collissions in removed nodes, workflows,...
// TODO: maybe can make general function to create GLOBAL IDs
func (c *CredentialServiceImpl) generateNewRequestID() string {
	newCredentialID := uuid.New().String()
	if newCredentialID == "" { // in case fail
		newCredentialID = uuid.New().String()
	}
	return newCredentialID
}

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
		log.Printf("ERROR | Cannot save Credential NOT added to dead letter %v", currentCredential)
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

func (c *CredentialServiceImpl) TransformWorkflow(currenteCredential *models.RequestExchangeCredential, workflow *models.Workflow) *models.Workflow {
	for i := 0; i < len(workflow.Nodes); i++ {
		if workflow.Nodes[i].ID == currenteCredential.NodeID {
			// TODO: in case changed btw i think dont changed
			workflow.Nodes[i].Data.CredentialData.ID = currenteCredential.ID
			// workflow.Nodes[i].Data.CredentialData.Data.ClientID = currenteCredential.Data.ClientID
			// workflow.Nodes[i].Data.CredentialData.Data.ClientSecret = currenteCredential.Data.ClientSecret
			workflow.Nodes[i].Data.CredentialData.Data = currenteCredential.Data
		}
	}
	return workflow
}
