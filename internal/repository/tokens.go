package repository

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/nadianeyl/nema-api/internal/domain"
)

type TokenRepository struct {
	DB db
}

func (r *TokenRepository) New(userID uuid.UUID, ttl time.Duration, scope string) (*domain.Token, error) {
	token, err := domain.GenerateToken(userID, ttl, scope)
	if err != nil {
		return nil, err
	}

	err = r.Insert(token)

	return token, err
}

func (r *TokenRepository) Insert(token *domain.Token) error {
	query := `
		INSERT INTO tokens (hash, user_id, expiry, scope)
		VALUES ($1, $2, $3, $4)`

	args := []any{token.Hash, token.UserID, token.Expiry, token.Scope}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := r.DB.ExecContext(ctx, query, args...)

	return err
}

func (r *TokenRepository) DeleteAllForUser(scope string, userID uuid.UUID) error {
	query := `
		DELETE FROM tokens
		WHERE scope = $1 AND user_id = $2`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := r.DB.ExecContext(ctx, query, scope, userID)

	return err
}
