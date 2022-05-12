package mtdb

import "database/sql"

func MigrateBlockDB(db *sql.DB, dbtype DatabaseType) error {
	var err error
	switch dbtype {
	case DATABASE_SQLITE:
		_, err = db.Exec(`CREATE TABLE IF NOT EXISTS blocks (pos INT PRIMARY KEY, data BLOB)`)
	case DATABASE_POSTGRES:
		_, err = db.Exec(`CREATE TABLE IF NOT EXISTS
			blocks (posX INT NOT NULL, posY INT NOT NULL, posZ INT NOT NULL, data BYTEA, PRIMARY KEY (posX,posY,posZ))`)
	}
	return err
}
