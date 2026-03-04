package to

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func newSqlite() (*sql.DB, error) {
	fmt.Println("to: Connecting to Sqlite db")

	wd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("i guess we don't know where tf we are??: %w", err)
	}
	dbFilePath := fmt.Sprintf("%s/database", wd)
	bdFileName := "database.sqlite3"

	_, err = os.Stat(dbFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(dbFilePath, os.ModePerm)
			if err != nil {
				return nil, fmt.Errorf("unable to create database directory: %w", err)
			}
		} else {
			fmt.Println("Couldn't check if the database directory exists - Going to continue like there is " +
				"nothing wrong, but things are probably broken...")
		}
	}

	sqliteURL := fmt.Sprintf("%s/%s", dbFilePath, bdFileName)

	db, err := sql.Open("sqlite3", sqliteURL)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to sqlite database: %w", err)
	}

	_ = db.Ping()

	return db, nil
}
