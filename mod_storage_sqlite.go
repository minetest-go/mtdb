package mtdb

import "database/sql"

type modStorageSqliteRepository struct {
	db *sql.DB
}

func (repo *modStorageSqliteRepository) Get(modname string, key []byte) (*ModStorageEntry, error) {
	row := repo.db.QueryRow("select modname,key,value from entries where modname = $1 and key = $2", modname, key)
	entry := &ModStorageEntry{}
	err := row.Scan(&entry.ModName, &entry.Key, &entry.Value)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return entry, err
}

func (repo *modStorageSqliteRepository) Create(entry *ModStorageEntry) error {
	_, err := repo.db.Exec("insert into entries(modname,key,value) values($1,$2,$3)", entry.ModName, entry.Key, entry.Value)
	return err
}

func (repo *modStorageSqliteRepository) Update(entry *ModStorageEntry) error {
	_, err := repo.db.Exec("update entries set value = $1 where modname = $2 and key = $3", entry.Value, entry.ModName, entry.Key)
	return err
}

func (repo *modStorageSqliteRepository) Delete(modname string, key []byte) error {
	_, err := repo.db.Exec("delete from entries where modname = $1 and key = $2", modname, key)
	return err
}
