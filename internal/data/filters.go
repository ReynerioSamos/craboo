package data

import (
	"github.com/ReynerioSamos/craboo/internal/validator"
)

// if an API does not pagination, it's not an API
// Also rate limiting, sorting, and graceful shutdown are the bare minimum of professional APIs
// The Filters type will contain the fields related to pagination
// and eventually the fields related to sorting
type Filters struct {
	Page     int // which page number does the client want
	PageSize int // how records per page
}

func ValidateFilters(v *validator.Validator, f Filters) {
	v.Check(f.Page > 0, "page", "must be greater than zero")
	v.Check(f.Page <= 500, "page", "must be a maximum of 500")
	v.Check(f.PageSize > 0, "page_size", "must be greater than zero")
	v.Check(f.PageSize <= 100, "page_size", "must be a maximum of 100")
}

// calculate how many records to send back
func (f Filters) limit() int {
	return f.PageSize
}

// calculate the offset so that we remember how many records have been sent
// and how many remail to be sent
func (f Filters) offset() int {
	return (f.Page - 1) * f.PageSize
}
