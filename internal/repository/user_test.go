package repository

import (
	"context"
	"testing"

	"github.com/SonDo580/bookmark-manager/internal/model"
	dbPkg "github.com/SonDo580/bookmark-manager/pkg/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestUserRepository_CreateUser(t *testing.T) {
	t.Parallel()

	const (
		testUsername       = "shadow"
		testEmail          = "shadow@x.com"
		testDisplayName    = "Shadow"
		testHashedPassword = "$2a$10$<hashed-password>"

		otherUsername = "sonic"
		otherEmail    = "sonic@x.com"
	)

	newUser := func(username, email string) *model.User {
		return &model.User{
			Username:    username,
			Email:       email,
			DisplayName: testDisplayName,
			Password:    testHashedPassword,
		}
	}

	testCases := []struct {
		name        string
		setupDB     func(ctx context.Context, repo UserRepository)
		user        *model.User
		expectedErr error
	}{
		{
			name:        "success",
			setupDB:     func(ctx context.Context, repo UserRepository) {},
			user:        newUser(testUsername, testEmail),
			expectedErr: nil,
		},
		{
			name: "fail: duplicate username",
			setupDB: func(ctx context.Context, repo UserRepository) {
				// Create a user with the same username in advance
				err := repo.CreateUser(ctx, newUser(testUsername, otherEmail))
				require.NoError(t, err)
			},
			user:        newUser(testUsername, testEmail),
			expectedErr: ErrDuplicateUsername,
		},
		{
			name: "fail: duplicate email",
			setupDB: func(ctx context.Context, repo UserRepository) {
				// Create a user with the same email in advance
				err := repo.CreateUser(ctx, newUser(otherUsername, testEmail))
				require.NoError(t, err)
			},
			user:        newUser(testUsername, testEmail),
			expectedErr: ErrDuplicateEmail,
		},
	}

	verifyCreatedUser := func(t *testing.T, ctx context.Context, db *gorm.DB, expected *model.User) {
		t.Helper()

		var stored model.User
		err := db.WithContext(ctx).
			Where("username = ?", expected.Username).
			First(&stored).Error
		require.NoError(t, err)

		assert.NotEmpty(t, stored.ID)
		assert.Equal(t, expected.Username, stored.Username)
		assert.Equal(t, expected.Email, stored.Email)
		assert.Equal(t, expected.DisplayName, stored.DisplayName)
		assert.Equal(t, expected.Password, stored.Password)
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()

			testDB := dbPkg.CreateTestDB(t)
			err := testDB.AutoMigrate(&model.User{})
			require.NoError(t, err)

			repo := NewUserRepository(testDB)
			tc.setupDB(ctx, repo)

			err = repo.CreateUser(ctx, tc.user)

			if tc.expectedErr != nil {
				assert.ErrorIs(t, err, tc.expectedErr)
			} else {
				assert.NoError(t, err)
				verifyCreatedUser(t, ctx, testDB, tc.user)
			}
		})
	}
}

func setupTestDB(t *testing.T, ctx context.Context, shouldSeed bool, fixtures []*model.User) *gorm.DB {
	t.Helper()

	testDB := dbPkg.CreateTestDB(t)
	err := testDB.AutoMigrate(&model.User{})
	require.NoError(t, err)

	if shouldSeed {
		for _, fixture := range fixtures {
			err = testDB.WithContext(ctx).Create(fixture).Error
			require.NoError(t, err)
		}
	}

	return testDB
}

func assertRetrievedUser(t *testing.T, retrievedUser *model.User, expectedUser *model.User) {
	t.Helper()

	require.NotNil(t, retrievedUser)
	assert.Equal(t, expectedUser.ID, retrievedUser.ID)
	assert.Equal(t, expectedUser.Username, retrievedUser.Username)
	assert.Equal(t, expectedUser.Email, retrievedUser.Email)
	assert.Equal(t, expectedUser.DisplayName, retrievedUser.DisplayName)
	assert.Equal(t, expectedUser.Password, retrievedUser.Password)
}

func TestUserRepository_GetUserByID(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name       string
		shouldSeed bool
	}{
		{
			name:       "success",
			shouldSeed: true,
		},
		{
			name:       "not found",
			shouldSeed: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			fixture := &model.User{
				ID:          "<user_id>",
				Username:    "shadow",
				Email:       "shadow@x.com",
				DisplayName: "Shadow",
				Password:    "<hashed-password>",
			}

			testDB := setupTestDB(t, ctx, tc.shouldSeed, []*model.User{fixture})
			repo := NewUserRepository(testDB)

			retrievedUser, err := repo.GetUserByID(ctx, fixture.ID)

			if !tc.shouldSeed {
				assert.ErrorIs(t, err, ErrUserNotFound)
				assert.Nil(t, retrievedUser)
			} else {
				require.NoError(t, err)
				assertRetrievedUser(t, retrievedUser, fixture)
			}
		})
	}
}

func TestUserRepository_GetUserByUsername(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name       string
		shouldSeed bool
	}{
		{
			name:       "success",
			shouldSeed: true,
		},
		{
			name:       "not found",
			shouldSeed: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			fixture := &model.User{
				ID:          "<user_id>",
				Username:    "shadow",
				Email:       "shadow@x.com",
				DisplayName: "Shadow",
				Password:    "<hashed-password>",
			}

			testDB := setupTestDB(t, ctx, tc.shouldSeed, []*model.User{fixture})
			repo := NewUserRepository(testDB)

			retrievedUser, err := repo.GetUserByUsername(ctx, fixture.Username)

			if !tc.shouldSeed {
				assert.ErrorIs(t, err, ErrUserNotFound)
				assert.Nil(t, retrievedUser)
			} else {
				require.NoError(t, err)
				assertRetrievedUser(t, retrievedUser, fixture)
			}
		})
	}
}

func TestUserRepository_UpdateUser(t *testing.T) {
	t.Parallel()

	const (
		user1ID          = "<user_1_id>"
		user1Username    = "shadow"
		user1Email       = "shadow@x.com"
		user1DisplayName = "Shadow"
		user1Password    = "<hashed-password-1>"

		user2ID          = "<user_2_id>"
		user2Username    = "sonic"
		user2Email       = "sonic@x.com"
		user2DisplayName = "Sonic"
		user2Password    = "<hashed-password-2>"

		otherUserID      = "<user_3_id>"
		otherDisplayName = "Silver"
		otherEmail       = "silver@x.com"
	)

	testCases := []struct {
		name           string
		targetID       string
		newDisplayName string
		newEmail       string
		expectedErr    error
	}{
		{
			name:           "success: same displayName and same email",
			targetID:       user1ID,
			newDisplayName: user1DisplayName,
			newEmail:       user1Email,
			expectedErr:    nil,
		},
		{
			name:           "success: different displayName and same email",
			targetID:       user1ID,
			newDisplayName: otherDisplayName,
			newEmail:       user1Email,
			expectedErr:    nil,
		},
		{
			name:           "success: different displayName and different email",
			targetID:       user1ID,
			newDisplayName: otherDisplayName,
			newEmail:       otherEmail,
			expectedErr:    nil,
		},
		{
			name:           "duplicate email: update user 1 with email of user 2",
			targetID:       user1ID,
			newDisplayName: otherDisplayName,
			newEmail:       user2Email,
			expectedErr:    ErrDuplicateEmail,
		},
		{
			name:           "not found: non-existent user ID",
			targetID:       otherUserID,
			newDisplayName: otherDisplayName,
			newEmail:       otherEmail,
			expectedErr:    ErrUserNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()

			user1 := &model.User{
				ID:          user1ID,
				Username:    user1Username,
				Email:       user1Email,
				DisplayName: user1DisplayName,
				Password:    user1Password,
			}

			user2 := &model.User{
				ID:          user2ID,
				Username:    user2Username,
				Email:       user2Email,
				DisplayName: user2DisplayName,
				Password:    user2Password,
			}

			testDB := setupTestDB(t, ctx, true, []*model.User{user1, user2})
			repo := NewUserRepository(testDB)

			err := repo.UpdateUser(ctx, tc.targetID, tc.newDisplayName, tc.newEmail)

			if tc.expectedErr != nil {
				assert.ErrorIs(t, err, tc.expectedErr)
			} else {
				require.NoError(t, err)

				targetUser, exists := map[string]*model.User{
					user1.ID: user1,
					user2.ID: user2,
				}[tc.targetID]
				require.True(t, exists)

				var storedUser model.User
				err = testDB.WithContext(ctx).
					Where("id = ?", tc.targetID).
					First(&storedUser).Error
				require.NoError(t, err)

				expectedUser := &model.User{
					ID:          targetUser.ID,
					Username:    targetUser.Username,
					Email:       tc.newEmail,
					DisplayName: tc.newDisplayName,
					Password:    targetUser.Password,
				}
				assertRetrievedUser(t, &storedUser, expectedUser)
			}
		})
	}
}
