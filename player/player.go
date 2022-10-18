package player

import (
	"database/sql"

	"github.com/minetest-go/mtdb/types"
)

func NewPlayerRepository(db *sql.DB, dbtype types.DatabaseType) PlayerRepository {
	switch dbtype {
	case types.DATABASE_SQLITE:
		return &PlayerSqliteRepository{db: db}
	case types.DATABASE_POSTGRES:
		return &PlayerPostgresRepository{db: db}
	}
	return nil
}
