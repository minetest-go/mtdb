package auth_test

import (
	"archive/zip"
	"bytes"
	"os"
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

	// search 1
	s_username := "te%"
	list, err := auth_repo.Search(&auth.AuthSearch{
		Usernamelike: &s_username,
	})
	assert.NoError(t, err)
	assert.NotNil(t, list)
	assert.Equal(t, 1, len(list))

	// search 2
	orderCol := auth.Name
	orderDir := auth.Ascending
	list, err = auth_repo.Search(&auth.AuthSearch{
		Usernamelike:   &s_username,
		OrderColumn:    &orderCol,
		OrderDirection: &orderDir,
	})
	assert.NoError(t, err)
	assert.NotNil(t, list)
	assert.Equal(t, 1, len(list))

	// search 3
	orderCol = auth.LastLogin
	orderDir = auth.Descending
	list, err = auth_repo.Search(&auth.AuthSearch{
		Usernamelike:   &s_username,
		OrderColumn:    &orderCol,
		OrderDirection: &orderDir,
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

	// update
	e.LastLogin = 456
	assert.NoError(t, auth_repo.Update(e))

	// check update
	e, err = auth_repo.GetByUsername("test")
	assert.NoError(t, err)
	assert.NotNil(t, e)
	assert.Equal(t, "test", e.Name)
	assert.Equal(t, "blah", e.Password)
	assert.Equal(t, 456, e.LastLogin)
	assert.NotNil(t, e.ID)

	// add priv
	priv := &auth.PrivilegeEntry{
		ID:        *e.ID,
		Privilege: "interact",
	}
	assert.NoError(t, priv_repo.Create(priv))

	// check update priv
	privs, err := priv_repo.GetByID(*e.ID)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(privs))

	// export
	buf := bytes.NewBuffer([]byte{})
	w := zip.NewWriter(buf)
	err = auth_repo.Export(w)
	assert.NoError(t, err)
	err = priv_repo.Export(w)
	assert.NoError(t, err)
	err = w.Close()
	assert.NoError(t, err)
	zipfile, err := os.CreateTemp(os.TempDir(), "auth.zip")
	assert.NoError(t, err)
	f, err := os.Create(zipfile.Name())
	assert.NoError(t, err)
	count, err = f.Write(buf.Bytes())
	assert.NoError(t, err)
	assert.True(t, count > 0)

	// delete priv
	assert.NoError(t, priv_repo.Delete(*e.ID, "interact"))

	// delete auth
	assert.NoError(t, auth_repo.Delete(*e.ID))

	// delete all
	assert.NoError(t, auth_repo.DeleteAll())

	// import
	z, err := zip.OpenReader(zipfile.Name())
	assert.NoError(t, err)
	err = auth_repo.Import(&z.Reader)
	assert.NoError(t, err)
	err = priv_repo.Import(&z.Reader)
	assert.NoError(t, err)

	// check import
	e, err = auth_repo.GetByUsername("test")
	assert.NoError(t, err)
	assert.NotNil(t, e)
	assert.Equal(t, "test", e.Name)
	assert.Equal(t, "blah", e.Password)
	assert.Equal(t, 456, e.LastLogin)
	assert.NotNil(t, e.ID)

	// check imported priv
	privs, err = priv_repo.GetByID(*e.ID)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(privs))

}
