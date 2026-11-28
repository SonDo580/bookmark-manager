package service

import (
	"context"
	"errors"
	"time"

	"github.com/SonDo580/bookmark-manager/internal/config"
	"github.com/SonDo580/bookmark-manager/internal/constants"
	"github.com/SonDo580/bookmark-manager/internal/dto"
	"github.com/SonDo580/bookmark-manager/internal/model"
	"github.com/SonDo580/bookmark-manager/internal/repository"
	"github.com/SonDo580/bookmark-manager/pkg/jwtutils"
	"github.com/SonDo580/bookmark-manager/pkg/utils"
	"github.com/golang-jwt/jwt/v5"
)

//go:generate mockery --name=UserService --filename user.go --outpkg service_mocks
type UserService interface {
	Register(ctx context.Context, input dto.RegisterRequest) (*model.User, error)
	Login(ctx context.Context, username, password string) (string, error)
	GetSelfInfoByID(ctx context.Context, id string) (*model.User, error)
	EditSelfInfo(ctx context.Context, id string, input dto.EditSelfInfoRequest) (updated bool, err error)
}

type userService struct {
	repo         repository.UserRepository
	jwtGenerator jwtutils.JWTGenerator
	cfg          *config.Config
}

func NewUserService(
	repo repository.UserRepository,
	jwtGenerator jwtutils.JWTGenerator,
	cfg *config.Config,
) UserService {
	return &userService{
		repo:         repo,
		jwtGenerator: jwtGenerator,
		cfg:          cfg,
	}
}

var (
	ErrDuplicateUsername = errors.New(constants.MsgUsernameAlreadyExists)
	ErrDuplicateEmail    = errors.New(constants.MsgEmailAlreadyExists)
	ErrWrongCredentials  = errors.New(constants.MsgWrongCredentials)
	ErrUserNotFound      = errors.New(constants.MsgUserNotFound)
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

func (s *userService) Login(ctx context.Context, username, password string) (string, error) {
	// Get user
	user, err := s.repo.GetUserByUsername(ctx, username)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrUserNotFound):
			return "", ErrWrongCredentials
		default:
			return "", err
		}
	}

	// Compare password
	passwordValid, err := utils.CompareHashAndPassword(user.Password, password)
	if err != nil {
		return "", err
	}
	if !passwordValid {
		return "", ErrWrongCredentials
	}

	// Generate and return token
	now := time.Now()
	return s.jwtGenerator.Generate(jwt.MapClaims{
		"sub": user.ID,
		"iat": now.Unix(),
		"exp": now.Add(s.cfg.TokenExp).Unix(),
	})
}

func (s *userService) GetSelfInfoByID(ctx context.Context, id string) (*model.User, error) {
	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrUserNotFound):
			return nil, ErrUserNotFound
		default:
			return nil, err
		}
	}

	return user, nil
}

func (s *userService) EditSelfInfo(ctx context.Context, id string, input dto.EditSelfInfoRequest) (updated bool, err error) {
	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return false, ErrUserNotFound
		}
		return false, err
	}

	if user.DisplayName == input.DisplayName && user.Email == input.Email {
		return false, nil // no changes
	}

	err = s.repo.UpdateUser(ctx, id, input.DisplayName, input.Email)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrUserNotFound):
			return false, ErrUserNotFound
		case errors.Is(err, repository.ErrDuplicateEmail):
			return false, ErrDuplicateEmail
		default:
			return false, err
		}
	}

	return true, nil
}
