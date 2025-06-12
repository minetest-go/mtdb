package block_test

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func getPostgresDB(t *testing.T) *sql.DB {
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

	db, err := sql.Open("postgres", connStr)
	assert.NoError(t, err)
	return db
}

type testingLogWriter struct {
	t *testing.T
}

func (l testingLogWriter) Write(b []byte) (n int, err error) {
	l.t.Log(string(b))
	return len(b), nil
}

func logToTesting(t *testing.T) {
	logrus.SetOutput(testingLogWriter{t})
	logrus.SetLevel(logrus.DebugLevel)
}

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
