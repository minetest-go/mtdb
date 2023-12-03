package wal_test

import (
	"database/sql"
	"os"
	"testing"

	_ "modernc.org/sqlite"

	"github.com/minetest-go/mtdb/wal"
	"github.com/stretchr/testify/assert"
)

func TestCheckJournalModeDelete(t *testing.T) {
	dbfile, err := os.CreateTemp(os.TempDir(), "auth.sqlite")
	assert.NoError(t, err)
	assert.NotNil(t, dbfile)
	copyFileContents("testdata/auth.sqlite", dbfile.Name())

	db, err := sql.Open("sqlite", "file:"+dbfile.Name()+"?mode=ro")
	assert.NoError(t, err)
	assert.Error(t, wal.EnableWAL(db))
}

func TestCheckJournalModeWal(t *testing.T) {
	dbfile, err := os.CreateTemp(os.TempDir(), "auth.wal.sqlite")
	assert.NoError(t, err)
	assert.NotNil(t, dbfile)
	copyFileContents("testdata/auth.wal.sqlite", dbfile.Name())

	db, err := sql.Open("sqlite", "file:"+dbfile.Name()+"?mode=ro")
	assert.NoError(t, err)
	assert.NoError(t, wal.EnableWAL(db))
}
