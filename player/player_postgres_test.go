package player_test

import (
	"testing"

	"github.com/minetest-go/mtdb/player"
	"github.com/minetest-go/mtdb/types"
	"github.com/stretchr/testify/assert"
)

func TestPostgresMigratePlayer(t *testing.T) {
	db, err := getPostgresDB(t)
	assert.NoError(t, err)

	assert.NoError(t, player.MigratePlayerDB(db, types.DATABASE_POSTGRES))
}

func TestPostgresPlayerRepo(t *testing.T) {
	db, err := getPostgresDB(t)
	assert.NoError(t, err)
	assert.NotNil(t, db)

	repo := player.NewPlayerRepository(db, types.DATABASE_POSTGRES)
	testRepository(t, repo)
}
