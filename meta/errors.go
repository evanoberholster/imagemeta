package meta

import "errors"

// Common Errors
var (
	// ErrInvalidHeader is an error for an Invalid Exif TiffHeader
	ErrInvalidHeader = errors.New("error TiffHeader is not valid")

	// ErrNoExif is an error for when no exif is found
	ErrNoExif = errors.New("error no Exif")

	// ErrBufLength
	ErrBufLength = errors.New("error buffer length insufficient")
)
