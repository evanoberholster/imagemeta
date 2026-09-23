package imagetype

import (
	"bytes"
	"encoding/binary"
	"testing"
)

// buildTiffDeep assembles a classic or BigTIFF header with a first IFD
// holding Make/Model ASCII tags (via string pool) and an optional inline
// DNGVersion tag. end is binary.LittleEndian or binary.BigEndian.
func buildTiffDeep(end binary.ByteOrder, big bool, manufacturer, model string, dngVersion bool) []byte {
	buf := make([]byte, 4096)
	put16 := func(off int, v uint16) { end.PutUint16(buf[off:], v) }
	put32 := func(off int, v uint32) { end.PutUint32(buf[off:], v) }
	put64 := func(off int, v uint64) { end.PutUint64(buf[off:], v) }

	entryCount := 2
	if dngVersion {
		entryCount = 3
	}
	var entriesAt, pool int
	if big {
		if end == binary.LittleEndian {
			copy(buf, []byte{0x49, 0x49, 0x2B, 0x00, 0x08, 0x00, 0x00, 0x00})
		} else {
			copy(buf, []byte{0x4D, 0x4D, 0x00, 0x2B, 0x00, 0x08, 0x00, 0x00})
		}
		put64(8, 16)
		put64(16, uint64(entryCount))
		entriesAt, pool = 24, 24+entryCount*20
	} else {
		if end == binary.LittleEndian {
			copy(buf, []byte{0x49, 0x49, 0x2A, 0x00})
		} else {
			copy(buf, []byte{0x4D, 0x4D, 0x00, 0x2A})
		}
		put32(4, 8)
		put16(8, uint16(entryCount))
		entriesAt, pool = 10, 10+entryCount*12
	}

	addASCII := func(idx int, tag uint16, s string) {
		// Mirror real writers: values that fit live inline, longer ones
		// go to the string pool. This exercises both reader paths.
		var e, inlineCap int
		if big {
			e, inlineCap = entriesAt+idx*20, 8
		} else {
			e, inlineCap = entriesAt+idx*12, 4
		}
		put16(e, tag)
		put16(e+2, tiffTypeASCII)
		if big {
			put64(e+4, uint64(len(s)+1))
		} else {
			put32(e+4, uint32(len(s)+1))
		}
		if len(s)+1 <= inlineCap {
			var v int
			if big {
				v = e + 12
			} else {
				v = e + 8
			}
			copy(buf[v:], s)
			return
		}
		if big {
			put64(e+12, uint64(pool))
		} else {
			put32(e+8, uint32(pool))
		}
		copy(buf[pool:], s)
		buf[pool+len(s)] = 0
		pool += len(s) + 1
	}
	addASCII(0, tiffTagMake, manufacturer)
	addASCII(1, tiffTagModel, model)
	if dngVersion {
		var e int
		if big {
			e = entriesAt + 2*20
			put16(e, tiffTagDNGVersion)
			put16(e+2, 7)
			put64(e+4, 4)
			copy(buf[e+12:], []byte{1, 7, 0, 0})
		} else {
			e = entriesAt + 2*12
			put16(e, tiffTagDNGVersion)
			put16(e+2, 7)
			put32(e+4, 4)
			copy(buf[e+8:], []byte{1, 7, 0, 0})
		}
	}
	return buf
}

func TestTiffDeepSubtype(t *testing.T) {
	t.Parallel()
	le, be := binary.LittleEndian, binary.BigEndian
	cases := []struct {
		name     string
		buf      []byte
		expected FileType
	}{
		{name: "Olympus", buf: buildTiffDeep(le, false, "OLYMPUS IMAGING CORP.", "E-M5MarkIII", false), expected: ImageORF},
		{name: "Sigma", buf: buildTiffDeep(le, false, "SIGMA", "fp", false), expected: ImageX3F},
		{name: "Pentax", buf: buildTiffDeep(le, false, "PENTAX", "K-3 Mark III", false), expected: ImagePEF},
		{name: "Ricoh", buf: buildTiffDeep(le, false, "RICOH IMAGING COMPANY, LTD.", "PENTAX K-1", false), expected: ImagePEF},
		{name: "Samsung", buf: buildTiffDeep(le, false, "SAMSUNG", "NX500", false), expected: ImageSRW},
		{name: "Minolta", buf: buildTiffDeep(le, false, "KONICA MINOLTA", "DYNAX 7D", false), expected: ImageMRW},
		{name: "Epson", buf: buildTiffDeep(le, false, "SEIKO EPSON CORP.", "R-D1", false), expected: ImageERF},
		{name: "Leaf", buf: buildTiffDeep(le, false, "Leaf", "Aptus 75", false), expected: ImageMOS},
		{name: "PhaseOne", buf: buildTiffDeep(le, false, "Phase One", "IQ4", false), expected: ImageIIQ},
		{name: "GoPro", buf: buildTiffDeep(le, false, "GoPro", "HERO11", false), expected: ImageGPR},
		{name: "Mamiya", buf: buildTiffDeep(le, false, "Mamiya", "ZD", false), expected: ImageMEF},
		{name: "NikonNRW", buf: buildTiffDeep(le, false, "NIKON CORPORATION", "COOLPIX P7800", false), expected: ImageNRW},
		{name: "NikonNEF", buf: buildTiffDeep(le, false, "NIKON CORPORATION", "D850", false), expected: ImageNEF},
		{name: "SonySR2", buf: buildTiffDeep(le, false, "SONY", "DSLR-A100", false), expected: ImageSR2},
		{name: "SonySRF", buf: buildTiffDeep(le, false, "SONY", "DSC-R1", false), expected: ImageSRF},
		{name: "SonyARW", buf: buildTiffDeep(le, false, "SONY", "ILCE-7M4", false), expected: ImageARW},
		{name: "KodakK25", buf: buildTiffDeep(le, false, "EASTMAN KODAK COMPANY", "DC25", false), expected: ImageK25},
		{name: "KodakDCR", buf: buildTiffDeep(le, false, "EASTMAN KODAK COMPANY", "DCS Pro 14n", false), expected: ImageDCR},
		{name: "KodakKDC", buf: buildTiffDeep(le, false, "EASTMAN KODAK COMPANY", "DC120", false), expected: ImageKDC},
		{name: "CanonFallsBack", buf: buildTiffDeep(le, false, "Canon", "EOS R5", false), expected: ImageTiff},
		{name: "LeicaFallsBack", buf: buildTiffDeep(le, false, "LEICA CAMERA AG", "M11", false), expected: ImageTiff},
		{name: "HasselbladFallsBack", buf: buildTiffDeep(le, false, "Hasselblad", "X2D", false), expected: ImageTiff},
		{name: "DNGVersion", buf: buildTiffDeep(le, false, "Apple", "iPhone 15", true), expected: ImageDNG},
		{name: "DNGVersionBE", buf: buildTiffDeep(be, false, "Google", "Pixel 8", true), expected: ImageDNG},
		{name: "BigEndian", buf: buildTiffDeep(be, false, "OLYMPUS IMAGING CORP.", "E-M1", false), expected: ImageORF},
		{name: "BigTIFF", buf: buildTiffDeep(le, true, "PENTAX", "645Z", false), expected: ImagePEF},
		{name: "BigTIFFDNG", buf: buildTiffDeep(be, true, "DJI", "Mavic", true), expected: ImageDNG},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got, err := Scan(bytes.NewReader(tc.buf))
			if err != nil {
				t.Fatalf("Scan() returned unexpected error: %v", err)
			}
			if got != tc.expected {
				t.Fatalf("Scan() = %s, expected %s", got, tc.expected)
			}
		})
	}
}

func TestTiffDeepSubtypeTruncated(t *testing.T) {
	t.Parallel()
	full := buildTiffDeep(binary.LittleEndian, false, "NIKON CORPORATION", "D850", false)
	for _, n := range []int{64, 100, 200} {
		got, err := Scan(bytes.NewReader(full[:n]))
		if err != nil {
			t.Fatalf("n=%d: Scan() returned unexpected error: %v", n, err)
		}
		// Truncated string pools must degrade to plain TIFF, never panic
		// or misreport.
		if got != ImageTiff && got != ImageNEF {
			t.Fatalf("n=%d: Scan() = %s, expected TIFF-family fallback", n, got)
		}
	}
}

func TestTiffMakeModelMalformed(t *testing.T) {
	t.Parallel()
	// Garbage entry count must not cause long loops or panics.
	buf := make([]byte, 4096)
	copy(buf, []byte{0x49, 0x49, 0x2A, 0x00})
	binary.LittleEndian.PutUint32(buf[4:], 8)
	binary.LittleEndian.PutUint16(buf[8:], 0xFFFF)
	if _, _, _, ok := tiffMakeModel(buf); ok {
		t.Fatal("tiffMakeModel succeeded on garbage entry count")
	}
	// Non-TIFF magic.
	if _, _, _, ok := tiffMakeModel(bytes.Repeat([]byte{0xFF}, 4096)); ok {
		t.Fatal("tiffMakeModel succeeded on non-TIFF data")
	}
}
