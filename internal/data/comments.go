package data

import (
	"time"

	"github.com/ReynerioSamos/craboo/internal/validator"
)

// each name begins with uppercase with so that they are exportable/public

// Make our JSON keys be displayed in all lowercase
// "-" means don't show this field
type Comment struct {
	ID        int64     `json:"id"`     // unique value for each comment
	Content   string    `json:"content` // the comment data
	Author    string    `json:"author`  // the person who wrote the comment
	CreatedAt time.Time `json:"-"`      // database timestamp
	Version   int32     `json:"version` // incremented on each update
}

func ValidateComment(v *validator.Validator, comment *Comment) {
	// check if content field is empty
	v.Check(comment.Content != "", "content", "must be provided")
	// check if author field is empty
	v.Check(comment.Author != "", "author", "must be provided")
	// check if the Content field is too long
	v.Check(len(comment.Content) <= 100, "content", "must not be more than 100 bytes long")
	//check if Author field is too long
	v.Check(len(comment.Author) <= 25, "author", "must not be more than 25 bytes long")
}
