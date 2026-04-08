package sqliteutils

import (
	"database/sql"
	"fmt"
)

func EnableAutoVacuum(db *sql.DB) error {
	_, err := db.Exec("pragma auto_vacuum = INCREMENTAL;")
	if err != nil {
		return fmt.Errorf("couldn't switch the db to incremental vacuum: %v", err)
	}

	return nil
}
