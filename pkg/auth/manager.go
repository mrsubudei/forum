package auth

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid"
)

type TokenManager interface {
	NewToken() (string, error)
	UpdateTTL() time.Time
	CheckTTL(TTL time.Time) bool
}

type Manager struct {
}

const SessionExpiredTime = 900 //seconds

func NewManager() (*Manager, error) {
	return &Manager{}, nil
}

func (m *Manager) NewToken() (string, error) {
	token, err := uuid.NewV4()
	if err != nil {
		return "", fmt.Errorf("auth - NewToken - NewV4: %w", err)
	}
	return fmt.Sprintf("%v", token), nil
}

func (m *Manager) UpdateTTL() time.Time {
	TTL := time.Now().Add(SessionExpiredTime * time.Second)
	return TTL
}

func (m *Manager) CheckTTL(TTL time.Time) bool {
	return TTL.Before(time.Now())
}
