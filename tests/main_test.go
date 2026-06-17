package tests

import (
    "database/sql"
    "os"
    "testing"
    _ "modernc.org/sqlite"
)

func TestMain(m *testing.M) {
    dbFile := DBFile
    if env := os.Getenv("TODO_DBFILE"); env != "" {
        dbFile = env
    }
    db, err := sql.Open("sqlite", dbFile)
    if err != nil {
        panic(err)
    }
    defer db.Close()
    _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS scheduler (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            date CHAR(8) NOT NULL DEFAULT '',
            title VARCHAR(255) NOT NULL,
            comment TEXT,
            repeat VARCHAR(128) NOT NULL DEFAULT ''
        );
        CREATE INDEX IF NOT EXISTS idx_date ON scheduler(date);
    `)
    if err != nil {
        panic(err)
    }
    os.Exit(m.Run())
}
