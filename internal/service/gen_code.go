package service

import (
	"crypto/rand"
	"math/big"
)

const DefaultCodeLength = 7
const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

//go:generate mockery --name=GenCode --filename gen_code.go --outpkg service_mocks
type GenCode interface {
	GenerateCode(length int) (string, error)
}

type genCode struct{}

func NewGenCode() GenCode {
	return &genCode{}
}

func (s *genCode) GenerateCode(length int) (string, error) {
	code := make([]byte, length)

	for i := range length {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}

		code[i] = charset[n.Int64()]
	}

	return string(code), nil
}
