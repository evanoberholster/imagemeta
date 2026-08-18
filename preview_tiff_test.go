package imagemeta

import (
	"bytes"
	"encoding/binary"
	"image"
	"image/jpeg"
	"testing"

	"github.com/pkg/errors"
)

func encodeTestJPEG(t *testing.T, width, height int) []byte {
	t.Helper()
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, image.NewRGBA(image.Rect(0, 0, width, height)), nil); err != nil {
		t.Fatalf("encode test jpeg: %v", err)
	}
	return buf.Bytes()
}

// buildTIFFWithPreviews assembles a minimal little-endian TIFF mimicking the
// NEF layout: IFD0 carries a thumbnail JPEG pair plus a SubIFDs pointer, the
// SubIFD carries the full-size preview JPEG ("JpgFromRaw").
func buildTIFFWithPreviews(thumb, preview []byte) []byte {
	le := binary.LittleEndian
	const (
		tagCompression = 0x0103
		tagSubIFDs     = 0x014a
		tagJPEGOffset  = 0x0201
		tagJPEGLength  = 0x0202
		typeShort      = 3
		typeLong       = 4
	)

	entry := func(id, typ uint16, value uint32) []byte {
		e := le.AppendUint16(nil, id)
		e = le.AppendUint16(e, typ)
		e = le.AppendUint32(e, 1)
		return le.AppendUint32(e, value)
	}

	buf := []byte{'I', 'I', 42, 0, 8, 0, 0, 0}
	// layout: header(8) ifd0(2+3*12+4) subifd(2+3*12+4) thumb preview
	ifd0Start := uint32(len(buf))
	subifdStart := ifd0Start + 2 + 3*12 + 4
	thumbStart := subifdStart + 2 + 3*12 + 4
	previewStart := thumbStart + uint32(len(thumb))

	ifd0 := le.AppendUint16(nil, 3)
	ifd0 = append(ifd0, entry(tagSubIFDs, typeLong, subifdStart)...)
	ifd0 = append(ifd0, entry(tagJPEGOffset, typeLong, thumbStart)...)
	ifd0 = append(ifd0, entry(tagJPEGLength, typeLong, uint32(len(thumb)))...)
	ifd0 = le.AppendUint32(ifd0, 0)

	subifd := le.AppendUint16(nil, 3)
	subifd = append(subifd, entry(tagCompression, typeShort, 6)...)
	subifd = append(subifd, entry(tagJPEGOffset, typeLong, previewStart)...)
	subifd = append(subifd, entry(tagJPEGLength, typeLong, uint32(len(preview)))...)
	subifd = le.AppendUint32(subifd, 0)

	buf = append(buf, ifd0...)
	buf = append(buf, subifd...)
	buf = append(buf, thumb...)
	return append(buf, preview...)
}

func TestTIFFPreviews(t *testing.T) {
	t.Parallel()

	thumb := encodeTestJPEG(t, 16, 8)
	preview := encodeTestJPEG(t, 64, 32)
	r := bytes.NewReader(buildTIFFWithPreviews(thumb, preview))

	previews, err := TIFFPreviews(r)
	if err != nil {
		t.Fatalf("TIFFPreviews: %v", err)
	}
	if len(previews) != 2 {
		t.Fatalf("len(previews) = %d, want 2", len(previews))
	}
	if previews[0].IFD != "SubIFD0" || previews[0].Length != uint32(len(preview)) {
		t.Fatalf("previews[0] = %+v, want the SubIFD0 preview of %d bytes", previews[0], len(preview))
	}
	if previews[1].IFD != "IFD0" || previews[1].Length != uint32(len(thumb)) {
		t.Fatalf("previews[1] = %+v, want the IFD0 thumbnail of %d bytes", previews[1], len(thumb))
	}

	got, err := ExtractPreview(r, previews[0])
	if err != nil {
		t.Fatalf("ExtractPreview: %v", err)
	}
	if !bytes.Equal(got, preview) {
		t.Fatal("ExtractPreview returned different bytes than the embedded preview")
	}
}

func TestPreviewTIFF(t *testing.T) {
	t.Parallel()

	thumb := encodeTestJPEG(t, 16, 8)
	preview := encodeTestJPEG(t, 64, 32)
	r := bytes.NewReader(buildTIFFWithPreviews(thumb, preview))

	got, err := PreviewTIFF(r)
	if err != nil {
		t.Fatalf("PreviewTIFF: %v", err)
	}
	if !bytes.Equal(got, preview) {
		t.Fatal("PreviewTIFF did not return the largest embedded preview")
	}
}

func TestPreviewTIFFNoPreview(t *testing.T) {
	t.Parallel()

	// both "previews" point at non-JPEG bytes: they are listed by their
	// tags but fail SOI validation, so no preview survives
	junk := bytes.Repeat([]byte{0x42}, 64)
	r := bytes.NewReader(buildTIFFWithPreviews(junk, junk))

	if _, err := PreviewTIFF(r); !errors.Is(err, ErrNoPreview) {
		t.Fatalf("PreviewTIFF error = %v, want ErrNoPreview", err)
	}
}

func TestExtractPreviewValidatesSOI(t *testing.T) {
	t.Parallel()

	thumb := encodeTestJPEG(t, 16, 8)
	preview := encodeTestJPEG(t, 64, 32)
	r := bytes.NewReader(buildTIFFWithPreviews(thumb, preview))

	if _, err := ExtractPreview(r, PreviewImage{IFD: "IFD0", Offset: 4, Length: 16}); !errors.Is(err, ErrNoPreview) {
		t.Fatalf("ExtractPreview error = %v, want ErrNoPreview", err)
	}
}

func TestIsRenderableJPEG(t *testing.T) {
	t.Parallel()

	baseline := encodeTestJPEG(t, 8, 8)
	if !isRenderableJPEG(baseline) {
		t.Fatal("isRenderableJPEG(baseline jpeg) = false, want true")
	}

	// lossless JPEG (SOF3) is what DNG raw sensor payloads use; it starts
	// with a regular SOI but common decoders cannot render it
	lossless := []byte{0xff, 0xd8, 0xff, 0xc3, 0x00, 0x0b, 8, 0, 16, 0, 16, 1, 0, 0x11, 0}
	if isRenderableJPEG(lossless) {
		t.Fatal("isRenderableJPEG(lossless jpeg) = true, want false")
	}

	if isRenderableJPEG([]byte{0x42, 0x42, 0x42, 0x42}) {
		t.Fatal("isRenderableJPEG(garbage) = true, want false")
	}
}

func TestExtractPreviewRejectsOutOfBoundsLength(t *testing.T) {
	t.Parallel()

	thumb := encodeTestJPEG(t, 16, 8)
	preview := encodeTestJPEG(t, 64, 32)
	data := buildTIFFWithPreviews(thumb, preview)
	r := bytes.NewReader(data)

	// a preview whose declared length runs past the file must not allocate it
	p := PreviewImage{IFD: "SubIFD0", Offset: 8, Length: 0xffffffff}
	if _, err := ExtractPreview(r, p); !errors.Is(err, ErrNoPreview) {
		t.Fatalf("ExtractPreview error = %v, want ErrNoPreview", err)
	}
}

func TestIsRenderableJPEGSegmentCap(t *testing.T) {
	t.Parallel()

	// SOI followed by more than maxJPEGSegments tiny APP1 segments hides the SOF
	buf := []byte{0xff, 0xd8}
	for n := 0; n < maxJPEGSegments+5; n++ {
		buf = append(buf, 0xff, 0xe1, 0x00, 0x04, 0x00, 0x00)
	}
	buf = append(buf, 0xff, 0xc0, 0x00, 0x0b, 0x08, 0, 16, 0, 16, 0x01, 0x01, 0x11, 0x00)
	if isRenderableJPEG(buf) {
		t.Fatal("isRenderableJPEG walked past the segment cap")
	}
}

func TestExtractPreviewMaxLengthOption(t *testing.T) {
	t.Parallel()

	thumb := encodeTestJPEG(t, 16, 8)
	preview := encodeTestJPEG(t, 64, 32)
	data := buildTIFFWithPreviews(thumb, preview)

	// the preview is larger than the configured limit
	small := WithMaxPreviewLength(int64(len(preview)) - 1)
	previews, err := TIFFPreviews(bytes.NewReader(data), small)
	if err != nil {
		t.Fatalf("TIFFPreviews: %v", err)
	}
	for _, p := range previews {
		if p.IFD == "SubIFD0" {
			t.Fatalf("TIFFPreviews surfaced a preview exceeding the max length: %+v", p)
		}
	}
	full := PreviewImage{IFD: "SubIFD0", Offset: 8, Length: uint32(len(preview))}
	if _, err := ExtractPreview(bytes.NewReader(data), full, small); !errors.Is(err, ErrNoPreview) {
		t.Fatalf("ExtractPreview error = %v, want ErrNoPreview", err)
	}
	// with no limit (default) it extracts fine
	if _, err := PreviewTIFF(bytes.NewReader(data)); err != nil {
		t.Fatalf("PreviewTIFF with default (unlimited): %v", err)
	}
}
