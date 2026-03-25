package db_test

import (
	"KitsuneSemCalda/KitsuneBot/internal/db"
	"os"
	"sync"
	"testing"

	"github.com/caiolandgraf/gest/v2/gest"
)

func TestConcurrentAccess(t *testing.T) {
	gest.Describe("Concurrent access tests").
		It("should handle concurrent inserts", func(t *gest.T) {
			tmpFile, err := os.CreateTemp("", "concurrent_*.db")
			t.Expect(err).ToBeNil()
			defer os.Remove(tmpFile.Name())
			tmpFile.Close()

			testDB, err := db.InitDB(tmpFile.Name())
			t.Expect(err).ToBeNil()
			defer testDB.Close()

			var wg sync.WaitGroup
			for i := 0; i < 10; i++ {
				wg.Add(1)
				go func(n int) {
					defer wg.Done()
					db.InsertMessage(testDB, "user", "message")
				}(i)
			}
			wg.Wait()

			count, err := db.GetMessageCount(testDB)
			t.Expect(err).ToBeNil()
			t.Expect(count).ToEqual(10)
		}).
		It("should handle concurrent reads and writes", func(t *gest.T) {
			tmpFile, err := os.CreateTemp("", "concurrent_*.db")
			t.Expect(err).ToBeNil()
			defer os.Remove(tmpFile.Name())
			tmpFile.Close()

			testDB, err := db.InitDB(tmpFile.Name())
			t.Expect(err).ToBeNil()
			defer testDB.Close()

			db.InsertMessage(testDB, "user1", "initial")

			var wg sync.WaitGroup
			errors := make(chan error, 20)

			for i := 0; i < 10; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					_, err := db.GetMessages(testDB)
					if err != nil {
						errors <- err
					}
				}()

				wg.Add(1)
				go func(n int) {
					defer wg.Done()
					err := db.InsertMessage(testDB, "user", "msg")
					if err != nil {
						errors <- err
					}
				}(i)
			}

			wg.Wait()
			close(errors)

			for err := range errors {
				t.Expect(err).ToBeNil()
			}

			count, err := db.GetMessageCount(testDB)
			t.Expect(err).ToBeNil()
			t.Expect(count).ToEqual(11)
		}).
		It("should handle concurrent deletes", func(t *gest.T) {
			tmpFile, err := os.CreateTemp("", "concurrent_*.db")
			t.Expect(err).ToBeNil()
			defer os.Remove(tmpFile.Name())
			tmpFile.Close()

			testDB, err := db.InitDB(tmpFile.Name())
			t.Expect(err).ToBeNil()
			defer testDB.Close()

			for i := 0; i < 5; i++ {
				db.InsertMessage(testDB, "user", "msg")
			}

			var wg sync.WaitGroup
			for i := 1; i <= 5; i++ {
				wg.Add(1)
				go func(id int) {
					defer wg.Done()
					db.DeleteMessage(testDB, id)
				}(i)
			}
			wg.Wait()

			count, err := db.GetMessageCount(testDB)
			t.Expect(err).ToBeNil()
			t.Expect(count).ToEqual(0)
		}).
		It("should handle concurrent updates", func(t *gest.T) {
			tmpFile, err := os.CreateTemp("", "concurrent_*.db")
			t.Expect(err).ToBeNil()
			defer os.Remove(tmpFile.Name())
			tmpFile.Close()

			testDB, err := db.InitDB(tmpFile.Name())
			t.Expect(err).ToBeNil()
			defer testDB.Close()

			db.InsertMessage(testDB, "user1", "original")

			var wg sync.WaitGroup
			for i := 0; i < 10; i++ {
				wg.Add(1)
				go func(n int) {
					defer wg.Done()
					db.UpdateMessage(testDB, 1, "updated", "")
				}(i)
			}
			wg.Wait()

			msg, err := db.GetMessageByID(testDB, 1)
			t.Expect(err).ToBeNil()
			t.Expect(msg.Content).ToEqual("updated")
		}).
		Run(t)
}
