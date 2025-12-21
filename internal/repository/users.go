package repository

import (
	"context"
	"database/sql"
	"time"
)

type UserRepository struct {
	DB *sql.DB
}

func (r *UserRepository) Insert(user *User) error {
	query := `
		INSERT INTO users (name, email, password_hash, activated, email_notifications_enabled)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at, version`

	args := []any{user.Name, user.Email, user.Password.hash, user.Activated, user.EmailNotificationsEnabled}

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
			return ErrDuplicateEmail
		default:
			return err
		}
	}

	return nil
}
