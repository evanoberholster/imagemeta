package isobmff

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"testing"

	"github.com/evanoberholster/imagemeta/imagetype"
	"github.com/evanoberholster/imagemeta/meta"
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

func TestReadBoxShortHeaderReturnsErrBufLength(t *testing.T) {
	r := NewReader(bytes.NewReader([]byte{0x00, 0x00, 0x00, 0x08}), nil, nil, nil)
	t.Cleanup(r.Close)

	_, err := r.readBox()
	if !errors.Is(err, ErrBufLength) {
		t.Fatalf("readBox() error = %v, want %v", err, ErrBufLength)
	}
	if errors.Is(err, io.EOF) {
		t.Fatalf("readBox() error = %v, should not be treated as EOF", err)
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

func TestBoxPeekNegativeReturnsError(t *testing.T) {
	b := box{remain: 4}
	if _, err := b.Peek(-1); !errors.Is(err, ErrBufLength) {
		t.Fatalf("Peek(-1) error = %v, want %v", err, ErrBufLength)
	}
}

func TestBoxDiscardNegativeReturnsError(t *testing.T) {
	b := box{remain: 4}
	if _, err := b.Discard(-1); !errors.Is(err, ErrBufLength) {
		t.Fatalf("Discard(-1) error = %v, want %v", err, ErrBufLength)
	}
}

func TestBoxPayloadOffsetUsesInt64Math(t *testing.T) {
	b := box{
		offset: 1 << 40,
		size:   32,
		remain: 8,
	}
	got := boxPayloadOffset(&b)
	want := uint64((int64(1) << 40) + 24)
	if got != want {
		t.Fatalf("boxPayloadOffset = %d, want %d", got, want)
	}
}

func TestInitMetadataGoalsPreviewCR3Only(t *testing.T) {
	noopPreview := func(_ io.Reader, _ meta.PreviewHeader) error { return nil }
	r := NewReader(bytes.NewReader(nil), nil, nil, noopPreview)
	t.Cleanup(r.Close)

	r.ftyp.MajorBrand = brandCrx
	r.initMetadataGoals()
	if r.hasGoal(metadataKindTHMB) || !r.hasGoal(metadataKindPRVW) {
		t.Fatal("expected PRVW-only goal for CR3 preview")
	}

	r.ftyp.MajorBrand = brandHeic
	r.initMetadataGoals()
	if r.hasGoal(metadataKindTHMB) || r.hasGoal(metadataKindPRVW) {
		t.Fatal("expected preview goals to be disabled for non-CR3")
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
	if r.offset != int64(len(data)) {
		t.Fatalf("offset = %d, want %d", r.offset, len(data))
	}
}

func TestReadMetadataSkipsMultipleJXLBoxes(t *testing.T) {
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
		// jxlc (to be skipped)
		0x00, 0x00, 0x00, 0x0C,
		'j', 'x', 'l', 'c',
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
	if r.offset != int64(len(data)) {
		t.Fatalf("offset = %d, want %d", r.offset, len(data))
	}
	if err := r.ReadMetadata(); !errors.Is(err, io.EOF) {
		t.Fatalf("ReadMetadata() EOF error = %v, want %v", err, io.EOF)
	}
}

func TestReadMetadataTruncatedTailReturnsErrBufLength(t *testing.T) {
	data := []byte{
		// ftyp
		0x00, 0x00, 0x00, 0x10,
		'f', 't', 'y', 'p',
		'a', 'v', 'i', 'f',
		'0', '0', '0', '1',
		// trailing truncated bytes (not a full box header)
		0x00, 0x00, 0x00,
	}

	r := NewReader(bytes.NewReader(data), nil, nil, nil)
	t.Cleanup(r.Close)

	if err := r.ReadFTYP(); err != nil {
		t.Fatalf("ReadFTYP() error = %v", err)
	}
	err := r.ReadMetadata()
	if !errors.Is(err, ErrBufLength) {
		t.Fatalf("ReadMetadata() error = %v, want %v", err, ErrBufLength)
	}
	if errors.Is(err, io.EOF) {
		t.Fatalf("ReadMetadata() error = %v, should not be EOF for truncated tail", err)
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
	if r.offset != int64(len(data)) {
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
	if r.offset != int64(len(data)) {
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
	if r.offset != int64(len(data)) {
		t.Fatalf("offset = %d, want %d", r.offset, len(data))
	}
}

func TestReadMetadataReadsMinimalExifHeader(t *testing.T) {
	data := []byte{
		// ftyp
		0x00, 0x00, 0x00, 0x10,
		'f', 't', 'y', 'p',
		'a', 'v', 'i', 'f',
		'0', '0', '0', '1',
		// Exif box with only TIFF header (8 bytes)
		0x00, 0x00, 0x00, 0x10,
		'E', 'x', 'i', 'f',
		'I', 'I', 0x2A, 0x00,
		0x08, 0x00, 0x00, 0x00,
	}

	r := NewReader(bytes.NewReader(data), nil, nil, nil)
	t.Cleanup(r.Close)

	if err := r.ReadFTYP(); err != nil {
		t.Fatalf("ReadFTYP() error = %v", err)
	}
	if err := r.ReadMetadata(); err != nil {
		t.Fatalf("ReadMetadata() error = %v", err)
	}
	if r.offset != int64(len(data)) {
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
	if r.offset != int64(len(data)) {
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

func TestReadMetaSkipsNonExifItemGraphBoxesForHEIF(t *testing.T) {
	// Malformed iprp/ipma payload that would fail if parsed.
	badIPMA := makeReaderTestBox("ipma", []byte{
		0x00, 0x00, 0x00, 0x00, // flags
		0x00, 0x00, // truncated entry_count
	})
	metaPayload := append([]byte{
		0x00, 0x00, 0x00, 0x00, // meta full box flags
	}, makeReaderTestBox("iprp", badIPMA)...)

	r := NewReader(bytes.NewReader(makeReaderTestBox("meta", metaPayload)), nil, nil, nil)
	t.Cleanup(r.Close)
	r.ftyp.MajorBrand = brandHeic

	b, err := r.readBox()
	if err != nil {
		t.Fatalf("readBox() error = %v", err)
	}
	if err := r.readMeta(&b); err != nil {
		t.Fatalf("readMeta() error = %v, want nil for HEIF skip path", err)
	}
}

func TestReadMetaParsesItemGraphBoxesForCR3(t *testing.T) {
	// Same malformed payload as above, but CR3 path should parse iprp/ipma and fail.
	badIPMA := makeReaderTestBox("ipma", []byte{
		0x00, 0x00, 0x00, 0x00, // flags
		0x00, 0x00, // truncated entry_count
	})
	metaPayload := append([]byte{
		0x00, 0x00, 0x00, 0x00, // meta full box flags
	}, makeReaderTestBox("iprp", badIPMA)...)

	r := NewReader(bytes.NewReader(makeReaderTestBox("meta", metaPayload)), nil, nil, nil)
	t.Cleanup(r.Close)
	r.ftyp.MajorBrand = brandCrx

	b, err := r.readBox()
	if err != nil {
		t.Fatalf("readBox() error = %v", err)
	}
	if err := r.readMeta(&b); !errors.Is(err, ErrBufLength) {
		t.Fatalf("readMeta() error = %v, want %v", err, ErrBufLength)
	}
}

func TestReadInfeVersion3ParsesItemID(t *testing.T) {
	payload := make([]byte, 0, 24)
	payload = append(payload, 0x03, 0x00, 0x00, 0x00) // version=3
	payload = binary.BigEndian.AppendUint32(payload, 0x1020)
	payload = binary.BigEndian.AppendUint16(payload, 0) // protection index
	payload = append(payload, 'E', 'x', 'i', 'f')       // item_type
	payload = append(payload, 0x00)                     // item_name
	data := makeReaderTestBox("infe", payload)          // no optional fields for Exif
	r := NewReader(bytes.NewReader(data), nil, nil, nil)
	t.Cleanup(r.Close)

	b, err := r.readBox()
	if err != nil {
		t.Fatalf("readBox() error: %v", err)
	}
	if err := r.readInfe(&b); err != nil {
		t.Fatalf("readInfe() error: %v", err)
	}
	if r.heic.exif.id != itemID(0x1020) {
		t.Fatalf("exif item ID = %d, want %d", r.heic.exif.id, 0x1020)
	}
}

func TestReadIinfVersion1WithNoEntries(t *testing.T) {
	payload := []byte{
		0x01, 0x00, 0x00, 0x00, // version=1
		0x00, 0x00, 0x00, 0x00, // entry_count=0
	}
	data := makeReaderTestBox("iinf", payload)

	r := NewReader(bytes.NewReader(data), nil, nil, nil)
	t.Cleanup(r.Close)

	b, err := r.readBox()
	if err != nil {
		t.Fatalf("readBox() error: %v", err)
	}
	if err := r.readIinf(&b); err != nil {
		t.Fatalf("readIinf() error: %v", err)
	}
	if b.remain != 0 {
		t.Fatalf("iinf remain = %d, want 0", b.remain)
	}
}

func TestHdlrTypeString(t *testing.T) {
	if got := hdlrPict.String(); got != "pict" {
		t.Fatalf("hdlrPict.String() = %q, want %q", got, "pict")
	}
	if got := hdlrVide.String(); got != "vide" {
		t.Fatalf("hdlrVide.String() = %q, want %q", got, "vide")
	}
	if got := hdlrMeta.String(); got != "meta" {
		t.Fatalf("hdlrMeta.String() = %q, want %q", got, "meta")
	}
	if got := hdlrUnknown.String(); got != "nnnn" {
		t.Fatalf("hdlrUnknown.String() = %q, want %q", got, "nnnn")
	}
}

func TestParseExtendedBoxSizeMaxInt64(t *testing.T) {
	buf := make([]byte, 16)
	binary.BigEndian.PutUint32(buf[:4], 1)
	copy(buf[4:8], []byte("mdat"))
	binary.BigEndian.PutUint64(buf[8:16], uint64(maxInt64Value))

	size, err := parseExtendedBoxSize(buf, typeMdat)
	if err != nil {
		t.Fatalf("parseExtendedBoxSize() error = %v", err)
	}
	if size != maxInt64Value {
		t.Fatalf("size = %d, want %d", size, maxInt64Value)
	}
}

func TestParseExtendedBoxSizeAboveInt64Fails(t *testing.T) {
	buf := make([]byte, 16)
	binary.BigEndian.PutUint32(buf[:4], 1)
	copy(buf[4:8], []byte("mdat"))
	binary.BigEndian.PutUint64(buf[8:16], uint64(maxInt64Value)+1)

	_, err := parseExtendedBoxSize(buf, typeMdat)
	if !errors.Is(err, errLargeBox) {
		t.Fatalf("parseExtendedBoxSize() error = %v, want %v", err, errLargeBox)
	}
}

func makeReaderTestBox(typ string, payload []byte) []byte {
	out := make([]byte, 8+len(payload))
	binary.BigEndian.PutUint32(out[:4], uint32(len(out)))
	copy(out[4:8], []byte(typ))
	copy(out[8:], payload)
	return out
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

func TestReadInfeMimeContentTypeTooLongReturnsError(t *testing.T) {
	mime := bytes.Repeat([]byte{'a'}, mimeContentTypeMaxLen+1)
	payload := make([]byte, 0, 4+2+2+4+1+len(mime)+1)
	payload = append(payload, 0x02, 0x00, 0x00, 0x00) // version=2
	payload = append(payload, 0x12, 0x34)             // item_ID
	payload = append(payload, 0x00, 0x00)             // protection_index
	payload = append(payload, 'm', 'i', 'm', 'e')     // item_type
	payload = append(payload, 0x00)                   // item_name
	payload = append(payload, mime...)                // content_type
	payload = append(payload, 0x00)                   // content_type terminator

	r := NewReader(bytes.NewReader(makeReaderTestBox("infe", payload)), nil, nil, nil)
	t.Cleanup(r.Close)

	b, err := r.readBox()
	if err != nil {
		t.Fatalf("readBox() error = %v", err)
	}
	if err := r.readInfe(&b); !errors.Is(err, ErrBoxStringTooLong) {
		t.Fatalf("readInfe() error = %v, want %v", err, ErrBoxStringTooLong)
	}
}

func TestReadInfeMimeXMPContentTypeSetsXMPID(t *testing.T) {
	payload := make([]byte, 0, 4+2+2+4+1+20)
	payload = append(payload, 0x02, 0x00, 0x00, 0x00)           // version=2
	payload = append(payload, 0x43, 0x21)                       // item_ID
	payload = append(payload, 0x00, 0x00)                       // protection_index
	payload = append(payload, 'm', 'i', 'm', 'e')               // item_type
	payload = append(payload, 0x00)                             // item_name
	payload = append(payload, []byte("application/rdf+xml")...) // content_type
	payload = append(payload, 0x00)                             // content_type terminator

	r := NewReader(bytes.NewReader(makeReaderTestBox("infe", payload)), nil, nil, nil)
	t.Cleanup(r.Close)
	r.setGoal(metadataKindXMP, true)

	b, err := r.readBox()
	if err != nil {
		t.Fatalf("readBox() error = %v", err)
	}
	if err := r.readInfe(&b); err != nil {
		t.Fatalf("readInfe() error = %v", err)
	}
	if r.heic.xml.id != itemID(0x4321) {
		t.Fatalf("xmp item ID = %d, want %d", r.heic.xml.id, 0x4321)
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

func TestResolveIlocExtentOffsetConstructionMethod1UsesIDAT(t *testing.T) {
	r := Reader{
		heic: heicMeta{
			idatData: offsetLength{
				offset: 100,
				length: 64,
			},
		},
	}
	ent := ilocEntry{
		id:                 1,
		baseOffset:         8,
		dataReferenceIndex: 0,
		constructionMethod: 1,
	}

	got, ok := r.resolveIlocExtentOffset(ent, 4)
	if !ok {
		t.Fatal("expected offset resolution to succeed")
	}
	if got != 112 {
		t.Fatalf("resolved offset = %d, want 112", got)
	}
}

func TestMetadataImageTypeFromMajorBrand(t *testing.T) {
	tests := []struct {
		name       string
		major      brand
		compatible []brand
		wantImage  imagetype.ImageType
	}{
		{name: "jxl", major: brandJxl, wantImage: imagetype.ImageJXL},
		{name: "avif", major: brandAvif, wantImage: imagetype.ImageAVIF},
		{name: "avis", major: brandAvis, wantImage: imagetype.ImageAVIF},
		{name: "heic", major: brandHeic, wantImage: imagetype.ImageHEIC},
		{name: "heif", major: brandHeif, wantImage: imagetype.ImageHEIF},
		{name: "cr3", major: brandCrx, wantImage: imagetype.ImageCR3},
		{name: "compatible-fallback-avif", major: brandUnknown, compatible: []brand{brandUnknown, brandAvif}, wantImage: imagetype.ImageAVIF},
		{name: "unknown-fallback", major: brandUnknown, wantImage: imagetype.ImageHEIF},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var compatible [maxBrandCount]brand
			copy(compatible[:], tt.compatible)
			r := Reader{
				ftyp: fileTypeBox{
					MajorBrand: tt.major,
					Compatible: compatible,
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
