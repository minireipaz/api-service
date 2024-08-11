package tests

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"minireipaz/pkg/infra/redisclient"
	"minireipaz/pkg/infra/tokenrepo"
)

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
				r.Client.Del(r.Ctx, "auth:token")
				tokenRepo.SetToken(nil)
			},
			expectedRes: nil,
			expectedErr: "no token found in redis",
		},

		{
			name: "token exists in redis and is valid",
			setup: func() {
				token := &tokenrepo.Token{
					ObtainedAt:  time.Now(),
					AccessToken: "valid-token-redis",
					TokenType:   "Bearer",
					ExpiresIn:   3600, // 1 hora
				}
				data, _ := json.Marshal(token)
				r.Set("auth:token", string(data))
			},
			expectedRes: &tokenrepo.Token{
				ObtainedAt:  time.Now(),
				AccessToken: "valid-token-redis",
				TokenType:   "Bearer",
				ExpiresIn:   3600,
			},
			expectedErr: "",
		},
		{
			name: "no token found in redis but in memory",
			setup: func() {
				r.Client.Del(r.Ctx, "auth:token")
			},
			expectedRes: &tokenrepo.Token{ // val from before test case
				ObtainedAt:  time.Now(),
				AccessToken: "valid-token-redis",
				TokenType:   "Bearer",
				ExpiresIn:   3600,
			},
			expectedErr: "",
		},

		{
			name: "token exists in memory and is valid",
			setup: func() {
				// Guardar un token válido en memoria
				tokenRepo.SaveToken(&tokenrepo.Token{
					ObtainedAt:  time.Now(),
					AccessToken: "valid-token",
					TokenType:   "Bearer",
					ExpiresIn:   3600, // 1 hora
				})
			},
			expectedRes: &tokenrepo.Token{
				ObtainedAt:  time.Now(),
				AccessToken: "valid-token",
				TokenType:   "Bearer",
				ExpiresIn:   3600,
			},
			expectedErr: "",
		},
		{
			name: "token exists in memory but is expired",
			setup: func() {
				// Guardar un token expirado en memoria
				tokenRepo.SaveToken(&tokenrepo.Token{
					ObtainedAt:  time.Now().Add(-2 * time.Hour),
					AccessToken: "expired-token",
					TokenType:   "Bearer",
					ExpiresIn:   3600, // 1 hora
				})
			},
			expectedRes: nil,
			expectedErr: "token expired",
		},
		{
			name: "token exists in redis but is expired",
			setup: func() {
				token := &tokenrepo.Token{
					ObtainedAt:  time.Now().Add(-2 * time.Hour),
					AccessToken: "expired-token-redis",
					TokenType:   "Bearer",
					ExpiresIn:   3600, // 1 hora
				}
				data, _ := json.Marshal(token)
				r.Set("auth:token", string(data))
			},
			expectedRes: nil,
			expectedErr: "token expired",
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

			token, err := tokenRepo.GetToken()

			if tt.expectedErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
				assert.Nil(t, token)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, token)
				assert.Equal(t, tt.expectedRes.AccessToken, token.AccessToken)
				assert.Equal(t, tt.expectedRes.TokenType, token.TokenType)
			}
		})
	}
}
