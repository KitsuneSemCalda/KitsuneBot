package db_test

import (
	"KitsuneSemCalda/KitsuneBot/internal/db"
	"os"
	"testing"
	"time"

	"github.com/caiolandgraf/gest/v2/gest"
)

func TestMainFlow(t *testing.T) {
	gest.Describe("Main flow integration tests").
		It("should setup and cleanup messages", func(t *gest.T) {
			tmpFile, err := os.CreateTemp("", "test_*.db")
			t.Expect(err).ToBeNil()
			defer os.Remove(tmpFile.Name())
			tmpFile.Close()

			err = db.Setup(tmpFile.Name())
			t.Expect(err).ToBeNil()
			defer db.DB.Close()

			db.InsertMessage(db.DB, "user1", "recent message")
			db.InsertMessage(db.DB, "user2", "old message")

			oldTimestamp := time.Now().AddDate(0, 0, -40)
			db.DB.Exec(
				"UPDATE messages SET timestamp = ? WHERE user = ?",
				oldTimestamp, "user2",
			)

			err = db.CleanupTTL(db.DB, 30)
			t.Expect(err).ToBeNil()

			count, err := db.GetMessageCount(db.DB)
			t.Expect(err).ToBeNil()
			t.Expect(count).ToEqual(1)

			msg, err := db.GetMessageByID(db.DB, 1)
			t.Expect(err).ToBeNil()
			t.Expect(msg.User).ToEqual("user1")
		}).
		It("should handle full message lifecycle", func(t *gest.T) {
			tmpFile, err := os.CreateTemp("", "test_*.db")
			t.Expect(err).ToBeNil()
			defer os.Remove(tmpFile.Name())
			tmpFile.Close()

			err = db.Setup(tmpFile.Name())
			t.Expect(err).ToBeNil()
			defer db.DB.Close()

			err = db.InsertMessage(db.DB, "user1", "original content")
			t.Expect(err).ToBeNil()

			msg, err := db.GetMessageByID(db.DB, 1)
			t.Expect(err).ToBeNil()
			t.Expect(msg.Content).ToEqual("original content")

			err = db.UpdateMessage(db.DB, 1, "updated content", "preprocessed")
			t.Expect(err).ToBeNil()

			msg, err = db.GetMessageByID(db.DB, 1)
			t.Expect(err).ToBeNil()
			t.Expect(msg.Content).ToEqual("updated content")
			t.Expect(msg.PreprocessedContent).ToEqual("preprocessed")

			err = db.DeleteMessage(db.DB, 1)
			t.Expect(err).ToBeNil()

			_, err = db.GetMessageByID(db.DB, 1)
			t.Expect(err).Not().ToBeNil()
		}).
		It("should filter messages by user correctly", func(t *gest.T) {
			tmpFile, err := os.CreateTemp("", "test_*.db")
			t.Expect(err).ToBeNil()
			defer os.Remove(tmpFile.Name())
			tmpFile.Close()

			err = db.Setup(tmpFile.Name())
			t.Expect(err).ToBeNil()
			defer db.DB.Close()

			db.InsertMessage(db.DB, "alice", "msg1")
			db.InsertMessage(db.DB, "bob", "msg2")
			db.InsertMessage(db.DB, "alice", "msg3")

			messages, err := db.GetMessagesByUser(db.DB, "alice")
			t.Expect(err).ToBeNil()
			t.Expect(messages).ToHaveLength(2)

			messages, err = db.GetMessagesByUser(db.DB, "bob")
			t.Expect(err).ToBeNil()
			t.Expect(messages).ToHaveLength(1)
		}).
		Run(t)
}
