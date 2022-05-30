package mtdb_test

import (
	"testing"

	"github.com/minetest-go/mtdb"
	"github.com/stretchr/testify/assert"
)

func testBlocksRepository(t *testing.T, block_repo mtdb.BlockRepository) {

	// cleanup
	assert.NoError(t, block_repo.Delete(0, 0, 0))

	// get nil
	block, err := block_repo.GetByPos(0, 0, 0)
	assert.NoError(t, err)
	assert.Nil(t, block)

	// create
	block = &mtdb.Block{
		PosX: 0,
		PosY: 0,
		PosZ: 0,
		Data: []byte{0x00, 0x01, 0x02},
	}
	assert.NoError(t, block_repo.Update(block))

	// get
	block, err = block_repo.GetByPos(0, 0, 0)
	assert.NoError(t, err)
	assert.NotNil(t, block)
	assert.Equal(t, 0, block.PosX)
	assert.Equal(t, 0, block.PosY)
	assert.Equal(t, 0, block.PosZ)
	assert.Equal(t, 3, len(block.Data))
	assert.Equal(t, uint8(0x00), block.Data[0])
	assert.Equal(t, uint8(0x01), block.Data[1])
	assert.Equal(t, uint8(0x02), block.Data[2])

	// delete
	assert.NoError(t, block_repo.Delete(0, 0, 0))

	// get nil
	block, err = block_repo.GetByPos(0, 0, 0)
	assert.NoError(t, err)
	assert.Nil(t, block)

}
