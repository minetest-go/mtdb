package mtdb

import (
	"database/sql"
	"fmt"
	"path"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"

	"github.com/minetest-go/mtdb/auth"
	"github.com/minetest-go/mtdb/block"
	"github.com/minetest-go/mtdb/mod_storage"
	"github.com/minetest-go/mtdb/player"
	"github.com/minetest-go/mtdb/types"
	"github.com/minetest-go/mtdb/wal"
	"github.com/minetest-go/mtdb/worldconfig"
	"github.com/sirupsen/logrus"
)

// Database connection context
type Context struct {
	Auth           *auth.AuthRepository
	Privs          *auth.PrivRepository
	Player         *player.PlayerRepository
	PlayerMetadata *player.PlayerMetadataRepository
	Blocks         block.BlockRepository
	ModStorage     mod_storage.ModStorageRepository
	open_databases []*sql.DB
}

// closes all database connections
func (ctx *Context) Close() {
	for _, db := range ctx.open_databases {
		db.Close()
	}
}

type connectMigrateOpts struct {
	Type             types.DatabaseType
	SQliteConnection string
	PSQLConnection   string
	MigrateFn        func(*sql.DB, types.DatabaseType) error
}

// connects to the configured database and migrates the schema
func connectAndMigrate(opts *connectMigrateOpts) (*sql.DB, error) {
	var datasource string
	var dbtype string

	logrus.WithFields(logrus.Fields{
		"db_type":     opts.Type,
		"sqlite_conn": opts.SQliteConnection,
		"pg_conn":     opts.PSQLConnection,
	}).Info("Connecting and migrating")

	switch opts.Type {
	case types.DATABASE_DUMMY:
		return nil, nil
	case types.DATABASE_POSTGRES:
		datasource = opts.PSQLConnection
		dbtype = "postgres"
	default:
		// default to sqlite
		datasource = fmt.Sprintf("%s?_timeout=15000&_journal=WAL&_sync=NORMAL&_cache=shared", opts.SQliteConnection)
		dbtype = "sqlite3"
	}

	if opts.Type == types.DATABASE_POSTGRES && datasource == "" {
		// pg connection unconfigured
		return nil, nil
	}

	db, err := sql.Open(dbtype, datasource)
	if err != nil {
		return nil, err
	}

	if opts.Type == types.DATABASE_SQLITE {
		// enable wal on sqlite databases
		err = wal.EnableWAL(db)
		if err != nil {
			return nil, err
		}
	}

	err = opts.MigrateFn(db, opts.Type)
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

	logrus.WithFields(logrus.Fields{
		"world_dir": world_dir,
		"world.mt":  wc,
	}).Debug("Creating new DB context")
	ctx := &Context{}

	// map
	dbtype := types.DatabaseType(wc[worldconfig.CONFIG_MAP_BACKEND])
	map_db, err := connectAndMigrate(&connectMigrateOpts{
		Type:             dbtype,
		SQliteConnection: path.Join(world_dir, "map.sqlite"),
		PSQLConnection:   wc[worldconfig.CONFIG_PSQL_MAP_CONNECTION],
		MigrateFn:        block.MigrateBlockDB,
	})
	if err != nil {
		return nil, err
	}
	if map_db != nil {
		ctx.Blocks = block.NewBlockRepository(map_db, dbtype)
		if ctx.Blocks == nil {
			return nil, fmt.Errorf("invalid repository dbtype: %v", dbtype)
		}
		ctx.open_databases = append(ctx.open_databases, map_db)
	}

	// auth/privs
	dbtype = types.DatabaseType(wc[worldconfig.CONFIG_AUTH_BACKEND])
	auth_db, err := connectAndMigrate(&connectMigrateOpts{
		Type:             dbtype,
		SQliteConnection: path.Join(world_dir, "auth.sqlite"),
		PSQLConnection:   wc[worldconfig.CONFIG_PSQL_AUTH_CONNECTION],
		MigrateFn:        auth.MigrateAuthDB,
	})
	if err != nil {
		return nil, err
	}
	if auth_db != nil {
		ctx.Auth = auth.NewAuthRepository(auth_db, dbtype)
		ctx.Privs = auth.NewPrivilegeRepository(auth_db, dbtype)
		ctx.open_databases = append(ctx.open_databases, auth_db)
	}

	// mod storage
	dbtype = types.DatabaseType(wc[worldconfig.CONFIG_STORAGE_BACKEND])
	mod_storage_db, err := connectAndMigrate(&connectMigrateOpts{
		Type:             dbtype,
		SQliteConnection: path.Join(world_dir, "mod_storage.sqlite"),
		PSQLConnection:   wc[worldconfig.CONFIG_PSQL_MOD_STORAGE_CONNECTION],
		MigrateFn:        mod_storage.MigrateModStorageDB,
	})
	if err != nil {
		return nil, err
	}
	if mod_storage_db != nil {
		ctx.ModStorage = mod_storage.NewModStorageRepository(mod_storage_db, dbtype)
		ctx.open_databases = append(ctx.open_databases, mod_storage_db)
	}

	// players
	dbtype = types.DatabaseType(wc[worldconfig.CONFIG_PLAYER_BACKEND])
	player_db, err := connectAndMigrate(&connectMigrateOpts{
		Type:             dbtype,
		SQliteConnection: path.Join(world_dir, "players.sqlite"),
		PSQLConnection:   wc[worldconfig.CONFIG_PSQL_PLAYER_CONNECTION],
		MigrateFn:        player.MigratePlayerDB,
	})
	if err != nil {
		return nil, err
	}
	if player_db != nil {
		ctx.Player = player.NewPlayerRepository(player_db, dbtype)
		ctx.PlayerMetadata = player.NewPlayerMetadataRepository(player_db, dbtype)
		ctx.open_databases = append(ctx.open_databases, player_db)
	}

	return ctx, nil
}

// creates just the connection to the block-repository
func NewBlockDB(world_dir string) (block.BlockRepository, error) {
	logrus.WithFields(logrus.Fields{"world_dir": world_dir}).Debug("Creating new Block-DB")

	wc, err := worldconfig.Parse(path.Join(world_dir, "world.mt"))
	if err != nil {
		return nil, err
	}

	// map
	dbtype := types.DatabaseType(wc[worldconfig.CONFIG_MAP_BACKEND])
	map_db, err := connectAndMigrate(&connectMigrateOpts{
		Type:             dbtype,
		SQliteConnection: path.Join(world_dir, "map.sqlite"),
		PSQLConnection:   wc[worldconfig.CONFIG_PSQL_MAP_CONNECTION],
		MigrateFn:        block.MigrateBlockDB,
	})
	if err != nil {
		return nil, err
	}
	if map_db != nil {
		return block.NewBlockRepository(map_db, dbtype), nil
	}
	return nil, nil
}
