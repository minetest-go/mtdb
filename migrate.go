package mtdb

import "database/sql"

func MigrateAuthSqlite(db *sql.DB) error {
	_, err := db.Exec(`
	CREATE TABLE if not exists
		auth (id INTEGER PRIMARY KEY AUTOINCREMENT,name VARCHAR(32) UNIQUE,password VARCHAR(512),last_login INTEGER);
	CREATE TABLE if not exists
		user_privileges (id INTEGER,privilege VARCHAR(32),PRIMARY KEY (id, privilege)CONSTRAINT fk_id FOREIGN KEY (id) REFERENCES auth (id) ON DELETE CASCADE);
	`)
	return err
}

func MigrateAuthPostgres(db *sql.DB) error {
	_, err := db.Exec(`
	CREATE TABLE if not exists
		auth (id SERIAL,name TEXT UNIQUE,password TEXT,last_login INT NOT NULL DEFAULT 0,PRIMARY KEY (id));
	CREATE TABLE if not exists
		user_privileges (id INT,privilege TEXT,PRIMARY KEY (id, privilege),CONSTRAINT fk_id FOREIGN KEY (id) REFERENCES auth (id) ON DELETE CASCADE);
	`)
	return err
}
