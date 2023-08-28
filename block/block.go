package block

import (
	"database/sql"

	"github.com/minetest-go/mtdb/types"
)

type Block struct {
	PosX int    `json:"x"`
	PosY int    `json:"y"`
	PosZ int    `json:"z"`
	Data []byte `json:"data"`
}

type BlockRepository interface {
	types.Backup
	GetByPos(x, y, z int) (*Block, error)
	Update(block *Block) error
	Delete(x, y, z int) error
	DeleteAll() error
	Vacuum() error
	Count() (int64, error)
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
