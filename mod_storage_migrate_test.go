package mtdb_test

import (
	"database/sql"
	"testing"

	"github.com/minetest-go/mtdb"
	"github.com/minetest-go/mtdb/types"
	"github.com/stretchr/testify/assert"
)

func TestMigrateModStorageSQlite(t *testing.T) {
	// open db
	db, err := sql.Open("sqlite", ":memory:")
	assert.NoError(t, err)

	assert.NoError(t, mtdb.MigrateModStorageDB(db, types.DATABASE_SQLITE))
}
