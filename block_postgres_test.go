package mtdb_test

import (
	"testing"

	"github.com/minetest-go/mtdb"
	"github.com/stretchr/testify/assert"
)

func TestPostgresBlocksRepo(t *testing.T) {
	db, err := getPostgresDB(t)
	assert.NoError(t, err)

	assert.NoError(t, mtdb.MigrateBlockDB(db, mtdb.DATABASE_POSTGRES))
	blocks_repo := mtdb.NewBlockRepository(db, mtdb.DATABASE_POSTGRES)
	testBlocksRepository(t, blocks_repo)
}
