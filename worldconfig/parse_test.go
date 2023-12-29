package worldconfig_test

import (
	"testing"

	"github.com/minetest-go/mtdb/worldconfig"
	"github.com/stretchr/testify/assert"
)

func TestInvalidFile(t *testing.T) {
	cfg, err := worldconfig.Parse("./testdata/non-existing.file")
	assert.Error(t, err)
	assert.Nil(t, cfg)
}

func TestParseSqlite(t *testing.T) {
	cfg, err := worldconfig.Parse("./testdata/world.mt.sqlite")
	assert.NoError(t, err)
	assert.Equal(t, worldconfig.BACKEND_SQLITE3, cfg[worldconfig.CONFIG_AUTH_BACKEND])
}

func TestParsePostgres(t *testing.T) {
	cfg, err := worldconfig.Parse("./testdata/world.mt.postgres")
	assert.NoError(t, err)
	assert.Equal(t, worldconfig.BACKEND_POSTGRES, cfg[worldconfig.CONFIG_AUTH_BACKEND])
	assert.Equal(t, "host=/var/run/postgresql user=postgres password=enter dbname=postgres", cfg[worldconfig.CONFIG_PSQL_AUTH_CONNECTION])
	assert.Equal(t, "host=postgres port=5432 user=postgres password=enter dbname=postgres", cfg[worldconfig.CONFIG_PSQL_MOD_STORAGE_CONNECTION])
}
