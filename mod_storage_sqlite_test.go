package mtdb_test

import (
	"database/sql"
	"os"
	"testing"

	"github.com/minetest-go/mtdb"
	"github.com/minetest-go/mtdb/types"
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
	repo := mtdb.NewModStorageRepository(db, types.DATABASE_SQLITE)
	assert.NotNil(t, repo)

	// existing entry
	entry, err := repo.Get("i3", []byte("data"))
	assert.NoError(t, err)
	assert.NotNil(t, entry)
	assert.Equal(t, []byte("return {[\"singleplayer\"] = {[\"waypoints\"] = {}}}"), entry.Value)

	// create
	entry = &mtdb.ModStorageEntry{
		ModName: "mymod",
		Key:     []byte("mykey"),
		Value:   []byte("myvalue"),
	}
	assert.NoError(t, repo.Create(entry))

	// update
	entry.Value = []byte("othervalue")
	assert.NoError(t, repo.Update(entry))

	entry2, err := repo.Get("mymod", []byte("mykey"))
	assert.NoError(t, err)
	assert.NotNil(t, entry2)
	assert.Equal(t, entry.Value, entry2.Value)

	// delete
	assert.NoError(t, repo.Delete("mymod", []byte("mykey")))
	entry, err = repo.Get("mymod", []byte("mykey"))
	assert.NoError(t, err)
	assert.Nil(t, entry)

}
