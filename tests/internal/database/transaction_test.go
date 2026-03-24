package db_test

import (
	"KitsuneSemCalda/KitsuneBot/internal/db"
	"testing"

	"github.com/caiolandgraf/gest/v2/gest"
)

func TestTransactions(t *testing.T) {
	gest.Describe("Transaction tests").
		It("should rollback on error", func(t *gest.T) {
			err := db.Setup(":memory:")
			t.Expect(err).ToBeNil()
			defer db.DB.Close()

			tx, err := db.DB.Begin()
			t.Expect(err).ToBeNil()

			_, err = tx.Exec("INSERT INTO messages (user, content, preprocessed_content) VALUES (?, ?, ?)", "user1", "msg1", "")
			t.Expect(err).ToBeNil()

			err = tx.Rollback()
			t.Expect(err).ToBeNil()

			count, err := db.GetMessageCount(db.DB)
			t.Expect(err).ToBeNil()
			t.Expect(count).ToEqual(0)
		}).
		It("should commit successfully", func(t *gest.T) {
			err := db.Setup(":memory:")
			t.Expect(err).ToBeNil()
			defer db.DB.Close()

			tx, err := db.DB.Begin()
			t.Expect(err).ToBeNil()

			_, err = tx.Exec("INSERT INTO messages (user, content, preprocessed_content) VALUES (?, ?, ?)", "user1", "msg1", "")
			t.Expect(err).ToBeNil()

			err = tx.Commit()
			t.Expect(err).ToBeNil()

			count, err := db.GetMessageCount(db.DB)
			t.Expect(err).ToBeNil()
			t.Expect(count).ToEqual(1)
		}).
		Run(t)
}
