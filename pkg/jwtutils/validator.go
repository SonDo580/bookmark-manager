package jwtutils

import (
	"crypto/rsa"
	"errors"
	"os"

	"github.com/SonDo580/bookmark-manager/internal/constants"
	"github.com/golang-jwt/jwt/v5"
)

//go:generate mockery --name=JWTValidator --filename validator.go --outpkg jwtutils_mocks
type JWTValidator interface {
	Validate(tokenStr string) (jwt.MapClaims, error)
}

type validator struct {
	publicKey *rsa.PublicKey
}

func NewJWTValidator(publicKeyPath string) (JWTValidator, error) {
	publicKeyData, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return nil, err
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyData)
	if err != nil {
		return nil, err
	}

	return &validator{
		publicKey: publicKey,
	}, nil
}

var (
	ErrInvalidToken = errors.New(constants.MsgInvalidToken)
	ErrExtractToken = errors.New(constants.MsgCannotExtractToken)
)

func (v *validator) Validate(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
		return v.publicKey, nil
	})

	if err != nil || !token.Valid {
		return nil, ErrInvalidToken
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims, nil
	}

	return nil, ErrExtractToken
}
