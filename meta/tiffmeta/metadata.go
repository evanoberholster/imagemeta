package tiffmeta

type Metadata struct {
	Header
}

func (m Metadata) Size() (width uint16, height uint16) {
	return 0, 0
}

func (m Metadata) TiffHeader() Header {
	return m.Header
}

func (m Metadata) XML() string {
	return ""
}
