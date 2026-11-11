package repository

import (
	"testing"
	"time"

	redisPkg "github.com/SonDo580/bookmark-manager/pkg/redis"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestURLStorage_StoreURL(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		key            string
		url            string
		exp            time.Duration
		trySecondStore bool
	}{
		{
			name: "no expiration",
			key:  "abcd1234",
			url:  "https://google.com",
			exp:  0,
		},
		{
			name: "with expiration",
			key:  "abcd1234",
			url:  "https://google.com",
			exp:  10 * time.Second,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()

			mockRedisServer, mockRedisClient := redisPkg.InitMockRedis(t)
			defer mockRedisClient.Close()

			testRepo := NewURLStorage(mockRedisClient)

			// First store: success
			success, err := testRepo.StoreURL(ctx, tc.key, tc.url, tc.exp)
			assert.NoError(t, err)
			assert.Equal(t, success, true)

			// Second store: should not store duplicate key
			success, err = testRepo.StoreURL(ctx, tc.key, tc.url, tc.exp)
			assert.NoError(t, err)
			assert.Equal(t, success, false)

			// Get with stored key
			res, err := mockRedisClient.Get(ctx, tc.key).Result()
			assert.NoError(t, err)
			assert.Equal(t, tc.url, res)

			// Get with non-existent key: should fail
			res, err = mockRedisClient.Get(ctx, tc.key+"x").Result()
			assert.ErrorIs(t, err, redis.Nil)

			// Simulate time passing
			if tc.exp > 0 {
				mockRedisServer.FastForward(tc.exp + 1*time.Second)

				// Get after stored key expires - should fail
				res, err = mockRedisClient.Get(ctx, tc.key).Result()
				assert.ErrorIs(t, err, redis.Nil)
			}
		})
	}
}

func TestURLStorage_GetURL(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		key  string
		url  string
		exp  time.Duration
	}{
		{
			name: "normal case",
			key:  "abc1234",
			url:  "https://google.com",
			exp:  0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()

			mockRedisServer, mockRedisClient := redisPkg.InitMockRedis(t)
			defer mockRedisClient.Close()

			testRepo := NewURLStorage(mockRedisClient)

			// Store
			err := mockRedisClient.Set(ctx, tc.key, tc.url, tc.exp).Err()
			require.NoError(t, err)

			// Get stored key
			url, err := testRepo.GetURL(ctx, tc.key)
			assert.NoError(t, err)
			assert.Equal(t, tc.url, url)

			// Get non-existent key
			url, err = testRepo.GetURL(ctx, tc.key+"x")
			assert.Empty(t, url)
			assert.ErrorIs(t, err, ErrURLNotFound)

			// Trigger other error (not redis.Nil)
			mockRedisServer.SetError("any")
			url, err = testRepo.GetURL(ctx, tc.key)
			assert.Empty(t, url)
			assert.Error(t, err)
			assert.NotErrorIs(t, err, ErrURLNotFound)
			mockRedisServer.SetError("")
		})
	}
}
