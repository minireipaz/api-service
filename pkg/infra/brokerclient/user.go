package brokerclient

import (
	"encoding/json"
	"log"
	"minireipaz/pkg/common"
	"minireipaz/pkg/domain/models"
	"time"
)

type UserKafkaRepository struct {
	client KafkaClient
}

func NewUserKafkaRepository(client KafkaClient) *UserKafkaRepository {
	return &UserKafkaRepository{
		client: client,
	}
}

func (u *UserKafkaRepository) CreateUser(user *models.SyncUserRequest) (sended bool) {
	userJSON, err := json.Marshal(user)
	if err != nil {
		log.Printf("ERROR | Cannot transform to JSON %v", err)
		return false
	}

	for i := 1; i < models.MaxAttempts; i++ {
		err = u.client.Produce("users.db.write", []byte(user.Sub), userJSON)
		if err == nil {
			return true
		}

		waitTime := common.RandomDuration(models.MaxRangeSleepDuration, models.MinRangeSleepDuration, i)
		log.Printf("ERROR | Cannot connect to Broker, attempt %d: %v. Retrying in %v", i, err, waitTime)
		time.Sleep(waitTime)
	}

	return false
}
