package mtdb

import "database/sql"

func MigrateModStorageDB(db *sql.DB, dbtype DatabaseType) error {
	var err error
	switch dbtype {
	case DATABASE_SQLITE:
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
