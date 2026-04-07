package to

import (
	"BlockMe/config"
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func newPostgres() (*sql.DB, error) {
	fmt.Println("to: Connecting to Postgres db")
	err := ensurePGDBExists()
	if err != nil {
		fmt.Printf("BlockMe was unable to ensure the '%s' database exists. Check out this error >>\n",
			config.Env.DbName)
		fmt.Printf("%s\n", err.Error())
		fmt.Printf("Maybe you didn't give us access to the default 'postgres' database? Or maybe you forgot "+
			"that you set the DB_TYPE environment variable to postgres but never actually set up a Postgres database "+
			"instance for us to use? Doesn't matter, we are going to throw a connection attempt at %s:%d/%s anyway to "+
			"see if it sticks\n", config.Env.DbHost, config.Env.DbPort, config.Env.DbName)
	}

	// Example: postgres://username:password@localhost:5432/database_name
	pgURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		config.Env.DbUser,
		config.Env.DbPass,
		config.Env.DbHost,
		config.Env.DbPort,
		config.Env.DbName,
	)

	db, err := sql.Open("pgx", pgURL)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to postgres database: %w", err)
	}

	return db, nil
}

// ensurePGDBExists attempts to create a connection to the default 'postgres' database.
func ensurePGDBExists() error {
	pgURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/postgres",
		config.Env.DbUser,
		config.Env.DbPass,
		config.Env.DbHost,
		config.Env.DbPort,
	)
	db, err := sql.Open("pgx", pgURL)
	if err != nil {
		return fmt.Errorf(
			"unable to ensure the existence of the %s database: %w",
			config.Env.DbName,
			err,
		)
	}
	defer db.Close()

	sqlStatement := `SELECT EXISTS (SELECT 1 FROM pg_database WHERE datname = $1);`
	var exists bool
	err = db.QueryRow(sqlStatement, config.Env.DbName).Scan(&exists)
	if err != nil {
		return fmt.Errorf(
			"unable to ensure the existence of the %s database: %w",
			config.Env.DbName,
			err,
		)
	}

	if !exists {
		// Can't use positional parameters ($1, $2...) for identifiers
		sqlStatement = fmt.Sprintf("CREATE DATABASE %s OWNER %s;", config.Env.DbName, config.Env.DbUser)
		_, err = db.Exec(sqlStatement)
		if err != nil {
			return fmt.Errorf("unable to create the %s database: %w", config.Env.DbName, err)
		}
	}
	return nil
}
