package player

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
	CreateOrUpdate(p *Player) error
}
