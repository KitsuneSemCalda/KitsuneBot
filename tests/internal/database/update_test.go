package db_test

import (
	"KitsuneSemCalda/KitsuneBot/internal/db"
	"testing"

	"github.com/caiolandgraf/gest/v2/gest"
)

func TestUpdateMessage(t *testing.T) {
	gest.Describe("UpdateMessage tests").
		It("should update message content", func(t *gest.T) {
			err := db.Setup(":memory:")
			t.Expect(err).ToBeNil()
			defer db.DB.Close()

			db.InsertMessage(db.DB, "user1", "original content")

			err = db.UpdateMessage(db.DB, 1, "updated content", "updated_preprocessed")
			t.Expect(err).ToBeNil()

			msg, err := db.GetMessageByID(db.DB, 1)
			t.Expect(err).ToBeNil()
			t.Expect(msg.Content).ToEqual("updated content")
			t.Expect(msg.PreprocessedContent).ToEqual("updated_preprocessed")
		}).
		It("should succeed even for non-existent ID (no rows affected)", func(t *gest.T) {
			err := db.Setup(":memory:")
			t.Expect(err).ToBeNil()
			defer db.DB.Close()

			err = db.UpdateMessage(db.DB, 999, "content", "preprocessed")
			t.Expect(err).ToBeNil()
		}).
		It("should update to empty content", func(t *gest.T) {
			err := db.Setup(":memory:")
			t.Expect(err).ToBeNil()
			defer db.DB.Close()

			db.InsertMessage(db.DB, "user1", "original")

			err = db.UpdateMessage(db.DB, 1, "", "")
			t.Expect(err).ToBeNil()

			msg, err := db.GetMessageByID(db.DB, 1)
			t.Expect(err).ToBeNil()
			t.Expect(msg.Content).ToEqual("")
		}).
		Run(t)
}
