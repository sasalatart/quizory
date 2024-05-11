package pagination

import "errors"

// Pagination contains the parameters needed for paginating a list of items.
type Pagination struct {
	Page     int
	PageSize int
}

// New creates a Pagination instance with the provided page and page size values. It defaults to
// page 0 and a page size of 25 if the provided values are nil.
func New(page, pageSize *int) Pagination {
	p := Pagination{Page: 0, PageSize: 25}
	if page != nil {
		p.Page = *page
	}
	if pageSize != nil {
		p.PageSize = *pageSize
	}
	return p
}

// Validate checks if a Pagination is valid, i.e. if its page is non-negative and its page size is
// positive and greater than 0.
func (p Pagination) Validate() error {
	if p.Page < 0 {
		return errors.New("page must be a non-negative integer")
	}
	if p.PageSize < 1 {
		return errors.New("page size must be a positive integer greater than 0")
	}
	return nil
}

// Offset returns the offset value that should be used when querying a database.
func (p Pagination) Offset() int {
	return p.Page * p.PageSize
}
