package xml

import "errors"

// Errors
var (
	// ErrNoXMP is returned when no XMP Root Tag is found.
	ErrNoXMP = errors.New("Error XMP not found")
	// ErrCloseTag is returned when a tag is closed.
	ErrCloseTag = errors.New("Error close tag")
)
