package auth_test

import (
	"testing"

	"github.com/minetest-go/mtdb/auth"
	"github.com/stretchr/testify/assert"
)

func testAuthRepository(t *testing.T, auth_repo *auth.AuthRepository, priv_repo *auth.PrivRepository) {
	// prepare test env
	e, err := auth_repo.GetByUsername("test")
	assert.NoError(t, err)
	if e != nil {
		assert.NoError(t, auth_repo.Delete(*e.ID))
	}

	// get by username (nonexistent)
	e, err = auth_repo.GetByUsername("test")
	assert.NoError(t, err)
	assert.Nil(t, e)

	// create new entry
	e = &auth.AuthEntry{
		Name:      "test",
		Password:  "blah",
		LastLogin: 123,
	}
	assert.NoError(t, auth_repo.Create(e))
	assert.NotNil(t, e.ID)

	// search
	s_username := "te%"
	list, err := auth_repo.Search(&auth.AuthSearch{
		Usernamelike: &s_username,
	})
	assert.NoError(t, err)
	assert.NotNil(t, list)
	assert.Equal(t, 1, len(list))

	// count
	count, err := auth_repo.Count(&auth.AuthSearch{
		Usernamelike: &s_username,
	})
	assert.NoError(t, err)
	assert.Equal(t, 1, count)

	// test duplicate
	auth2 := &auth.AuthEntry{
		Name:     "test",
		Password: "123",
	}
	assert.Error(t, auth_repo.Create(auth2))

	e.LastLogin = 456
	assert.NoError(t, auth_repo.Update(e))

	e, err = auth_repo.GetByUsername("test")
	assert.NoError(t, err)
	assert.NotNil(t, e)
	assert.Equal(t, "test", e.Name)
	assert.Equal(t, "blah", e.Password)
	assert.Equal(t, 456, e.LastLogin)
	assert.NotNil(t, e.ID)

	priv := &auth.PrivilegeEntry{
		ID:        *e.ID,
		Privilege: "interact",
	}
	assert.NoError(t, priv_repo.Create(priv))

	privs, err := priv_repo.GetByID(*e.ID)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(privs))

	assert.NoError(t, priv_repo.Delete(*e.ID, "interact"))

	assert.NoError(t, auth_repo.Delete(*e.ID))
}