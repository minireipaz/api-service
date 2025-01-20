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

func (c *CredentialServiceImpl) CreateCredential(credentialFrontend *models.RequestCreateCredential) (*models.RequestExchangeCredential, error) {
	// obtain data from DB
	dbCredential := c.GetCredentialByID(&credentialFrontend.Sub, &credentialFrontend.ID)
	if dbCredential.Status != 200 {
		return nil, fmt.Errorf("cannot saved")
	}

	transformedCredential := c.mergeDBCredentail(dbCredential, credentialFrontend)
	if transformedCredential == nil {
		return nil, fmt.Errorf("cannot saved")
	}
	switch transformedCredential.Type { // TODO: refactor
	case "googlesheets":
		authURL := c.googleOAuthRepo.GenerateAuthURL(transformedCredential, &transformedCredential.CredentialCreatedNew)
		credentialFrontend.Data.RedirectURL = *authURL
	default:
		return nil, errors.New("ERROR | currentCredential Type not supported")
	}

	return transformedCredential, nil
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

func (c *CredentialServiceImpl) GetCredentialByID(userID *string, credentialID *string) (response *models.ResponseGetCredential) {
	credentials, err := c.credentialHTTP.GetCredentialByID(userID, credentialID, 1) // 1000 credential as max
	if err != nil {
		log.Printf("ERROR | getcredentialbyid %v by userid %s credential id %s", err, *userID, *credentialID)
		return &models.ResponseGetCredential{
			Status: http.StatusBadRequest,
			Error:  err.Error(),
		}
	}

	return &models.ResponseGetCredential{
		Status:      http.StatusOK,
		Credentials: credentials,
	}
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

func (c *CredentialServiceImpl) mergeDBCredentail(dbCredential *models.ResponseGetCredential, credentialFrontend *models.RequestCreateCredential) (transformedCredential *models.RequestExchangeCredential) {
	if len(*dbCredential.Credentials) == 0 {
		// new credential
		if credentialFrontend.ID == "none" || !strings.HasPrefix(credentialFrontend.ID, "credential_") {
			newCredential := true
			transformedCredential = c.transformNewCredentialToBackend(credentialFrontend, newCredential)
		} else {
			// ban user hackerman
			log.Printf("CRITICAL | lock user")
		}
	}

	if len(*dbCredential.Credentials) == 1 {
		// update credential ???
		newCredential := false
		transformedCredential = c.transformExistedCredentialToBackend(credentialFrontend, dbCredential, newCredential)
	} else if len(*dbCredential.Credentials) > 1 {
		return nil
	}

	return transformedCredential
}

func (c *CredentialServiceImpl) CreateTokenCredential(credentialFrontend *models.RequestCreateCredential) (saved bool, transformedCredentialID *string, err error) {
	// obtain data from DB
	dbCredential := c.GetCredentialByID(&credentialFrontend.Sub, &credentialFrontend.ID)
	if dbCredential.Status != 200 {
		return false, nil, fmt.Errorf("cannot saved")
	}

	transformedCredential := c.mergeDBCredentail(dbCredential, credentialFrontend)
	if transformedCredential == nil {
		return false, nil, fmt.Errorf("cannot saved")
	}

	switch transformedCredential.Type { // TODO: refactor
	case "notiontoken":
		saved, err = c.saveTokenCredential(transformedCredential)
	default:
		return false, nil, errors.New("ERROR | currentCredential Type not supported")
	}
	// ID only necessary when created new one
	return saved, &transformedCredential.ID, err
}

// TODO: can be merged two funcitions transformNewCredentialToBackend and transformExistedCredentialToBackend
func (c *CredentialServiceImpl) transformNewCredentialToBackend(credentialFrontend *models.RequestCreateCredential, createdNew bool) *models.RequestExchangeCredential {
	now := time.Now().UTC()
	newID := c.generateNewIDCredential(credentialFrontend)
	requestID := c.generateNewRequestID()
	exchangeCredential := &models.RequestExchangeCredential{
		ID:                   newID,
		RequestID:            requestID,
		CredentialCreatedNew: createdNew,
		RevokedAt:            &models.CustomTime{Time: models.TimeDefault},
		LastUsedAt:           &models.CustomTime{Time: models.TimeDefault},
		ExpiresAt:            &models.CustomTime{Time: now.AddDate(1, 0, 0)}, // TODO: 1 year - better control for lifetime,
		UpdatedAt:            &models.CustomTime{Time: now},
		CreatedAt:            &models.CustomTime{Time: now},
		NodeID:               credentialFrontend.NodeID,
		Sub:                  credentialFrontend.Sub,
		WorkflowID:           credentialFrontend.WorkflowID,
		Type:                 credentialFrontend.Type,
		Name:                 credentialFrontend.Name,
		Data:                 credentialFrontend.Data,
		Version:              1,
		IsActive:             true,
	}
	return exchangeCredential
}

func (c *CredentialServiceImpl) transformExistedCredentialToBackend(credential *models.RequestCreateCredential, dbCredential *models.ResponseGetCredential, newCredential bool) *models.RequestExchangeCredential {
	now := time.Now().UTC()
	// already checked for len
	dbCredentialCurrent := *dbCredential.Credentials
	// TODO: better
	revokedat := dbCredentialCurrent[0].RevokedAt
	lastusedat := dbCredentialCurrent[0].LastUsedAt
	expiresat := &models.CustomTime{Time: now.AddDate(1, 0, 0)} // TODO: 1 year - better control for lifetime
	updatedat := &models.CustomTime{Time: now}
	createdat := dbCredentialCurrent[0].CreatedAt
	newRequestID := c.generateNewRequestID()

	exchangeCredential := &models.RequestExchangeCredential{
		RequestID:            newRequestID,
		CredentialCreatedNew: newCredential,
		RevokedAt:            revokedat,
		LastUsedAt:           lastusedat,
		ExpiresAt:            expiresat,
		UpdatedAt:            updatedat,
		CreatedAt:            createdat,
		NodeID:               credential.NodeID,
		Sub:                  dbCredentialCurrent[0].Sub,
		WorkflowID:           credential.WorkflowID,
		ID:                   dbCredentialCurrent[0].ID,
		Type:                 credential.Type,
		Name:                 credential.Name,
		Data:                 credential.Data,
		Version:              1,
		IsActive:             true,
	}
	return exchangeCredential
}

func (c *CredentialServiceImpl) saveTokenCredential(transformedCredential *models.RequestExchangeCredential) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), models.MaxTimeoutContext)
	defer cancel() // checkout
	// TODO: maybe insert in retry template check if it's inserted locked
	locked, err := c.insertLocker(transformedCredential)
	if err != nil {
		log.Printf("ERROR | cannot insert lock in savetokencredential %v", err)
		return false, fmt.Errorf("error cannot lock key %v", transformedCredential)
	}

	if !locked { // not error perse
		return false, fmt.Errorf("wait 5 seconds")
	}

	sended := false
	err = c.retryTemplateWithError(ctx, models.MaxAttempts, func() error {
		var lastError error
		sended = c.saveCredentialExchange(&transformedCredential.Data.Token, &transformedCredential.Data.TokenRefresh, &transformedCredential.ExpiresAt.Time, transformedCredential)
		if !sended {
			log.Printf("ERROR | Cannot save Credential NOT added to dead letter %v", transformedCredential)
			lastError = fmt.Errorf("ERROR | Cannot save Credential NOT added to dead letter %v", transformedCredential)
		}
		return lastError
	})

	return sended, err
}
