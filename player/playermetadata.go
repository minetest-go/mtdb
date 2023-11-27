package player

import (
	"archive/zip"
	"bufio"
	"bytes"
	"database/sql"
	"encoding/json"
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

func (r *PlayerMetadataRepository) Export(z *zip.Writer) error {
	w, err := z.Create("playermetadata.json")
	if err != nil {
		return err
	}
	enc := json.NewEncoder(w)

	rows, err := r.db.Query("select player from player_metadata")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		name := ""
		err = rows.Scan(&name)
		if err != nil {
			return err
		}

		list, err := r.GetPlayerMetadata(name)
		if err != nil {
			return err
		}

		for _, entry := range list {
			err = enc.Encode(entry)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *PlayerMetadataRepository) Import(z *zip.Reader) error {
	f, err := z.Open("playermetadata.json")
	if err != nil {
		return err
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		dc := json.NewDecoder(bytes.NewReader(sc.Bytes()))
		e := &PlayerMetadata{}
		err = dc.Decode(e)
		if err != nil {
			return err
		}

		err = r.SetPlayerMetadata(e)
		if err != nil {
			return err
		}
	}

	return nil
}
