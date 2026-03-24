package db

import (
	"time"
)

type Message struct {
	ID                  int
	User                string
	Content             string
	PreprocessedContent string
	Timestamp           time.Time
}
