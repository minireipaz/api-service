package redisclient_test

import (
	"fmt"
	"minireipaz/pkg/config"
	"minireipaz/pkg/domain/models"
	"minireipaz/pkg/infra/redisclient"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	config.LoadEnvs("../../..")
	os.Exit(m.Run())
}

func Test_InsertUser(t *testing.T) {
	r := redisclient.NewRedisClient()
	userRepo := redisclient.NewUserRedisRepository(r)

	tests := []struct {
		name           string
		user           *models.SyncUserRequest
		setup          func()
		cleanup        func()
		wantInserted   bool
		wantLockExists bool
		wantUserExists bool
		wantErr        bool
	}{
		{
			name: "insertion fails due to unexpected error",
			user: &models.SyncUserRequest{
				Sub:         "1111111111",
				AccessToken: "valid-token",
			},
			setup: func() {
				// Simulate a Redis connection error
				r.Client.Close()
			},
			cleanup: func() {
				// Reconnect to Redis
				r = redisclient.NewRedisClient()
				userRepo = redisclient.NewUserRedisRepository(r)
			},
			wantInserted:   false,
			wantLockExists: false,
			wantUserExists: false,
			wantErr:        true,
		},
		{
			name: "repeated insertion attempt",
			user: &models.SyncUserRequest{
				Sub:         "7777777777",
				AccessToken: "valid-token",
			},
			setup: func() {
				// Primero insertamos el usuario
				userRepo.InsertUser(&models.SyncUserRequest{
					Sub:         "7777777777",
					AccessToken: "valid-token",
				})
			},
			cleanup: func() {
				r.Client.Del(r.Ctx, fmt.Sprintf("lock:user:%s", "7777777777"))
				r.Client.Del(r.Ctx, fmt.Sprintf("users:%s", "7777777777"))
			},
			wantInserted:   false,
			wantLockExists: true,
			wantUserExists: false,
			wantErr:        true,
		},

		{
			name: "successful user insertion",
			user: &models.SyncUserRequest{
				Sub:         "1234567890",
				AccessToken: "valid-token",
			},
			setup: func() {},
			cleanup: func() {
				r.Client.Del(r.Ctx, fmt.Sprintf("lock:user:%s", "1234567890"))
				r.Client.Del(r.Ctx, fmt.Sprintf("users:%s", "1234567890"))
			},
			wantInserted:   true,
			wantLockExists: false,
			wantUserExists: false,
			wantErr:        false,
		},

		// -------------
		{
			name: "concurrent insertion attempt",
			user: &models.SyncUserRequest{
				Sub:         "9999999999",
				AccessToken: "valid-token",
			},
			setup: func() {
				// Simular una inserción concurrente creando el lock manualmente
				lockKey := fmt.Sprintf("lock:user:%s", "9999999999")
				r.Client.SetNX(r.Ctx, lockKey, "dummy", 20*time.Second)
			},
			cleanup: func() {
				r.Client.Del(r.Ctx, fmt.Sprintf("lock:user:%s", "9999999999"))
				r.Client.Del(r.Ctx, fmt.Sprintf("users:%s", "9999999999"))
			},
			wantInserted:   false,
			wantLockExists: true,
			wantUserExists: false,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			defer tt.cleanup()
			inserted, lockExists, userExists, err := userRepo.InsertUser(tt.user)
			if tt.wantErr {
				if tt.wantInserted != inserted || tt.wantLockExists != lockExists || tt.wantUserExists != userExists || err == nil {
					t.Errorf(`Expected wantInserted %v, got %v | wantLockExists %v, got %v | wantUserExists %v, got %v`,
						tt.wantInserted, inserted,
						tt.wantLockExists, lockExists,
						tt.wantUserExists, userExists,
					)
				}
			} else {
				assert.NoError(t, err)
				assert.True(t, inserted)

				// Verify that the user was actually inserted
				exists, err := userRepo.CheckUserExist(tt.user)
				assert.NoError(t, err)
				assert.True(t, exists)
			}
		})
	}
}
