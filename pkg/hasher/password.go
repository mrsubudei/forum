package hasher

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type PasswordHasher interface {
	Hash(password string) (string, error)
	CheckPassword(existHashed, entered string) error
}

type BcryptHasher struct {
}

func NewBcryptHasher() *BcryptHasher {
	return &BcryptHasher{}
}

func (b *BcryptHasher) Hash(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("hasher - Hash - GenerateFromPassword: %w", err)
	}
	return string(hashed), nil
}

func (b *BcryptHasher) CheckPassword(existHashed, entered string) error {
	return bcrypt.CompareHashAndPassword([]byte(existHashed), []byte(entered))
}
