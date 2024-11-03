package auth

import (
	"database/sql"

	"github.com/minetest-go/mtdb/types"
)

type PrivilegeEntry struct {
	ID        int64  `json:"id"`
	Privilege string `json:"privilege"`
}

type PrivRepository struct {
	db     *sql.DB
	dbtype types.DatabaseType
}

func NewPrivilegeRepository(db *sql.DB, dbtype types.DatabaseType) *PrivRepository {
	return &PrivRepository{db: db, dbtype: dbtype}
}

func (repo *PrivRepository) GetByID(id int64) ([]*PrivilegeEntry, error) {
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

func (repo *PrivRepository) Create(entry *PrivilegeEntry) error {
	_, err := repo.db.Exec("insert into user_privileges(id,privilege) values($1,$2)", entry.ID, entry.Privilege)
	return err
}

func (repo *PrivRepository) Delete(id int64, privilege string) error {
	_, err := repo.db.Exec("delete from user_privileges where id = $1 and privilege = $2", id, privilege)
	return err
}
