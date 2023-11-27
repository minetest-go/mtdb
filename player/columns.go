package player

import (
	"github.com/minetest-go/mtdb/types"
)

func getColumns(dbtype types.DatabaseType) []string {
	switch dbtype {
	case types.DATABASE_SQLITE:
		return []string{
			"name", "pitch", "yaw",
			"posx", "posy", "posz",
			"hp", "breath",
			"strftime('%s', creation_date)",
			"strftime('%s', modification_date)",
		}
	case types.DATABASE_POSTGRES:
		return []string{
			"name", "pitch", "yaw",
			"posx", "posy", "posz",
			"hp", "breath",
			"extract(epoch from creation_date)::int",
			"extract(epoch from modification_date)::int",
		}
	default:
		return nil
	}
}

func scanPlayer(scan func(v ...any) error) (*Player, error) {
	p := &Player{}
	err := scan(&p.Name, &p.Pitch, &p.Yaw, &p.PosX, &p.PosY, &p.PosZ, &p.HP, &p.Breath, &p.CreationDate, &p.ModificationDate)
	return p, err
}
