package mtdb

import (
	"database/sql"

	"github.com/minetest-go/mtdb/types"
)

type Block struct {
	PosX int
	PosY int
	PosZ int
	Data []byte
}

type BlockRepository interface {
	GetByPos(x, y, z int) (*Block, error)
	Update(block *Block) error
	Delete(x, y, z int) error
}

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
