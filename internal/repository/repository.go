package repository

import (
	"database/sql"

	"github.com/nadianeyl/nema-api/internal/jsonlog"
)

type Repositories struct {
	Users UserRepository
}

func NewRepositories(db *sql.DB, logger *jsonlog.Logger) Repositories {
	return Repositories{
		Users: UserRepository{DB: db},
	}
}
