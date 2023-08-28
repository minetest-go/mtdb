package block_test

import (
	"archive/zip"
	"bytes"
	"os"
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
