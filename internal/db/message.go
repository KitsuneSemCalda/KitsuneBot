package db

import (
	"time"
)

type Message struct {
	ID                  int
	User                string
	Content             string
	PreprocessedContent string
	Timesamp            time.Time
}
