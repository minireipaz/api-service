package brokerclient

import (
	"encoding/json"
	"minireipaz/pkg/domain/models"
)

type UserKafkaRepository struct {
	client KafkaClient
}

func NewUserKafkaRepository(client KafkaClient) *UserKafkaRepository {
	return &UserKafkaRepository{
		client: client,
	}
}

func (u *UserKafkaRepository) Create(user *models.Users) (created bool, exist bool) {
	userJSON, err := json.Marshal(user)
	if err != nil {
		return false, false
	}

	err = u.client.Produce("users.db.write", []byte(user.Stub), userJSON)
	if err != nil {
		return false, false
	}

	return true, false
}
