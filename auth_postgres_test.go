package mtdb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPostgresDB(t *testing.T) {
	db, err := getPostgresDB(t)
	assert.NoError(t, err)

	assert.NoError(t, MigrateAuthDB(db, DATABASE_POSTGRES))

	auth_repo := NewAuthRepository(db, DATABASE_POSTGRES)
	priv_repo := NewPrivilegeRepository(db, DATABASE_POSTGRES)

	testAuthRepository(t, auth_repo, priv_repo)
}
