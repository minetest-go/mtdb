package block

import (
	"database/sql"

	"github.com/sirupsen/logrus"
)

type postgresBlockRepository struct {
	db *sql.DB
}

func (repo *postgresBlockRepository) GetByPos(x, y, z int) (*Block, error) {
	rows, err := repo.db.Query("select posX,posY,posZ,data from blocks where posX=$1 and posY=$2 and posZ=$3", x, y, z)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if !rows.Next() {
		return nil, nil
	}
	entry := &Block{}
	err = rows.Scan(&entry.PosX, &entry.PosY, &entry.PosZ, &entry.Data)
	return entry, err
}

func (repo *postgresBlockRepository) Iterator(x, y, z int) (chan *Block, error) {
	rows, err := repo.db.Query(`
		SELECT posX, posY, posZ, data
		FROM blocks
		WHERE posX >= $1 AND posY >= $2 AND posZ >= $3
		ORDER BY posX, posY, posZ
		`, x, y, z)
	if err != nil {
		return nil, err
	}

	l := logrus.WithField("iterating_from", []int{x, y, z})
	ch := make(chan *Block)
	count := int64(0)

	// Spawn go routine to fetch rows and send to channel
	go func() {
		defer close(ch)
		defer rows.Close()

		l.Debug("Retrieving database rows")
		for rows.Next() {
			// Debug progress while fetching rows every 100's
			count++
			if count%100 == 0 {
				l.Debugf("Retrieved %d records so far", count)
			}
			// Fetch and send to channel
			b := &Block{}
			if err = rows.Scan(&b.PosX, &b.PosY, &b.PosZ, &b.Data); err != nil {
				l.Errorf("Failed to read next item from iterator: %v", err)
			}
			ch <- b
		}
		l.Debug("Iterator finished, closing up rows and channel")
	}()

	// Return channel to main component
	return ch, nil
}

func (repo *postgresBlockRepository) Update(block *Block) error {
	_, err := repo.db.Exec("insert into blocks(posX,posY,posZ,data) values($1,$2,$3,$4) ON CONFLICT ON CONSTRAINT blocks_pkey DO UPDATE SET data = $4",
		block.PosX, block.PosY, block.PosZ, block.Data)
	return err
}

func (repo *postgresBlockRepository) Delete(x, y, z int) error {
	_, err := repo.db.Exec("delete from blocks where posX=$1 and posY=$2 and posZ=$3", x, y, z)
	return err
}

func (repo *postgresBlockRepository) Vacuum() error {
	_, err := repo.db.Exec("vacuum")
	return err
}
