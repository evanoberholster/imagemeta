package xmp

import "errors"

// Errors
var (
	// ErrNoXMP is returned when no XMP Root Tag is found.
	ErrNoXMP = errors.New("Error XMP not found")
	// ErrPropertyNotSet
	ErrPropertyNotSet = errors.New("Error property not set")
)
