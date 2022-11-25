package auth

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid"
)

type TokenManager interface {
	NewToken() (string, error)
	UpdateTTL() time.Time
	CheckTTLExpired(TTL time.Time) (bool, error)
}

type Manager struct {
}

const (
	SessionExpiredTime = 1800 //+n seconds
	TimeFormat         = "2006-01-02 15:04:05"
)

func NewManager() *Manager {
	return &Manager{}
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

func (m *Manager) CheckTTLExpired(TTL time.Time) (bool, error) {
	now := time.Now()
	formatted := now.Format(TimeFormat)
	timeNow, err := time.Parse(TimeFormat, formatted)
	if err != nil {
		return false, fmt.Errorf("auth - CheckTTLExpired - Parse: %w", err)
	}
	return TTL.Before(timeNow), nil
}
