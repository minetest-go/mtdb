package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/minetest-go/mtdb"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

var help = flag.Bool("help", false, "shows the help")
var show_version = flag.Bool("version", false, "shows the version")
var migrate = flag.Bool("migrate", false, "just migrates the database schemas and exit")

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
