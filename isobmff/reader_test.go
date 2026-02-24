package isobmff

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/evanoberholster/imagemeta/imagetype"
)

func TestReadBoxEightByteHeader(t *testing.T) {
	data := []byte{
		0x00, 0x00, 0x00, 0x08, // size
		'f', 'r', 'e', 'e', // type
	}

	r := NewReader(bytes.NewReader(data), nil, nil, nil)
	t.Cleanup(r.Close)

	b, err := r.readBox()
	if err != nil {
		t.Fatalf("readBox error: %v", err)
	}
	if b.boxType != typeFree {
		t.Fatalf("boxType = %v, want %v", b.boxType, typeFree)
	}
	if b.size != 8 {
		t.Fatalf("size = %d, want 8", b.size)
	}
	if b.remain != 0 {
		t.Fatalf("remain = %d, want 0", b.remain)
	}
}

func TestReadInnerBoxEightByteHeader(t *testing.T) {
	data := []byte{
		0x00, 0x00, 0x00, 0x10, // outer size
		'm', 'o', 'o', 'v', // outer type
		0x00, 0x00, 0x00, 0x08, // inner size
		'f', 'r', 'e', 'e', // inner type
	}

	r := NewReader(bytes.NewReader(data), nil, nil, nil)
	t.Cleanup(r.Close)

	outer, err := r.readBox()
	if err != nil {
		t.Fatalf("readBox outer error: %v", err)
	}

	inner, ok, err := outer.readInnerBox()
	if err != nil {
		t.Fatalf("readInnerBox error: %v", err)
	}
	if !ok {
		t.Fatal("expected inner box")
	}
	if inner.boxType != typeFree {
		t.Fatalf("inner type = %v, want %v", inner.boxType, typeFree)
	}
	if inner.remain != 0 {
		t.Fatalf("inner remain = %d, want 0", inner.remain)
	}

	_, ok, err = outer.readInnerBox()
	if err != nil {
		t.Fatalf("readInnerBox second call error: %v", err)
	}
	if ok {
		t.Fatal("expected no more inner boxes")
	}
}

func TestDiscardUsesSeekForLargeSkips(t *testing.T) {
	payload := make([]byte, seekDiscardThreshold+256)
	source := &countingReadSeeker{Reader: bytes.NewReader(payload)}

	r := NewReader(source, nil, nil, nil)
	t.Cleanup(r.Close)

	discarded, err := r.discard(seekDiscardThreshold)
	if err != nil {
		t.Fatalf("discard error: %v", err)
	}
	if discarded != seekDiscardThreshold {
		t.Fatalf("discarded = %d, want %d", discarded, seekDiscardThreshold)
	}
	if source.seekCalls == 0 {
		t.Fatal("expected seek path to be used")
	}
	if source.readCalls != 0 {
		t.Fatalf("read calls = %d, want 0", source.readCalls)
	}
}

func TestResetPreservesBufferSize(t *testing.T) {
	r := NewReader(bytes.NewReader([]byte("0123456789abcdef")), nil, nil, nil)
	t.Cleanup(r.Close)

	r.reset(bytes.NewReader([]byte("abcdefghijklmno1")))

	if got := r.br.Size(); got != bufReaderSize {
		t.Fatalf("buffer size = %d, want %d", got, bufReaderSize)
	}
}

func TestReadMetadataContinuesAfterJXLP(t *testing.T) {
	data := []byte{
		// ftyp
		0x00, 0x00, 0x00, 0x10,
		'f', 't', 'y', 'p',
		'a', 'v', 'i', 'f',
		'0', '0', '0', '1',
		// jxlp (to be skipped)
		0x00, 0x00, 0x00, 0x0C,
		'j', 'x', 'l', 'p',
		0x00, 0x00, 0x00, 0x00,
		// meta (empty full box, version+flags=0)
		0x00, 0x00, 0x00, 0x0C,
		'm', 'e', 't', 'a',
		0x00, 0x00, 0x00, 0x00,
	}

	r := NewReader(bytes.NewReader(data), nil, nil, nil)
	t.Cleanup(r.Close)

	if err := r.ReadFTYP(); err != nil {
		t.Fatalf("ReadFTYP() error = %v", err)
	}
	if err := r.ReadMetadata(); err != nil {
		t.Fatalf("ReadMetadata() error = %v", err)
	}
	if r.offset != len(data) {
		t.Fatalf("offset = %d, want %d", r.offset, len(data))
	}
}

func TestReadMetadataReadsTopLevelExifBox(t *testing.T) {
	data := []byte{
		// ftyp
		0x00, 0x00, 0x00, 0x10,
		'f', 't', 'y', 'p',
		'a', 'v', 'i', 'f',
		'0', '0', '0', '1',
		// Exif box with TIFF header payload
		0x00, 0x00, 0x00, 0x18,
		'E', 'x', 'i', 'f',
		'I', 'I', 0x2A, 0x00, // TIFF byte-order + marker
		0x08, 0x00, 0x00, 0x00, // IFD0 offset
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
	}

	r := NewReader(bytes.NewReader(data), nil, nil, nil)
	t.Cleanup(r.Close)

	if err := r.ReadFTYP(); err != nil {
		t.Fatalf("ReadFTYP() error = %v", err)
	}
	if err := r.ReadMetadata(); err != nil {
		t.Fatalf("ReadMetadata() error = %v", err)
	}
	if r.offset != len(data) {
		t.Fatalf("offset = %d, want %d", r.offset, len(data))
	}
}

func TestReadMetadataReadsExifWithAPP1Prefix(t *testing.T) {
	data := []byte{
		// ftyp
		0x00, 0x00, 0x00, 0x10,
		'f', 't', 'y', 'p',
		'a', 'v', 'i', 'f',
		'0', '0', '0', '1',
		// Exif box with APP1-style Exif\0\0 prefix
		0x00, 0x00, 0x00, 0x1E,
		'E', 'x', 'i', 'f',
		'E', 'x', 'i', 'f', 0x00, 0x00,
		'I', 'I', 0x2A, 0x00,
		0x08, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
	}

	r := NewReader(bytes.NewReader(data), nil, nil, nil)
	t.Cleanup(r.Close)

	if err := r.ReadFTYP(); err != nil {
		t.Fatalf("ReadFTYP() error = %v", err)
	}
	if err := r.ReadMetadata(); err != nil {
		t.Fatalf("ReadMetadata() error = %v", err)
	}
	if r.offset != len(data) {
		t.Fatalf("offset = %d, want %d", r.offset, len(data))
	}
}

func TestReadMetadataReadsExifWithTIFFOffsetPrefix(t *testing.T) {
	data := []byte{
		// ftyp
		0x00, 0x00, 0x00, 0x10,
		'f', 't', 'y', 'p',
		'a', 'v', 'i', 'f',
		'0', '0', '0', '1',
		// Exif box with 4-byte TIFF offset prefix
		0x00, 0x00, 0x00, 0x1C,
		'E', 'x', 'i', 'f',
		0x00, 0x00, 0x00, 0x00, // TIFF starts immediately after this offset field
		'I', 'I', 0x2A, 0x00,
		0x08, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
	}

	r := NewReader(bytes.NewReader(data), nil, nil, nil)
	t.Cleanup(r.Close)

	if err := r.ReadFTYP(); err != nil {
		t.Fatalf("ReadFTYP() error = %v", err)
	}
	if err := r.ReadMetadata(); err != nil {
		t.Fatalf("ReadMetadata() error = %v", err)
	}
	if r.offset != len(data) {
		t.Fatalf("offset = %d, want %d", r.offset, len(data))
	}
}

func TestReadMetadataSkipsUnknownTopLevelBox(t *testing.T) {
	data := []byte{
		// ftyp
		0x00, 0x00, 0x00, 0x10,
		'f', 't', 'y', 'p',
		'a', 'v', 'i', 'f',
		'0', '0', '0', '1',
		// unknown box with payload
		0x00, 0x00, 0x00, 0x0C,
		'a', 'b', 'c', 'd',
		0x01, 0x02, 0x03, 0x04,
		// meta (empty full box)
		0x00, 0x00, 0x00, 0x0C,
		'm', 'e', 't', 'a',
		0x00, 0x00, 0x00, 0x00,
	}

	r := NewReader(bytes.NewReader(data), nil, nil, nil)
	t.Cleanup(r.Close)

	if err := r.ReadFTYP(); err != nil {
		t.Fatalf("ReadFTYP() error = %v", err)
	}
	if err := r.ReadMetadata(); err != nil {
		t.Fatalf("ReadMetadata() first call error = %v", err)
	}
	if err := r.ReadMetadata(); err != nil {
		t.Fatalf("ReadMetadata() second call error = %v", err)
	}
	if r.offset != len(data) {
		t.Fatalf("offset = %d, want %d", r.offset, len(data))
	}
}

func TestReadMetadataReturnsEOF(t *testing.T) {
	data := []byte{
		// ftyp
		0x00, 0x00, 0x00, 0x10,
		'f', 't', 'y', 'p',
		'a', 'v', 'i', 'f',
		'0', '0', '0', '1',
		// meta (empty full box)
		0x00, 0x00, 0x00, 0x0C,
		'm', 'e', 't', 'a',
		0x00, 0x00, 0x00, 0x00,
	}

	r := NewReader(bytes.NewReader(data), nil, nil, nil)
	t.Cleanup(r.Close)

	if err := r.ReadFTYP(); err != nil {
		t.Fatalf("ReadFTYP() error = %v", err)
	}
	if err := r.ReadMetadata(); err != nil {
		t.Fatalf("ReadMetadata() first call error = %v", err)
	}
	if err := r.ReadMetadata(); !errors.Is(err, io.EOF) {
		t.Fatalf("ReadMetadata() EOF error = %v, want %v", err, io.EOF)
	}
}

func TestReadInfeTruncatedReturnsError(t *testing.T) {
	data := []byte{
		0x00, 0x00, 0x00, 0x14, // size
		'i', 'n', 'f', 'e', // type
		0x02, 0x00, 0x00, 0x00, // version=2
		0x00, 0x01, // item_ID
		0x00, 0x00, // protection index
		'm', 'i', 'm', 'e', // item_type
	}

	r := NewReader(bytes.NewReader(data), nil, nil, nil)
	t.Cleanup(r.Close)

	b, err := r.readBox()
	if err != nil {
		t.Fatalf("readBox() error = %v", err)
	}
	if err := r.readInfe(&b); !errors.Is(err, ErrBufLength) {
		t.Fatalf("readInfe() error = %v, want %v", err, ErrBufLength)
	}
}

func TestReadIlocUnsupportedFieldSizeReturnsError(t *testing.T) {
	data := []byte{
		0x00, 0x00, 0x00, 0x1A, // size
		'i', 'l', 'o', 'c', // type
		0x00, 0x00, 0x00, 0x00, // version=0
		0x31, 0x00, // offset_size=3 (unsupported), length_size=1, base_offset_size=0
		0x00, 0x01, // item_count=1
		0x00, 0x01, // item_ID
		0x00, 0x00, // data_reference_index
		0x00, 0x01, // extent_count
		0x00, 0x00, 0x00, // extent_offset (3 bytes)
		0x05, // extent_length
	}

	r := NewReader(bytes.NewReader(data), nil, nil, nil)
	t.Cleanup(r.Close)

	b, err := r.readBox()
	if err != nil {
		t.Fatalf("readBox() error = %v", err)
	}
	if err := r.readIloc(&b); !errors.Is(err, ErrUnsupportedFieldSize) {
		t.Fatalf("readIloc() error = %v, want %v", err, ErrUnsupportedFieldSize)
	}
}

func TestMetadataImageTypeFromMajorBrand(t *testing.T) {
	tests := []struct {
		name      string
		major     Brand
		wantImage imagetype.ImageType
	}{
		{name: "jxl", major: brandJxl, wantImage: imagetype.ImageJXL},
		{name: "avif", major: brandAvif, wantImage: imagetype.ImageAVIF},
		{name: "avis", major: brandAvis, wantImage: imagetype.ImageAVIF},
		{name: "heic", major: brandHeic, wantImage: imagetype.ImageHEIC},
		{name: "heif", major: brandHeif, wantImage: imagetype.ImageHEIF},
		{name: "cr3", major: brandCrx, wantImage: imagetype.ImageCR3},
		{name: "unknown-fallback", major: brandUnknown, wantImage: imagetype.ImageHEIF},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Reader{
				ftyp: FileTypeBox{
					MajorBrand: tt.major,
				},
			}
			if got := r.metadataImageType(); got != tt.wantImage {
				t.Fatalf("metadataImageType() = %v, want %v", got, tt.wantImage)
			}
		})
	}
}

type countingReadSeeker struct {
	*bytes.Reader
	readCalls int
	seekCalls int
}

func (c *countingReadSeeker) Read(p []byte) (int, error) {
	c.readCalls++
	return c.Reader.Read(p)
}

func (c *countingReadSeeker) Seek(offset int64, whence int) (int64, error) {
	c.seekCalls++
	return c.Reader.Seek(offset, whence)
}
