package mtdb_test

import (
	"testing"

	"github.com/minetest-go/mtdb"
	"github.com/stretchr/testify/assert"
)

func TestPostgresMigratePlayer(t *testing.T) {
	db, err := getPostgresDB(t)
	assert.NoError(t, err)

	assert.NoError(t, mtdb.MigratePlayerDB(db, mtdb.DATABASE_POSTGRES))
}

func TestPostgresPlayerRepo(t *testing.T) {
	db, err := getPostgresDB(t)
	assert.NoError(t, err)
	assert.NotNil(t, db)
	// TODO
}
