package mod_storage

import (
	"archive/zip"
	"bufio"
	"bytes"
	"database/sql"
	"encoding/json"
)

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

func (repo *modStorageSqliteRepository) Export(z *zip.Writer) error {
	w, err := z.Create("mod_storage.json")
	if err != nil {
		return err
	}
	enc := json.NewEncoder(w)

	rows, err := repo.db.Query("select modname,key,value from entries")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		e := &ModStorageEntry{}
		err = rows.Scan(&e.ModName, &e.Key, &e.Value)
		if err != nil {
			return err
		}

		err = enc.Encode(e)
		if err != nil {
			return err
		}
	}

	return nil
}

func (repo *modStorageSqliteRepository) Import(z *zip.Reader) error {
	f, err := z.Open("mod_storage.json")
	if err != nil {
		return err
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		dc := json.NewDecoder(bytes.NewReader(sc.Bytes()))
		e := &ModStorageEntry{}
		err = dc.Decode(e)
		if err != nil {
			return err
		}

		_, err := repo.db.Exec("insert or replace into entries(modname,key,value) values($1,$2,$3)", e.ModName, e.Key, e.Value)
		if err != nil {
			return err
		}
	}

	return nil
}
