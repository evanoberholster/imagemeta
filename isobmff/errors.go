package isobmff

import "errors"

// Errors
var (
	ErrBrandNotSupported        = errors.New("error brand not supported")
	ErrBufLength                = errors.New("insufficient buffer length")
	ErrItemTypeWS               = errors.New("itemType doesn't end on whitespace")
	ErrRemainLengthInsufficient = errors.New("remain length insufficient")
	ErrFlagsLength              = errors.New("failed to read 4 bytes of Flags")
	ErrItemTypeLength           = errors.New("insufficient itemType Length")
	errLargeBox                 = errors.New("unexpectedly large box")
	errUintSize                 = errors.New("invalid uintn read size")
	ErrWrongBoxType             = errors.New("error wrong box type")

	// ErrInfeVersionNotSupported is returned when an infe box with an unsupported was found.
	ErrInfeVersionNotSupported = errors.New("infe box version not supported")
)
