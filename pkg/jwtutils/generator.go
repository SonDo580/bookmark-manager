package jwtutils

import (
	"crypto/rsa"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

//go:generate mockery --name=JWTGenerator --filename generator.go --outpkg jwtutils_mocks
type JWTGenerator interface {
	Generate(payload jwt.MapClaims) (string, error)
}

type generator struct {
	privateKey *rsa.PrivateKey
}

func NewJWTGenerator(privateKeyPath string) (JWTGenerator, error) {
	privateKeyData, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, err
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyData)
	if err != nil {
		return nil, err
	}

	return &generator{
		privateKey: privateKey,
	}, nil
}

func (g *generator) Generate(payload jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, payload)
	return token.SignedString(g.privateKey)
}
