package mtdb

import (
	"database/sql"
	"fmt"
	"os"
	"testing"
)

func getPostgresDB(t *testing.T) (*sql.DB, error) {
	if os.Getenv("PGHOST") == "" {
		t.SkipNow()
	}

	connStr := fmt.Sprintf(
		"user=%s password=%s port=%s host=%s dbname=%s sslmode=disable",
		os.Getenv("PGUSER"),
		os.Getenv("PGPASSWORD"),
		os.Getenv("PGPORT"),
		os.Getenv("PGHOST"),
		os.Getenv("PGDATABASE"))

	return sql.Open("postgres", connStr)
}
