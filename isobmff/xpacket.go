package isobmff

import "bytes"

const xpacketProbeLength = 512

var (
	xpacketPIStart = []byte("<?xpacket")
	xmpMetaStart   = []byte("<x:xmpmeta")
	rdfStart       = []byte("<rdf:RDF")
	utf8BOM        = []byte{0xEF, 0xBB, 0xBF}
)

// evaluateXPacketHeader inspects a small prefix of the payload for common XMP markers.
// It does not consume bytes from the source.
func evaluateXPacketHeader(b *box) (h XPacketHeader, err error) {
	payloadOffset := b.offset + int(b.size) - b.remain
	h.Offset = uint64(payloadOffset)
	if b.remain > int(^uint32(0)) {
		h.Length = ^uint32(0)
	} else {
		h.Length = uint32(b.remain)
	}
	if b.remain == 0 {
		return h, nil
	}

	probeLen := b.remain
	if probeLen > xpacketProbeLength {
		probeLen = xpacketProbeLength
	}
	buf, err := b.Peek(probeLen)
	if err != nil {
		return h, err
	}

	probe := bytes.TrimLeft(buf, "\x00\t\r\n ")
	if bytes.HasPrefix(probe, utf8BOM) {
		probe = probe[len(utf8BOM):]
	}

	h.HasXPacketPI = bytes.HasPrefix(probe, xpacketPIStart) || bytes.Contains(probe, xpacketPIStart)
	h.HasXMPMeta = bytes.Contains(probe, xmpMetaStart) || bytes.Contains(probe, rdfStart)
	return h, nil
}
