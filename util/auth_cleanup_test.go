package util_test

import (
	"database/sql"
	"testing"

	"github.com/minetest-go/mtdb/auth"
	"github.com/minetest-go/mtdb/player"
	"github.com/minetest-go/mtdb/types"
	"github.com/minetest-go/mtdb/util"
	"github.com/stretchr/testify/assert"

	_ "github.com/mattn/go-sqlite3"
)

func TestAuthCleanup(t *testing.T) {

	// open db
	db, err := sql.Open("sqlite3", ":memory:")
	assert.NoError(t, err)

	assert.NoError(t, auth.MigrateAuthDB(db, types.DATABASE_SQLITE))
	auth_repo := auth.NewAuthRepository(db, types.DATABASE_SQLITE)
	assert.NotNil(t, auth_repo)

	assert.NoError(t, player.MigratePlayerDB(db, types.DATABASE_SQLITE))
	player_repo := player.NewPlayerRepository(db, types.DATABASE_SQLITE)
	assert.NotNil(t, player_repo)

	assert.NoError(t, auth_repo.Create(&auth.AuthEntry{
		Name: "p1",
	}))

	// only auth entry exists
	assert.NoError(t, auth_repo.Create(&auth.AuthEntry{
		Name: "p2",
	}))

	assert.NoError(t, player_repo.CreateOrUpdate(&player.Player{
		Name: "p1",
	}))

	// verify before
	auth_count, err := auth_repo.Count(&auth.AuthSearch{})
	assert.NoError(t, err)
	assert.Equal(t, 2, auth_count)

	// cleanup
	count, err := util.AuthCleanup(auth_repo, player_repo)
	assert.NoError(t, err)
	assert.Equal(t, 1, count)

	// verify after
	auth_count, err = auth_repo.Count(&auth.AuthSearch{})
	assert.NoError(t, err)
	assert.Equal(t, 1, auth_count)

}
