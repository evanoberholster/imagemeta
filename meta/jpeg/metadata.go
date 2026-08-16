package jpeg

import "errors"

// Metadata contains metadata stored directly in JPEG marker segments.
type Metadata struct {
	SOF       SOF
	JFIF      *JFIF
	CIFF      *CIFF
	MPF       *MPF
	ICC       *ICCProfile
	Photoshop *Photoshop
	IPTC      *IPTC
	Adobe     *Adobe

	iccChunks map[uint8][]byte
	iccTotal  uint8
}

// SOF stores primary image dimensions from a JPEG Start Of Frame marker.
type SOF struct {
	Marker          string
	EncodingProcess uint8
	BitsPerSample   uint8
	Width           uint16
	Height          uint16
	ColorComponents uint8
}

func (m *Metadata) finish() error {
	if m == nil {
		return nil
	}
	if len(m.iccChunks) > 0 {
		return m.finishICC()
	}
	return nil
}

func errShortSegment(name string) error {
	return errors.New("jpeg: short " + name + " segment")
}
