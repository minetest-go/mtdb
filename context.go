package mtdb

import (
	"database/sql"
	"path"

	_ "github.com/lib/pq"
	"github.com/minetest-go/mtdb/auth"
	"github.com/minetest-go/mtdb/block"
	"github.com/minetest-go/mtdb/mod_storage"
	"github.com/minetest-go/mtdb/player"
	"github.com/minetest-go/mtdb/types"
	"github.com/minetest-go/mtdb/wal"
	"github.com/minetest-go/mtdb/worldconfig"
	_ "modernc.org/sqlite"
)

type Context struct {
	Auth           *auth.AuthRepository
	Privs          *auth.PrivRepository
	Player         player.PlayerRepository
	Blocks         block.BlockRepository
	ModStorage     mod_storage.ModStorageRepository
	map_db         *sql.DB
	player_db      *sql.DB
	auth_db        *sql.DB
	mod_storage_db *sql.DB
}

// closes all database connections
func (ctx *Context) Close() {
	ctx.map_db.Close()
	ctx.player_db.Close()
	ctx.auth_db.Close()
	ctx.mod_storage_db.Close()
}

// parses the "world.mt" file in the world-dir and creates a new context
func New(world_dir string) (*Context, error) {
	wc, err := worldconfig.Parse(path.Join(world_dir, "world.mt"))
	if err != nil {
		return nil, err
	}
	ctx := &Context{}

	//TODO: refactor/minimize repetitive code

	// create map repos
	switch wc[worldconfig.CONFIG_MAP_BACKEND] {
	case worldconfig.BACKEND_SQLITE3:
		map_db, err := sql.Open("sqlite", path.Join(world_dir, "map.sqlite"))
		if err != nil {
			return nil, err
		}

		err = wal.EnableWAL(map_db)
		if err != nil {
			return nil, err
		}

		err = block.MigrateBlockDB(map_db, types.DATABASE_SQLITE)
		if err != nil {
			return nil, err
		}

		ctx.Blocks = block.NewBlockRepository(map_db, types.DATABASE_SQLITE)
		ctx.map_db = map_db

	case worldconfig.BACKEND_POSTGRES:
		map_db, err := sql.Open("postgres", wc[worldconfig.CONFIG_PSQL_MAP_CONNECTION])
		if err != nil {
			return nil, err
		}

		err = block.MigrateBlockDB(map_db, types.DATABASE_POSTGRES)
		if err != nil {
			return nil, err
		}

		ctx.Blocks = block.NewBlockRepository(map_db, types.DATABASE_POSTGRES)
		ctx.map_db = map_db
	}

	// create auth repos
	switch wc[worldconfig.CONFIG_AUTH_BACKEND] {
	case worldconfig.BACKEND_SQLITE3:
		auth_db, err := sql.Open("sqlite", path.Join(world_dir, "auth.sqlite"))
		if err != nil {
			return nil, err
		}

		err = wal.EnableWAL(auth_db)
		if err != nil {
			return nil, err
		}

		err = auth.MigrateAuthDB(auth_db, types.DATABASE_SQLITE)
		if err != nil {
			return nil, err
		}

		ctx.Auth = auth.NewAuthRepository(auth_db, types.DATABASE_SQLITE)
		ctx.Privs = auth.NewPrivilegeRepository(auth_db, types.DATABASE_SQLITE)
		ctx.auth_db = auth_db

	case worldconfig.BACKEND_POSTGRES:
		auth_db, err := sql.Open("postgres", wc[worldconfig.CONFIG_PSQL_AUTH_CONNECTION])
		if err != nil {
			return nil, err
		}

		err = auth.MigrateAuthDB(auth_db, types.DATABASE_POSTGRES)
		if err != nil {
			return nil, err
		}

		ctx.Auth = auth.NewAuthRepository(auth_db, types.DATABASE_POSTGRES)
		ctx.Privs = auth.NewPrivilegeRepository(auth_db, types.DATABASE_POSTGRES)
		ctx.auth_db = auth_db

	}

	// mod storage
	switch wc[worldconfig.CONFIG_STORAGE_BACKEND] {
	case worldconfig.BACKEND_SQLITE3:
		mod_storage_db, err := sql.Open("sqlite", path.Join(world_dir, "mod_storage.sqlite"))
		if err != nil {
			return nil, err
		}

		err = mod_storage.MigrateModStorageDB(mod_storage_db, types.DATABASE_SQLITE)
		if err != nil {
			return nil, err
		}

		ctx.ModStorage = mod_storage.NewModStorageRepository(mod_storage_db, types.DATABASE_SQLITE)
		ctx.mod_storage_db = mod_storage_db
	}

	// players
	switch wc[worldconfig.CONFIG_PLAYER_BACKEND] {
	case worldconfig.BACKEND_SQLITE3:
		player_db, err := sql.Open("sqlite", path.Join(world_dir, "players.sqlite"))
		if err != nil {
			return nil, err
		}

		err = mod_storage.MigrateModStorageDB(player_db, types.DATABASE_SQLITE)
		if err != nil {
			return nil, err
		}

		ctx.Player = player.NewPlayerRepository(player_db, types.DATABASE_SQLITE)
		ctx.player_db = player_db
	}

	return ctx, nil
}
