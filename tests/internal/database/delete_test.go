package db_test

import (
	"KitsuneSemCalda/KitsuneBot/internal/db"
	"testing"

	"github.com/caiolandgraf/gest/v2/gest"
)

func TestDeleteMessage(t *testing.T) {
	gest.Describe("DeleteMessage tests").
		It("should delete message by ID", func(t *gest.T) {
			err := db.Setup(":memory:")
			t.Expect(err).ToBeNil()
			defer db.DB.Close()

			db.InsertMessage(db.DB, "user1", "msg1")
			db.InsertMessage(db.DB, "user2", "msg2")

			err = db.DeleteMessage(db.DB, 1)
			t.Expect(err).ToBeNil()

			messages, err := db.GetMessages(db.DB)
			t.Expect(err).ToBeNil()
			t.Expect(messages).ToHaveLength(1)
			t.Expect(messages[0].ID).ToEqual(2)
		}).
		It("should not affect other messages", func(t *gest.T) {
			err := db.Setup(":memory:")
			t.Expect(err).ToBeNil()
			defer db.DB.Close()

			db.InsertMessage(db.DB, "user1", "msg1")
			db.InsertMessage(db.DB, "user2", "msg2")

			db.DeleteMessage(db.DB, 1)

			msg, err := db.GetMessageByID(db.DB, 2)
			t.Expect(err).ToBeNil()
			t.Expect(msg.Content).ToEqual("msg2")
		}).
		It("should return success for non-existent ID", func(t *gest.T) {
			err := db.Setup(":memory:")
			t.Expect(err).ToBeNil()
			defer db.DB.Close()

			err = db.DeleteMessage(db.DB, 999)
			t.Expect(err).ToBeNil()
		}).
		It("should delete all messages", func(t *gest.T) {
			err := db.Setup(":memory:")
			t.Expect(err).ToBeNil()
			defer db.DB.Close()

			db.InsertMessage(db.DB, "user1", "msg1")
			db.InsertMessage(db.DB, "user2", "msg2")

			db.DeleteMessage(db.DB, 1)
			db.DeleteMessage(db.DB, 2)

			messages, err := db.GetMessages(db.DB)
			t.Expect(err).ToBeNil()
			t.Expect(messages).ToHaveLength(0)
		}).
		Run(t)
}

func TestGetMessageCount(t *testing.T) {
	gest.Describe("GetMessageCount tests").
		It("should return 0 for empty database", func(t *gest.T) {
			err := db.Setup(":memory:")
			t.Expect(err).ToBeNil()
			defer db.DB.Close()

			count, err := db.GetMessageCount(db.DB)
			t.Expect(err).ToBeNil()
			t.Expect(count).ToEqual(0)
		}).
		It("should return correct count after inserts", func(t *gest.T) {
			err := db.Setup(":memory:")
			t.Expect(err).ToBeNil()
			defer db.DB.Close()

			db.InsertMessage(db.DB, "user1", "msg1")
			db.InsertMessage(db.DB, "user2", "msg2")
			db.InsertMessage(db.DB, "user3", "msg3")

			count, err := db.GetMessageCount(db.DB)
			t.Expect(err).ToBeNil()
			t.Expect(count).ToEqual(3)
		}).
		It("should update count after delete", func(t *gest.T) {
			err := db.Setup(":memory:")
			t.Expect(err).ToBeNil()
			defer db.DB.Close()

			db.InsertMessage(db.DB, "user1", "msg1")
			db.InsertMessage(db.DB, "user2", "msg2")
			db.DeleteMessage(db.DB, 1)

			count, err := db.GetMessageCount(db.DB)
			t.Expect(err).ToBeNil()
			t.Expect(count).ToEqual(1)
		}).
		Run(t)
}
