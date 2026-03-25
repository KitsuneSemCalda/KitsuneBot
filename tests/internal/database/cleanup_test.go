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
		It("should delete exactly at boundary (7 days ago)", func(t *gest.T) {
			err := db.Setup(":memory:")
			t.Expect(err).ToBeNil()
			defer db.DB.Close()

			boundaryTime := time.Now().AddDate(0, 0, -7).Add(-1 * time.Second)
			db.DB.Exec(
				"INSERT INTO messages (user, content, preprocessed_content, timestamp) VALUES (?, ?, ?, ?)",
				"user1", "old", "", boundaryTime,
			)

			err = db.CleanupTTL(db.DB, 7)
			t.Expect(err).ToBeNil()

			count, err := db.GetMessageCount(db.DB)
			t.Expect(err).ToBeNil()
			t.Expect(count).ToEqual(0)
		}).
		It("should keep message exactly at TTL boundary", func(t *gest.T) {
			err := db.Setup(":memory:")
			t.Expect(err).ToBeNil()
			defer db.DB.Close()

			boundaryTime := time.Now().AddDate(0, 0, -7)
			db.DB.Exec(
				"INSERT INTO messages (user, content, preprocessed_content, timestamp) VALUES (?, ?, ?, ?)",
				"user1", "boundary", "", boundaryTime,
			)

			err = db.CleanupTTL(db.DB, 7)
			t.Expect(err).ToBeNil()

			count, err := db.GetMessageCount(db.DB)
			t.Expect(err).ToBeNil()
			t.Expect(count).ToEqual(1)
		}).
		It("should handle very large TTL (365 days)", func(t *gest.T) {
			err := db.Setup(":memory:")
			t.Expect(err).ToBeNil()
			defer db.DB.Close()

			oldTimestamp := time.Now().AddDate(0, 0, -100)
			db.DB.Exec(
				"INSERT INTO messages (user, content, preprocessed_content, timestamp) VALUES (?, ?, ?, ?)",
				"user1", "old", "", oldTimestamp,
			)

			err = db.CleanupTTL(db.DB, 365)
			t.Expect(err).ToBeNil()

			count, err := db.GetMessageCount(db.DB)
			t.Expect(err).ToBeNil()
			t.Expect(count).ToEqual(1)
		}).
		It("should delete all when TTL is negative", func(t *gest.T) {
			err := db.Setup(":memory:")
			t.Expect(err).ToBeNil()
			defer db.DB.Close()

			db.InsertMessage(db.DB, "user1", "msg1")

			err = db.CleanupTTL(db.DB, -1)
			t.Expect(err).ToBeNil()

			count, err := db.GetMessageCount(db.DB)
			t.Expect(err).ToBeNil()
			t.Expect(count).ToEqual(0)
		}).
		Run(t)
}
