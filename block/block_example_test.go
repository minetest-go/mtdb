package block_test

import (
	"database/sql"
	"fmt"

	"github.com/minetest-go/mtdb/block"
	"github.com/minetest-go/mtdb/types"
)

func ExampleNewBlockRepository() {
	// create a new db with given arguments
	db, err := sql.Open("sqlite3", "./map.sqlite?_journal=WAL")
	if err != nil {
		panic(err)
	}

	// optionally: create or migrate db-schema
	block.MigrateBlockDB(db, types.DATABASE_SQLITE)

	// create a new repository instance
	repo := block.NewBlockRepository(db, types.DATABASE_SQLITE)

	// fetch the mapblock at position 0,0,0
	block, err := repo.GetByPos(0, 0, 0)
	if err != nil {
		panic(err)
	}

	// dump the raw and unparsed binary data
	// use the github.com/minetest-go/mapparser project to parse the actual content
	fmt.Printf("Mapblock content: %s\n", block.Data)
}
