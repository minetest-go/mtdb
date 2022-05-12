package mtdb

import (
	"database/sql"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSqliteBlockRepo(t *testing.T) {
	// open db
	dbfile, err := os.CreateTemp(os.TempDir(), "map.sqlite")
	assert.NoError(t, err)
	assert.NotNil(t, dbfile)
	db, err := sql.Open("sqlite", "file:"+dbfile.Name())
	assert.NoError(t, err)
	assert.NoError(t, EnableWAL(db))

	assert.NoError(t, MigrateBlockDB(db, DATABASE_SQLITE))
	blocks_repo := NewBlockRepository(db, DATABASE_SQLITE)
	testBlocksRepository(t, blocks_repo)
}
