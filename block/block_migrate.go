package block

import (
	"database/sql"

	"github.com/minetest-go/mtdb/types"
)

func MigrateBlockDB(db *sql.DB, dbtype types.DatabaseType) error {
	var err error
	switch dbtype {
	case types.DATABASE_SQLITE:
		_, err = db.Exec(`CREATE TABLE IF NOT EXISTS blocks (pos INT PRIMARY KEY, data BLOB)`)
	case types.DATABASE_POSTGRES:
		_, err = db.Exec(`CREATE TABLE IF NOT EXISTS
			blocks (posX INT NOT NULL, posY INT NOT NULL, posZ INT NOT NULL, data BYTEA, PRIMARY KEY (posX,posY,posZ))`)
	}
	return err
}
