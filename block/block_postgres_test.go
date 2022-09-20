package block_test

import (
	"testing"

	"github.com/minetest-go/mtdb/block"
	"github.com/minetest-go/mtdb/types"
	"github.com/stretchr/testify/assert"
)

func TestPostgresBlocksRepo(t *testing.T) {
	db, err := getPostgresDB(t)
	assert.NoError(t, err)

	assert.NoError(t, block.MigrateBlockDB(db, types.DATABASE_POSTGRES))
	blocks_repo := block.NewBlockRepository(db, types.DATABASE_POSTGRES)
	testBlocksRepository(t, blocks_repo)
}
