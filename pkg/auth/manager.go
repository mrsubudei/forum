package auth

import (
	"fmt"
	"time"

	"forum/internal/config"

	"github.com/gofrs/uuid"
)

type TokenManager interface {
	NewToken() (string, error)
	UpdateTTL() time.Time
	CheckTTLExpired(TTL time.Time) (bool, error)
}

type Manager struct {
	Cfg config.Config
}

const TimeFormat = "2006-01-02 15:04:05"

func NewManager(cfg config.Config) *Manager {
	return &Manager{
		Cfg: cfg,
	}
}

func (m *Manager) NewToken() (string, error) {
	token, err := uuid.NewV4()
	if err != nil {
		return "", fmt.Errorf("auth - NewToken - NewV4: %w", err)
	}
	return fmt.Sprintf("%v", token), nil
}

func (m *Manager) UpdateTTL() time.Time {
	TTL := time.Now().Add(time.Duration(m.Cfg.TokenManager.SessionExpiringTime * int(time.Second)))
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
