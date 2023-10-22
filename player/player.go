package player

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/minetest-go/mtdb/types"
)

func NewPlayerRepository(db *sql.DB, dbtype types.DatabaseType) *PlayerRepository {
	return &PlayerRepository{db: db, dbtype: dbtype}
}

type PlayerRepository struct {
	db     *sql.DB
	dbtype types.DatabaseType
}

func (r *PlayerRepository) GetPlayer(name string) (*Player, error) {
	q := fmt.Sprintf("select %s from player where name = $1", strings.Join(getColumns(r.dbtype), ","))

	row := r.db.QueryRow(q, name)
	p, err := scanPlayer(row.Scan)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return p, nil
}

func (r *PlayerRepository) CreateOrUpdate(p *Player) error {
	var q string
	switch r.dbtype {
	case types.DATABASE_SQLITE:
		q = `
			insert or replace into
			player(
				name,pitch,yaw,
				posx,posy,posz,
				hp,breath,
				creation_date,
				modification_date
			)
			values(
				$1,$2,$3,
				$4,$5,$6,
				$7,$8,
				datetime($9, 'unixepoch'),
				datetime($10, 'unixepoch')
			)
		`
	case types.DATABASE_POSTGRES:
		q = `
			insert into player(
				name,pitch,yaw,
				posx,posy,posz,
				hp,breath,
				creation_date,
				modification_date
			)
			values(
				$1,$2,$3,
				$4,$5,$6,
				$7,$8,
				to_timestamp($9),
				to_timestamp($10)
			)
			on conflict (name) do update
			set
				pitch = EXCLUDED.pitch,
				yaw = EXCLUDED.yaw,
				posx = EXCLUDED.posx,
				posy = EXCLUDED.posy,
				posz = EXCLUDED.posz,
				hp = EXCLUDED.hp,
				breath = EXCLUDED.breath,
				creation_date = EXCLUDED.creation_date,
				modification_date = EXCLUDED.modification_date
		`
	default:
		return errors.New("invalid dbtype")
	}

	_, err := r.db.Exec(q, p.Name, p.Pitch, p.Yaw, p.PosX, p.PosY, p.PosZ, p.HP, p.Breath, p.CreationDate, p.ModificationDate)
	return err
}

func (r *PlayerRepository) RemovePlayer(name string) error {
	_, err := r.db.Exec("delete from player where name = $1", name)
	return err
}
