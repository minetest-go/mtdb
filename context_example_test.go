package mtdb_test

import (
	"fmt"

	"github.com/minetest-go/mtdb"
)

func ExampleContext() {
	// create a new context (an object with all repositories) from a minetest world-directory
	ctx, err := mtdb.New("/xy")
	if err != nil {
		panic(err)
	}

	// retrieve an auth entry of the "admin" user
	entry, err := ctx.Auth.GetByUsername("admin")
	if err != nil {
		panic(err)
	}

	fmt.Printf("Admin password-entry: %s\n", entry.Password)

	// retrieve admin's privileges
	privs, err := ctx.Privs.GetByID(*entry.ID)
	for _, priv := range privs {
		fmt.Printf(" + %s\n", priv.Privilege)
	}

	// get a mapblock from the database
	block, err := ctx.Blocks.GetByPos(0, 0, 0)
	if err != nil {
		panic(err)
	}
	// dump the raw and unparsed binary data
	// use the github.com/minetest-go/mapparser project to parse the actual content
	fmt.Printf("Mapblock content: %s\n", block.Data)
}

func ExampleNewBlockDB() {
	// create a new block-repository with the backend-connection specified in the world.mt file in the given directory
	blocks, err := mtdb.NewBlockDB("/my-world-dir")
	if err != nil {
		panic(err)
	}

	// fetch the mapblock at position 0,0,0
	block, err := blocks.GetByPos(0, 0, 0)
	if err != nil {
		panic(err)
	}

	// dump the raw and unparsed binary data
	// use the github.com/minetest-go/mapparser project to parse the actual content
	fmt.Printf("Mapblock content: %s\n", block.Data)
}
