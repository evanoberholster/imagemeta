package imagemeta

import (
	"bytes"
	"encoding/binary"
	"strings"
	"testing"

	"github.com/evanoberholster/imagemeta/imagetype"
)

// buildBigTIFF assembles a minimal little-endian BigTIFF (DNG 1.7) container
// exercising every value-mapping path: an at-offset ASCII string (Make), an
// inline ASCII string (Model), an inline 8-byte RATIONAL (XResolution) and two
// numeric-scalar offsets, one carried as a 64-bit LONG8 (JPEGInterchangeFormat).
func buildBigTIFF(makeStr, model string, xresNum uint32, preview []byte) []byte {
	le := binary.LittleEndian
	const (
		tagMake       = 0x010F
		tagModel      = 0x0110
		tagXReson     = 0x011A
		tagJPEGOffset = 0x0201
		tagJPEGLength = 0x0202
		typeASCII     = 2
		typeLong      = 4
		typeRational  = 5
		typeLong8     = 16
	)

	// entry: id[2] type[2] count[8] value/offset[8, left-justified]
	entry := func(id, typ uint16, count uint64, value []byte) []byte {
		e := le.AppendUint16(nil, id)
		e = le.AppendUint16(e, typ)
		e = le.AppendUint64(e, count)
		slot := make([]byte, 8)
		copy(slot, value)
		return append(e, slot...)
	}

	const entryCount = 5
	// header(16) + count(8) + entries(5*20) + next(8)
	ifd0Start := uint64(16)
	dataStart := ifd0Start + 8 + entryCount*20 + 8
	makeBytes := append([]byte(makeStr), 0)
	makeOffset := dataStart
	previewStart := dataStart + uint64(len(makeBytes))

	rational := le.AppendUint32(nil, xresNum)
	rational = le.AppendUint32(rational, 1) // denominator

	ifd0 := le.AppendUint64(nil, entryCount)
	ifd0 = append(ifd0, entry(tagMake, typeASCII, uint64(len(makeBytes)), le.AppendUint64(nil, makeOffset))...)
	ifd0 = append(ifd0, entry(tagModel, typeASCII, uint64(len(model)+1), append([]byte(model), 0))...)
	ifd0 = append(ifd0, entry(tagXReson, typeRational, 1, rational)...)
	ifd0 = append(ifd0, entry(tagJPEGOffset, typeLong8, 1, le.AppendUint64(nil, previewStart))...)
	ifd0 = append(ifd0, entry(tagJPEGLength, typeLong, 1, le.AppendUint32(nil, uint32(len(preview))))...)
	ifd0 = le.AppendUint64(ifd0, 0) // no next IFD

	buf := []byte{'I', 'I', 43, 0, 8, 0, 0, 0}
	buf = le.AppendUint64(buf, ifd0Start)
	buf = append(buf, ifd0...)
	buf = append(buf, makeBytes...)
	return append(buf, preview...)
}

func TestBigTIFFImageTypeDetection(t *testing.T) {
	preview := encodeTestJPEG(t, 8, 8)
	data := buildBigTIFF("Make", "Model", 300, preview)
	it, err := imagetype.Buf(data)
	if err != nil {
		t.Fatalf("imagetype.Buf: %v", err)
	}
	if it != imagetype.ImageTiff {
		t.Fatalf("imagetype.Buf: want ImageTiff, got %s", it)
	}
}

func TestBigTIFFDecodeMetadata(t *testing.T) {
	preview := encodeTestJPEG(t, 8, 8)
	data := buildBigTIFF("TestMake", "Cam", 300, preview)

	ex, err := DecodeTiff(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("DecodeTiff: %v", err)
	}
	// Make is read from a file offset; imagemeta canonicalizes it to lower case.
	if !strings.EqualFold(ex.IFD0.Make, "TestMake") {
		t.Errorf("Make: want %q (any case), got %q (at-offset ASCII)", "TestMake", ex.IFD0.Make)
	}
	if ex.IFD0.Model != "Cam" {
		t.Errorf("Model: want %q, got %q (inline ASCII)", "Cam", ex.IFD0.Model)
	}
	// The RATIONAL is packed inline in the 8-byte BigTIFF value field: this is
	// the path that needs the full 8 inline bytes, not just the low four.
	if ex.IFD0.XResolution != 300 {
		t.Errorf("XResolution: want 300, got %v (inline 8-byte RATIONAL)", ex.IFD0.XResolution)
	}
}

func TestBigTIFFPreviewExtraction(t *testing.T) {
	preview := encodeTestJPEG(t, 32, 16)
	data := buildBigTIFF("Make", "Model", 72, preview)
	r := bytes.NewReader(data)

	previews, err := TIFFPreviews(r)
	if err != nil {
		t.Fatalf("TIFFPreviews: %v", err)
	}
	if len(previews) != 1 {
		t.Fatalf("want 1 preview, got %d", len(previews))
	}
	got, err := ExtractPreview(r, previews[0])
	if err != nil {
		t.Fatalf("ExtractPreview: %v", err)
	}
	if !bytes.Equal(got, preview) {
		t.Fatalf("preview bytes mismatch: want %d bytes, got %d", len(preview), len(got))
	}
}

func TestBigTIFFRejectsBadOffsetSize(t *testing.T) {
	preview := encodeTestJPEG(t, 8, 8)
	data := buildBigTIFF("Make", "Model", 300, preview)
	data[4] = 4 // offset size must be 8; a malformed header must not be parsed

	if _, err := DecodeTiff(bytes.NewReader(data)); err == nil {
		t.Fatal("DecodeTiff: expected error for malformed BigTIFF offset size")
	}
}
