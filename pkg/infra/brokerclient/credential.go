package brokerclient

import (
	"encoding/json"
	"log"
	"minireipaz/pkg/common"
	"minireipaz/pkg/domain/models"
	"time"
)

type CredentialCommand struct {
	Type       string                    `json:"type,omitempty"`
	Credential *models.CredentialPayload `json:"credential"`
	Timestamp  time.Time                 `json:"timestamp,omitempty"`
}

type CredentialKafkaRepository struct {
	client KafkaClient
}

func NewCredentialKafkaRepository(client KafkaClient) *CredentialKafkaRepository {
	return &CredentialKafkaRepository{
		client: client,
	}
}

func (c *CredentialKafkaRepository) CreateCredential(token, refresh *string, expire *time.Time, stateInfo *models.RequestExchangeCredential) (sended bool) {
	payload := c.credentialToPayload(stateInfo, token, refresh, expire)
	if payload == nil {
		return false
	}
	typeCommand := CommandTypeCreate
	if !stateInfo.CredentialCreatedNew {
		typeCommand = CommandTypeUpdate
	}
	command := CredentialCommand{
		Credential: payload,
		Type:       typeCommand, // TODO: mmmm
		Timestamp:  time.Now().UTC(),
	}
	sended = c.PublishCommand(command, stateInfo.ID)
	return sended
}

// use sync.pool in serverless not necessary
// TODO: marked as optimiced if it's necessary
func (c *CredentialKafkaRepository) credentialToPayload(stateInfo *models.RequestExchangeCredential, token, refresh *string, expire *time.Time) *models.CredentialPayload {
	// now := models.CustomTime{
	// 	Time: time.Now().UTC(),
	// }
	customExpire := models.CustomTime{
		Time: *expire,
	}

	stateInfo.Data.Token = *token
	stateInfo.Data.TokenRefresh = *refresh

	stateInfo.ExpiresAt = &customExpire // needed values?
	// stateInfo.CreatedAt = &now
	// stateInfo.UpdatedAt = &now
	// stateInfo.LastUsedAt = &now
	// stateInfo.RevokedAt = &models.CustomTime{Time: models.TimeDefault}
	// stateInfo.Version = 1
	// stateInfo.IsActive = true

	dataCredential, err := json.Marshal(stateInfo.Data)
	if err != nil {
		log.Printf("ERROR | Cannot convert to json %v", stateInfo.Data)
		return nil
	}

	payload := &models.CredentialPayload{
		RequestExchangeCredential: *stateInfo,
		Data:                      string(dataCredential),
	}
	return payload
}

func (c *CredentialKafkaRepository) PublishCommand(credentialCommand CredentialCommand, key string) bool {
	command, err := json.Marshal(credentialCommand)
	if err != nil {
		log.Printf("ERROR | Cannot transform to JSON %v", err)
		return false
	}

	for i := 1; i < models.MaxAttempts; i++ {
		err = c.client.Produce("credentials.command", []byte(key), command)
		if err == nil {
			return true
		}

		waitTime := common.RandomDuration(models.MaxRangeSleepDuration, models.MinRangeSleepDuration, i)
		log.Printf("ERROR | Cannot connect to Broker, attempt %d: %v. Retrying in %v", i, err, waitTime)
		time.Sleep(waitTime)
	}

	return false
}
