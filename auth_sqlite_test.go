package mtdb

import (
	"database/sql"
	"io"
	"os"
	"testing"

	_ "modernc.org/sqlite"

	"github.com/stretchr/testify/assert"
)

func copyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}

func TestEmptySQliteRepo(t *testing.T) {
	// init stuff
	dbfile, err := os.CreateTemp(os.TempDir(), "auth.empty.sqlite")
	assert.NoError(t, err)

	// open db
	db, err := sql.Open("sqlite", "file:"+dbfile.Name())
	assert.NoError(t, err)
	repo := NewAuthRepository(db)
	assert.NotNil(t, repo)

	// existing entry
	entry, err := repo.GetByUsername("test")
	assert.Error(t, err)
	assert.Nil(t, entry)
}

func TestSQliteRepo(t *testing.T) {
	// init stuff
	dbfile, err := os.CreateTemp(os.TempDir(), "auth.sqlite")
	assert.NoError(t, err)
	assert.NotNil(t, dbfile)
	copyFileContents("testdata/auth.wal.sqlite", dbfile.Name())

	// open db
	db, err := sql.Open("sqlite", "file:"+dbfile.Name())
	assert.NoError(t, err)
	repo := NewAuthRepository(db)
	assert.NotNil(t, repo)

	// existing entry
	entry, err := repo.GetByUsername("test")
	assert.NoError(t, err)
	assert.NotNil(t, entry)
	assert.Equal(t, "test", entry.Name)
	assert.Equal(t, "#1#TxqLUa/uEJvZzPc3A0xwpA#oalXnktlS0bskc7bccsoVTeGwgAwUOyYhhceBu7wAyITkYjCtrzcDg6W5Co5V+oWUSG13y7TIoEfIg6rafaKzAbwRUC9RVGCeYRIUaa0hgEkIe9VkDmpeQ/kfF8zT8p7prOcpyrjWIJR+gmlD8Bf1mrxoPoBLDbvmxkcet327kQ9H4EMlIlv+w3XCufoPGFQ1UrfWiVqqK8dEmt/ldLPfxiK1Rg8MkwswEekymP1jyN9Cpq3w8spVVcjsxsAzI5M7QhSyqMMrIThdgBsUqMBOCULdV+jbRBBiA/ClywtZ8vvBpN9VGqsQuhmQG0h5x3fqPyR2XNdp9Ocm3zHBoJy/w", entry.Password)
	assert.Equal(t, int64(2), *entry.ID)
	assert.Equal(t, 1649603232, entry.LastLogin)

	// non-existing entry
	entry, err = repo.GetByUsername("bogus")
	assert.NoError(t, err)
	assert.Nil(t, entry)

	// create entry
	new_entry := &AuthEntry{
		Name:      "createduser",
		Password:  "blah",
		LastLogin: 456,
	}
	assert.NoError(t, repo.Create(new_entry))
	assert.NotNil(t, new_entry.ID)

	// check newly created entry
	entry, err = repo.GetByUsername("createduser")
	assert.NoError(t, err)
	assert.NotNil(t, entry)
	assert.Equal(t, new_entry.Name, entry.Name)
	assert.Equal(t, new_entry.Password, entry.Password)
	assert.Equal(t, *new_entry.ID, *entry.ID)
	assert.Equal(t, new_entry.LastLogin, entry.LastLogin)

	// change things
	new_entry.Name = "x"
	new_entry.Password = "y"
	new_entry.LastLogin = 123
	assert.NoError(t, repo.Update(new_entry))
	entry, err = repo.GetByUsername("x")
	assert.NoError(t, err)
	assert.NotNil(t, entry)
	assert.Equal(t, new_entry.Name, entry.Name)
	assert.Equal(t, new_entry.Password, entry.Password)
	assert.Equal(t, *new_entry.ID, *entry.ID)
	assert.Equal(t, new_entry.LastLogin, entry.LastLogin)

	// remove new user
	assert.NoError(t, repo.Delete(*new_entry.ID))
	entry, err = repo.GetByUsername("x")
	assert.NoError(t, err)
	assert.Nil(t, entry)

}

func TestSQlitePrivRepo(t *testing.T) {
	// init stuff
	dbfile, err := os.CreateTemp(os.TempDir(), "auth.sqlite")
	assert.NoError(t, err)
	assert.NotNil(t, dbfile)
	copyFileContents("testdata/auth.wal.sqlite", dbfile.Name())

	// open db
	db, err := sql.Open("sqlite", "file:"+dbfile.Name())
	assert.NoError(t, err)
	repo := NewPrivilegeRepository(db)
	assert.NotNil(t, repo)

	// read privs
	list, err := repo.GetByID(2)
	assert.NoError(t, err)
	assert.NotNil(t, list)
	assert.Equal(t, 2, len(list))

	privs := make(map[string]bool)
	for _, e := range list {
		privs[e.Privilege] = true
	}
	assert.True(t, privs["interact"])
	assert.True(t, privs["shout"])

	// create
	assert.NoError(t, repo.Create(&PrivilegeEntry{ID: 2, Privilege: "stuff"}))

	// verify
	list, err = repo.GetByID(2)
	assert.NoError(t, err)
	assert.NotNil(t, list)
	assert.Equal(t, 3, len(list))

	privs = make(map[string]bool)
	for _, e := range list {
		privs[e.Privilege] = true
	}
	assert.True(t, privs["interact"])
	assert.True(t, privs["shout"])
	assert.True(t, privs["stuff"])

	// delete
	assert.NoError(t, repo.Delete(2, "stuff"))

	list, err = repo.GetByID(2)
	assert.NoError(t, err)
	assert.NotNil(t, list)
	assert.Equal(t, 2, len(list))

	privs = make(map[string]bool)
	for _, e := range list {
		privs[e.Privilege] = true
	}
	assert.True(t, privs["interact"])
	assert.True(t, privs["shout"])
}
