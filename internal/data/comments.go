package data

import (
	"context"
	"database/sql"
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

// A CommentModel expects a connection pool
type CommentModel struct {
	DB *sql.DB
}

// Insert a new row in the commetns table
// Expects a pointer to the actual
func (c CommentModel) Insert(comment *Comment) error {
	// the SQL query to be executed against the database table
	query := `
		INSERT INTO comments (content, author)
		VALUES ($1, $2)
		RETURNING id, created_at, version
		`
	// the actual values to replace $1, and $2
	args := []any{comment.Content, comment.Author}

	// Create a context with a 3-second timeout. No database
	// operation should take more than 3 seconds or we will quit it
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// executre the query against the comments database table. We ask for the
	// id, created_at, and the version to be sent back to us which we will use
	// to update the Comment struct later on

	return c.DB.QueryRowContext(ctx, query, args...).Scan(&comment.ID, &comment.CreatedAt, &comment.Version)
}
