package mtdb

import "database/sql"

func MigratePlayerDB(db *sql.DB, dbtype DatabaseType) error {
	var err error
	switch dbtype {
	case DATABASE_SQLITE:
		_, err = db.Exec(`
			CREATE TABLE IF NOT EXISTS player(
				name VARCHAR(50) NOT NULL,
				pitch NUMERIC(11, 4) NOT NULL,
				yaw NUMERIC(11, 4) NOT NULL,
				posX NUMERIC(11, 4) NOT NULL,
				posY NUMERIC(11, 4) NOT NULL,
				posZ NUMERIC(11, 4) NOT NULL,
				hp INT NOT NULL,
				breath INT NOT NULL,
				creation_date DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
				modification_date DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
				PRIMARY KEY (name)
			);
			CREATE TABLE IF NOT EXISTS player_metadata (
				player VARCHAR(50) NOT NULL,
				metadata VARCHAR(256) NOT NULL,
				value TEXT,
				PRIMARY KEY(player, metadata),
				FOREIGN KEY (player) REFERENCES player (name) ON DELETE CASCADE
			);
			CREATE TABLE IF NOT EXISTS player_inventories (
				player VARCHAR(50) NOT NULL,
				inv_id INT NOT NULL,
				inv_width INT NOT NULL,
				inv_name TEXT NOT NULL DEFAULT '',
				inv_size INT NOT NULL,
				PRIMARY KEY(player, inv_id),
				FOREIGN KEY (player) REFERENCES player (name) ON DELETE CASCADE
			);
			CREATE TABLE IF NOT EXISTS player_inventory_items (
				player VARCHAR(50) NOT NULL,
				inv_id INT NOT NULL,
				slot_id INT NOT NULL,
				item TEXT NOT NULL DEFAULT '',
				PRIMARY KEY(player, inv_id, slot_id),
				FOREIGN KEY (player) REFERENCES player (name) ON DELETE CASCADE
			);
		`)
	case DATABASE_POSTGRES:
		_, err = db.Exec(`
			CREATE TABLE IF NOT EXISTS player (
				name VARCHAR(60) NOT NULL,
				pitch NUMERIC(15, 7) NOT NULL,
				yaw NUMERIC(15, 7) NOT NULL,
				posX NUMERIC(15, 7) NOT NULL,
				posY NUMERIC(15, 7) NOT NULL,
				posZ NUMERIC(15, 7) NOT NULL,
				hp INT NOT NULL,
				breath INT NOT NULL,
				creation_date TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
				modification_date TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
				PRIMARY KEY (name)
			);
			CREATE TABLE IF NOT EXISTS player_inventories (
				player VARCHAR(60) NOT NULL,
				inv_id INT NOT NULL,
				inv_width INT NOT NULL,
				inv_name TEXT NOT NULL DEFAULT '',
				inv_size INT NOT NULL,
				PRIMARY KEY(player, inv_id),
				CONSTRAINT player_inventories_fkey FOREIGN KEY (player) REFERENCES 
				player (name) ON DELETE CASCADE
			);
			CREATE TABLE IF NOT EXISTS player_inventory_items (
				player VARCHAR(60) NOT NULL,
				inv_id INT NOT NULL,
				slot_id INT NOT NULL,
				item TEXT NOT NULL DEFAULT '',
				PRIMARY KEY(player, inv_id, slot_id),
				CONSTRAINT player_inventory_items_fkey FOREIGN KEY (player) REFERENCES 
				player (name) ON DELETE CASCADE
			);
			CREATE TABLE IF NOT EXISTS player_metadata (
				player VARCHAR(60) NOT NULL,
				attr VARCHAR(256) NOT NULL,
				value TEXT,
				PRIMARY KEY(player, attr),
				CONSTRAINT player_metadata_fkey FOREIGN KEY (player) REFERENCES 
				player (name) ON DELETE CASCADE
			);
		`)
	}
	return err
}
