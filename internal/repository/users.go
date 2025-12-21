package repository

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"time"

	"github.com/nadianeyl/nema-api/internal/domain"
)

type UserRepository struct {
	DB *sql.DB
}

func (r *UserRepository) Insert(user *domain.User) error {
	query := `
		INSERT INTO users (name, email, password_hash, activated, email_notifications_enabled)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at, version`

	args := []any{user.Name, user.Email, user.Password.Hash, user.Activated, user.EmailNotificationsEnabled}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := r.DB.QueryRowContext(ctx, query, args...).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Version,
	)

	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return domain.ErrDuplicateEmail
		default:
			return err
		}
	}

	return nil
}

func (r *UserRepository) GetByEmail(email string) (*domain.User, error) {
	query := `
		SELECT id, name, email, password_hash, activated, email_notifications_enabled, created_at, updated_at, version
		FROM users
		WHERE email = $1`

	var user domain.User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := r.DB.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password.Hash,
		&user.Activated,
		&user.EmailNotificationsEnabled,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, domain.ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (r UserRepository) Update(user *domain.User) error {
	query := `
		UPDATE users
		SET name = $1, email = $2, password_hash = $3, activated = $4, email_notifications_enabled = $5, updated_at = $6, version = version + 1
		WHERE id = $7 AND version = $8
		RETURNING updated_at, version`

	args := []any{
		user.Name,
		user.Email,
		user.Password.Hash,
		user.Activated,
		user.EmailNotificationsEnabled,
		time.Now(),
		user.ID,
		user.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := r.DB.QueryRowContext(ctx, query, args...).Scan(&user.UpdatedAt, &user.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return domain.ErrDuplicateEmail
		case errors.Is(err, sql.ErrNoRows):
			return domain.ErrEditConflict
		default:
			return err
		}
	}

	return nil
}

func (r UserRepository) GetForToken(tokenScope, tokenPlaintext string) (*domain.User, error) {
	tokenHash := sha256.Sum256([]byte(tokenPlaintext))

	query := `
		SELECT users.id, users.name, users.email, users.password_hash, users.activated, users.email_notifications_enabled, users.created_at, users.updated_at, users.version
		FROM users
		INNER JOIN tokens
		ON users.id = tokens.user_id
		WHERE tokens.hash = $1
		AND tokens.scope = $2
		AND tokens.expiry > $3`

	args := []any{tokenHash[:], tokenScope, time.Now()}

	var user domain.User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := r.DB.QueryRowContext(ctx, query, args...).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password.Hash,
		&user.Activated,
		&user.EmailNotificationsEnabled,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, domain.ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}
