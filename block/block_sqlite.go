package block

import (
	"archive/zip"
	"bufio"
	"bytes"
	"database/sql"
	"encoding/json"
)

type sqliteBlockRepository struct {
	db *sql.DB
}

// https://bitbucket.org/s_l_teichmann/mtsatellite/src/e1bf980a2b278c570b3f44f9452c9c087558acb3/common/coords.go?at=default&fileviewer=file-view-default
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

func (repo *sqliteBlockRepository) GetByPos(x, y, z int) (*Block, error) {
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

func (repo *sqliteBlockRepository) Update(block *Block) error {
	pos := CoordToPlain(block.PosX, block.PosY, block.PosZ)
	_, err := repo.db.Exec("replace into blocks(pos,data) values($1, $2)", pos, block.Data)
	return err
}

func (repo *sqliteBlockRepository) Delete(x, y, z int) error {
	pos := CoordToPlain(x, y, z)
	_, err := repo.db.Exec("delete from blocks where pos=$1", pos)
	return err
}

func (repo *sqliteBlockRepository) DeleteAll() error {
	_, err := repo.db.Exec("delete from blocks")
	return err
}

func (repo *sqliteBlockRepository) Vacuum() error {
	_, err := repo.db.Exec("vacuum")
	return err
}

func (repo *sqliteBlockRepository) Count() (int64, error) {
	row := repo.db.QueryRow("select count(*) from blocks")
	count := int64(0)
	err := row.Scan(&count)
	return count, err
}

func (r *sqliteBlockRepository) Export(z *zip.Writer) error {
	w, err := z.Create("blocks.json")
	if err != nil {
		return err
	}
	enc := json.NewEncoder(w)

	rows, err := r.db.Query("select pos from blocks")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var pos int64
		err = rows.Scan(&pos)
		if err != nil {
			return err
		}

		x, y, z := PlainToCoord(pos)
		block, err := r.GetByPos(x, y, z)
		if err != nil {
			return err
		}

		err = enc.Encode(block)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *sqliteBlockRepository) Import(z *zip.Reader) error {
	f, err := z.Open("blocks.json")
	if err != nil {
		return err
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		dc := json.NewDecoder(bytes.NewReader(sc.Bytes()))
		e := &Block{}
		err = dc.Decode(e)
		if err != nil {
			return err
		}

		err = r.Update(e)
		if err != nil {
			return err
		}
	}

	return nil
}
