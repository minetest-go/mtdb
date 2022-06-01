package mtdb_test

import (
	"database/sql"
	"os"
	"testing"

	"github.com/minetest-go/mtdb"
	"github.com/stretchr/testify/assert"
)

func TestModStorageSQliteRepo(t *testing.T) {
	// init stuff
	dbfile, err := os.CreateTemp(os.TempDir(), "mod_storage.sqlite")
	assert.NoError(t, err)
	assert.NotNil(t, dbfile)
	copyFileContents("testdata/mod_storage.sqlite", dbfile.Name())

	// open db
	db, err := sql.Open("sqlite", "file:"+dbfile.Name())
	assert.NoError(t, err)
	repo := mtdb.NewModStorageRepository(db, mtdb.DATABASE_SQLITE)
	assert.NotNil(t, repo)

	// existing entry
	entry, err := repo.Get("i3", []byte("data"))
	assert.NoError(t, err)
	assert.NotNil(t, entry)
	assert.Equal(t, []byte("return {[\"singleplayer\"] = {[\"waypoints\"] = {}}}"), entry.Value)

	// create
	assert.NoError(t, repo.Create(&mtdb.ModStorageEntry{
		ModName: "mymod",
		Key:     []byte("mykey"),
		Value:   []byte("myvalue"),
	}))

	//TODO: update/delete
}
