package repository

import (
	"context"
	"errors"

	"github.com/SonDo580/bookmark-manager/internal/constants"
	"github.com/SonDo580/bookmark-manager/internal/model"
	dbPkg "github.com/SonDo580/bookmark-manager/pkg/db"
	"gorm.io/gorm"
)

//go:generate mockery --name=UserRepository --filename user.go --outpkg repository_mocks
type UserRepository interface {
	CreateUser(ctx context.Context, user *model.User) error
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)
	GetUserByID(ctx context.Context, id string) (*model.User, error)
	UpdateUser(ctx context.Context, id, displayName, email string) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

var (
	ErrDuplicateUsername = errors.New(constants.MsgUsernameAlreadyExists)
	ErrDuplicateEmail    = errors.New(constants.MsgEmailAlreadyExists)
	ErrUserNotFound      = errors.New(constants.MsgUserNotFound)
)

var (
	duplicateUsernameFilter = dbPkg.UniqueViolationFilter("uq_users_username", "users.username", ErrDuplicateUsername)
	duplicateEmailFilter    = dbPkg.UniqueViolationFilter("uq_users_email", "users.email", ErrDuplicateEmail)
	userNotFoundFilter      = dbPkg.NoRecordFilter(ErrUserNotFound)
)

func (r *userRepository) CreateUser(ctx context.Context, user *model.User) error {
	err := r.db.WithContext(ctx).Create(user).Error

	if err != nil {
		return dbPkg.CatchDBError(
			[]dbPkg.ErrorFilter{
				duplicateUsernameFilter,
				duplicateEmailFilter,
			},
			err,
		)
	}

	return nil
}

func (r *userRepository) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	user := &model.User{}

	err := r.db.WithContext(ctx).Where("username = ?", username).First(user).Error
	if err != nil {
		return nil, dbPkg.CatchDBError(
			[]dbPkg.ErrorFilter{userNotFoundFilter},
			err,
		)
	}

	return user, nil
}

func (r *userRepository) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	user := &model.User{}

	err := r.db.WithContext(ctx).Where("id = ?", id).First(user).Error
	if err != nil {
		return nil, dbPkg.CatchDBError(
			[]dbPkg.ErrorFilter{userNotFoundFilter},
			err,
		)
	}

	return user, nil
}

func (r *userRepository) UpdateUser(ctx context.Context, id, displayName, email string) error {
	result := r.db.WithContext(ctx).
		Model(&model.User{}).
		Where("id = ?", id).
		Select("DisplayName", "Email").
		Updates(model.User{
			DisplayName: displayName,
			Email:       email,
		})

	if result.Error != nil {
		return dbPkg.CatchDBError(
			[]dbPkg.ErrorFilter{duplicateEmailFilter},
			result.Error,
		)
	}

	if result.RowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}
