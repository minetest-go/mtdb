package mtdb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPostgresBlocksRepo(t *testing.T) {
	db, err := getPostgresDB(t)
	assert.NoError(t, err)

	assert.NoError(t, MigrateBlockDB(db, DATABASE_POSTGRES))
	blocks_repo := NewBlockRepository(db, DATABASE_POSTGRES)
	testBlocksRepository(t, blocks_repo)
}
