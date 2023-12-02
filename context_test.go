package mtdb_test

import (
	"archive/zip"
	"bytes"
	"fmt"
	"os"
	"path"
	"testing"

	_ "github.com/lib/pq"
	_ "modernc.org/sqlite"

	"github.com/minetest-go/mtdb"
	"github.com/stretchr/testify/assert"
)

func repoSmokeTests(t *testing.T, repos *mtdb.Context) {
	p, err := repos.Player.GetPlayer("nonexistent")
	assert.NoError(t, err)
	assert.Nil(t, p)

	m, err := repos.PlayerMetadata.GetPlayerMetadata("nonexistent")
	assert.NoError(t, err)
	assert.NotNil(t, m)
	assert.Equal(t, 0, len(m))

	a, err := repos.Auth.GetByUsername("nonexistent")
	assert.NoError(t, err)
	assert.Nil(t, a)
}

func TestNoConfig(t *testing.T) {
	tmpdir, err := os.MkdirTemp(os.TempDir(), "mtdb")
	assert.NoError(t, err)

	repos, err := mtdb.New(tmpdir)
	assert.Error(t, err)
	assert.Nil(t, repos)
}

func TestNewSqlite(t *testing.T) {
	tmpdir := os.TempDir()
	contents := `
backend = sqlite3
auth_backend = sqlite3
player_backend = sqlite3
mod_storage_backend = sqlite3
	`
	err := os.WriteFile(path.Join(tmpdir, "world.mt"), []byte(contents), 0644)
	assert.NoError(t, err)

	repos, err := mtdb.New(tmpdir)
	assert.NoError(t, err)
	assert.NotNil(t, repos)
	assert.NotNil(t, repos.Auth)
	assert.NotNil(t, repos.Privs)
	assert.NotNil(t, repos.Blocks)
	assert.NotNil(t, repos.Player)
	assert.NotNil(t, repos.PlayerMetadata)
	assert.NotNil(t, repos.ModStorage)

	repoSmokeTests(t, repos)
}

func TestExportImportSqlite(t *testing.T) {
	tmpdir := os.TempDir()
	contents := `
backend = sqlite3
auth_backend = sqlite3
player_backend = sqlite3
mod_storage_backend = sqlite3
	`
	err := os.WriteFile(path.Join(tmpdir, "world.mt"), []byte(contents), 0644)
	assert.NoError(t, err)

	repos, err := mtdb.New(tmpdir)
	assert.NoError(t, err)

	// export
	buf := bytes.NewBuffer([]byte{})
	w := zip.NewWriter(buf)
	err = repos.Export(w)
	assert.NoError(t, err)
	err = w.Close()
	assert.NoError(t, err)
	zipfile, err := os.CreateTemp(os.TempDir(), "dump.zip")
	assert.NoError(t, err)
	f, err := os.Create(zipfile.Name())
	assert.NoError(t, err)
	count, err := f.Write(buf.Bytes())
	assert.NoError(t, err)
	assert.True(t, count > 0)

	// import
	z, err := zip.OpenReader(zipfile.Name())
	assert.NoError(t, err)
	err = repos.Import(&z.Reader)
	assert.NoError(t, err)
}

func TestNewSqliteWithDummyMap(t *testing.T) {
	tmpdir := os.TempDir()
	contents := `
backend = dummy
auth_backend = sqlite3
player_backend = sqlite3
mod_storage_backend = sqlite3
	`
	err := os.WriteFile(path.Join(tmpdir, "world.mt"), []byte(contents), 0644)
	assert.NoError(t, err)

	repos, err := mtdb.New(tmpdir)
	assert.NoError(t, err)
	assert.NotNil(t, repos)
	assert.NotNil(t, repos.Auth)
	assert.NotNil(t, repos.Privs)
	assert.Nil(t, repos.Blocks)
	assert.NotNil(t, repos.Player)
	assert.NotNil(t, repos.PlayerMetadata)
	assert.NotNil(t, repos.ModStorage)

	repoSmokeTests(t, repos)
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
pgsql_player_connection = ` + connStr + `
	`
	err := os.WriteFile(path.Join(tmpdir, "world.mt"), []byte(contents), 0644)
	assert.NoError(t, err)

	repos, err := mtdb.New(tmpdir)
	assert.NoError(t, err)
	assert.NotNil(t, repos)
	assert.NotNil(t, repos.Auth)
	assert.NotNil(t, repos.Privs)
	assert.NotNil(t, repos.Blocks)
	assert.NotNil(t, repos.Player)
	assert.NotNil(t, repos.PlayerMetadata)
	//assert.NotNil(t, repos.ModStorage)

	repoSmokeTests(t, repos)
}
