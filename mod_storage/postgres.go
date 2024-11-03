package mod_storage

import (
	"database/sql"
)

type modStoragePostgresRepository struct {
	db *sql.DB
}

func (repo *modStoragePostgresRepository) Get(modname string, key []byte) (*ModStorageEntry, error) {
	row := repo.db.QueryRow("select modname,key,value from mod_storage where modname = $1 and key = $2", modname, key)
	entry := &ModStorageEntry{}
	err := row.Scan(&entry.ModName, &entry.Key, &entry.Value)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return entry, err
}

func (repo *modStoragePostgresRepository) Create(entry *ModStorageEntry) error {
	_, err := repo.db.Exec("insert into mod_storage(modname,key,value) values($1,$2,$3)", entry.ModName, entry.Key, entry.Value)
	return err
}

func (repo *modStoragePostgresRepository) Update(entry *ModStorageEntry) error {
	_, err := repo.db.Exec("update mod_storage set value = $1 where modname = $2 and key = $3", entry.Value, entry.ModName, entry.Key)
	return err
}

func (repo *modStoragePostgresRepository) Delete(modname string, key []byte) error {
	_, err := repo.db.Exec("delete from mod_storage where modname = $1 and key = $2", modname, key)
	return err
}

func (repo *modStoragePostgresRepository) Count() (int64, error) {
	row := repo.db.QueryRow("select count(*) from mod_storage")
	count := int64(0)
	err := row.Scan(&count)
	return count, err
}
