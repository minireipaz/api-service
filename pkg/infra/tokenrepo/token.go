package tokenrepo

import (
	"encoding/json"
	"fmt"
	"log"

	"minireipaz/pkg/config"
	"minireipaz/pkg/domain/models"
	"minireipaz/pkg/infra/redisclient"
	"sync"
	"time"
)

type Token struct {
	ObtainedAt  time.Time     `json:"obtained_at"`
	AccessToken *string       `json:"access_token"`
	TokenType   string        `json:"token_type"`
	ExpiresIn   time.Duration `json:"expires_in"`
}

type TokenRepository struct {
	mu          sync.RWMutex
	redisClient *redisclient.RedisClient
	key         string
	token       *Token
}

func NewTokenRepository(redisClient *redisclient.RedisClient) *TokenRepository {
	return &TokenRepository{
		redisClient: redisClient,
		key:         "serviceuser:token",
	}
}

func (r *TokenRepository) GetToken() (*Token, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.token != nil {
		if r.isExpired(*r.token) {
			return nil, fmt.Errorf("token expired")
		}
		return r.token, nil
	}

	data, err := r.redisClient.Get(r.key)
	if err != nil {
		return nil, err
	}
	if data == "" { // Not exist key in redis
		return nil, fmt.Errorf("no token found in redis")
	}

	var token Token
	err = json.Unmarshal([]byte(data), &token)
	if err != nil {
		return nil, err
	}

	if r.isExpired(token) {
		return nil, fmt.Errorf("token expired")
	}

	r.token = &token
	return r.token, nil
}

func (r *TokenRepository) isExpired(token Token) bool {
	if config.GetEnv("ROTATE_SERVICE_USER_TOKEN", "n") == "y" {
		if time.Now().UTC().After(token.ObtainedAt.Add(token.ExpiresIn * time.Second)) {
			return true
		}
	}
	return false
}

func (r *TokenRepository) SaveToken(accessToken *string, expiresIn *time.Duration) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	token := Token{
		AccessToken: accessToken,
		ExpiresIn:   *expiresIn - models.SaveOffset, // -10 seconds
		ObtainedAt:  time.Now().UTC(),
	}

	data, err := json.Marshal(token)
	if err != nil {
		return err
	}

	err = r.redisClient.WatchToken(string(data), r.key, (token.ExpiresIn)*time.Second)
	if err == nil {
		r.token = &token
		return nil
	}

	log.Printf("ERROR | Failed to save token, %v", err)
	return err
}

func (r *TokenRepository) SetToken(token *Token) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.token = token
}
