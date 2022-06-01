package mtdb

import "database/sql"

// internal name: "entries"
type ModStorageEntry struct {
	ModName string `json:"modname"`
	Key     []byte `json:"key"`
	Value   []byte `json:"value"`
}

type ModStorageRepository interface {
	Get(modname string, key []byte) (*ModStorageEntry, error)
	Create(entry *ModStorageEntry) error
	Update(entry *ModStorageEntry) error
	Delete(modname string, key []byte) error
}

func NewModStorageRepository(db *sql.DB, dbtype DatabaseType) ModStorageRepository {
	switch dbtype {
	case DATABASE_SQLITE:
		return &modStorageSqliteRepository{db: db}
	default:
		return nil
	}
}
