package repository

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	ScopeActivation     = "activation"
	ScopeAuthentication = "authentication"
)

var (
	ErrDuplicateEmail         = errors.New("duplicate email")
	ErrRecordNotFound         = errors.New("record not found")
	ErrEditConflict           = errors.New("edit conflict")
	ErrInvalidOrExpiredToken  = errors.New("invalid or expired token")
	ErrInvalidAuthCredentials = errors.New("invalid authentication credentials")
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
	plaintext *string
	hash      []byte
}

func (p *password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}

	p.plaintext = &plaintextPassword
	p.hash = hash

	return nil
}

func (p *password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintextPassword))
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

type Token struct {
	Plaintext string
	Hash      []byte
	UserID    uuid.UUID
	Expiry    time.Time
	Scope     string
}

func generateToken(userID uuid.UUID, ttl time.Duration, scope string) (*Token, error) {
	token := &Token{
		UserID: userID,
		Expiry: time.Now().Add(ttl),
		Scope:  scope,
	}

	randomBytes := make([]byte, 16)

	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}

	token.Plaintext = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)

	hash := sha256.Sum256([]byte(token.Plaintext))
	token.Hash = hash[:]

	return token, nil
}
