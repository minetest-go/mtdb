Minetest database repositories for the `golang` ecosystem

![](https://github.com/minetest-go/mtdb/workflows/test/badge.svg)
[![Coverage Status](https://coveralls.io/repos/github/minetest-go/mtdb/badge.svg)](https://coveralls.io/github/minetest-go/mtdb)

Docs: https://pkg.go.dev/github.com/minetest-go/mtdb

# Features

* Read and write users/privs from and to the `auth` database
* Read and write from and to the `map` database
* Read and write from the `mod_storage` database

Supported databases:

* Sqlite3 (auth,blocks,mod_storage)
* Postgres (auth.blocks)

# License

Code: **MIT**