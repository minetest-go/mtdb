package mtdb

import (
	"database/sql"

	"github.com/minetest-go/mtdb/types"
)

func MigrateModStorageDB(db *sql.DB, dbtype types.DatabaseType) error {
	var err error
	switch dbtype {
	case types.DATABASE_SQLITE:
		_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS entries (
			modname TEXT NOT NULL,
			key BLOB NOT NULL,
			value BLOB NOT NULL,
			PRIMARY KEY (modname, key)
		)`)
	}
	return err
}
