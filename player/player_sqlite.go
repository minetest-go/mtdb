package player

import (
	"database/sql"
	"errors"
)

type PlayerSqliteRepository struct {
	db *sql.DB
}

func (r *PlayerSqliteRepository) GetPlayer(name string) (*Player, error) {
	q := `
		select name,pitch,yaw,
			posx,posy,posz,
			hp,breath,
			strftime('%s', creation_date),strftime('%s', modification_date)
		from player
		where name = $1
	`
	row := r.db.QueryRow(q, name)
	p := &Player{}
	err := row.Scan(&p.Name, &p.Pitch, &p.Yaw, &p.PosX, &p.PosY, &p.PosZ, &p.HP, &p.Breath, &p.CreationDate, &p.ModificationDate)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return p, nil
}

func (r *PlayerSqliteRepository) CreateOrUpdate(p *Player) error {
	q := `
		insert or replace into
		player(
			name,pitch,yaw,
			posx,posy,posz,
			hp,breath,
			creation_date,modification_date
		)
		values(
			$1,$2,$3,
			$4,$5,$6,
			$7,$8,
			datetime($9, 'unixepoch'),datetime($10, 'unixepoch')
		)
	`
	_, err := r.db.Exec(q, p.Name, p.Pitch, p.Yaw, p.PosX, p.PosY, p.PosZ, p.HP, p.Breath, p.CreationDate, p.ModificationDate)
	return err
}
