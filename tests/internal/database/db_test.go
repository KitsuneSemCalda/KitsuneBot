package db_test

import (
	"KitsuneSemCalda/KitsuneBot/internal/db"
	"testing"
	"time"

	"github.com/caiolandgraf/gest/v2/gest"
)

func TestDB(t *testing.T) {
	s := gest.Describe("Database tests")

	s.It("should insert a message and retrieve it", func(t *gest.T) {
		err := db.Setup(":memory:")
		t.Expect(err).ToBeNil()
		defer db.DB.Close()

		res, err := db.DB.Exec(
			"INSERT INTO messages (user, content, preprocessed_content) VALUES (?, ?, ?)",
			"user1", "hello world", "hello_world",
		)
		t.Expect(err).ToBeNil()

		id, _ := res.LastInsertId()

		// Recuperar mensagem
		row := db.DB.QueryRow("SELECT content FROM messages WHERE id = ?", id)
		var content string
		err = row.Scan(&content)
		t.Expect(err).ToBeNil()
		t.Expect(content).ToBe("hello world")
	})

	s.It("should delete messages older than TTL", func(t *gest.T) {
		err := db.Setup(":memory:")
		t.Expect(err).ToBeNil()
		defer db.DB.Close()

		// Inserir mensagem antiga
		oldTimestamp := time.Now().AddDate(0, 0, -10) // 10 dias atrás
		_, err = db.DB.Exec(
			"INSERT INTO messages (user, content, preprocessed_content, timestamp) VALUES (?, ?, ?, ?)",
			"user2", "old message", "old_message", oldTimestamp,
		)
		t.Expect(err).ToBeNil()

		// Cleanup TTL de 7 dias
		err = db.CleanupTTL(db.DB, 7)
		t.Expect(err).ToBeNil()

		// Verificar se foi deletada
		row := db.DB.QueryRow("SELECT COUNT(*) FROM messages WHERE content = ?", "old message")
		var count int
		err = row.Scan(&count)
		t.Expect(err).ToBeNil()
		t.Expect(count).ToBe(0)
	})

	s.Run(t)
}
