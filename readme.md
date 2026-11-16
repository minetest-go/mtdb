Luanti database repositories for the `golang` ecosystem

![](https://github.com/luanti-golang/mtdb/workflows/test/badge.svg)

Docs: https://pkg.go.dev/github.com/luanti-golang/mtdb

# Features

* Read and write users/privs from and to the `auth` database
* Read and write player-data and metadata from and to the `player` database
* Read and write from and to the `map` (blocks) database
* Read and write from the `mod_storage` database

Supported databases:

* Sqlite3 (auth,player,blocks,mod_storage)
* Postgres (auth,player,blocks)

# License

Code: **MIT**