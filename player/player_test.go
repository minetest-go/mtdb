package player_test

import (
	"archive/zip"
	"bytes"
	"database/sql"
	"fmt"
	"math/rand"
	"os"
	"testing"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"

	"github.com/minetest-go/mtdb/player"
	"github.com/minetest-go/mtdb/types"
	"github.com/stretchr/testify/assert"
)

func testRepository(t *testing.T, repo *player.PlayerRepository) {
	assert.NotNil(t, repo)

	p1 := &player.Player{
		Name:             fmt.Sprintf("player-%d", rand.Int()),
		Yaw:              1.4,
		Pitch:            0.9,
		PosX:             1200.0,
		PosY:             2300.0,
		PosZ:             4500.0,
		HP:               10,
		Breath:           9,
		CreationDate:     12000,
		ModificationDate: 15000,
	}

	assert.NoError(t, repo.CreateOrUpdate(p1))
	assert.NoError(t, repo.CreateOrUpdate(p1))

	p2, err := repo.GetPlayer(p1.Name)
	assert.NoError(t, err)
	assert.NotNil(t, p2)
	assert.Equal(t, p1.Breath, p2.Breath)
	assert.Equal(t, p1.HP, p2.HP)
	assert.Equal(t, p1.CreationDate, p2.CreationDate)
	assert.Equal(t, p1.ModificationDate, p2.ModificationDate)
	assert.Equal(t, p1.Name, p2.Name)
}

func TestSQliteMigratePlayer(t *testing.T) {
	// open db
	db, err := sql.Open("sqlite3", ":memory:")
	assert.NoError(t, err)

	assert.NoError(t, player.MigratePlayerDB(db, types.DATABASE_SQLITE))
}

func TestSqlitePlayerRepo(t *testing.T) {
	// init stuff
	dbfile, err := os.CreateTemp(os.TempDir(), "players.sqlite")
	assert.NoError(t, err)
	assert.NotNil(t, dbfile)
	copyFileContents("testdata/players.sqlite", dbfile.Name())

	// open db
	db, err := sql.Open("sqlite3", "file:"+dbfile.Name())
	assert.NoError(t, err)
	assert.NoError(t, player.MigratePlayerDB(db, types.DATABASE_SQLITE))
	repo := player.NewPlayerRepository(db, types.DATABASE_SQLITE)
	assert.NotNil(t, repo)

	// count
	player_count, err := repo.Count()
	assert.NoError(t, err)
	assert.Equal(t, int64(1), player_count)

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

	// export
	buf := bytes.NewBuffer([]byte{})
	w := zip.NewWriter(buf)
	err = repo.Export(w)
	assert.NoError(t, err)
	err = w.Close()
	assert.NoError(t, err)
	zipfile, err := os.CreateTemp(os.TempDir(), "player.zip")
	assert.NoError(t, err)
	f, err := os.Create(zipfile.Name())
	assert.NoError(t, err)
	count, err := f.Write(buf.Bytes())
	assert.NoError(t, err)
	assert.True(t, count > 0)

	// import
	z, err := zip.OpenReader(zipfile.Name())
	assert.NoError(t, err)
	err = repo.Import(&z.Reader)
	assert.NoError(t, err)
}

func TestSQlitePlayerRepo2(t *testing.T) {
	// open db
	db, err := sql.Open("sqlite3", ":memory:")
	assert.NoError(t, err)

	assert.NoError(t, player.MigratePlayerDB(db, types.DATABASE_SQLITE))
	repo := player.NewPlayerRepository(db, types.DATABASE_SQLITE)
	testRepository(t, repo)
}

func TestPostgresMigratePlayer(t *testing.T) {
	db, err := getPostgresDB(t)
	assert.NoError(t, err)

	assert.NoError(t, player.MigratePlayerDB(db, types.DATABASE_POSTGRES))
}

func TestPostgresPlayerRepo(t *testing.T) {
	db, err := getPostgresDB(t)
	assert.NoError(t, err)
	assert.NotNil(t, db)

	repo := player.NewPlayerRepository(db, types.DATABASE_POSTGRES)
	testRepository(t, repo)
}
