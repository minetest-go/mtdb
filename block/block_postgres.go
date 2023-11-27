package block

import (
	"archive/zip"
	"bufio"
	"bytes"
	"database/sql"

	"encoding/json"

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
				return
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

func (repo *postgresBlockRepository) Count() (int64, error) {
	row := repo.db.QueryRow("select count(*) from blocks")
	count := int64(0)
	err := row.Scan(&count)
	return count, err
}

func (r *postgresBlockRepository) Export(z *zip.Writer) error {
	w, err := z.Create("blocks.json")
	if err != nil {
		return err
	}
	enc := json.NewEncoder(w)

	rows, err := r.db.Query("select posX,posY,posZ,data from blocks")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		block := &Block{}
		err = rows.Scan(&block.PosX, &block.PosY, &block.PosZ, &block.Data)
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

func (r *postgresBlockRepository) Import(z *zip.Reader) error {
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
