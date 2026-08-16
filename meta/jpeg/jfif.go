package jpeg

// JFIF stores fields from an APP0 JFIF segment.
type JFIF struct {
	MajorVersion   uint8
	MinorVersion   uint8
	ResolutionUnit uint8
	XResolution    uint16
	YResolution    uint16
}

func parseJFIF(payload []byte) (*JFIF, error) {
	if len(payload) < 14 || !isJFIFPayload(payload) {
		return nil, errShortSegment("JFIF")
	}
	return &JFIF{
		MajorVersion:   payload[5],
		MinorVersion:   payload[6],
		ResolutionUnit: payload[7],
		XResolution:    jpegEndian.Uint16(payload[8:10]),
		YResolution:    jpegEndian.Uint16(payload[10:12]),
	}, nil
}
