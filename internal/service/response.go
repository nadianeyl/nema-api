package service

import (
	"time"

	"github.com/google/uuid"
)

type UserResponse struct {
	ID                        uuid.UUID `json:"id"`
	Name                      string    `json:"name"`
	Email                     string    `json:"email"`
	Activated                 bool      `json:"activated"`
	EmailNotificationsEnabled bool      `json:"email_notifications_enabled"`
	CreatedAt                 time.Time `json:"created_at"`
	UpdatedAt                 time.Time `json:"updated_at"`
}

type TokenResponse struct {
	Plaintext string    `json:"token"`
	Expiry    time.Time `json:"expiry"`
}
