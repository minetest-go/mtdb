package player

import (
	"database/sql"
	"errors"

	"github.com/minetest-go/mtdb/types"
)

func NewPlayerMetadataRepository(db *sql.DB, dbtype types.DatabaseType) *PlayerMetadataRepository {
	return &PlayerMetadataRepository{db: db, dbtype: dbtype}
}

type PlayerMetadataRepository struct {
	db     *sql.DB
	dbtype types.DatabaseType
}

func (r *PlayerMetadataRepository) GetPlayerMetadata(name string) ([]*PlayerMetadata, error) {
	var q string
	switch r.dbtype {
	case types.DATABASE_SQLITE:
		q = "select player,metadata,value from player_metadata where player = $1"
	case types.DATABASE_POSTGRES:
		q = "select player,attr,value from player_metadata where player = $1"
	default:
		return nil, errors.New("invalid dbtype")
	}

	rows, err := r.db.Query(q, name)
	if err != nil {
		return nil, err
	}
	list := make([]*PlayerMetadata, 0)
	for rows.Next() {
		pm := &PlayerMetadata{}
		err = rows.Scan(&pm.Player, &pm.Metadata, &pm.Value)
		if err != nil {
			return nil, err
		}
		list = append(list, pm)
	}
	return list, rows.Close()
}

func (r *PlayerMetadataRepository) SetPlayerMetadata(md *PlayerMetadata) error {
	var q string
	switch r.dbtype {
	case types.DATABASE_SQLITE:
		q = "insert or replace into player_metadata(player,metadata,value) values($1,$2,$3)"
	case types.DATABASE_POSTGRES:
		q = "insert into player_metadata(player,attr,value) values($1,$2,$3) on conflict (player,attr) do update set value = EXCLUDED.value"
	default:
		return errors.New("invalid dbtype")
	}

	_, err := r.db.Exec(q, md.Player, md.Metadata, md.Value)
	return err
}
