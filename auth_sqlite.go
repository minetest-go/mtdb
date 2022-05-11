package mtdb

import (
	"database/sql"
	"errors"
)

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
