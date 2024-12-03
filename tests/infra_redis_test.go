package tests

import (
	"context"
	"minireipaz/pkg/domain/models"
	"minireipaz/pkg/infra/redisclient"
	"sync"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

var (
	testRedisClient *redis.Client
)

func TestSingleRequest(t *testing.T) {
	actionID := "unique-actionid-1"
	newAction := &models.RequestGoogleAction{
		ActionID:  actionID,
		RequestID: "unique-requestid-1_" + time.Now().UTC().Format(models.LayoutTimestamp),
	}

	testRedisClient = redisclient.NewRedisClient().Client
	rc := &redisclient.RedisClient{Client: testRedisClient}
	repo := redisclient.NewActionsRepository(rc)
	defer testRedisClient.Del(context.Background(), repo.GetActionsGlobalAll(), "lock:"+actionID)
	created, existed, err := repo.Create(newAction)

	assert.NoError(t, err)
	assert.True(t, created)
	assert.False(t, existed)
	// Verify action exists in Redis
	exists, err := testRedisClient.HExists(context.Background(), repo.GetActionsGlobalAll(), actionID).Result()
	assert.NoError(t, err)
	assert.True(t, exists)
}

func TestConcurrentRequests(t *testing.T) {
	actionID := "concurrent-actionid"
	newAction := &models.RequestGoogleAction{
		ActionID:  actionID,
		RequestID: "concurrent-requestid-1_" + time.Now().UTC().Format(models.LayoutTimestamp),
	}
	testRedisClient = redisclient.NewRedisClient().Client
	rc := &redisclient.RedisClient{Client: testRedisClient}
	repo := redisclient.NewActionsRepository(rc)
	defer testRedisClient.Del(context.Background(), repo.GetActionsGlobalAll(), "lock:"+actionID)
	var wg sync.WaitGroup
	numRoutines := 10
	wg.Add(numRoutines)
	for i := 0; i < numRoutines; i++ {
		go func() {
			defer wg.Done()
			repo.Create(newAction)
		}()
	}
	wg.Wait()
	// Verify only one action exists in Redis
	exists, err := testRedisClient.HExists(context.Background(), repo.GetActionsGlobalAll(), actionID).Result()
	assert.NoError(t, err)
	assert.True(t, exists)
	// Check the number of entries with this actionID
	keys, err := testRedisClient.HKeys(context.Background(), repo.GetActionsGlobalAll()).Result()
	assert.NoError(t, err)
	count := 0
	for _, key := range keys {
		if key == actionID {
			count++
		}
	}
	assert.Equal(t, 1, count)
}

func TestRetryMechanism(t *testing.T) {
	actionID := "retry-actionid"
	newAction := &models.RequestGoogleAction{
		ActionID:  actionID,
		RequestID: "unique-requestid-1_" + time.Now().UTC().Format(models.LayoutTimestamp),
	}
	testRedisClient = redisclient.NewRedisClient().Client
	rc := &redisclient.RedisClient{Client: testRedisClient}
	repo := redisclient.NewActionsRepository(rc)
	defer testRedisClient.Del(context.Background(), repo.GetActionsGlobalAll(), "lock:"+actionID)
	// Simulate a scenario where the first attempt fails
	// For example, by introducing a delay or conflicting operations
	// This is a simplified simulation; in real scenarios, you might need more complex setups
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		repo.Create(newAction)
	}()
	go func() {
		defer wg.Done()
		repo.Create(newAction)
	}()
	wg.Wait()
	// Verify only one action exists in Redis
	exists, err := testRedisClient.HExists(context.Background(), repo.GetActionsGlobalAll(), actionID).Result()
	assert.NoError(t, err)
	assert.True(t, exists)
	// Check the number of entries with this actionID
	keys, err := testRedisClient.HKeys(context.Background(), repo.GetActionsGlobalAll()).Result()
	assert.NoError(t, err)
	count := 0
	for _, key := range keys {
		if key == actionID {
			count++
		}
	}
	assert.Equal(t, 1, count)
}

func TestLockExpiry(t *testing.T) {
	actionID := "expiry-actionid"
	newAction := &models.RequestGoogleAction{
		ActionID:  actionID,
		RequestID: "unique-requestid-1_" + time.Now().UTC().Format(models.LayoutTimestamp),
	}
	testRedisClient = redisclient.NewRedisClient().Client
	rc := &redisclient.RedisClient{Client: testRedisClient}
	repo := redisclient.NewActionsRepository(rc)
	defer testRedisClient.Del(context.Background(), repo.GetActionsGlobalAll(), "lock:"+actionID)
	// Create an action with a lock that expires after a short time
	created, existed, err := repo.Create(newAction)
	assert.NoError(t, err)
	assert.True(t, created)
	assert.False(t, existed)
	// Wait for the lock to expire
	time.Sleep(models.MaxTimeForLocks + 1*time.Second)
	// Create the action again
	createdAgain, existedAgain, err := repo.Create(newAction)
	assert.NoError(t, err)
	assert.True(t, createdAgain)
	assert.False(t, existedAgain)
	// Verify only one action exists in Redis
	exists, err := testRedisClient.HExists(context.Background(), repo.GetActionsGlobalAll(), actionID).Result()
	assert.NoError(t, err)
	assert.True(t, exists)
	// Check the number of entries with this actionID
	keys, err := testRedisClient.HKeys(context.Background(), repo.GetActionsGlobalAll()).Result()
	assert.NoError(t, err)
	count := 0
	for _, key := range keys {
		if key == actionID {
			count++
		}
	}
	assert.Equal(t, 1, count)
}

func TestErrorHandling(t *testing.T) {
	actionID := "error-actionid"
	newAction := &models.RequestGoogleAction{
		ActionID:  actionID,
		RequestID: "unique-requestid-1_" + time.Now().UTC().Format(models.LayoutTimestamp),
	}
	testRedisClient = redisclient.NewRedisClient().Client
	rc := &redisclient.RedisClient{Client: testRedisClient}
	repo := redisclient.NewActionsRepository(rc)
	// Introduce an error in Redis operations
	// For example, by closing the Redis connection
	testRedisClient.Close()
	created, existed, err := repo.Create(newAction)
	assert.Error(t, err)
	assert.False(t, created)
	assert.False(t, existed)
	// Reopen the connection
	testRedisClient = redisclient.NewRedisClient().Client
	// Verify that no duplicate action was created
	exists, err := testRedisClient.HExists(context.Background(), repo.GetActionsGlobalAll(), actionID).Result()
	assert.NoError(t, err)
	assert.False(t, exists)
}
