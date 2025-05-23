package mod_storage_test

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"

	"github.com/minetest-go/mtdb/mod_storage"
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
	db, err := sql.Open("sqlite3", "file:"+dbfile.Name())
	assert.NoError(t, err)
	repo := mod_storage.NewModStorageRepository(db, types.DATABASE_SQLITE)
	assert.NotNil(t, repo)

	// existing entry
	entry, err := repo.Get("i3", []byte("data"))
	assert.NoError(t, err)
	assert.NotNil(t, entry)
	assert.Equal(t, []byte("return {[\"singleplayer\"] = {[\"waypoints\"] = {}}}"), entry.Value)

	// create
	entry = &mod_storage.ModStorageEntry{
		ModName: "mymod",
		Key:     []byte("mykey"),
		Value:   []byte("myvalue"),
	}
	assert.NoError(t, repo.Create(entry))

	// count
	entry_count, err := repo.Count()
	assert.NoError(t, err)
	assert.Equal(t, int64(3), entry_count)

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
