package auth

import (
	"errors"
	"time"
)

type TokenManager interface {
	NewToken(userId string, ttl time.Duration) (string, error)
}

type Manager struct {
	signingKey string
}

func NewManager(signingKey string) (*Manager, error) {
	if signingKey == "" {
		return nil, errors.New("empty signing key")
	}

	return &Manager{signingKey: signingKey}, nil
}

func (m *Manager) NewToken(userId string, ttl time.Duration) (string, error) {
	var err error
	return "", err
}
