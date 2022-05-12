Minetest database repositories for the `golang` ecosystem

![](https://github.com/minetest-go/mtdb/workflows/test/badge.svg)
[![Coverage Status](https://coveralls.io/repos/github/minetest-go/mtdb/badge.svg)](https://coveralls.io/github/minetest-go/mtdb)

Docs: https://pkg.go.dev/github.com/minetest-go/mtdb

# Features

* Read and write users/privs from and to the `auth` database
* Read and write from and to the `map` database

Supported databases:

* Sqlite3
* Postgres

# Usage

Read from an existing `auth.sqlite` database:
```golang
import (
    "database/sql"
    _ "modernc.org/sqlite"
    "github/minetest-go/mtdb"
)

func main() {
    auth_db, err := sql.Open("sqlite", "file:auth.sqlite")
    if err != nil {
        panic(err)
    }

    // Enable the wal mode for concurrent access
    err = mtdb.EnableWAL(auth_db)
    if err != nil {
        panic(err)
    }

    auth_repo := mtdb.NewAuthRepository(auth_db, mtdb.DATABASE_SQLITE)
    priv_repo := mtdb.NewPrivilegeRepository(auth_db, mtdb.DATABASE_SQLITE)

    // Read a user
    admin_user, err := auth_repo.GetByUsername("admin")
    if err != nil {
        panic(err)
    }

    fmt.Printf("User: %s, Last-login: %d\n", admin_user.Name, admin_user.LastLogin)

    // read privileges
    admin_privs, err := priv_repo.GetByID(*admin_user.ID)
    for _, priv := range admin_privs {
        fmt.Printf("+ Priv: %s\n", priv.Privilege)
    }
}
```


# License

Code: **MIT**