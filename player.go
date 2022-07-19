package mtdb

import (
	"database/sql"
	"errors"
)

type Player struct {
	Name             string  `json:"name"`
	Pitch            float64 `json:"pitch"`
	Yaw              float64 `json:"yaw"`
	PosX             float64 `json:"posx"`
	PosY             float64 `json:"posy"`
	PosZ             float64 `json:"posz"`
	HP               int     `json:"hp"`
	Breath           int     `json:"breath"`
	CreationDate     int64   `json:"creation_date"`
	ModificationDate int64   `json:"modification_date"`
}

type PlayerMetadata struct {
	Player   string `json:"player"`
	Metadata string `json:"metadata"`
	Value    string `json:"value"`
}

type PlayerInventories struct {
	Player   string `json:"player"`
	InvID    int    `json:"inv_id"`
	InvWidth int    `json:"inv_width"`
	InvName  string `json:"inv_name"`
	InvSize  int    `json:"inv_size"`
}

type PlayerInventoryItems struct {
	Player string `json:"player"`
	InvID  int    `json:"inv_id"`
	SlotID int    `json:"slot_id"`
	Item   string `json:"item"`
}

type PlayerRepository interface {
	GetPlayer(name string) (*Player, error)
}

func NewPlayerRepository(db *sql.DB, dbtype DatabaseType) PlayerRepository {
	switch dbtype {
	case DATABASE_SQLITE:
		return &PlayerSqliteRepository{db: db}
	}
	return nil
}

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
