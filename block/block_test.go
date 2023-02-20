package block_test

import (
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
	negX, negY, negZ := block.NodeToBlock(-32000, -32000, -32000)
	posX, posY, posZ := block.NodeToBlock(32000, 32000, 32000)

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
