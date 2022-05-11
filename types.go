package mtdb

type DatabaseType string

const (
	DATABASE_SQLITE   DatabaseType = "sqlite"
	DATABASE_POSTGRES DatabaseType = "postgres"
)
