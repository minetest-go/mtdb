package mtdb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPostgresDB(t *testing.T) {
	db, err := getPostgresDB(t)
	assert.NoError(t, err)

	assert.NoError(t, MigrateAuthPostgres(db))

	// delete existing data
	_, err = db.Exec("delete from auth;")
	assert.NoError(t, err)

	auth_repo := NewAuthRepository(db, DATABASE_POSTGRES)
	assert.NotNil(t, auth_repo)
	priv_repo := NewPrivilegeRepository(db, DATABASE_POSTGRES)
	assert.NotNil(t, priv_repo)

	user := &AuthEntry{
		Name:     "test",
		Password: "dummy",
	}

	assert.NoError(t, auth_repo.Create(user))
	assert.NotNil(t, user.ID)

	entry, err := auth_repo.GetByUsername("test")
	assert.NoError(t, err)
	assert.NotNil(t, entry)

	priv := &PrivilegeEntry{
		ID:        *user.ID,
		Privilege: "interact",
	}

	assert.NoError(t, priv_repo.Create(priv))
}
