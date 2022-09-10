package auth_test

import (
	"testing"

	"github.com/minetest-go/mtdb/auth"
	"github.com/minetest-go/mtdb/types"
	"github.com/stretchr/testify/assert"
)

func TestPostgresAuthRepo(t *testing.T) {
	db, err := getPostgresDB(t)
	assert.NoError(t, err)

	assert.NoError(t, auth.MigrateAuthDB(db, types.DATABASE_POSTGRES))

	auth_repo := auth.NewAuthRepository(db, types.DATABASE_POSTGRES)
	priv_repo := auth.NewPrivilegeRepository(db, types.DATABASE_POSTGRES)

	testAuthRepository(t, auth_repo, priv_repo)
}
