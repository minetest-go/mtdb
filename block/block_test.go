package block_test

import (
	"testing"

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

}
