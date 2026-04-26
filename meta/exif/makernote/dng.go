package makernote

import "github.com/evanoberholster/imagemeta/meta/exif/tag"

// DNG groups Adobe DNG extension tags and Adobe private data records.
type DNG struct {
	DNGVersion         [8]byte
	DNGBackwardVersion [8]byte

	CameraModel         string
	OriginalRawFileName string
	ProfileName         string

	DNGVersionCount         uint8
	DNGBackwardVersionCount uint8

	BestQualityScale tag.RationalU
	AdobeData        DNGAdobeData
}

// DNGAdobeData stores selected information from IFD0 tag 0xc634
// (DNGAdobeData / Adobe private data).
//
// ExifTool parses this as an "Adobe\0" record stream. We currently model the
// overall record count plus the Adobe-mutated maker-note record details needed
// to rebase and parse MakN data.
type DNGAdobeData struct {
	RecordCount             uint8
	MakerNoteOriginalOffset uint32
	MakerNoteRecordLength   uint32
}
