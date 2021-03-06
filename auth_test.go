package mtdb_test

import (
	"testing"

	"github.com/minetest-go/mtdb"
	"github.com/stretchr/testify/assert"
)

func testAuthRepository(t *testing.T, auth_repo *mtdb.AuthRepository, priv_repo *mtdb.PrivRepository) {
	// prepare test env
	auth, err := auth_repo.GetByUsername("test")
	assert.NoError(t, err)
	if auth != nil {
		assert.NoError(t, auth_repo.Delete(*auth.ID))
	}

	// get by username (nonexistent)
	auth, err = auth_repo.GetByUsername("test")
	assert.NoError(t, err)
	assert.Nil(t, auth)

	// create new entry
	auth = &mtdb.AuthEntry{
		Name:      "test",
		Password:  "blah",
		LastLogin: 123,
	}
	assert.NoError(t, auth_repo.Create(auth))
	assert.NotNil(t, auth.ID)

	// search
	s_username := "te%"
	list, err := auth_repo.Search(&mtdb.AuthSearch{
		Usernamelike: &s_username,
	})
	assert.NoError(t, err)
	assert.NotNil(t, list)
	assert.Equal(t, 1, len(list))

	// count
	count, err := auth_repo.Count(&mtdb.AuthSearch{
		Usernamelike: &s_username,
	})
	assert.NoError(t, err)
	assert.Equal(t, 1, count)

	// test duplicate
	auth2 := &mtdb.AuthEntry{
		Name:     "test",
		Password: "123",
	}
	assert.Error(t, auth_repo.Create(auth2))

	auth.LastLogin = 456
	assert.NoError(t, auth_repo.Update(auth))

	auth, err = auth_repo.GetByUsername("test")
	assert.NoError(t, err)
	assert.NotNil(t, auth)
	assert.Equal(t, "test", auth.Name)
	assert.Equal(t, "blah", auth.Password)
	assert.Equal(t, 456, auth.LastLogin)
	assert.NotNil(t, auth.ID)

	priv := &mtdb.PrivilegeEntry{
		ID:        *auth.ID,
		Privilege: "interact",
	}
	assert.NoError(t, priv_repo.Create(priv))

	privs, err := priv_repo.GetByID(*auth.ID)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(privs))

	assert.NoError(t, priv_repo.Delete(*auth.ID, "interact"))

	assert.NoError(t, auth_repo.Delete(*auth.ID))
}
