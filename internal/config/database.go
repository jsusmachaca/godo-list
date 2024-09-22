package config

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func GetConnection() (*sql.DB, error) {
	var err error
	db, err := sql.Open("sqlite3", "db.sqlite3")
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func Migrate(db *sql.DB) error {
	const create string = `
	CREATE TABLE IF NOT EXISTS tasks (
	id VARCHAR NOT NULL PRIMARY KEY,
	name VARCHAR NOT NULL,
	done BOOL DEFAULT (true)
	);`

	_, err := db.Exec(create)
	if err != nil {
		return err
	}
	return nil
}
