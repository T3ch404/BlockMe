package to

// 'to' is short for "BlockMe Database Handler"

import (
	"BlockMe/config"
	"database/sql"
	"fmt"
)

var DB *sql.DB

func Setup() error {
	fmt.Printf("\nInitializing database connection...\n")

	var err error
	switch config.Env.DbType {
	case "postgres":
		// Setup and return Postgres database connection
		DB, err = newPostgres()
	case "sqlite":
		// Setup and return Sqlite database connection
		DB, err = newSqlite()
	default:
		fmt.Println("DB_TYPE is either invalid or not set. Continuing with sqlite")
		DB, err = newSqlite()
	}
	if err != nil {
		return err
	}

	err = ensureTablesExist()
	if err != nil {
		return err
	}

	fmt.Println("to: Successfully connected to database")
	return nil
}

func ensureTablesExist() error {
	fmt.Printf("to: Checking database tables...\n")

	dbType := config.Env.DbType
	if dbType == "" {
		dbType = "sqlite"
	}

	var sqlStatements []string
	if dbType == "sqlite" {
		sqlStatements = append(sqlStatements, `CREATE TABLE IF NOT EXISTS ip_blocklist (ip TEXT PRIMARY KEY, timestamp TIMESTAMP NOT NULL);`)
	} else if dbType == "postgres" {
		sqlStatements = append(sqlStatements, `CREATE TABLE IF NOT EXISTS ip_blocklist (ip INET PRIMARY KEY, timestamp TIMESTAMPTZ NOT NULL);`)
	}

	for _, statement := range sqlStatements {
		_, err := DB.Exec(statement)
		if err != nil {
			return fmt.Errorf("could not create required table(s): %w", err)
		}
	}

	return nil
}
