package imagemeta

type TiffMetadata struct {
	TiffHeader
}

func (tm TiffMetadata) Size() (width uint16, height uint16) {
	return 0, 0
}

func (tm TiffMetadata) Header() TiffHeader {
	return tm.TiffHeader
}

func (tm TiffMetadata) XMP() string {
	return ""
}