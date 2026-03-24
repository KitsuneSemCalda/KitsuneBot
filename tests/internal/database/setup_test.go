package db_test

import (
	"KitsuneSemCalda/KitsuneBot/internal/db"
	"testing"

	"github.com/caiolandgraf/gest/v2/gest"
)

func TestInitDB(t *testing.T) {
	gest.Describe("InitDB tests").
		It("should create database with correct schema", func(t *gest.T) {
			testDB, err := db.InitDB(":memory:")
			t.Expect(err).ToBeNil()
			defer testDB.Close()

			row := testDB.QueryRow("SELECT sql FROM sqlite_master WHERE type='table' AND name='messages'")
			var schema string
			err = row.Scan(&schema)
			t.Expect(err).ToBeNil()
			t.Expect(schema).ToContain("id INTEGER PRIMARY KEY")
			t.Expect(schema).ToContain("user TEXT NOT NULL")
			t.Expect(schema).ToContain("content TEXT NOT NULL")
			t.Expect(schema).ToContain("preprocessed_content")
			t.Expect(schema).ToContain("timestamp")
		}).
		It("should use CREATE TABLE IF NOT EXISTS", func(t *gest.T) {
			testDB, err := db.InitDB(":memory:")
			t.Expect(err).ToBeNil()

			err = testDB.Close()
			t.Expect(err).ToBeNil()

			testDB2, err := db.InitDB(":memory:")
			t.Expect(err).ToBeNil()
			defer testDB2.Close()

			count, err := db.GetMessageCount(testDB2)
			t.Expect(err).ToBeNil()
			t.Expect(count).ToEqual(0)
		}).
		Run(t)
}

func TestSetup(t *testing.T) {
	gest.Describe("Setup tests").
		It("should setup database successfully", func(t *gest.T) {
			err := db.Setup(":memory:")
			t.Expect(err).ToBeNil()
			t.Expect(db.DB).Not().ToBeNil()
			db.DB.Close()
		}).
		Run(t)
}
