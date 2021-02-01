package xmp

import "errors"

// Errors
var (
	// ErrNoXMP is returned when no XMP Root Tag is found.
	ErrNoXMP = errors.New("error XMP not found")
	// ErrPropertyNotSet
	ErrPropertyNotSet = errors.New("error property not set")
)
