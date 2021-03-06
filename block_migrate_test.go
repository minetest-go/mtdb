package mtdb_test

import (
	"database/sql"
	"testing"

	_ "modernc.org/sqlite"

	"github.com/minetest-go/mtdb"
	"github.com/stretchr/testify/assert"
)

func TestMigrateBlockSQlite(t *testing.T) {
	// open db
	db, err := sql.Open("sqlite", ":memory:")
	assert.NoError(t, err)

	assert.NoError(t, mtdb.MigrateBlockDB(db, mtdb.DATABASE_SQLITE))
}

func TestMigrateBlockPostgres(t *testing.T) {
	db, err := getPostgresDB(t)
	assert.NoError(t, err)

	assert.NoError(t, mtdb.MigrateBlockDB(db, mtdb.DATABASE_POSTGRES))
}
