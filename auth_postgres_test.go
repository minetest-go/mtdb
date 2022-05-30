package mtdb_test

import (
	"testing"

	"github.com/minetest-go/mtdb"
	"github.com/stretchr/testify/assert"
)

func TestPostgresAuthRepo(t *testing.T) {
	db, err := getPostgresDB(t)
	assert.NoError(t, err)

	assert.NoError(t, mtdb.MigrateAuthDB(db, mtdb.DATABASE_POSTGRES))

	auth_repo := mtdb.NewAuthRepository(db, mtdb.DATABASE_POSTGRES)
	priv_repo := mtdb.NewPrivilegeRepository(db, mtdb.DATABASE_POSTGRES)

	testAuthRepository(t, auth_repo, priv_repo)
}
