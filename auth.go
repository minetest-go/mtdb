package mtdb

import (
	"database/sql"
)

type AuthEntry struct {
	ID        *int64 `json:"id"`
	Name      string `json:"name"`
	Password  string `json:"password"`
	LastLogin int    `json:"last_login"`
}

type DBAuthRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) *DBAuthRepository {
	return &DBAuthRepository{db: db}
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
	result, err := repo.db.Exec("insert into auth(id,name,password,last_login) values($1,$2,$3,$4)", entry.ID, entry.Name, entry.Password, entry.LastLogin)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	// assign returned id
	entry.ID = &id
	return nil
}

func (repo *DBAuthRepository) Update(entry *AuthEntry) error {
	_, err := repo.db.Exec("update auth set name = $1, password = $2, last_login = $3 where id = $4", entry.Name, entry.Password, entry.LastLogin, entry.ID)
	return err
}

func (repo *DBAuthRepository) Delete(id int64) error {
	_, err := repo.db.Exec("delete from auth where id = $1", id)
	return err
}
