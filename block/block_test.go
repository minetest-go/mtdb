package block_test

import (
	"archive/zip"
	"bytes"
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/minetest-go/mtdb/block"
	"github.com/stretchr/testify/assert"
)

func testBlocksRepository(t *testing.T, block_repo block.BlockRepository) {

	// cleanup
	assert.NoError(t, block_repo.Delete(0, 0, 0))

	// get nil
	b, err := block_repo.GetByPos(0, 0, 0)
	assert.NoError(t, err)
	assert.Nil(t, b)

	// create
	b = &block.Block{
		PosX: 0,
		PosY: 0,
		PosZ: 0,
		Data: []byte{0x00, 0x01, 0x02},
	}
	assert.NoError(t, block_repo.Update(b))

	// count
	blocks, err := block_repo.Count()
	assert.NoError(t, err)
	assert.Equal(t, int64(1), blocks)

	// get
	b, err = block_repo.GetByPos(0, 0, 0)
	assert.NoError(t, err)
	assert.NotNil(t, b)
	assert.Equal(t, 0, b.PosX)
	assert.Equal(t, 0, b.PosY)
	assert.Equal(t, 0, b.PosZ)
	assert.Equal(t, 3, len(b.Data))
	assert.Equal(t, uint8(0x00), b.Data[0])
	assert.Equal(t, uint8(0x01), b.Data[1])
	assert.Equal(t, uint8(0x02), b.Data[2])

	// export
	buf := bytes.NewBuffer([]byte{})
	w := zip.NewWriter(buf)
	err = block_repo.Export(w)
	assert.NoError(t, err)
	err = w.Close()
	assert.NoError(t, err)
	zipfile, err := os.CreateTemp(os.TempDir(), "blocks.zip")
	assert.NoError(t, err)
	f, err := os.Create(zipfile.Name())
	assert.NoError(t, err)
	count, err := f.Write(buf.Bytes())
	assert.NoError(t, err)
	assert.True(t, count > 0)

	// delete
	assert.NoError(t, block_repo.Delete(0, 0, 0))

	// get nil
	b, err = block_repo.GetByPos(0, 0, 0)
	assert.NoError(t, err)
	assert.Nil(t, b)

	// vacuum
	assert.NoError(t, block_repo.Vacuum())

	// import
	z, err := zip.OpenReader(zipfile.Name())
	assert.NoError(t, err)
	err = block_repo.Import(&z.Reader)
	assert.NoError(t, err)
}

func testBlocksRepositoryIterator(t *testing.T, blocks_repo block.BlockRepository) {
	negX, negY, negZ := block.AsBlockPos(-32000, -32000, -32000)
	posX, posY, posZ := block.AsBlockPos(32000, 32000, 32000)

	testData := []block.Block{
		{negX, negY, negZ, []byte("negative")},
		{0, 0, 0, []byte("zero")},
		{posX, posY, posZ, []byte("positive")},
	}
	for i := range testData {
		b := testData[i]
		blocks_repo.Update(&b)
	}
	// logrus.SetLevel(logrus.DebugLevel)

	// Helper function to loop over all channel data
	// TODO(ronoaldo): simplify this so we don't have
	// to work so verbose. Perhaps implement a wrapper
	// to the channel?
	consumeAll := func(it chan *block.Block) int {
		t.Logf("consumeAll: fetching data from iterator")
		count := 0
		for {
			select {
			case b, ok := <-it:
				if !ok {
					return count
				}
				t.Logf("consumeAll: got %#v", b)
				count++
			case <-time.After(3 * time.Second):
				t.Errorf("consumeAll: timed out")
				return count
			}
			if count > 10 {
				panic("consumeAll: too many items returned from channel")
			}
		}
	}

	// Fetch from neg -> pos, retrieves all three blocks
	it, err := blocks_repo.Iterator(negX, negY, negZ)
	if assert.NoError(t, err) {
		assert.Equal(t, 3, consumeAll(it))
	}

	// Fetch from zero -> pos, retrieves two blocks
	it, err = blocks_repo.Iterator(0, 0, 0)
	if assert.NoError(t, err) {
		assert.Equal(t, 2, consumeAll(it))
	}

	// Fetch from zero +1 -> pos, retrieves only one
	it, err = blocks_repo.Iterator(0, 0, 1)
	if assert.NoError(t, err) {
		assert.Equal(t, 1, consumeAll(it))
	}

	// Fetch from 2000,2000,2000, retrieves zero blocks
	it, err = blocks_repo.Iterator(posX+1, posY+1, posZ+1)
	if assert.NoError(t, err) {
		assert.Equal(t, 0, consumeAll(it))
	}
}

func testIteratorErrorHandling(t *testing.T, blocks_repo block.BlockRepository, db *sql.DB, mockDataCorruption string) {
	setUp := func() {
		if err := blocks_repo.Update(&block.Block{1, 2, 3, []byte("default:stone")}); err != nil {
			t.Fatalf("setUp: error loading test data: %v", err)
		}

		// Forcing an error during iterator loop
		if _, err := db.Exec(mockDataCorruption); err != nil {
			t.Fatalf("Error renaming column: %v", err)
		}
	}

	tearDown := func() {
		if _, err := db.Exec("DROP TABLE blocks"); err != nil {
			t.Fatalf("tearDown: error resetting test db: %v", err)
		}
	}

	setUp()
	defer tearDown()

	ch, err := blocks_repo.Iterator(0, 0, 0)
	if err != nil {
		t.Fatalf("Error loading the iterator: %v", err)
	}

	count := 0
	for b := range ch {
		t.Logf("Block: %v", b)
		count++
	}

	assert.Equal(t, 0, count, "should not return any blocks when data is corrupted")
}
