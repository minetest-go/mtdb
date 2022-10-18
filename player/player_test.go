package player_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/minetest-go/mtdb/player"
	"github.com/stretchr/testify/assert"
)

func testRepository(t *testing.T, repo player.PlayerRepository) {
	assert.NotNil(t, repo)

	p1 := &player.Player{
		Name:             fmt.Sprintf("player-%d", rand.Int()),
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
	assert.NoError(t, repo.CreateOrUpdate(p1))

	p2, err := repo.GetPlayer(p1.Name)
	assert.NoError(t, err)
	assert.NotNil(t, p2)
	assert.Equal(t, p1.Breath, p2.Breath)
	assert.Equal(t, p1.HP, p2.HP)
	assert.Equal(t, p1.CreationDate, p2.CreationDate)
	assert.Equal(t, p1.ModificationDate, p2.ModificationDate)
	assert.Equal(t, p1.Name, p2.Name)
}
