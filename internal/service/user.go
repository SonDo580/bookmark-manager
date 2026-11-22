package service

import (
	"context"
	"errors"

	"github.com/SonDo580/bookmark-manager/internal/constants"
	"github.com/SonDo580/bookmark-manager/internal/dto"
	"github.com/SonDo580/bookmark-manager/internal/model"
	"github.com/SonDo580/bookmark-manager/internal/repository"
	"github.com/SonDo580/bookmark-manager/pkg/utils"
)

//go:generate mockery --name=UserService --filename user.go --outpkg service_mocks
type UserService interface {
	Register(ctx context.Context, input dto.RegisterRequest) (*model.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

var (
	ErrDuplicateUsername = errors.New(constants.MsgUsernameAlreadyExists)
	ErrDuplicateEmail    = errors.New(constants.MsgEmailAlreadyExists)
)

func (s *userService) Register(ctx context.Context, input dto.RegisterRequest) (*model.User, error) {
	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Username:    input.Username,
		Password:    hashedPassword,
		DisplayName: input.DisplayName,
		Email:       input.Email,
	}

	err = s.repo.CreateUser(ctx, user)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrDuplicateUsername):
			return nil, ErrDuplicateUsername
		case errors.Is(err, repository.ErrDuplicateEmail):
			return nil, ErrDuplicateEmail
		default:
			return nil, err
		}
	}

	return user, nil
}
