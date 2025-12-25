package repository

import "database/sql"

type Repositories struct {
	Users    UserRepository
	Tokens   TokenRepository
	Accounts AccountRepository
}

func NewRepositories(db *sql.DB) Repositories {
	return Repositories{
		Users:    UserRepository{DB: db},
		Tokens:   TokenRepository{DB: db},
		Accounts: AccountRepository{DB: db},
	}
}
