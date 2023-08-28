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

func NewPlayerRepository(db *sql.DB, dbtype types.DatabaseType) *PlayerRepository {
	return &PlayerRepository{db: db, dbtype: dbtype}
}

type PlayerRepository struct {
	db     *sql.DB
	dbtype types.DatabaseType
}

func (r *PlayerRepository) GetPlayer(name string) (*Player, error) {
	var q string
	switch r.dbtype {
	case types.DATABASE_SQLITE:
		q = `
			select name,pitch,yaw,
				posx,posy,posz,
				hp,breath,
				strftime('%s', creation_date),strftime('%s', modification_date)
			from player
			where name = $1
		`
	case types.DATABASE_POSTGRES:
		q = `
			select name,pitch,yaw,
				posx,posy,posz,
				hp,breath,
				extract(epoch from creation_date)::int,extract(epoch from modification_date)::int
			from player
			where name = $1
		`
	default:
		return nil, errors.New("invalid dbtype")
	}

	row := r.db.QueryRow(q, name)
	p := &Player{}
	err := row.Scan(&p.Name, &p.Pitch, &p.Yaw, &p.PosX, &p.PosY, &p.PosZ, &p.HP, &p.Breath, &p.CreationDate, &p.ModificationDate)
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
				creation_date,modification_date
			)
			values(
				$1,$2,$3,
				$4,$5,$6,
				$7,$8,
				datetime($9, 'unixepoch'),datetime($10, 'unixepoch')
			)
		`
	case types.DATABASE_POSTGRES:
		q = `
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

func (r *PlayerRepository) Count() (int64, error) {
	row := r.db.QueryRow("select count(*) from player")
	count := int64(0)
	err := row.Scan(&count)
	return count, err
}

func (r *PlayerRepository) Export(z *zip.Writer) error {
	w, err := z.Create("player.json")
	if err != nil {
		return err
	}
	enc := json.NewEncoder(w)

	rows, err := r.db.Query("select name from player")
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

		player, err := r.GetPlayer(name)
		if err != nil {
			return err
		}

		err = enc.Encode(player)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *PlayerRepository) Import(z *zip.Reader) error {
	f, err := z.Open("player.json")
	if err != nil {
		return err
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		dc := json.NewDecoder(bytes.NewReader(sc.Bytes()))
		e := &Player{}
		err = dc.Decode(e)
		if err != nil {
			return err
		}

		err = r.CreateOrUpdate(e)
		if err != nil {
			return err
		}
	}

	return nil
}
