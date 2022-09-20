package types

type DatabaseType string

const (
	DATABASE_SQLITE   DatabaseType = "sqlite3"
	DATABASE_POSTGRES DatabaseType = "postgresql"
)
