package block_test

import (
	"database/sql"
	"testing"

	"github.com/minetest-go/mtdb/block"
	"github.com/minetest-go/mtdb/types"
	"github.com/stretchr/testify/assert"
)

func setupPostgress(t *testing.T) (block.BlockRepository, *sql.DB) {
	db := getPostgresDB(t)

	// Cleanup any previous data
	db.Exec("delete from blocks")

	assert.NoError(t, block.MigrateBlockDB(db, types.DATABASE_POSTGRES))
	blocks_repo, err := block.NewBlockRepository(db, types.DATABASE_POSTGRES)
	if err != nil {
		panic(err)
	}

	assert.NotNil(t, blocks_repo)
	return blocks_repo, db
}

func TestPostgresBlocksRepo(t *testing.T) {
	r, _ := setupPostgress(t)
	testBlocksRepository(t, r)
}

func TestPostgresMaxConnections(t *testing.T) {
	r, db := setupPostgress(t)

	var maxConnections int
	row := db.QueryRow("show max_connections;")
	err := row.Scan(&maxConnections)
	assert.NoError(t, err)
	t.Logf("Testing against %v max connections", maxConnections)

	fakeBlock := block.Block{
		PosX: 1,
		PosY: 1,
		PosZ: 1,
		Data: []byte("test"),
	}
	if err := r.Update(&fakeBlock); err != nil {
		t.Fatalf("Error setting up test case: %v", err)
	}

	// Run more than max_connections query operations in a loop
	count := 0
	for i := 0; i < maxConnections*2; i++ {
		b, err := r.GetByPos(1, 1, 1)
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

func TestPostgresIterator(t *testing.T) {
	logToTesting(t)
	blocks_repo, _ := setupPostgress(t)
	testBlocksRepositoryIterator(t, blocks_repo)
}

func TestPostgresIteratorBatches(t *testing.T) {
	logToTesting(t)

	oldSize := block.IteratorBatchSize
	setUp := func() {
		block.IteratorBatchSize = 1
	}
	tearDown := func() {
		block.IteratorBatchSize = oldSize
	}

	setUp()
	defer tearDown()
	blocks_repo, _ := setupPostgress(t)
	testBlocksRepositoryIterator(t, blocks_repo)
}

func TestPostgresIteratorErrorHandling(t *testing.T) {
	blocks_repo, db := setupPostgress(t)
	defer db.Close()

	testIteratorErrorHandling(t, blocks_repo, db, `
		ALTER TABLE blocks ALTER COLUMN posX TYPE float;
		UPDATE blocks SET posX = 18446744073709551615;`)
}

func TestPostgresIteratorCloser(t *testing.T) {
	r, _ := setupPostgress(t)
	testIteratorClose(t, r)
}
