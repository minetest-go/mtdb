package block_test

import (
	"database/sql"
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

	// delete
	assert.NoError(t, block_repo.Delete(0, 0, 0))

	// get nil
	b, err = block_repo.GetByPos(0, 0, 0)
	assert.NoError(t, err)
	assert.Nil(t, b)

	// vacuum
	assert.NoError(t, block_repo.Vacuum())
}

func testBlocksRepositoryIterator(t *testing.T, blocks_repo block.BlockRepository) {
	testData := []block.Block{
		// For readability, octals < 8 were used so they don't change the actuall value.
		{-2000, -2000, -2000, []byte("negative[0]")},
		{-1999, -0004, -2000, []byte("negative[1]")},
		{+0000, +0000, +0000, []byte("zero")},
		{-2001, +0000, +0001, []byte("negative[2]")},
		{+0001, +0002, +2000, []byte("positive[0]")},
		{+2000, +2000, +2000, []byte("positive[1]")},
	}
	setUp := func() {
		for i := range testData {
			b := testData[i]
			blocks_repo.Update(&b)
		}
	}
	tearDown := func() {
		for _, b := range testData {
			blocks_repo.Delete(b.PosX, b.PosY, b.PosZ)
		}
	}

	setUp()
	defer tearDown()

	// Helper function to loop over all channel data
	// TODO(ronoaldo): simplify this so we don't have
	// to work so verbose. Perhaps implement a wrapper
	// to the channel?
	consumeAll := func(tc string, it chan *block.Block) int {
		t.Logf("Test Case: %s", tc)
		t.Logf("consumeAll: fetching data from iterator")
		count := 0
		for {
			select {
			case b, ok := <-it:
				if !ok {
					return count
				}
				t.Logf("consumeAll: got %v", b)
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

	type testCase struct {
		x        int
		y        int
		z        int
		name     string
		expected int
	}

	// Sort order should be (z, y, x) to keep consistency with how sqlite
	// in minetest works (using a single pos field with z+y+x summed and byte-shifted)
	testCases := []testCase{
		{-2001, -2001, -2001, "starting from -2000,-2000,-2000 should return 6", 6},
		{-2000, -2000, -2000, "starting from -2000,-2000,-2000 should return 5", 5},
		{-0001, -0001, -0001, "starting from -0001,-0001,-0001 should return 4", 4},
		{+0000, +0000, +0000, "starting from +0000,+0000,+0000 should return 3", 3},
		{+0000, +0000, +0001, "starting from +0000,+0000,+0001 should return 2", 2},
		{+0000, +0000, +1999, "starting from +2000,+2000,+1999 should return 2", 2},
		{+1999, +1999, +1999, "starting from +1999,+1999,+1999 should return 2", 2},
		{+2000, +2000, +2000, "starting from +2000,+2000,+2000 should return 0", 0},
		{+2000, +2000, +2001, "starting from +2000,+2000,+2001 should return 0", 0},
		{+2000, +2001, +2000, "starting from +2000,+2001,+2000 should return 0", 0},
		{+2001, +2000, +2000, "starting from +2001,+2000,+2000 should return 0", 0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			it, _, err := blocks_repo.Iterator(tc.x, tc.y, tc.z)
			if assert.NoError(t, err) {
				assert.Equal(t, tc.expected, consumeAll(tc.name, it))
			}
		})
	}
}

func testIteratorErrorHandling(t *testing.T, blocks_repo block.BlockRepository, db *sql.DB, mockDataCorruption string) {
	logToTesting(t)
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

	ch, _, err := blocks_repo.Iterator(0, 0, 0)
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

func testIteratorClose(t *testing.T, r block.BlockRepository) {
	logToTesting(t)

	// setUp: Generates 1000+ blocks
	for x := -10; x <= 10; x += 2 {
		for y := -10; y <= 10; y += 2 {
			for z := -10; z <= 10; z += 2 {
				r.Update(&block.Block{x, y, z, []byte("default:stone")})
			}
		}
	}

	it, cl, err := r.Iterator(0, 0, 0)
	assert.NoError(t, err, "no error should be returned when initializing iterator")
	assert.NotNil(t, cl, "closer should not be nil")

	count := 0
	for b := range it {
		t.Logf("Block received: %v", b)
		assert.NotNil(t, b, "should not return a nil block from iterator")
		count++

		if count >= 10 {
			t.Logf("Closing the bridge at %d", count)
			assert.NoError(t, cl.Close(), "closer should not have any errors")
			break
		}
	}

	totalCount, err := r.Count()
	assert.NoError(t, err, "should not return error when counting")

	t.Logf("Retrieved %d blocks from a total of %d", count, totalCount)
}
