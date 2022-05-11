package mtdb

import (
	"database/sql"
)

type PrivilegeEntry struct {
	ID        int64  `json:"id"`
	Privilege string `json:"privilege"`
}

type DBPrivRepository struct {
	db     *sql.DB
	dbtype DatabaseType
}

func NewPrivilegeRepository(db *sql.DB, dbtype DatabaseType) *DBPrivRepository {
	return &DBPrivRepository{db: db, dbtype: dbtype}
}

func (repo *DBPrivRepository) GetByID(id int64) ([]*PrivilegeEntry, error) {
	rows, err := repo.db.Query("select id,privilege from user_privileges where id = $1", id)
	if err != nil {
		return nil, err
	}
	list := make([]*PrivilegeEntry, 0)
	for rows.Next() {
		entry := &PrivilegeEntry{}
		err = rows.Scan(&entry.ID, &entry.Privilege)
		if err != nil {
			return nil, err
		}
		list = append(list, entry)
	}
	return list, nil
}

func (repo *DBPrivRepository) Create(entry *PrivilegeEntry) error {
	_, err := repo.db.Exec("insert into user_privileges(id,privilege) values($1,$2)", entry.ID, entry.Privilege)
	return err
}

func (repo *DBPrivRepository) Delete(id int64, privilege string) error {
	_, err := repo.db.Exec("delete from user_privileges where id = $1 and privilege = $2", id, privilege)
	return err
}
