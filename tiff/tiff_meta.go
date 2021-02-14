package tiff

type TiffMetadata struct {
	header Header
}

func (tm TiffMetadata) Size() (width uint16, height uint16) {
	return 0, 0
}

func (tm TiffMetadata) Header() Header {
	return tm.header
}

func (tm TiffMetadata) XMP() string {
	return ""
}
