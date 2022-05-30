package worldconfig_test

import (
	"fmt"
	"testing"

	"github.com/minetest-go/mtdb/worldconfig"
	"github.com/stretchr/testify/assert"
)

func TestParseSqlite(t *testing.T) {
	cfg, err := worldconfig.Parse("./testdata/world.mt.sqlite")
	assert.NoError(t, err)
	if cfg[worldconfig.CONFIG_AUTH_BACKEND] != worldconfig.BACKEND_SQLITE3 {
		t.Fatal("not sqlite3")
	}
}

func TestParsePostgres(t *testing.T) {
	cfg, err := worldconfig.Parse("./testdata/world.mt.postgres")
	assert.NoError(t, err)
	fmt.Println(cfg)
	if cfg[worldconfig.CONFIG_AUTH_BACKEND] != worldconfig.BACKEND_POSTGRES {
		t.Fatal("not postgres")
	}

	if cfg[worldconfig.CONFIG_PSQL_AUTH_CONNECTION] != "host=/var/run/postgresql user=postgres password=enter dbname=postgres" {
		t.Fatal("param err")
	}
}
