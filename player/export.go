package player

import (
	"archive/zip"
	"bufio"
	"bytes"
	"encoding/json"
)

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
