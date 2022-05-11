package mtdb

import (
	"database/sql"
	"errors"
)

func MigrateAuth(db *sql.DB, dbtype string) error {
	var err error
	switch dbtype {
	case "postgres":
		_, err = db.Exec(`
		CREATE TABLE if not exists
			auth (id SERIAL,name TEXT UNIQUE,password TEXT,last_login INT NOT NULL DEFAULT 0,PRIMARY KEY (id));
		CREATE TABLE if not exists
			user_privileges (id INT,privilege TEXT,PRIMARY KEY (id, privilege),CONSTRAINT fk_id FOREIGN KEY (id) REFERENCES auth (id) ON DELETE CASCADE);
		`)
	case "sqlite":
		_, err = db.Exec(`
		CREATE TABLE if not exists
			auth (id INTEGER PRIMARY KEY AUTOINCREMENT,name VARCHAR(32) UNIQUE,password VARCHAR(512),last_login INTEGER);
		CREATE TABLE if not exists
			user_privileges (id INTEGER,privilege VARCHAR(32),PRIMARY KEY (id, privilege)CONSTRAINT fk_id FOREIGN KEY (id) REFERENCES auth (id) ON DELETE CASCADE);
		`)
	}
	return err
}

type AuthEntry struct {
	ID        *int64 `json:"id"`
	Name      string `json:"name"`
	Password  string `json:"password"`
	LastLogin int    `json:"last_login"`
}

type DBAuthRepository struct {
	db     *sql.DB
	dbtype DatabaseType
}

func NewAuthRepository(db *sql.DB, dbtype DatabaseType) *DBAuthRepository {
	return &DBAuthRepository{db: db, dbtype: dbtype}
}

func (repo *DBAuthRepository) GetByUsername(username string) (*AuthEntry, error) {
	rows, err := repo.db.Query("select id,name,password,last_login from auth where name = $1", username)
	if err != nil {
		return nil, err
	}
	if !rows.Next() {
		return nil, nil
	}
	entry := &AuthEntry{}
	err = rows.Scan(&entry.ID, &entry.Name, &entry.Password, &entry.LastLogin)
	return entry, err
}

func (repo *DBAuthRepository) Create(entry *AuthEntry) error {
	var err error
	switch repo.dbtype {
	case DATABASE_POSTGRES:
		rows, err := repo.db.Query("insert into auth(name,password,last_login) values($1,$2,$3) returning id", entry.Name, entry.Password, entry.LastLogin)
		if err != nil {
			return err
		}
		if !rows.Next() {
			return errors.New("no id returned")
		}
		err = rows.Scan(&entry.ID)
	case DATABASE_SQLITE:
		res, err := repo.db.Exec("insert into auth(name,password,last_login) values($1,$2,$3)", entry.Name, entry.Password, entry.LastLogin)
		if err != nil {
			return err
		}
		id, err := res.LastInsertId()
		entry.ID = &id
	}
	return err
}

func (repo *DBAuthRepository) Update(entry *AuthEntry) error {
	_, err := repo.db.Exec("update auth set name = $1, password = $2, last_login = $3 where id = $4", entry.Name, entry.Password, entry.LastLogin, entry.ID)
	return err
}

func (repo *DBAuthRepository) Delete(id int64) error {
	_, err := repo.db.Exec("delete from auth where id = $1", id)
	return err
}
