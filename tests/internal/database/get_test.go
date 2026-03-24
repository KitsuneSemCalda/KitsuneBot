package db_test

import (
	"KitsuneSemCalda/KitsuneBot/internal/db"
	"testing"
	"time"

	"github.com/caiolandgraf/gest/v2/gest"
)

func TestGetMessages(t *testing.T) {
	gest.Describe("GetMessages tests").
		It("should retrieve all messages", func(t *gest.T) {
			err := db.Setup(":memory:")
			t.Expect(err).ToBeNil()
			defer db.DB.Close()

			db.InsertMessage(db.DB, "user1", "msg1")
			db.InsertMessage(db.DB, "user2", "msg2")
			db.InsertMessage(db.DB, "user3", "msg3")

			messages, err := db.GetMessages(db.DB)
			t.Expect(err).ToBeNil()
			t.Expect(messages).ToHaveLength(3)
		}).
		It("should return messages ordered by timestamp DESC", func(t *gest.T) {
			err := db.Setup(":memory:")
			t.Expect(err).ToBeNil()
			defer db.DB.Close()

			ts1 := time.Now().Add(-2 * time.Hour)
			ts2 := time.Now().Add(-1 * time.Hour)
			ts3 := time.Now()

			db.DB.Exec("INSERT INTO messages (user, content, preprocessed_content, timestamp) VALUES (?, ?, ?, ?)", "user1", "msg1", "", ts1)
			db.DB.Exec("INSERT INTO messages (user, content, preprocessed_content, timestamp) VALUES (?, ?, ?, ?)", "user2", "msg2", "", ts2)
			db.DB.Exec("INSERT INTO messages (user, content, preprocessed_content, timestamp) VALUES (?, ?, ?, ?)", "user3", "msg3", "", ts3)

			messages, err := db.GetMessages(db.DB)
			t.Expect(err).ToBeNil()
			t.Expect(messages[0].User).ToEqual("user3")
			t.Expect(messages[2].User).ToEqual("user1")
		}).
		It("should return empty slice when no messages", func(t *gest.T) {
			err := db.Setup(":memory:")
			t.Expect(err).ToBeNil()
			defer db.DB.Close()

			messages, err := db.GetMessages(db.DB)
			t.Expect(err).ToBeNil()
			t.Expect(messages).ToHaveLength(0)
		}).
		It("should handle query error when db is closed", func(t *gest.T) {
			err := db.Setup(":memory:")
			t.Expect(err).ToBeNil()

			db.DB.Close()

			_, err = db.GetMessages(db.DB)
			t.Expect(err).Not().ToBeNil()
		}).
		Run(t)
}

func TestGetMessageByID(t *testing.T) {
	gest.Describe("GetMessageByID tests").
		It("should retrieve message by ID", func(t *gest.T) {
			err := db.Setup(":memory:")
			t.Expect(err).ToBeNil()
			defer db.DB.Close()

			db.InsertMessage(db.DB, "user1", "test message")

			msg, err := db.GetMessageByID(db.DB, 1)
			t.Expect(err).ToBeNil()
			t.Expect(msg.ID).ToEqual(1)
			t.Expect(msg.User).ToEqual("user1")
			t.Expect(msg.Content).ToEqual("test message")
		}).
		It("should return error for non-existent ID", func(t *gest.T) {
			err := db.Setup(":memory:")
			t.Expect(err).ToBeNil()
			defer db.DB.Close()

			_, err = db.GetMessageByID(db.DB, 999)
			t.Expect(err).Not().ToBeNil()
		}).
		It("should return error when db is closed", func(t *gest.T) {
			err := db.Setup(":memory:")
			t.Expect(err).ToBeNil()

			db.DB.Close()

			_, err = db.GetMessageByID(db.DB, 1)
			t.Expect(err).Not().ToBeNil()
		}).
		Run(t)
}

func TestGetMessagesByUser(t *testing.T) {
	gest.Describe("GetMessagesByUser tests").
		It("should retrieve messages by user", func(t *gest.T) {
			err := db.Setup(":memory:")
			t.Expect(err).ToBeNil()
			defer db.DB.Close()

			db.InsertMessage(db.DB, "user1", "msg1")
			db.InsertMessage(db.DB, "user1", "msg2")
			db.InsertMessage(db.DB, "user2", "msg3")

			messages, err := db.GetMessagesByUser(db.DB, "user1")
			t.Expect(err).ToBeNil()
			t.Expect(messages).ToHaveLength(2)
		}).
		It("should return empty for user with no messages", func(t *gest.T) {
			err := db.Setup(":memory:")
			t.Expect(err).ToBeNil()
			defer db.DB.Close()

			db.InsertMessage(db.DB, "user1", "msg1")

			messages, err := db.GetMessagesByUser(db.DB, "nonexistent")
			t.Expect(err).ToBeNil()
			t.Expect(messages).ToHaveLength(0)
		}).
		It("should return messages ordered by timestamp DESC", func(t *gest.T) {
			err := db.Setup(":memory:")
			t.Expect(err).ToBeNil()
			defer db.DB.Close()

			ts1 := time.Now().Add(-1 * time.Hour)
			ts2 := time.Now()

			db.DB.Exec("INSERT INTO messages (user, content, preprocessed_content, timestamp) VALUES (?, ?, ?, ?)", "user1", "msg1", "", ts1)
			db.DB.Exec("INSERT INTO messages (user, content, preprocessed_content, timestamp) VALUES (?, ?, ?, ?)", "user1", "msg2", "", ts2)

			messages, err := db.GetMessagesByUser(db.DB, "user1")
			t.Expect(err).ToBeNil()
			t.Expect(messages[0].Content).ToEqual("msg2")
			t.Expect(messages[1].Content).ToEqual("msg1")
		}).
		Run(t)
}
