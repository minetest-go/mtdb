package mtdb

import (
	"database/sql"
	"fmt"
)

type AuthEntry struct {
	ID        *int64 `json:"id"`
	Name      string `json:"name"`
	Password  string `json:"password"`
	LastLogin int    `json:"last_login"`
}

func NewAuthRepository(db *sql.DB, dbtype DatabaseType) *AuthRepository {
	return &AuthRepository{db: db}
}

type AuthRepository struct {
	db *sql.DB
}

func (repo *AuthRepository) GetByUsername(username string) (*AuthEntry, error) {
	row := repo.db.QueryRow("select id,name,password,last_login from auth where name = $1", username)
	entry := &AuthEntry{}
	err := row.Scan(&entry.ID, &entry.Name, &entry.Password, &entry.LastLogin)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return entry, err
}

type AuthSearch struct {
	Usernamelike *string `json:"usernamelike"`
	Username     *string `json:"username"`
}

func (repo *AuthRepository) buildWhereClause(fields string, s *AuthSearch) (string, []interface{}) {
	q := `select ` + fields + ` from auth where true `
	args := make([]interface{}, 0)
	i := 1

	if s.Username != nil {
		q += fmt.Sprintf(" and name = $%d", i)
		args = append(args, *s.Username)
		i++
	}

	if s.Usernamelike != nil {
		q += fmt.Sprintf(" and name like $%d", i)
		args = append(args, *s.Usernamelike)
		i++
	}

	// limit result length to 1000 per default
	q += " limit 1000"

	return q, args
}

func (repo *AuthRepository) Search(s *AuthSearch) ([]*AuthEntry, error) {
	q, args := repo.buildWhereClause("id,name,password,last_login", s)
	rows, err := repo.db.Query(q, args...)
	if err != nil {
		return nil, err
	}
	list := make([]*AuthEntry, 0)
	for rows.Next() {
		entry := &AuthEntry{}
		err := rows.Scan(&entry.ID, &entry.Name, &entry.Password, &entry.LastLogin)
		if err != nil {
			return nil, err
		}
		list = append(list, entry)
	}
	return list, nil
}

func (repo *AuthRepository) Count(s *AuthSearch) (int, error) {
	q, args := repo.buildWhereClause("count(*)", s)
	row := repo.db.QueryRow(q, args...)
	count := 0
	err := row.Scan(&count)
	return count, err
}

func (repo *AuthRepository) Create(entry *AuthEntry) error {
	row := repo.db.QueryRow("insert into auth(name,password,last_login) values($1,$2,$3) returning id", entry.Name, entry.Password, entry.LastLogin)
	return row.Scan(&entry.ID)
}

func (repo *AuthRepository) Update(entry *AuthEntry) error {
	_, err := repo.db.Exec("update auth set name = $1, password = $2, last_login = $3 where id = $4", entry.Name, entry.Password, entry.LastLogin, entry.ID)
	return err
}

func (repo *AuthRepository) Delete(id int64) error {
	_, err := repo.db.Exec("delete from auth where id = $1", id)
	return err
}
