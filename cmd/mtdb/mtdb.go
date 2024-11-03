package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/minetest-go/mtdb"
	"github.com/minetest-go/mtdb/worldconfig"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

var help = flag.Bool("help", false, "shows the help")
var show_version = flag.Bool("version", false, "shows the version")
var migrate = flag.Bool("migrate", false, "just migrates the database schemas and exit")
var init_world = flag.Bool("init", false, "initialize world.mt with defaults if it does not exist")

func main() {
	flag.Parse()

	if *help {
		flag.Usage()
		return
	}

	if *show_version {
		fmt.Printf("mtdb %s, commit %s, built at %s\n", version, commit, date)
	}

	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	if *init_world {
		worldmt_file := path.Join(wd, "world.mt")
		_, err := os.Stat(worldmt_file)
		if errors.Is(err, os.ErrNotExist) {
			err := os.WriteFile(worldmt_file, []byte(worldconfig.DEFAULT_CONFIG), 0644)
			if err != nil {
				panic(err)
			}
		}
	}

	ctx, err := mtdb.New(wd)
	if err != nil {
		panic(err)
	}
	defer ctx.Close()

	if *migrate {
		// already migrated at this point
		fmt.Println("Databases migrated / initialized")
		return
	}

}
