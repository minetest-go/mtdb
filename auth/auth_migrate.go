package auth

import (
	"database/sql"

	"github.com/minetest-go/mtdb/types"
)

func MigrateAuthDB(db *sql.DB, dbtype types.DatabaseType) error {
	var err error
	switch dbtype {
	case types.DATABASE_SQLITE:
		_, err = db.Exec(`
		CREATE TABLE if not exists
			auth (id INTEGER PRIMARY KEY AUTOINCREMENT,name VARCHAR(32) UNIQUE,password VARCHAR(512),last_login INTEGER);
		CREATE TABLE if not exists
			user_privileges (id INTEGER,privilege VARCHAR(32),PRIMARY KEY (id, privilege)CONSTRAINT fk_id FOREIGN KEY (id) REFERENCES auth (id) ON DELETE CASCADE);
		`)
	case types.DATABASE_POSTGRES:
		_, err = db.Exec(`
		CREATE TABLE if not exists
			auth (id SERIAL,name TEXT UNIQUE,password TEXT,last_login INT NOT NULL DEFAULT 0,PRIMARY KEY (id));
		CREATE TABLE if not exists
			user_privileges (id INT,privilege TEXT,PRIMARY KEY (id, privilege),CONSTRAINT fk_id FOREIGN KEY (id) REFERENCES auth (id) ON DELETE CASCADE);
		`)
	}
	return err
}
