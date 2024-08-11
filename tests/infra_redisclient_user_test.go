package tests

import (
	"errors"
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
	config.LoadEnvs("../")
	os.Exit(m.Run())
}

func TestUserRedisRepository_CheckUserExist(t *testing.T) {
	// Crear un cliente Redis real
	r := redisclient.NewRedisClient()

	// Crear un repositorio de usuarios
	userRepo := redisclient.NewUserRedisRepository(r)

	// Definir un usuario de prueba
	user := &models.SyncUserRequest{
		Sub: "test-user",
	}

	tests := []struct {
		name          string
		setup         func()
		expectedErr   error
		expectedExist bool
	}{
		{
			name: "user exists",
			setup: func() {
				// Set up Redis to simulate the user existing
				r.Set(fmt.Sprintf("users:%s", user.Sub), "some_value")
			},
			expectedErr:   nil,
			expectedExist: true,
		},
		{
			name: "user does not exist",
			setup: func() {
				// Ensure the key does not exist in Redis
				r.Client.Del(r.Ctx, fmt.Sprintf("users:%s", user.Sub))
			},
			expectedErr:   nil,
			expectedExist: false,
		},
		{
			name: "redis error on exists check",
			setup: func() {
				// Simulate an error in Redis (e.g., by shutting down the Redis server)
				r.Client.Close() // This will force an error when Exists is called
			},
			expectedErr:   errors.New("ERROR | Cannot check if user exist test-user. More than 10 intents"),
			expectedExist: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			exists, err := userRepo.CheckUserExist(user)

			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedExist, exists)
		})
	}
}

func TestUserRedisRepository_CheckLockExist(t *testing.T) {
	// Crear un cliente Redis real
	r := redisclient.NewRedisClient()

	// Crear un repositorio de usuarios
	userRepo := redisclient.NewUserRedisRepository(r)

	// Definir un usuario de prueba
	user := &models.SyncUserRequest{
		Sub: "test-user",
	}

	tests := []struct {
		name        string
		setup       func()
		expectedErr error
		expectedRes bool
	}{
		{
			name: "lock exists",
			setup: func() {
				// Set up Redis to simulate the lock existing
				lockKey := fmt.Sprintf("lock:users:%s", user.Sub)
				r.Set(lockKey, "locked")
			},
			expectedErr: nil,
			expectedRes: true,
		},
		{
			name: "lock does not exist",
			setup: func() {
				// Ensure the lock key does not exist in Redis
				lockKey := fmt.Sprintf("lock:users:%s", user.Sub)
				r.Client.Del(r.Ctx, lockKey)
			},
			expectedErr: nil,
			expectedRes: false,
		},
		{
			name: "redis error on lock check",
			setup: func() {
				// Simulate an error in Redis (e.g., by shutting down the Redis server)
				r.Client.Close() // This will force an error when Exists is called
			},
			expectedErr: errors.New("ERROR | Cannot check if exist lock for user test-user. More than 10 intents"),
			expectedRes: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			exists, err := userRepo.CheckLockExist(user)

			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedRes, exists)
		})
	}
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

func TestUserRedisRepository_AddLock(t *testing.T) {
	// Crear un cliente Redis real
	r := redisclient.NewRedisClient()

	// Crear un repositorio de usuarios
	userRepo := redisclient.NewUserRedisRepository(r)

	// Definir un usuario de prueba
	user := &models.SyncUserRequest{
		Sub: "test-user",
	}

	tests := []struct {
		name        string
		setup       func()
		expectedErr error
		expectedRes bool
	}{
		{
			name: "successful lock creation",
			setup: func() {
				// Eliminar cualquier lock existente para asegurar que se pueda crear uno nuevo
				lockKey := fmt.Sprintf("lock:user:%s", user.Sub)
				r.Client.Del(r.Ctx, lockKey)
			},
			expectedErr: nil,
			expectedRes: true,
		},
		{
			name: "lock already exists",
			setup: func() {
				// Crear un lock previo para simular un conflicto en la creación
				lockKey := fmt.Sprintf("lock:user:%s", user.Sub)
				r.Set(lockKey, "dummy")
			},
			expectedErr: nil,
			expectedRes: false,
		},
		{
			name: "redis error on lock creation",
			setup: func() {
				r.Client.Close()
			},
			expectedErr: errors.New("ERROR | Cannot create lock for user test-user. More than 10 intents"),
			expectedRes: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			locked, err := userRepo.AddLock(user)

			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedRes, locked)
		})
	}
}

func TestUserRedisRepository_RemoveLock(t *testing.T) {
	r := redisclient.NewRedisClient()

	userRepo := redisclient.NewUserRedisRepository(r)

	user := &models.SyncUserRequest{
		Sub: "test-user",
	}

	tests := []struct {
		name        string
		setup       func()
		expectedRes bool
	}{
		{
			name: "successful lock removal",
			setup: func() {
				lockKey := fmt.Sprintf("lock:user:%s", user.Sub)
				r.Set(lockKey, "dummy")
			},
			expectedRes: true,
		},
		{
			name: "lock already removed",
			setup: func() {
				lockKey := fmt.Sprintf("lock:user:%s", user.Sub)
				r.Client.Del(r.Ctx, lockKey)
			},
			expectedRes: true,
		},
		{
			name: "redis error during lock removal",
			setup: func() {
				r.Client.Close()
			},
			expectedRes: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			removed := userRepo.RemoveLock(user)

			assert.Equal(t, tt.expectedRes, removed)
		})
	}
}
