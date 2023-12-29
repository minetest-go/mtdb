package worldconfig

import (
	"bufio"
	"os"
	"strings"
)

const (
	BACKEND_SQLITE3  = "sqlite3"
	BACKEND_FILES    = "files"
	BACKEND_POSTGRES = "postgresql"
)

const (
	CONFIG_MAP_BACKEND                 = "backend"
	CONFIG_AUTH_BACKEND                = "auth_backend"
	CONFIG_STORAGE_BACKEND             = "mod_storage_backend"
	CONFIG_PLAYER_BACKEND              = "player_backend"
	CONFIG_PSQL_PLAYER_CONNECTION      = "pgsql_player_connection"
	CONFIG_PSQL_MAP_CONNECTION         = "pgsql_connection"
	CONFIG_PSQL_AUTH_CONNECTION        = "pgsql_auth_connection"
	CONFIG_PSQL_MOD_STORAGE_CONNECTION = "pgsql_mod_storage_connection"
)

const DEFAULT_CONFIG = `
mod_storage_backend = sqlite3
auth_backend = sqlite3
player_backend = sqlite3
backend = sqlite3
`

func Parse(filename string) (map[string]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	cfg := make(map[string]string)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		sepIndex := strings.Index(line, "=")
		if sepIndex < 0 {
			continue
		}

		valueStr := strings.Trim(line[sepIndex+1:], " ")
		keyStr := strings.Trim(line[:sepIndex], " ")

		cfg[keyStr] = valueStr
	}

	return cfg, nil
}
