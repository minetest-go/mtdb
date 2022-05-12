package mtdb

import "database/sql"

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

func NewBlockRepository(db *sql.DB, dbtype DatabaseType) BlockRepository {
	switch dbtype {
	case DATABASE_POSTGRES:
		return &postgresBlockRepository{db: db}
	case DATABASE_SQLITE:
		return &sqliteBlockRepository{db: db}
	default:
		return nil
	}
}
