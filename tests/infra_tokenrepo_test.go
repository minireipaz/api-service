package tests

import (
	"encoding/json"
	"testing"
	"time"

	"minireipaz/pkg/infra/redisclient"
	"minireipaz/pkg/infra/tokenrepo"

	"github.com/stretchr/testify/assert"
)

var validTokenRedis = "valid-token-redis"
var expiredTokenRedis = "expired-token-redis"
var expiredToken = "expired-token"
var validToken = "valid-token"
var secondsExpired = time.Second * 7600

func TestTokenRepository_GetToken(t *testing.T) {
	r := redisclient.NewRedisClient()
	tokenRepo := tokenrepo.NewTokenRepository(r)
	tests := []struct {
		name        string
		setup       func()
		expectedRes *tokenrepo.Token
		expectedErr string
	}{

		{
			name: "no token found in redis and in memory",
			setup: func() {
				r.Client.Del(r.Ctx, "serviceuser_backend:token")
				tokenRepo.SetToken(nil)
			},
			expectedRes: nil,
			expectedErr: "no token found in redis",
		},

		{
			name: "token exists in redis and is valid",
			setup: func() {
				token := &tokenrepo.Token{
					ObtainedAt:  time.Now().UTC(),
					AccessToken: &validTokenRedis,
					TokenType:   "Bearer",
					ExpiresIn:   3600, // 1 hora
				}
				data, _ := json.Marshal(token)
				r.Set("serviceuser_backend:token", string(data))
			},
			expectedRes: &tokenrepo.Token{
				ObtainedAt:  time.Now().UTC(),
				AccessToken: &validTokenRedis,
				TokenType:   "Bearer",
				ExpiresIn:   3600,
			},
			expectedErr: "",
		},
		{
			name: "no token found in redis but in memory",
			setup: func() {
				r.Client.Del(r.Ctx, "serviceuser_backend:token")
			},
			expectedRes: &tokenrepo.Token{ // val from before test case
				ObtainedAt:  time.Now().UTC(),
				AccessToken: &validTokenRedis,
				TokenType:   "Bearer",
				ExpiresIn:   3600,
			},
			expectedErr: "",
		},

		// {
		// 	name: "token exists in memory and is valid",
		// 	setup: func() {
		// 		tokenRepo.SaveServiceUserToken(&validToken, &secondsExpired)
		// 	},
		// 	expectedRes: &tokenrepo.Token{
		// 		ObtainedAt:  time.Now().UTC(),
		// 		AccessToken: &validToken,
		// 		TokenType:   "Bearer",
		// 		ExpiresIn:   3600,
		// 	},
		// 	expectedErr: "",
		// },
    // {
		// 	name: "token exists in redis but is expired",
		// 	setup: func() {
		// 		token := &tokenrepo.Token{
		// 			ObtainedAt:  time.Now().UTC().Add(-2 * time.Hour),
		// 			AccessToken: &expiredTokenRedis,
		// 			TokenType:   "Bearer",
		// 			ExpiresIn:   3600, // 1 hora
		// 		}
		// 		data, _ := json.Marshal(token)
		// 		r.Set("serviceuser_backend:token", string(data))
		// 	},
		// 	expectedRes: &tokenrepo.Token{
		// 		ObtainedAt:  time.Now().UTC().Add(-2 * time.Hour),
		// 		AccessToken: &expiredTokenRedis,
		// 		TokenType:   "Bearer",
		// 		ExpiresIn:   3600, // 1 hora
		// 	},
		// 	expectedErr: "",
		// },
		{
			name: "token exists in memory",
			setup: func() {
				tokenRepo.SaveServiceUserToken(&expiredToken, &secondsExpired)
			},
			expectedRes: nil,
			expectedErr: "",
		},

		{
			name: "error fetching token from redis",
			setup: func() {
				tokenRepo.SetToken(nil)
				r.Client.Close()
			},
			expectedRes: nil,
			expectedErr: "redis: client is closed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			token, err := tokenRepo.GetServiceUserToken()

			if tt.expectedErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
				assert.Nil(t, token)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, token)
				if tt.expectedRes != nil {
					assert.Equal(t, tt.expectedRes.AccessToken, token.AccessToken)
					assert.Equal(t, tt.expectedRes.TokenType, token.TokenType)
				}
			}
		})
	}
}

func TestTokenRepository_SaveToken(t *testing.T) {
	r := redisclient.NewRedisClient()
	tokenRepo := tokenrepo.NewTokenRepository(r)

	tests := []struct {
		name        string
		setup       func()
		token       *tokenrepo.Token
		expectedErr string
	}{
		{
			name: "successfully save token",
			setup: func() {
			},
			token: &tokenrepo.Token{
				ObtainedAt:  time.Now().UTC(),
				AccessToken: &validToken,
				TokenType:   "Bearer",
				ExpiresIn:   3600, // 1 hora
			},
			expectedErr: "",
		},
		{
			name: "error saving token",
			setup: func() {
				r.Client.Close()
			},
			token: &tokenrepo.Token{
				ObtainedAt:  time.Now().UTC(),
				AccessToken: &validToken, //"valid-token",
				TokenType:   "Bearer",
				ExpiresIn:   3600,
			},
			expectedErr: "redis: client is closed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			err := tokenRepo.SaveServiceUserToken(tt.token.AccessToken, &tt.token.ExpiresIn)

			if tt.expectedErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
