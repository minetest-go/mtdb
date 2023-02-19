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

func TestMaxConnections(t *testing.T) {
	db, err := getPostgresDB(t)
	assert.NoError(t, err)

	assert.NoError(t, block.MigrateBlockDB(db, types.DATABASE_POSTGRES))
	blocks_repo := block.NewBlockRepository(db, types.DATABASE_POSTGRES)
	assert.NotNil(t, blocks_repo)

	var maxConnections int
	row := db.QueryRow("show max_connections;")
	err = row.Scan(&maxConnections)
	assert.NoError(t, err)
	t.Logf("Testing against %v max connections", maxConnections)

	fakeBlock := block.Block{
		PosX: 1,
		PosY: 1,
		PosZ: 1,
		Data: []byte("test"),
	}
	if err := blocks_repo.Update(&fakeBlock); err != nil {
		t.Fatalf("Error setting up test case: %v", err)
	}

	// Run more than max_connections query operations in a loop
	count := 0
	for i := 0; i < maxConnections*2; i++ {
		b, err := blocks_repo.GetByPos(1, 1, 1)
		count++
		if b != nil && count%10 == 0 {
			t.Logf("Executed %d operations. b=%v", count, b)
		}
		if err != nil {
			t.Errorf("Unexpected error after %d operations: %v", count, err)
			break
		}
	}
}
