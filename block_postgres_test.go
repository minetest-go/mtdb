package mtdb_test

import (
	"testing"

	"github.com/minetest-go/mtdb"
	"github.com/minetest-go/mtdb/types"
	"github.com/stretchr/testify/assert"
)

func TestPostgresBlocksRepo(t *testing.T) {
	db, err := getPostgresDB(t)
	assert.NoError(t, err)

	assert.NoError(t, mtdb.MigrateBlockDB(db, types.DATABASE_POSTGRES))
	blocks_repo := mtdb.NewBlockRepository(db, types.DATABASE_POSTGRES)
	testBlocksRepository(t, blocks_repo)
}
