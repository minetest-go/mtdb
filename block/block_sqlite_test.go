package block_test

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"

	"github.com/minetest-go/mtdb/block"
	"github.com/minetest-go/mtdb/types"
	"github.com/minetest-go/mtdb/wal"
	"github.com/stretchr/testify/assert"
)

func setupSqlite(t *testing.T) (block.BlockRepository, *sql.DB) {
	dbfile, err := os.CreateTemp(os.TempDir(), "map.sqlite")
	assert.NoError(t, err)
	assert.NotNil(t, dbfile)
	db, err := sql.Open("sqlite3", "file:"+dbfile.Name())
	assert.NoError(t, err)
	assert.NoError(t, wal.EnableWAL(db))

	assert.NoError(t, block.MigrateBlockDB(db, types.DATABASE_SQLITE))
	blocks_repo, err := block.NewBlockRepository(db, types.DATABASE_SQLITE)
	assert.NoError(t, err)

	return blocks_repo, db
}

func TestSqliteBlockRepo(t *testing.T) {
	// open db
	r, _ := setupSqlite(t)
	defer r.Close()
	testBlocksRepository(t, r)
}

func TestSqliteIterator(t *testing.T) {
	r, _ := setupSqlite(t)
	defer r.Close()
	testBlocksRepositoryIterator(t, r)
}

func TestSqliteIteratorErrorHandling(t *testing.T) {
	r, db := setupSqlite(t)
	defer db.Close()
	defer r.Close()

	testIteratorErrorHandling(t, r, db, `UPDATE blocks SET pos = 18446744073709551615;`)
}

func TestSqliteIteratorCloser(t *testing.T) {
	r, _ := setupSqlite(t)
	defer r.Close()
	testIteratorClose(t, r)
}

func TestCoordToPlain(t *testing.T) {
	nodes := []struct {
		x, y, z int
	}{
		{0, 0, 0},
		{1, -1, 1},
		{-1, -1, -1},

		{30912, 30912, 30912},
		{-30912, -30912, -30912},
	}

	for i, tc := range nodes {
		t.Logf("Test case: #%d", i)

		x1, y1, z1 := block.AsBlockPos(tc.x, tc.y, tc.z)
		pos := block.CoordToPlain(x1, y1, z1)
		x2, y2, z2 := block.PlainToCoord(pos)

		t.Logf("in=%v,%v,%v => pos=%v => out=%v,%v,%v", x1, y1, z1, pos, x2, y2, z2)
		if x1 != x2 || y1 != y2 || z1 != z2 {
			t.Errorf("Unexpected coord returned from pos:"+
				"x=%v,y=%v=z=%v => x=%v, y=%v, z=%v", x1, y1, z1, x2, y2, z2)
		}
	}
}

func TestBlockRepoLegacy(t *testing.T) {
	dbfile, err := os.CreateTemp(os.TempDir(), "map.sqlite")
	assert.NoError(t, err)
	assert.NotNil(t, dbfile)
	copyFileContents("testdata/map_legacy_column.sqlite", dbfile.Name())

	db, err := sql.Open("sqlite3", fmt.Sprintf("file:%s", dbfile.Name()))
	assert.NoError(t, err)
	assert.NoError(t, wal.EnableWAL(db))

	repo, err := block.NewBlockRepository(db, types.DATABASE_SQLITE)
	assert.NoError(t, err)
	assert.NotNil(t, repo)

	b, err := repo.GetByPos(0, 0, 0)
	assert.NoError(t, err)
	assert.NotNil(t, b)
	assert.True(t, len(b.Data) > 0)
}

func TestBlockRepoMultiColumn(t *testing.T) {
	dbfile, err := os.CreateTemp(os.TempDir(), "map.sqlite")
	assert.NoError(t, err)
	assert.NotNil(t, dbfile)
	copyFileContents("testdata/map_multi_pos_column.sqlite", dbfile.Name())

	db, err := sql.Open("sqlite3", fmt.Sprintf("file:%s", dbfile.Name()))
	assert.NoError(t, err)
	assert.NoError(t, wal.EnableWAL(db))

	repo, err := block.NewBlockRepository(db, types.DATABASE_SQLITE)
	assert.NoError(t, err)
	assert.NotNil(t, repo)

	b, err := repo.GetByPos(0, 0, 0)
	assert.NoError(t, err)
	assert.NotNil(t, b)
	assert.True(t, len(b.Data) > 0)
}
