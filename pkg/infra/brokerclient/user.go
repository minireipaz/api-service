package brokerclient

import (
	"encoding/json"
	"log"
	"minireipaz/pkg/common"
	"minireipaz/pkg/domain/models"
	"time"
)

const (
	offset     = 1 * time.Second
	timedrift  = 500 * time.Millisecond
	maxIntents = 10
)

type UserKafkaRepository struct {
	client KafkaClient
}

func NewUserKafkaRepository(client KafkaClient) *UserKafkaRepository {
	return &UserKafkaRepository{
		client: client,
	}
}

func (u *UserKafkaRepository) CreateUser(user *models.Users) (sended bool) {
	userJSON, err := json.Marshal(user)
	if err != nil {
		log.Printf("ERROR | Cannot transform to JSON %v", err)
		return false
	}

	for i := 1; i < maxIntents; i++ {
		err = u.client.Produce("users.db.write", []byte(user.Sub), userJSON)
		if err == nil {
			return true
		}

		waitTime := common.RandomDuration(models.MaxSleepDuration, models.MinSleepDuration, i)
		log.Printf("ERROR | Cannot connect to Broker, attempt %d: %v. Retrying in %v", i, err, waitTime)
		time.Sleep(waitTime)
	}

	return false
}
