// Package meta contains meta types for image metadata
package meta

import "io"

// Reader that is compatible with imagemeta
type Reader interface {
	io.ReaderAt
	io.ReadSeeker
}
