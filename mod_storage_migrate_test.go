package mtdb_test

import (
	"database/sql"
	"testing"

	"github.com/minetest-go/mtdb"
	"github.com/stretchr/testify/assert"
)

func TestMigrateModStorageSQlite(t *testing.T) {
	// open db
	db, err := sql.Open("sqlite", ":memory:")
	assert.NoError(t, err)

	assert.NoError(t, mtdb.MigrateModStorageDB(db, mtdb.DATABASE_SQLITE))
}
