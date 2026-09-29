package db

import (
	"database/sql"
	"errors"
	"fmt"
	"io/fs"
	"os"

	_ "modernc.org/sqlite"
)

const schema = `
CREATE TABLE scheduler (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	date CHAR(8) NOT NULL DEFAULT '',
	title VARCHAR(256) NOT NULL DEFAULT '',
	comment TEXT NOT NULL DEFAULT '',
	repeat VARCHAR(128) NOT NULL DEFAULT ''
);
CREATE INDEX scheduler_date ON scheduler (date);
`

var db *sql.DB

func Init(dbFile string) error {
	install := false 
	_, err := os.Stat(dbFile)
	if errors.Is(err, fs.ErrNotExist) {
		install = true
	} else if err != nil {
		return fmt.Errorf("check db file: %w", err)
	}

	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return fmt.Errorf("open db: %w", err)
	}

	if install {
		if _, err := db.Exec(schema); err != nil {
			db.Close()
			return fmt.Errorf("create schema: %w", err)
		}
	}

	return nil
}

func Close() error {
	if db == nil {
		return nil
	}

	return db.Close()
}