package player_test

import (
	"database/sql"
	"testing"

	"github.com/minetest-go/mtdb/player"
	"github.com/minetest-go/mtdb/types"
	"github.com/stretchr/testify/assert"
)

func ref[T any](s T) *T {
	return &s
}

func testPlayerSearch(t *testing.T, repo *player.PlayerRepository) {
	assert.NotNil(t, repo)

	p1 := &player.Player{
		Name:             "player1",
		Yaw:              1.4,
		Pitch:            0.9,
		PosX:             1200.0,
		PosY:             2300.0,
		PosZ:             4500.0,
		HP:               10,
		Breath:           9,
		CreationDate:     12000,
		ModificationDate: 15000,
	}
	assert.NoError(t, repo.CreateOrUpdate(p1))

	p2 := &player.Player{
		Name:             "player2",
		Yaw:              1.4,
		Pitch:            0.9,
		PosX:             1200.0,
		PosY:             2300.0,
		PosZ:             4500.0,
		HP:               10,
		Breath:           9,
		CreationDate:     12000,
		ModificationDate: 10000,
	}
	assert.NoError(t, repo.CreateOrUpdate(p2))

	// search all
	res, err := repo.Search(&player.PlayerSearch{})
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, 2, len(res))

	// search all with limit
	res, err = repo.Search(&player.PlayerSearch{Limit: ref(1)})
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, 1, len(res))

	// search all with limit, order by mod-date
	res, err = repo.Search(&player.PlayerSearch{
		Limit:          ref(1),
		OrderColumn:    ref(player.ModificationDate),
		OrderDirection: ref(player.Ascending),
	})
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, 1, len(res))
	assert.Equal(t, "player2", res[0].Name)

	// count all
	c, err := repo.Count(&player.PlayerSearch{})
	assert.NoError(t, err)
	assert.Equal(t, 2, c)

	// search by name
	res, err = repo.Search(&player.PlayerSearch{
		Name: ref("player1"),
	})
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, 1, len(res))
	assert.Equal(t, "player1", res[0].Name)

	// search by namelike
	res, err = repo.Search(&player.PlayerSearch{
		Namelike: ref("player%"),
	})
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, 2, len(res))

	// search by namelike (no match)
	res, err = repo.Search(&player.PlayerSearch{
		Namelike: ref("xy%"),
	})
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, 0, len(res))

	// delete
	assert.NoError(t, repo.RemovePlayer("player1"))
	assert.NoError(t, repo.RemovePlayer("player2"))
}

func TestPlayerSearchSQLite(t *testing.T) {
	// open db
	db, err := sql.Open("sqlite3", ":memory:")
	assert.NoError(t, err)

	assert.NoError(t, player.MigratePlayerDB(db, types.DATABASE_SQLITE))
	repo := player.NewPlayerRepository(db, types.DATABASE_SQLITE)
	testPlayerSearch(t, repo)
}

func TestPlayerSearchPostgres(t *testing.T) {
	// open db
	db, err := getPostgresDB(t)
	assert.NoError(t, err)
	assert.NotNil(t, db)
	db.Exec("delete from player")

	assert.NoError(t, player.MigratePlayerDB(db, types.DATABASE_POSTGRES))
	repo := player.NewPlayerRepository(db, types.DATABASE_POSTGRES)
	testPlayerSearch(t, repo)
}
