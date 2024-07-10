package database_models

// PostgresEnv represents the PostgreSQL environment
type PostgresEnv struct {
	Addr     string
	Database string
	User     string
	Password string
}
