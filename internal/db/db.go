package db

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path)

	if err != nil {
		return nil, err
	}

	createTable := `
    CREATE TABLE IF NOT EXISTS messages (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        user TEXT NOT NULL,
        content TEXT NOT NULL,
        preprocessed_content TEXT,
        timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
    );`

	_, err = db.Exec(createTable)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func CleanupTTL(db *sql.DB, days int) error {
	cut_off := time.Now().AddDate(0, 0, -days).Format("2006-01-02 15:04:05")
	_, err := db.Exec("DELETE FROM messages WHERE timestamp < ?", cut_off)
	return err
}

func Setup(path string) error {
	var err error
	DB, err = InitDB(path)
	return err
}
