package mtdb

import "database/sql"

type AuthEntry struct {
	ID        *int64 `json:"id"`
	Name      string `json:"name"`
	Password  string `json:"password"`
	LastLogin int    `json:"last_login"`
}

type AuthRepository interface {
	GetByUsername(username string) (*AuthEntry, error)
	Create(entry *AuthEntry) error
	Update(entry *AuthEntry) error
	Delete(id int64) error
}

func NewAuthRepository(db *sql.DB, dbtype DatabaseType) AuthRepository {
	switch dbtype {
	case DATABASE_SQLITE:
		return &SqliteAuthRepository{db: db}
	case DATABASE_POSTGRES:
		return &PostgresAuthRepository{db: db}
	default:
		return nil
	}
}
