package player

import (
	"database/sql"
	"errors"
)

type PlayerPostgresRepository struct {
	db *sql.DB
}

func (r *PlayerPostgresRepository) GetPlayer(name string) (*Player, error) {
	q := `
		select name,pitch,yaw,
			posx,posy,posz,
			hp,breath,
			extract(epoch from creation_date)::int,extract(epoch from modification_date)::int
		from player
		where name = $1
	`
	row := r.db.QueryRow(q, name)
	if row.Err() != nil {
		return nil, row.Err()
	}
	p := &Player{}
	err := row.Scan(&p.Name, &p.Pitch, &p.Yaw, &p.PosX, &p.PosY, &p.PosZ, &p.HP, &p.Breath, &p.CreationDate, &p.ModificationDate)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return p, nil
}

func (r *PlayerPostgresRepository) CreateOrUpdate(p *Player) error {
	q := `
		insert into player(
			name,pitch,yaw,
			posx,posy,posz,
			hp,breath,
			creation_date,modification_date
		)
		values(
			$1,$2,$3,
			$4,$5,$6,
			$7,$8,
			to_timestamp($9),to_timestamp($10)
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
	_, err := r.db.Exec(q, p.Name, p.Pitch, p.Yaw, p.PosX, p.PosY, p.PosZ, p.HP, p.Breath, p.CreationDate, p.ModificationDate)
	return err
}
