package block_test

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
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

type testingLogWriter struct {
	t *testing.T
}

func (l testingLogWriter) Write(b []byte) (n int, err error) {
	l.t.Logf(string(b))
	return len(b), nil
}

func logToTesting(t *testing.T) {
	logrus.SetOutput(testingLogWriter{t})
	logrus.SetLevel(logrus.DebugLevel)
}
