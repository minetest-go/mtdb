package util

import (
	"fmt"

	"github.com/minetest-go/mtdb/auth"
	"github.com/minetest-go/mtdb/player"
)

func AuthCleanup(auth_repo *auth.AuthRepository, player_repo *player.PlayerRepository) (int, error) {
	count := 0

	auth_entries, err := auth_repo.Search(&auth.AuthSearch{})
	if err != nil {
		return 0, fmt.Errorf("error fetching auth entries: %v", err)
	}
	for _, auth_entry := range auth_entries {
		player_entry, err := player_repo.GetPlayer(auth_entry.Name)
		if err != nil {
			return 0, fmt.Errorf("error fetching player entry '%s': %v", auth_entry.Name, err)
		}

		if player_entry == nil {
			err = auth_repo.Delete(*auth_entry.ID)
			if err != nil {
				return 0, fmt.Errorf("error removing player entry '%s': %v", auth_entry.Name, err)
			}
			count++
		}
	}

	return count, nil
}
