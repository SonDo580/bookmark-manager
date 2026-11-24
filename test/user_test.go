package test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/SonDo580/bookmark-manager/internal/api"
	"github.com/SonDo580/bookmark-manager/internal/constants"
	"github.com/SonDo580/bookmark-manager/internal/dto"
	"github.com/SonDo580/bookmark-manager/internal/model"
	dbPkg "github.com/SonDo580/bookmark-manager/pkg/db"
	redisPkg "github.com/SonDo580/bookmark-manager/pkg/redis"
)

func TestUser_Register(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	cfg := getDefaultConfig()

	const (
		correctMethod   = http.MethodPost
		correctEndpoint = "/v1/users/register"

		testUsername    = "shadow"
		testEmail       = "shadow@x.com"
		testPassword    = "examplePass"
		testDisplayName = "Shadow"

		otherUsername       = "sonic"
		otherEmail          = "sonic@x.com"
		otherHashedPassword = "$2a$10$<hashed-password>"
	)

	newInput := func(username, email string) dto.RegisterRequest {
		return dto.RegisterRequest{
			Username:    username,
			Email:       email,
			Password:    testPassword,
			DisplayName: testDisplayName,
		}
	}
	newDefaultInput := func() dto.RegisterRequest {
		return newInput(testUsername, testEmail)
	}

	testCases := []struct {
		name             string
		setupDB          func(ctx context.Context, t *testing.T, db *gorm.DB)
		setupTestRequest func(t *testing.T) *http.Request
		expectedStatus   int
		expectedErrMsg   string
	}{
		{
			name:    "success",
			setupDB: func(ctx context.Context, t *testing.T, db *gorm.DB) {},
			setupTestRequest: func(t *testing.T) *http.Request {
				return newJSONRequest(t, correctMethod, correctEndpoint, newDefaultInput())
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:    "fail: invalid request body (lack required fields)",
			setupDB: func(ctx context.Context, t *testing.T, db *gorm.DB) {},
			setupTestRequest: func(t *testing.T) *http.Request {
				return newJSONRequest(t, correctMethod, correctEndpoint, dto.RegisterRequest{})
			},
			expectedStatus: http.StatusBadRequest,
			expectedErrMsg: constants.MsgInvalidInput,
		},
		{
			name: "fail: duplicate username",
			setupDB: func(ctx context.Context, t *testing.T, db *gorm.DB) {
				err := db.WithContext(ctx).Create(&model.User{
					Username:    testUsername, // same as default input
					Email:       otherEmail,
					Password:    otherHashedPassword,
					DisplayName: testDisplayName,
				}).Error
				require.NoError(t, err)
			},
			setupTestRequest: func(t *testing.T) *http.Request {
				return newJSONRequest(t, correctMethod, correctEndpoint, newDefaultInput())
			},
			expectedStatus: http.StatusConflict,
			expectedErrMsg: constants.MsgUsernameAlreadyExists,
		},
		{
			name: "fail: duplicate email",
			setupDB: func(ctx context.Context, t *testing.T, db *gorm.DB) {
				err := db.WithContext(ctx).Create(&model.User{
					Username:    otherUsername,
					Email:       testEmail, // same as default input
					Password:    otherHashedPassword,
					DisplayName: testDisplayName,
				}).Error
				require.NoError(t, err)
			},
			setupTestRequest: func(t *testing.T) *http.Request {
				return newJSONRequest(t, correctMethod, correctEndpoint, newDefaultInput())
			},
			expectedStatus: http.StatusConflict,
			expectedErrMsg: constants.MsgEmailAlreadyExists,
		},
		{
			name: "fail: internal server error",
			setupDB: func(ctx context.Context, t *testing.T, db *gorm.DB) {
				// Drop the "users" table
				err := db.WithContext(ctx).Migrator().DropTable(&model.User{})
				require.NoError(t, err)
			},
			setupTestRequest: func(t *testing.T) *http.Request {
				return newJSONRequest(t, correctMethod, correctEndpoint, newDefaultInput())
			},
			expectedStatus: http.StatusInternalServerError,
			expectedErrMsg: constants.MsgInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			_, mockRedisClient := redisPkg.InitMockRedis(t)
			defer mockRedisClient.Close()

			testDB := dbPkg.CreateTestDB(t)
			err := testDB.AutoMigrate(&model.User{})
			require.NoError(t, err)

			ctx := t.Context()
			tc.setupDB(ctx, t, testDB)

			testAPIEngine := api.NewEngine(cfg, mockRedisClient, testDB)
			rec := httptest.NewRecorder()
			req := tc.setupTestRequest(t)
			testAPIEngine.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedStatus, rec.Code)

			switch tc.expectedStatus {
			case http.StatusCreated:
				var resp dto.RegisterResponse
				err := json.Unmarshal(rec.Body.Bytes(), &resp)
				require.NoError(t, err)

				assert.NotEmpty(t, resp.Message)
				assert.NotEmpty(t, resp.Data.ID)
				assert.Empty(t, resp.Data.Password)
				assert.Equal(t, testUsername, resp.Data.Username)
				assert.Equal(t, testEmail, resp.Data.Email)
				assert.Equal(t, testDisplayName, resp.Data.DisplayName)
			default:
				var resp dto.ErrorResponse
				err := json.Unmarshal(rec.Body.Bytes(), &resp)
				require.NoError(t, err)

				assert.Equal(t, tc.expectedErrMsg, resp.Message)
			}
		})
	}
}
