package mtdb

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	_ "github.com/lib/pq"

	"github.com/stretchr/testify/assert"
)

func IsDatabaseAvailable() bool {
	return os.Getenv("PGHOST") != ""
}

func TestPostgresDB(t *testing.T) {
	if !IsDatabaseAvailable() {
		t.SkipNow()
	}

	connStr := fmt.Sprintf(
		"user=%s password=%s port=%s host=%s dbname=%s sslmode=disable",
		os.Getenv("PGUSER"),
		os.Getenv("PGPASSWORD"),
		os.Getenv("PGPORT"),
		os.Getenv("PGHOST"),
		os.Getenv("PGDATABASE"))

	db, err := sql.Open("postgres", connStr)
	assert.NoError(t, err)

	assert.NoError(t, MigrateAuth(db, "postgres"))

	// delete existing data
	_, err = db.Exec("delete from auth;")
	assert.NoError(t, err)

	repo := NewAuthRepository(db, DATABASE_POSTGRES)
	assert.NotNil(t, repo)

	user := &AuthEntry{
		Name:     "test",
		Password: "dummy",
	}

	assert.NoError(t, repo.Create(user))

	entry, err := repo.GetByUsername("test")
	assert.NoError(t, err)
	assert.NotNil(t, entry)

}
