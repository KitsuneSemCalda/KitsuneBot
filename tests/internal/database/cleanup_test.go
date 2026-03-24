package db_test

import (
	"KitsuneSemCalda/KitsuneBot/internal/db"
	"testing"
	"time"

	"github.com/caiolandgraf/gest/v2/gest"
)

func TestCleanupTTL(t *testing.T) {
	gest.Describe("CleanupTTL tests").
		It("should delete messages older than TTL", func(t *gest.T) {
			err := db.Setup(":memory:")
			t.Expect(err).ToBeNil()
			defer db.DB.Close()

			oldTimestamp := time.Now().AddDate(0, 0, -10)
			db.DB.Exec(
				"INSERT INTO messages (user, content, preprocessed_content, timestamp) VALUES (?, ?, ?, ?)",
				"user1", "old message", "", oldTimestamp,
			)

			err = db.CleanupTTL(db.DB, 7)
			t.Expect(err).ToBeNil()

			count, err := db.GetMessageCount(db.DB)
			t.Expect(err).ToBeNil()
			t.Expect(count).ToEqual(0)
		}).
		It("should keep messages newer than TTL", func(t *gest.T) {
			err := db.Setup(":memory:")
			t.Expect(err).ToBeNil()
			defer db.DB.Close()

			db.InsertMessage(db.DB, "user1", "recent message")

			err = db.CleanupTTL(db.DB, 7)
			t.Expect(err).ToBeNil()

			count, err := db.GetMessageCount(db.DB)
			t.Expect(err).ToBeNil()
			t.Expect(count).ToEqual(1)
		}).
		It("should keep messages when TTL is 0 (only deletes past messages)", func(t *gest.T) {
			err := db.Setup(":memory:")
			t.Expect(err).ToBeNil()
			defer db.DB.Close()

			db.InsertMessage(db.DB, "user1", "msg1")
			db.InsertMessage(db.DB, "user2", "msg2")

			err = db.CleanupTTL(db.DB, 0)
			t.Expect(err).ToBeNil()

			count, err := db.GetMessageCount(db.DB)
			t.Expect(err).ToBeNil()
			t.Expect(count).ToEqual(2)
		}).
		It("should keep all when no messages are old enough", func(t *gest.T) {
			err := db.Setup(":memory:")
			t.Expect(err).ToBeNil()
			defer db.DB.Close()

			db.InsertMessage(db.DB, "user1", "msg1")

			err = db.CleanupTTL(db.DB, 30)
			t.Expect(err).ToBeNil()

			count, err := db.GetMessageCount(db.DB)
			t.Expect(err).ToBeNil()
			t.Expect(count).ToEqual(1)
		}).
		It("should handle empty database", func(t *gest.T) {
			err := db.Setup(":memory:")
			t.Expect(err).ToBeNil()
			defer db.DB.Close()

			err = db.CleanupTTL(db.DB, 7)
			t.Expect(err).ToBeNil()

			count, err := db.GetMessageCount(db.DB)
			t.Expect(err).ToBeNil()
			t.Expect(count).ToEqual(0)
		}).
		Run(t)
}
