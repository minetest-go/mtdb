package block

import (
	"database/sql"
	"fmt"
	"math"

	"github.com/minetest-go/mtdb/types"
)

type Block struct {
	PosX int    `json:"x"`
	PosY int    `json:"y"`
	PosZ int    `json:"z"`
	Data []byte `json:"data"`
}

func (b *Block) String() string {
	if b == nil {
		return "nil"
	}
	v := b.Data
	if len(b.Data) > 20 {
		v = b.Data[:20]
	}
	return fmt.Sprintf("Block{X: %d, Y: %d, Z: %d, data: \"%v\"}", b.PosX, b.PosY, b.PosZ, string(v))
}

// BlockRepository implementes data access layer for the Minetest map data.
// All positions are in mapblock coordinates, as described here:
// https://github.com/minetest/minetest/blob/master/doc/lua_api.md#mapblock-coordinates
type BlockRepository interface {
	// GetByPost returns the map block at positions X,Y,Z.
	GetByPos(x, y, z int) (*Block, error)

	// Iterator returns a channel to fetch all data from the starting position
	// X,Y,Z (exclusive), with the map blocks sorted by position ascending.
	// Sorting is done by Z, Y, X to keep consistency with Sqlite map format.
	Iterator(x, y, z int) (chan *Block, types.Closer, error)

	// Update upserts the provided map block in the database, using the position
	// as key.
	Update(block *Block) error

	// Delete removes the map block from the database denoted by the x,y,z
	// coordinates.
	Delete(x, y, z int) error

	// Vacuum executes the storage layer vacuum command. Useful to reclaim
	// storage space if not done automatically by the backend.
	Vacuum() error

	// Count returns the total number of stored blocks in the map database.
	Count() (int64, error)

	// Close gracefully finishes the connection with the database backend.
	Close() error
}

// NewBlockRepository initializes the connection with the appropriate database
// backend and returns the BlockRepository implementation suited for it.
func NewBlockRepository(db *sql.DB, dbtype types.DatabaseType) BlockRepository {
	switch dbtype {
	case types.DATABASE_POSTGRES:
		return &postgresBlockRepository{db: db}
	case types.DATABASE_SQLITE:
		return &sqliteBlockRepository{db: db}
	default:
		return nil
	}
}

// AsBlockPos converts the coordinates from the given Node into the equivalent
// Block position. Each block contains 16x16x16 nodes.
func AsBlockPos(x, y, z int) (int, int, int) {
	pos := func(x int) int { return int(math.Floor(float64(x) / 16.0)) }
	return pos(x), pos(y), pos(z)
}
