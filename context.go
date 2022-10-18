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
	Player         *player.PlayerRepository
	PlayerMetadata *player.PlayerMetadataRepository
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

func connectAndMigrate(t types.DatabaseType, sqliteConn, psqlConn string, migFn func(*sql.DB, types.DatabaseType) error) (*sql.DB, error) {
	var datasource string
	var dbtype string

	switch t {
	case types.DATABASE_POSTGRES:
		datasource = psqlConn
		dbtype = "postgres"
	default:
		// default to sqlite
		datasource = sqliteConn
		dbtype = "sqlite"
	}

	if t == types.DATABASE_POSTGRES && datasource == "" {
		// pg connection unconfigured
		return nil, nil
	}

	db, err := sql.Open(string(dbtype), datasource)
	if err != nil {
		return nil, err
	}

	if t == types.DATABASE_SQLITE {
		// enable wal on sqlite databases
		err = wal.EnableWAL(db)
		if err != nil {
			return nil, err
		}
	}

	err = migFn(db, t)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// parses the "world.mt" file in the world-dir and creates a new context
func New(world_dir string) (*Context, error) {
	wc, err := worldconfig.Parse(path.Join(world_dir, "world.mt"))
	if err != nil {
		return nil, err
	}
	ctx := &Context{}

	// map
	dbtype := types.DatabaseType(wc[worldconfig.CONFIG_MAP_BACKEND])
	ctx.map_db, err = connectAndMigrate(
		dbtype,
		path.Join(world_dir, "map.sqlite"),
		wc[worldconfig.CONFIG_PSQL_MAP_CONNECTION],
		block.MigrateBlockDB,
	)
	if err != nil {
		return nil, err
	}
	if ctx.map_db != nil {
		ctx.Blocks = block.NewBlockRepository(ctx.map_db, types.DATABASE_SQLITE)
	}

	// auth/privs
	dbtype = types.DatabaseType(wc[worldconfig.CONFIG_AUTH_BACKEND])
	ctx.auth_db, err = connectAndMigrate(
		dbtype,
		path.Join(world_dir, "auth.sqlite"),
		wc[worldconfig.CONFIG_PSQL_AUTH_CONNECTION],
		auth.MigrateAuthDB,
	)
	if err != nil {
		return nil, err
	}
	if ctx.auth_db != nil {
		ctx.Auth = auth.NewAuthRepository(ctx.auth_db, dbtype)
		ctx.Privs = auth.NewPrivilegeRepository(ctx.auth_db, dbtype)
	}

	// mod storage
	dbtype = types.DatabaseType(wc[worldconfig.CONFIG_STORAGE_BACKEND])
	ctx.mod_storage_db, err = connectAndMigrate(
		dbtype,
		path.Join(world_dir, "mod_storage.sqlite"),
		"not implemented",
		mod_storage.MigrateModStorageDB,
	)
	if err != nil {
		return nil, err
	}
	if ctx.mod_storage_db != nil {
		ctx.ModStorage = mod_storage.NewModStorageRepository(ctx.mod_storage_db, dbtype)
	}

	// players
	dbtype = types.DatabaseType(wc[worldconfig.CONFIG_PLAYER_BACKEND])
	ctx.player_db, err = connectAndMigrate(
		dbtype,
		path.Join(world_dir, "players.sqlite"),
		wc[worldconfig.CONFIG_PSQL_PLAYER_CONNECTION],
		player.MigratePlayerDB,
	)
	if err != nil {
		return nil, err
	}
	if ctx.player_db != nil {
		ctx.Player = player.NewPlayerRepository(ctx.player_db, dbtype)
		ctx.PlayerMetadata = player.NewPlayerMetadataRepository(ctx.player_db, dbtype)
	}

	return ctx, nil
}
