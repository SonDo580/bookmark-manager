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
)

func (r *userRepository) CreateUser(ctx context.Context, user *model.User) error {
	err := r.db.WithContext(ctx).Create(user).Error

	if err != nil {
		return dbPkg.CatchDBError([]dbPkg.ErrorFilter{
			dbPkg.PGUniqueViolationFilter("idx_users_username", ErrDuplicateUsername),
			dbPkg.PGUniqueViolationFilter("idx_users_email", ErrDuplicateEmail),
		}, err)
	}

	return nil
}
