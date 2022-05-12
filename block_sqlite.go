package mtdb

import "database/sql"

type SqliteBlockRepository struct {
	db *sql.DB
}

//https://bitbucket.org/s_l_teichmann/mtsatellite/src/e1bf980a2b278c570b3f44f9452c9c087558acb3/common/coords.go?at=default&fileviewer=file-view-default
const (
	numBitsPerComponent = 12
	modulo              = 1 << numBitsPerComponent
	maxPositive         = modulo / 2
	minValue            = -1 << (numBitsPerComponent - 1)
	maxValue            = 1<<(numBitsPerComponent-1) - 1

	MinPlainCoord = -34351347711
)

func CoordToPlain(x, y, z int) int64 {
	return int64(z)<<(2*numBitsPerComponent) +
		int64(y)<<numBitsPerComponent +
		int64(x)
}

func unsignedToSigned(i int16) int {
	if i < maxPositive {
		return int(i)
	}
	return int(i - maxPositive*2)
}

// To match C++ code.
func pythonModulo(i int16) int16 {
	const mask = modulo - 1
	if i >= 0 {
		return i & mask
	}
	return modulo - -i&mask
}

func PlainToCoord(i int64) (x, y, z int) {
	x = unsignedToSigned(pythonModulo(int16(i)))
	i = (i - int64(x)) >> numBitsPerComponent
	y = unsignedToSigned(pythonModulo(int16(i)))
	i = (i - int64(y)) >> numBitsPerComponent
	z = unsignedToSigned(pythonModulo(int16(i)))
	return x, y, z
}

func (repo *SqliteBlockRepository) GetByPos(x, y, z int) (*Block, error) {
	pos := CoordToPlain(x, y, z)
	rows, err := repo.db.Query("select pos,data from blocks where pos=$1", pos)
	if err != nil {
		return nil, err
	}
	if !rows.Next() {
		return nil, nil
	}
	entry := &Block{}
	err = rows.Scan(&pos, &entry.Data)
	entry.PosX, entry.PosY, entry.PosZ = PlainToCoord(pos)
	return entry, err
}

func (repo *SqliteBlockRepository) Update(block *Block) error {
	pos := CoordToPlain(block.PosX, block.PosY, block.PosZ)
	_, err := repo.db.Exec("replace into blocks(pos,data) values($1, $2)", pos, block.Data)
	return err
}

func (repo *SqliteBlockRepository) Delete(x, y, z int) error {
	pos := CoordToPlain(x, y, z)
	_, err := repo.db.Exec("delete from blocks where pos=$1", pos)
	return err
}
