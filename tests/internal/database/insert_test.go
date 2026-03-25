package db_test

import (
	"KitsuneSemCalda/KitsuneBot/internal/db"
	"strings"
	"testing"

	"github.com/caiolandgraf/gest/v2/gest"
)

func TestInsertMessage(t *testing.T) {
	gest.Describe("InsertMessage tests").
		It("should insert message with all fields", func(t *gest.T) {
			err := db.Setup(":memory:")
			t.Expect(err).ToBeNil()
			defer db.DB.Close()

			err = db.InsertMessage(db.DB, "user1", "hello world")
			t.Expect(err).ToBeNil()

			messages, err := db.GetMessages(db.DB)
			t.Expect(err).ToBeNil()
			t.Expect(messages).ToHaveLength(1)
			t.Expect(messages[0].User).ToEqual("user1")
			t.Expect(messages[0].Content).ToEqual("hello world")
		}).
		It("should insert multiple messages", func(t *gest.T) {
			err := db.Setup(":memory:")
			t.Expect(err).ToBeNil()
			defer db.DB.Close()

			err = db.InsertMessage(db.DB, "user1", "msg1")
			t.Expect(err).ToBeNil()
			err = db.InsertMessage(db.DB, "user2", "msg2")
			t.Expect(err).ToBeNil()
			err = db.InsertMessage(db.DB, "user3", "msg3")
			t.Expect(err).ToBeNil()

			messages, err := db.GetMessages(db.DB)
			t.Expect(err).ToBeNil()
			t.Expect(messages).ToHaveLength(3)
		}).
		It("should insert message with empty content", func(t *gest.T) {
			err := db.Setup(":memory:")
			t.Expect(err).ToBeNil()
			defer db.DB.Close()

			err = db.InsertMessage(db.DB, "user1", "")
			t.Expect(err).ToBeNil()

			messages, err := db.GetMessages(db.DB)
			t.Expect(err).ToBeNil()
			t.Expect(messages).ToHaveLength(1)
			t.Expect(messages[0].Content).ToEqual("")
		}).
		It("should insert message with special characters", func(t *gest.T) {
			err := db.Setup(":memory:")
			t.Expect(err).ToBeNil()
			defer db.DB.Close()

			special := "Hello 🦊 + emoji & symbols <script>alert('xss')</script>"
			err = db.InsertMessage(db.DB, "user_1", special)
			t.Expect(err).ToBeNil()

			msg, err := db.GetMessageByID(db.DB, 1)
			t.Expect(err).ToBeNil()
			t.Expect(msg.Content).ToEqual(special)
		}).
		It("should insert message with unicode characters", func(t *gest.T) {
			err := db.Setup(":memory:")
			t.Expect(err).ToBeNil()
			defer db.DB.Close()

			unicode := "こんにちは世界 مرحبا بالعالم 🎉"
			err = db.InsertMessage(db.DB, "user1", unicode)
			t.Expect(err).ToBeNil()

			msg, err := db.GetMessageByID(db.DB, 1)
			t.Expect(err).ToBeNil()
			t.Expect(msg.Content).ToEqual(unicode)
		}).
		It("should insert message with very long content", func(t *gest.T) {
			err := db.Setup(":memory:")
			t.Expect(err).ToBeNil()
			defer db.DB.Close()

			longContent := strings.Repeat("a", 10000)
			err = db.InsertMessage(db.DB, "user1", longContent)
			t.Expect(err).ToBeNil()

			msg, err := db.GetMessageByID(db.DB, 1)
			t.Expect(err).ToBeNil()
			t.Expect(len(msg.Content)).ToEqual(10000)
		}).
		It("should allow empty user (SQLite behavior)", func(t *gest.T) {
			err := db.Setup(":memory:")
			t.Expect(err).ToBeNil()
			defer db.DB.Close()

			err = db.InsertMessage(db.DB, "", "content")
			t.Expect(err).ToBeNil()

			messages, err := db.GetMessages(db.DB)
			t.Expect(err).ToBeNil()
			t.Expect(messages[0].User).ToEqual("")
		}).
		It("should handle user with special characters", func(t *gest.T) {
			err := db.Setup(":memory:")
			t.Expect(err).ToBeNil()
			defer db.DB.Close()

			err = db.InsertMessage(db.DB, "user'db", "content")
			t.Expect(err).ToBeNil()

			msg, err := db.GetMessageByID(db.DB, 1)
			t.Expect(err).ToBeNil()
			t.Expect(msg.User).ToEqual("user'db")
		}).
		It("should store preprocessed_content as empty by default", func(t *gest.T) {
			err := db.Setup(":memory:")
			t.Expect(err).ToBeNil()
			defer db.DB.Close()

			err = db.InsertMessage(db.DB, "user1", "hello")
			t.Expect(err).ToBeNil()

			msg, err := db.GetMessageByID(db.DB, 1)
			t.Expect(err).ToBeNil()
			t.Expect(msg.PreprocessedContent).ToEqual("")
		}).
		Run(t)
}
