Minetest database repositories for the `golang` ecosystem

![](https://github.com/minetest-go/mtdb/workflows/test/badge.svg)
[![Coverage Status](https://coveralls.io/repos/github/minetest-go/mtdb/badge.svg)](https://coveralls.io/github/minetest-go/mtdb)

Docs: https://pkg.go.dev/github.com/minetest-go/mtdb

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