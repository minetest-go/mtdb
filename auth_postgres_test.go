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
	repo := NewAuthRepository(db)
	assert.NotNil(t, repo)

	entry, err := repo.GetByUsername("test")
	assert.NoError(t, err)
	assert.NotNil(t, entry)

}
