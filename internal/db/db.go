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

func InsertMessage(db *sql.DB, user, content string) error {
	_, err := db.Exec("INSERT INTO messages (user, content, preprocessed_content) VALUES (?, ?, ?)", user, content, "")
	return err
}

func GetMessages(db *sql.DB) ([]Message, error) {
	rows, err := db.Query("SELECT id, user, content, preprocessed_content, timestamp FROM messages ORDER BY timestamp DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var msg Message
		if err := rows.Scan(&msg.ID, &msg.User, &msg.Content, &msg.PreprocessedContent, &msg.Timestamp); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	return messages, nil
}

func GetMessagesByUser(db *sql.DB, user string) ([]Message, error) {
	rows, err := db.Query("SELECT id, user, content, preprocessed_content, timestamp FROM messages WHERE user = ? ORDER BY timestamp DESC", user)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var messages []Message

	for rows.Next() {
		var msg Message
		if err := rows.Scan(&msg.ID, &msg.User, &msg.Content, &msg.PreprocessedContent, &msg.Timestamp); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	return messages, nil
}

func DeleteMessage(db *sql.DB, id int) error {
	_, err := db.Exec("DELETE FROM messages WHERE id = ?", id)
	return err
}

func GetMessageByID(db *sql.DB, id int) (*Message, error) {
	row := db.QueryRow("SELECT id, user, content, preprocessed_content, timestamp FROM messages WHERE id = ?", id)
	var msg Message
	err := row.Scan(&msg.ID, &msg.User, &msg.Content, &msg.PreprocessedContent, &msg.Timestamp)
	if err != nil {
		return nil, err
	}
	return &msg, nil
}

func UpdateMessage(db *sql.DB, id int, content, preprocessedContent string) error {
	_, err := db.Exec("UPDATE messages SET content = ?, preprocessed_content = ? WHERE id = ?", content, preprocessedContent, id)
	return err
}

func GetMessageCount(db *sql.DB) (int, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM messages").Scan(&count)
	return count, err
}
