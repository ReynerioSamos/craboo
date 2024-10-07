package data

import (
	"time"
)

// each name begins with uppercase with so that they are exportable/public

type Comment struct {
	ID        int64     // unique value for each comment
	COntent   string    // the comment data
	Author    string    // the person who wrote the comment
	CreatedAt time.Time // database timestamp
	Version   int32     // incremented on each update
}
