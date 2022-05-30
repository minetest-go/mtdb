package mtdb

import (
	"database/sql"
	"errors"
)

type postgresAuthRepository struct {
	db *sql.DB
}

func (repo *postgresAuthRepository) GetByUsername(username string) (*AuthEntry, error) {
	row := repo.db.QueryRow("select id,name,password,last_login from auth where name = $1", username)
	entry := &AuthEntry{}
	err := row.Scan(&entry.ID, &entry.Name, &entry.Password, &entry.LastLogin)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return entry, err
}

func (repo *postgresAuthRepository) Create(entry *AuthEntry) error {
	rows, err := repo.db.Query("insert into auth(name,password,last_login) values($1,$2,$3) returning id", entry.Name, entry.Password, entry.LastLogin)
	if err != nil {
		return err
	}
	if !rows.Next() {
		return errors.New("no id returned")
	}
	err = rows.Scan(&entry.ID)
	return err
}

func (repo *postgresAuthRepository) Update(entry *AuthEntry) error {
	_, err := repo.db.Exec("update auth set name = $1, password = $2, last_login = $3 where id = $4", entry.Name, entry.Password, entry.LastLogin, entry.ID)
	return err
}

func (repo *postgresAuthRepository) Delete(id int64) error {
	_, err := repo.db.Exec("delete from auth where id = $1", id)
	return err
}
