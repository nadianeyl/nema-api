package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID                        uuid.UUID
	Name                      string
	Email                     string
	Password                  password
	Activated                 bool
	EmailNotificationsEnabled bool
	CreatedAt                 time.Time
	UpdatedAt                 time.Time
	Version                   int
}

type password struct {
	Plaintext *string
	Hash      []byte
}

func (p *password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}

	p.Plaintext = &plaintextPassword
	p.Hash = hash

	return nil
}

func (p *password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.Hash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}
