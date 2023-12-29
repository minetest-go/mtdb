package mod_storage_test

import (
	"testing"

	"github.com/minetest-go/mtdb/mod_storage"
	"github.com/minetest-go/mtdb/types"
	"github.com/stretchr/testify/assert"
)

func TestModStoragePostgresRepo(t *testing.T) {
	// open db
	db := getPostgresDB(t)
	repo := mod_storage.NewModStorageRepository(db, types.DATABASE_POSTGRES)
	assert.NotNil(t, repo)

	// cleanup
	_, err := db.Exec("delete from mod_storage")
	assert.NoError(t, err)

	// create
	entry := &mod_storage.ModStorageEntry{
		ModName: "mymod",
		Key:     []byte("mykey"),
		Value:   []byte("myvalue"),
	}
	assert.NoError(t, repo.Create(entry))

	// count
	entry_count, err := repo.Count()
	assert.NoError(t, err)
	assert.Equal(t, int64(1), entry_count)

	// delete
	assert.NoError(t, repo.Delete("mymod", []byte("mykey")))
	entry, err = repo.Get("mymod", []byte("mykey"))
	assert.NoError(t, err)
	assert.Nil(t, entry)
}
