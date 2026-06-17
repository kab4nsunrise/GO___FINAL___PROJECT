package db

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

var db *sql.DB

const schema = `
CREATE TABLE IF NOT EXISTS scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT '',
    title VARCHAR(255) NOT NULL,
    comment TEXT,
    repeat VARCHAR(128) NOT NULL DEFAULT ''
);
CREATE INDEX IF NOT EXISTS idx_date ON scheduler(date);
`

func Init(dbFile string) error {
	var err error
	db, err = sql.Open("sqlite", dbFile)
	if err != nil {
		log.Printf("Failed to open DB: %v", err)
		return err
	}
	_, err = db.Exec(schema)
	if err != nil {
		log.Printf("Failed to create table: %v", err)
		return err
	}
	log.Println("Database initialized successfully")
	return nil
}

func GetDB() *sql.DB {
	return db
}
