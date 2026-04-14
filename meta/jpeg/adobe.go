package jpeg

// Adobe stores fields from the APP14 Adobe DCT segment.
type Adobe struct {
	DCTEncodeVersion uint16
	APP14Flags0      uint16
	APP14Flags1      uint16
	ColorTransform   uint8
}

func parseAdobe(payload []byte) (*Adobe, error) {
	if len(payload) < 12 || !isAdobePayload(payload) {
		return nil, errShortSegment("Adobe")
	}
	return &Adobe{
		DCTEncodeVersion: jpegEndian.Uint16(payload[5:7]),
		APP14Flags0:      jpegEndian.Uint16(payload[7:9]),
		APP14Flags1:      jpegEndian.Uint16(payload[9:11]),
		ColorTransform:   payload[11],
	}, nil
}
