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

type partialUser struct {
	ID          string
	Username    string
	Password    string
	Email       string
	DisplayName string
}

func TestUserRepository_CreateUser(t *testing.T) {
	t.Parallel()

	const (
		testUsername       = "shadow"
		testEmail          = "shadow@x.com"
		testDisplayName    = "Shadow"
		testHashedPassword = "$2a$10$<hashed-password>"
	)

	const (
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

		var stored partialUser
		findErr := db.WithContext(ctx).
			Table("users").
			Where("username = ?", expected.Username).
			First(&stored).Error
		require.NoError(t, findErr)

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
