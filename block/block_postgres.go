package block

import (
	"archive/zip"
	"bufio"
	"bytes"
	"database/sql"

	"encoding/json"

	"github.com/minetest-go/mtdb/types"
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

// IteratorBatchSize is a default value to be used while batching over Iterators.
var IteratorBatchSize = 4096

// iteratorQuery uses the keyset pagination as described in this SO question
// https://stackoverflow.com/q/57504274
//
// The WHERE (args) >= (params) works in a way to keep the sorting during comparision
// without the need to implement this as an offset/cursor
var iteratorQuery = `-- Query blocks to iterate over
SELECT posX, posY, posZ, data
FROM blocks
WHERE (posZ, posY, posX) > ($3, $2, $1)
ORDER BY posZ, posY, posX
LIMIT $4`

func (repo *postgresBlockRepository) Iterator(x, y, z int) (chan *Block, types.Closer, error) {
	// logging setup
	l := logrus.WithField("iterating_from", []int{x, y, z})
	count := int64(0)
	page := 0
	pageSize := 0

	l.Debug("running database query")
	rows, err := repo.db.Query(iteratorQuery, x, y, z, IteratorBatchSize)
	if err != nil {
		return nil, nil, err
	}

	ch := make(chan *Block, IteratorBatchSize)
	done := make(types.WhenDone, 1)
	lastPos := Block{}

	// Spawn go routine to fetch rows and send to channel
	go func() {
		defer close(ch)
		defer rows.Close()

		l.Debug("retrieving database rows ...")
		for {
			select {
			case <-done:
				// We can now return, we are done
				l.Debugf("iterator closed by caller; finishing up...")
				return
			default:
				if rows.Next() {
					// Debug progress while fetching rows every 100's
					count++
					pageSize++
					if count%100 == 0 {
						l.Debugf("retrieved %d records so far", count)
					}
					// Fetch and send to channel
					b := &Block{}
					if err = rows.Scan(&b.PosX, &b.PosY, &b.PosZ, &b.Data); err != nil {
						l.Errorf("failed to read next item from iterator: %v; aborting", err)
						return
					}
					lastPos.PosX, lastPos.PosY, lastPos.PosZ = b.PosX, b.PosY, b.PosZ
					ch <- b
				} else {
					f := logrus.Fields{"last_pos": lastPos, "last_page_size": pageSize, "page": page}
					if pageSize > 0 {
						page++
						// If the previous batch is > 0, restart the query from last position
						if err := rows.Close(); err != nil {
							l.WithField("err", err).Warning("error closing previous batch")
						}
						rows, err = repo.db.Query(iteratorQuery, lastPos.PosX, lastPos.PosY, lastPos.PosZ, IteratorBatchSize)
						if err != nil {
							l.WithField("err", err).Warning("error restarting query")
							return
						}
						pageSize = 0
						l.WithFields(f).Debug("batch finished, restarting next batch")
					} else {
						l.WithFields(f).Debug("iterator finished, closing up rows and channel")
						return
					}
				}
			}
		}
	}()

	// Return channel to main component
	return ch, done, nil
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

func (r *postgresBlockRepository) Close() error {
	return r.db.Close()
}
