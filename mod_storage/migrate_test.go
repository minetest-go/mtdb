package mod_storage_test

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"

	"github.com/minetest-go/mtdb/mod_storage"
	"github.com/minetest-go/mtdb/types"
	"github.com/stretchr/testify/assert"
)

func TestMigrateModStorageSQlite(t *testing.T) {
	// open db
	db, err := sql.Open("sqlite3", ":memory:")
	assert.NoError(t, err)

	assert.NoError(t, mod_storage.MigrateModStorageDB(db, types.DATABASE_SQLITE))
}

func TestMigrateModStoragePostgres(t *testing.T) {
	db := getPostgresDB(t)
	assert.NoError(t, mod_storage.MigrateModStorageDB(db, types.DATABASE_POSTGRES))
}
