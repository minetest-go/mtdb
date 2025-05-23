package block

import (
	"database/sql"

	"github.com/minetest-go/mtdb/types"
	"github.com/sirupsen/logrus"
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
	defer rows.Close()
	if !rows.Next() {
		return nil, nil
	}
	entry := &Block{}
	err = rows.Scan(&pos, &entry.Data)
	entry.PosX, entry.PosY, entry.PosZ = PlainToCoord(pos)
	return entry, err
}

func (repo *sqliteBlockRepository) Iterator(x, y, z int) (chan *Block, types.Closer, error) {
	pos := CoordToPlain(x, y, z)
	rows, err := repo.db.Query(`
		SELECT pos, data
		FROM blocks
		WHERE pos > $1
		ORDER BY pos
		`, pos)
	if err != nil {
		return nil, nil, err
	}

	l := logrus.
		WithField("iterating_from", []int{x, y, z}).
		WithField("pos", pos)
	ch := make(chan *Block)
	done := make(types.WhenDone, 1)
	count := int64(0)

	// Spawn go routine to fetch rows and send to channel
	go func() {
		defer close(ch)
		defer rows.Close()

		l.Debug("Retrieving database rows")
		for {
			select {
			case <-done:
				l.Debugf("Iterator closed by caller. Finishing up...")
				return
			default:
				if rows.Next() {
					// Debug progress while fetching rows every 100's
					count++
					if count%100 == 0 {
						l.Debugf("Retrieved %d records so far", count)
					}
					// Fetch and send to channel
					b := &Block{}
					if err = rows.Scan(&pos, &b.Data); err != nil {
						l.Errorf("Failed to read next item from iterator: %v", err)
						return
					}
					b.PosX, b.PosY, b.PosZ = PlainToCoord(pos)
					ch <- b
				} else {
					l.Debug("Iterator finished, closing up rows and channel")
					return
				}
			}
		}
	}()

	// Return channel to main component
	return ch, done, nil
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

func (r *sqliteBlockRepository) Close() error {
	return r.db.Close()
}
