package mtdb

import (
	"database/sql"
	"path"

	_ "github.com/lib/pq"
	"github.com/minetest-go/mtdb/worldconfig"
	_ "modernc.org/sqlite"
)

type Context struct {
	Auth    AuthRepository
	Privs   *PrivRepository
	Blocks  BlockRepository
	map_db  *sql.DB
	auth_db *sql.DB
}

// closes all database connections
func (ctx *Context) Close() {
	ctx.map_db.Close()
	ctx.auth_db.Close()
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

		err = EnableWAL(map_db)
		if err != nil {
			return nil, err
		}

		err = MigrateBlockDB(map_db, DATABASE_SQLITE)
		if err != nil {
			return nil, err
		}

		ctx.Blocks = NewBlockRepository(map_db, DATABASE_SQLITE)
		ctx.map_db = map_db

	case worldconfig.BACKEND_POSTGRES:
		map_db, err := sql.Open("postgres", wc[worldconfig.CONFIG_PSQL_MAP_CONNECTION])
		if err != nil {
			return nil, err
		}

		err = MigrateBlockDB(map_db, DATABASE_POSTGRES)
		if err != nil {
			return nil, err
		}

		ctx.Blocks = NewBlockRepository(map_db, DATABASE_POSTGRES)
		ctx.map_db = map_db
	}

	// create auth repos
	switch wc[worldconfig.CONFIG_AUTH_BACKEND] {
	case worldconfig.BACKEND_SQLITE3:
		auth_db, err := sql.Open("sqlite", path.Join(world_dir, "auth.sqlite"))
		if err != nil {
			return nil, err
		}

		err = EnableWAL(auth_db)
		if err != nil {
			return nil, err
		}

		err = MigrateAuthDB(auth_db, DATABASE_SQLITE)
		if err != nil {
			return nil, err
		}

		ctx.Auth = NewAuthRepository(auth_db, DATABASE_SQLITE)
		ctx.Privs = NewPrivilegeRepository(auth_db, DATABASE_SQLITE)
		ctx.auth_db = auth_db

	case worldconfig.BACKEND_POSTGRES:
		auth_db, err := sql.Open("postgres", wc[worldconfig.CONFIG_PSQL_AUTH_CONNECTION])
		if err != nil {
			return nil, err
		}

		err = MigrateAuthDB(auth_db, DATABASE_POSTGRES)
		if err != nil {
			return nil, err
		}

		ctx.Auth = NewAuthRepository(auth_db, DATABASE_POSTGRES)
		ctx.Privs = NewPrivilegeRepository(auth_db, DATABASE_POSTGRES)
		ctx.auth_db = auth_db

	}

	return ctx, nil
}
