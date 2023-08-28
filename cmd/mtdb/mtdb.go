package main

import (
	"archive/zip"
	"errors"
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/minetest-go/mtdb"
	"github.com/minetest-go/mtdb/auth"
	"github.com/minetest-go/mtdb/worldconfig"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

var help = flag.Bool("help", false, "shows the help")
var show_version = flag.Bool("version", false, "shows the version")
var show_stats = flag.Bool("stats", false, "show database statistics")
var migrate = flag.Bool("migrate", false, "just migrates the database schemas and exit")
var export = flag.String("export", "", "exports the database to the given zip file")
var import_file = flag.String("import", "", "imports the database from a given zip file")
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

	if *show_stats {
		block_count, err := ctx.Blocks.Count()
		if err != nil {
			panic(err)
		}
		auth_count, err := ctx.Auth.Count(&auth.AuthSearch{})
		if err != nil {
			panic(err)
		}
		player_count, err := ctx.Player.Count()
		if err != nil {
			panic(err)
		}
		modstorage_count, err := ctx.ModStorage.Count()
		if err != nil {
			panic(err)
		}

		fmt.Printf("Blocks: %d, Auth-entries: %d, Player-entries: %d, Modstorage-entries: %d\n", block_count, auth_count, player_count, modstorage_count)
	}

	if *export != "" {
		fmt.Printf("Exporting database to '%s'\n", *export)
		f, err := os.Create(*export)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		z := zip.NewWriter(f)
		defer z.Close()

		err = ctx.Export(z)
		if err != nil {
			panic(err)
		}

		fmt.Println("Database exported")
	}

	if *import_file != "" {
		fmt.Printf("Importing database from '%s'\n", *import_file)

		z, err := zip.OpenReader(*import_file)
		if err != nil {
			panic(err)
		}
		defer z.Close()

		err = ctx.Import(&z.Reader)
		if err != nil {
			panic(err)
		}

		fmt.Println("Database imported")
	}

}
