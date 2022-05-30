package mtdb_test

import (
	"database/sql"
	"testing"

	_ "modernc.org/sqlite"

	"github.com/minetest-go/mtdb"
	"github.com/stretchr/testify/assert"
)

func TestMigrateAuthSQlite(t *testing.T) {
	// open db
	db, err := sql.Open("sqlite", ":memory:")
	assert.NoError(t, err)

	assert.NoError(t, mtdb.MigrateAuthDB(db, mtdb.DATABASE_SQLITE))
}

func TestMigrateAuthPostgres(t *testing.T) {
	db, err := getPostgresDB(t)
	assert.NoError(t, err)

	assert.NoError(t, mtdb.MigrateAuthDB(db, mtdb.DATABASE_POSTGRES))
}
