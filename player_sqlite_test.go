package mtdb_test

import (
	"database/sql"
	"os"
	"testing"

	"github.com/minetest-go/mtdb"
	"github.com/stretchr/testify/assert"
)

func TestSQliteMigratePlayer(t *testing.T) {
	// open db
	db, err := sql.Open("sqlite", ":memory:")
	assert.NoError(t, err)

	assert.NoError(t, mtdb.MigratePlayerDB(db, mtdb.DATABASE_SQLITE))
}

func TestSqlitePlayerRepo(t *testing.T) {
	// init stuff
	dbfile, err := os.CreateTemp(os.TempDir(), "players.sqlite")
	assert.NoError(t, err)
	assert.NotNil(t, dbfile)
	copyFileContents("testdata/players.sqlite", dbfile.Name())

	// open db
	db, err := sql.Open("sqlite", "file:"+dbfile.Name())
	assert.NoError(t, err)
	assert.NoError(t, mtdb.MigratePlayerDB(db, mtdb.DATABASE_SQLITE))
	repo := mtdb.NewPlayerRepository(db, mtdb.DATABASE_SQLITE)
	assert.NotNil(t, repo)

	// existing entry
	p, err := repo.GetPlayer("singleplayer")
	assert.NoError(t, err)
	assert.NotNil(t, p)
	assert.Equal(t, "singleplayer", p.Name)
	assert.InDelta(t, 7.17000007629395, p.Pitch, 0.01)
	assert.InDelta(t, 272.760009765625, p.Yaw, 0.01)
	assert.InDelta(t, -1631.63000488281, p.PosX, 0.01)
	assert.InDelta(t, 196.210006713867, p.PosY, 0.01)
	assert.InDelta(t, -430.440002441406, p.PosZ, 0.01)
	assert.Equal(t, 20, p.HP)
	assert.Equal(t, 10, p.Breath)
	assert.Equal(t, int64(1648301850), p.CreationDate)
	assert.Equal(t, int64(1652728478), p.ModificationDate)

	// non-existing entry
	p, err = repo.GetPlayer("dummy")
	assert.NoError(t, err)
	assert.Nil(t, p)
}
