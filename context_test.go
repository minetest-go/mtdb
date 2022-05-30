package mtdb_test

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/minetest-go/mtdb"
	"github.com/stretchr/testify/assert"
)

func TestNewSqlite(t *testing.T) {
	tmpdir := os.TempDir()
	contents := `
backend = sqlite3
auth_backend = sqlite3
player_backend = sqlite3
	`
	err := os.WriteFile(path.Join(tmpdir, "world.mt"), []byte(contents), 0644)
	assert.NoError(t, err)

	repos, err := mtdb.New(tmpdir)
	assert.NoError(t, err)
	assert.NotNil(t, repos)
	assert.NotNil(t, repos.Auth)
	assert.NotNil(t, repos.Privs)
	assert.NotNil(t, repos.Blocks)
}

func TestNewPostgres(t *testing.T) {
	if os.Getenv("PGHOST") == "" {
		t.SkipNow()
	}

	connStr := fmt.Sprintf(
		"user=%s password=%s port=%s host=%s dbname=%s sslmode=disable",
		os.Getenv("PGUSER"),
		os.Getenv("PGPASSWORD"),
		os.Getenv("PGPORT"),
		os.Getenv("PGHOST"),
		os.Getenv("PGDATABASE"))

	tmpdir := os.TempDir()
	contents := `
backend = postgresql
pgsql_connection = ` + connStr + `
auth_backend = postgresql
pgsql_auth_connection = ` + connStr + `
player_backend = postgresql
	`
	err := os.WriteFile(path.Join(tmpdir, "world.mt"), []byte(contents), 0644)
	assert.NoError(t, err)

	repos, err := mtdb.New(tmpdir)
	assert.NoError(t, err)
	assert.NotNil(t, repos)
	assert.NotNil(t, repos.Auth)
	assert.NotNil(t, repos.Privs)
	assert.NotNil(t, repos.Blocks)
}
