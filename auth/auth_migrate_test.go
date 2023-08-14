package auth_test

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"

	"github.com/minetest-go/mtdb/auth"
	"github.com/minetest-go/mtdb/types"
	"github.com/stretchr/testify/assert"
)

func TestMigrateAuthSQlite(t *testing.T) {
	// open db
	db, err := sql.Open("sqlite3", ":memory:")
	assert.NoError(t, err)

	assert.NoError(t, auth.MigrateAuthDB(db, types.DATABASE_SQLITE))
}

func TestMigrateAuthPostgres(t *testing.T) {
	db, err := getPostgresDB(t)
	assert.NoError(t, err)

	assert.NoError(t, auth.MigrateAuthDB(db, types.DATABASE_POSTGRES))
}
