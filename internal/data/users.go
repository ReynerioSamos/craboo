package data

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"time"

	"github.com/ReynerioSamos/craboo/internal/validator"
)

type User struct {
	ID        int64     `json:"id"`       // unique value for each user
	Email     string    `json:"email"`    // the email of user
	Fullname  string    `json:"fullname"` // the full name of user
	CreatedAt time.Time `json:"-"`        // database timestamp
}

func ValidateUser(v *validator.Validator, user *User) {
	// check if email field is empty
	v.Check(user.Email != "", "email", "must be provided")

	// check if fullname field is empty
	v.Check(user.Fullname != "", "fullname", "must be provided")

	// check if the email field is too long
	v.Check(len(user.Email) <= 254, "email", "must not be more than 254 bytes long")

	// validation for email formatt using regex
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	v.Check(emailRegex.MatchString(user.Email), "email", "must be a valid email address")

	//check if fullname field is too long
	v.Check(len(user.Fullname) <= 50, "fullname", "must not be more than 50 bytes long")
}

// A UserModel expects a connection pool
type UserModel struct {
	DB *sql.DB
}

// Insert a new row in the users table
// Expects a pointer to the actual user
func (u UserModel) Insert(user *User) error {
	// the SQL query to be executed against the database table
	query := `
		INSERT INTO users (email, fullname)
		VALUES ($1, $2)
		RETURNING id, created_at
		`
	// the actual values to replace $1, and $2
	args := []any{user.Email, user.Fullname}

	// Create a context with a 3-second timeout. No database
	// operation should take more than 3 seconds or we will quit it
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// execute the query against the users database table. We ask for the
	// id, created_at, and the version to be sent back to us which we will use
	// to update the user struct later on

	return u.DB.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.CreatedAt)
}

// Get a specific user from the users table
func (u UserModel) Get(id int64) (*User, error) {
	// check if the id is valid
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	// the SQL query to be executed against the database table
	query := `
		SELECT id, created_at, email, fullname
		FROM users
		WHERE id = $1
		`
	// declare a variable of type user to store the returned user
	var user User

	// Set a 3-second context/time
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := u.DB.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Email,
		&user.Fullname,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &user, nil
}

func (u UserModel) Update(user *User) error {
	// The SQL query to be executed against the database table
	query := `
		UPDATE users
		SET email = $1, fullname = $2
		WHERE id = $3
		RETURNING email, fullname
		`

	args := []any{user.Email, user.Fullname, user.ID}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return u.DB.QueryRowContext(ctx, query, args...).Scan(&user.Email, &user.Fullname)
}

func (u UserModel) Delete(id int64) error {
	// check if the id is valid
	if id < 1 {
		return ErrRecordNotFound
	}
	// the SQL query to be executied against the database table
	query := `
		DELETE FROM users
		WHERE id = $1	
		`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// ExecContext does not return any rows unlike QueryRowContext.
	// It only returns information about the query execution
	// such as how many rows were affected
	result, err := u.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	// Were any rows deleted?
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	// Probably a wrong id was provided or the client is trying to delete an already deleted user
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}
