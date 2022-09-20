package block_test

import (
	"database/sql"
	"os"
	"testing"

	_ "modernc.org/sqlite"

	"github.com/minetest-go/mtdb"
	"github.com/minetest-go/mtdb/block"
	"github.com/minetest-go/mtdb/types"
	"github.com/stretchr/testify/assert"
)

func TestSqliteBlockRepo(t *testing.T) {
	// open db
	dbfile, err := os.CreateTemp(os.TempDir(), "map.sqlite")
	assert.NoError(t, err)
	assert.NotNil(t, dbfile)
	db, err := sql.Open("sqlite", "file:"+dbfile.Name())
	assert.NoError(t, err)
	assert.NoError(t, mtdb.EnableWAL(db))

	assert.NoError(t, block.MigrateBlockDB(db, types.DATABASE_SQLITE))
	blocks_repo := block.NewBlockRepository(db, types.DATABASE_SQLITE)
	testBlocksRepository(t, blocks_repo)
}
