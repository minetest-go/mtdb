package mtdb

import (
	"database/sql"
	"testing"

	_ "modernc.org/sqlite"

	"github.com/stretchr/testify/assert"
)

func TestMigrateBlocksSQlite(t *testing.T) {
	// open db
	db, err := sql.Open("sqlite", ":memory:")
	assert.NoError(t, err)

	assert.NoError(t, MigrateBlockDB(db, DATABASE_SQLITE))
}

func TestMigrateBlocksPostgres(t *testing.T) {
	db, err := getPostgresDB(t)
	assert.NoError(t, err)

	assert.NoError(t, MigrateBlockDB(db, DATABASE_POSTGRES))
}
