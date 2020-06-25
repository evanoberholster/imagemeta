package jpegmeta

import (
	"errors"
	"io"
)

// Errors
var (
	ErrNoJPEGMarker = errors.New("No JPEG Marker")
	ErrEndOfImage   = errors.New("End of Image")
)

// Scan -
// TODO: Write tests
func Scan(reader io.Reader) (m Metadata, err error) {
	defer func() {
		if state := recover(); state != nil {
			err = state.(error)
		}
	}()
	m = newMetadata(reader)
	var buf []byte
	for {
		if buf, err = m.br.Peek(16); err != nil {
			if err == io.EOF {
				err = ErrNoJPEGMarker
				return
			}
			panic(err)
		}

		if !isMarkerFirstByte(buf) {
			if err = m.discard(1); err != nil {
				panic(err)
			}
			continue
		}

		if err := m.scanMarkers(buf); err == nil {
			continue
		}

		break
	}
	return
}
