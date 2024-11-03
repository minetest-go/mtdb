package player_test

import (
	"database/sql"
	"os"
	"testing"

	"github.com/minetest-go/mtdb/player"
	"github.com/minetest-go/mtdb/types"
	"github.com/stretchr/testify/assert"
)

func testPlayerMetadata(t *testing.T, repo *player.PlayerMetadataRepository, prepo *player.PlayerRepository) {
	assert.NotNil(t, repo)
	assert.NoError(t, prepo.RemovePlayer("singleplayer"))
	assert.NoError(t, prepo.CreateOrUpdate(&player.Player{Name: "singleplayer"}))

	assert.NoError(t, repo.SetPlayerMetadata(&player.PlayerMetadata{Player: "singleplayer", Metadata: "x", Value: "y"}))
	mdlist, err := repo.GetPlayerMetadata("singleplayer")
	assert.NoError(t, err)
	assert.NotNil(t, mdlist)
	assert.Equal(t, 1, len(mdlist))
	assert.Equal(t, "singleplayer", mdlist[0].Player)
	assert.Equal(t, "x", mdlist[0].Metadata)
	assert.Equal(t, "y", mdlist[0].Value)
}

func TestPlayerMetadataSQlite(t *testing.T) {
	dbfile, err := os.CreateTemp(os.TempDir(), "playermetadata.sqlite")
	assert.NoError(t, err)
	db, err := sql.Open("sqlite3", dbfile.Name())
	assert.NoError(t, err)

	assert.NoError(t, player.MigratePlayerDB(db, types.DATABASE_SQLITE))
	repo := player.NewPlayerMetadataRepository(db, types.DATABASE_SQLITE)
	prepo := player.NewPlayerRepository(db, types.DATABASE_SQLITE)
	testPlayerMetadata(t, repo, prepo)
}

func TestPlayerMetadataPostgres(t *testing.T) {
	db, err := getPostgresDB(t)
	assert.NoError(t, err)

	assert.NoError(t, player.MigratePlayerDB(db, types.DATABASE_POSTGRES))
	repo := player.NewPlayerMetadataRepository(db, types.DATABASE_POSTGRES)
	prepo := player.NewPlayerRepository(db, types.DATABASE_POSTGRES)
	testPlayerMetadata(t, repo, prepo)
}
