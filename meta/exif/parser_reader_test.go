package exif

import (
	"bufio"
	"bytes"
	"errors"
	"math"
	"testing"

	"github.com/evanoberholster/imagemeta/meta"
	"github.com/evanoberholster/imagemeta/meta/exif/makernote"
	"github.com/evanoberholster/imagemeta/meta/exif/tag"
	"github.com/evanoberholster/imagemeta/meta/utils"
	"github.com/rs/zerolog"
)

func TestMergeIFD0CoreFields(t *testing.T) {
	t.Parallel()

	dst := Exif{}
	dst.IFD0.ImageWidth = 100
	dst.IFD0.TileWidth = 0

	src := Exif{}
	src.IFD0.ImageWidth = 200
	src.IFD0.ImageHeight = 300
	src.IFD0.TileWidth = 512
	src.IFD0.TileLength = 256
	src.IFD0.RowsPerStrip = 42

	mergeIFD0CoreFields(&dst, src)

	if dst.IFD0.ImageWidth != 100 {
		t.Fatalf("ImageWidth overwritten: got %d, want 100", dst.IFD0.ImageWidth)
	}
	if dst.IFD0.ImageHeight != 300 {
		t.Fatalf("ImageHeight not copied: got %d, want 300", dst.IFD0.ImageHeight)
	}
	if dst.IFD0.TileWidth != 512 || dst.IFD0.TileLength != 256 {
		t.Fatalf("tile fields not copied: got (%d,%d)", dst.IFD0.TileWidth, dst.IFD0.TileLength)
	}
	if dst.IFD0.RowsPerStrip != 42 {
		t.Fatalf("RowsPerStrip not copied: got %d, want 42", dst.IFD0.RowsPerStrip)
	}

	mergeIFD0CoreFields(nil, src) // no panic
}

func TestTagTypeFor(t *testing.T) {
	t.Parallel()

	if got := tagTypeFor(tag.IFD0, tag.TagExifIFDPointer, tag.TypeLong); got != tag.TypeIfd {
		t.Fatalf("tagTypeFor(ifd0, exif ptr, long) = %v, want %v", got, tag.TypeIfd)
	}
	if got := tagTypeFor(tag.IFD0, tag.TagGPSIFDPointer, tag.TypeUndefined); got != tag.TypeIfd {
		t.Fatalf("tagTypeFor(ifd0, gps ptr, undef) = %v, want %v", got, tag.TypeIfd)
	}
	if got := tagTypeFor(tag.ExifIFD, tag.TagMakerNote, tag.TypeLong); got != tag.TypeIfd {
		t.Fatalf("tagTypeFor(exif, makernote, long) = %v, want %v", got, tag.TypeIfd)
	}
	if got := tagTypeFor(tag.IFD0, tag.TagMake, tag.TypeLong); got != tag.TypeLong {
		t.Fatalf("tagTypeFor(ifd0, make, long) = %v, want %v", got, tag.TypeLong)
	}
	if got := tagTypeFor(tag.ExifIFD, tag.TagMakerNote, tag.TypeShort); got != tag.TypeShort {
		t.Fatalf("tagTypeFor(exif, makernote, short) = %v, want %v", got, tag.TypeShort)
	}
}

func TestTagUsesIfdType(t *testing.T) {
	t.Parallel()

	if !tagUsesIfdType(tag.IFD0, tag.TagExifIFDPointer) {
		t.Fatal("tagUsesIfdType(ifd0, exif ptr) = false, want true")
	}
	if !tagUsesIfdType(tag.IFD0, tag.TagGPSIFDPointer) {
		t.Fatal("tagUsesIfdType(ifd0, gps ptr) = false, want true")
	}
	if !tagUsesIfdType(tag.ExifIFD, tag.TagMakerNote) {
		t.Fatal("tagUsesIfdType(exif, makernote) = false, want true")
	}
	if tagUsesIfdType(tag.IFD0, tag.TagMake) {
		t.Fatal("tagUsesIfdType(ifd0, make) = true, want false")
	}
	if tagUsesIfdType(tag.GPSIFD, tag.TagMakerNote) {
		t.Fatal("tagUsesIfdType(gps, makernote) = true, want false")
	}
}

func TestTagFromBuffer(t *testing.T) {
	t.Parallel()

	d := tag.NewDirectory(utils.LittleEndian, tag.IFD0, 0, 0x100, 0x10)
	var buf [12]byte
	d.ByteOrder.PutUint16(buf[0:2], uint16(tag.TagExifIFDPointer))
	d.ByteOrder.PutUint16(buf[2:4], uint16(tag.TypeLong))
	d.ByteOrder.PutUint32(buf[4:8], 1)
	d.ByteOrder.PutUint32(buf[8:12], 0x20)

	entry, err := tagFromBuffer(d, buf[:])
	if err != nil {
		t.Fatalf("tagFromBuffer() unexpected error: %v", err)
	}
	if entry.ID != tag.TagExifIFDPointer {
		t.Fatalf("entry.ID = %v, want %v", entry.ID, tag.TagExifIFDPointer)
	}
	if entry.Type != tag.TypeIfd {
		t.Fatalf("entry.Type = %v, want %v", entry.Type, tag.TypeIfd)
	}
	if entry.UnitCount != 1 {
		t.Fatalf("entry.UnitCount = %d, want 1", entry.UnitCount)
	}
	if entry.ValueOffset != 0x30 {
		t.Fatalf("entry.ValueOffset = 0x%x, want 0x30", entry.ValueOffset)
	}

	dEmbedded := tag.NewDirectory(utils.LittleEndian, tag.MakerNoteIFD, 0, 0, 0x853a)
	var embedded [12]byte
	dEmbedded.ByteOrder.PutUint16(embedded[0:2], uint16(makernote.TagNikonMakerNoteVersion))
	dEmbedded.ByteOrder.PutUint16(embedded[2:4], uint16(tag.TypeUndefined))
	dEmbedded.ByteOrder.PutUint32(embedded[4:8], 4)
	copy(embedded[8:12], []byte("0211"))

	entry, err = tagFromBuffer(dEmbedded, embedded[:])
	if err != nil {
		t.Fatalf("tagFromBuffer() embedded unexpected error: %v", err)
	}
	if got, want := entry.ValueOffset, uint32(0x31313230); got != want {
		t.Fatalf("embedded entry.ValueOffset = 0x%x, want 0x%x", got, want)
	}

	d2 := tag.NewDirectory(utils.LittleEndian, tag.IFD0, 0, 0, 0)
	var invalid [12]byte
	d2.ByteOrder.PutUint16(invalid[0:2], uint16(tag.TagMake))
	d2.ByteOrder.PutUint16(invalid[2:4], 0) // TypeUnknown
	d2.ByteOrder.PutUint32(invalid[4:8], 1)
	d2.ByteOrder.PutUint32(invalid[8:12], 0)

	_, err = tagFromBuffer(d2, invalid[:])
	if !errors.Is(err, tag.ErrTagTypeNotValid) {
		t.Fatalf("tagFromBuffer() error = %v, want %v", err, tag.ErrTagTypeNotValid)
	}
}

func TestParseDirectoryTagHeadersBulkMatchesPerEntry(t *testing.T) {
	t.Parallel()

	var payload [36]byte
	// Embedded LONG: ImageWidth=4000.
	utils.LittleEndian.PutUint16(payload[0:2], uint16(tag.TagImageWidth))
	utils.LittleEndian.PutUint16(payload[2:4], uint16(tag.TypeLong))
	utils.LittleEndian.PutUint32(payload[4:8], 1)
	utils.LittleEndian.PutUint32(payload[8:12], 4000)
	// Embedded LONG: ImageLength=3000.
	utils.LittleEndian.PutUint16(payload[12:14], uint16(tag.TagImageLength))
	utils.LittleEndian.PutUint16(payload[14:16], uint16(tag.TypeLong))
	utils.LittleEndian.PutUint32(payload[16:20], 1)
	utils.LittleEndian.PutUint32(payload[20:24], 3000)
	// Non-embedded LONG[2]: queued tag.
	utils.LittleEndian.PutUint16(payload[24:26], uint16(tag.TagStripOffsets))
	utils.LittleEndian.PutUint16(payload[26:28], uint16(tag.TypeLong))
	utils.LittleEndian.PutUint32(payload[28:32], 2)
	utils.LittleEndian.PutUint32(payload[32:36], 0x2000)

	directory := tag.NewDirectory(utils.LittleEndian, tag.IFD0, 0, 0, 0)

	perEntry := NewReader(zerolog.Nop())
	defer perEntry.Close()
	perEntryRaw := bytes.NewReader(payload[:])
	perEntry.Reset(bufio.NewReaderSize(perEntryRaw, len(payload)))
	if err := perEntry.parseDirectoryTagHeadersPerEntry(directory, 3); err != nil {
		t.Fatalf("parseDirectoryTagHeadersPerEntry() error = %v", err)
	}

	bulk := NewReader(zerolog.Nop())
	defer bulk.Close()
	bulkRaw := bytes.NewReader(payload[:])
	bulk.Reset(bufio.NewReaderSize(bulkRaw, len(payload)))
	if err := bulk.parseDirectoryTagHeadersBulk(directory, 3); err != nil {
		t.Fatalf("parseDirectoryTagHeadersBulk() error = %v", err)
	}

	if perEntry.Exif.IFD0.ImageWidth != bulk.Exif.IFD0.ImageWidth {
		t.Fatalf("ImageWidth mismatch: per-entry=%d bulk=%d", perEntry.Exif.IFD0.ImageWidth, bulk.Exif.IFD0.ImageWidth)
	}
	if perEntry.Exif.IFD0.ImageHeight != bulk.Exif.IFD0.ImageHeight {
		t.Fatalf("ImageHeight mismatch: per-entry=%d bulk=%d", perEntry.Exif.IFD0.ImageHeight, bulk.Exif.IFD0.ImageHeight)
	}
	if perEntry.state.len != bulk.state.len {
		t.Fatalf("queued tag count mismatch: per-entry=%d bulk=%d", perEntry.state.len, bulk.state.len)
	}
	if perEntry.state.len != 1 {
		t.Fatalf("queued tag count = %d, want 1", perEntry.state.len)
	}
	if perEntry.state.tag[0] != bulk.state.tag[0] {
		t.Fatalf("queued tag mismatch: per-entry=%+v bulk=%+v", perEntry.state.tag[0], bulk.state.tag[0])
	}
}

func TestParseDirectoryTagHeadersBulkTrustedEmbeddedBaseOffset(t *testing.T) {
	t.Parallel()

	var payload [12]byte
	utils.LittleEndian.PutUint16(payload[0:2], uint16(makernote.TagNikonMakerNoteVersion))
	utils.LittleEndian.PutUint16(payload[2:4], uint16(tag.TypeUndefined))
	utils.LittleEndian.PutUint32(payload[4:8], 4)
	copy(payload[8:12], []byte("0211"))

	directory := tag.NewDirectory(utils.LittleEndian, tag.MakerNoteIFD, 0, 0, 0x853a)
	r := NewReader(zerolog.Nop())
	defer r.Close()
	r.Reset(bufio.NewReaderSize(bytes.NewReader(payload[:]), len(payload)))
	r.Exif.CameraMakeID = makernote.CameraMakeNikon
	r.Exif.MakerNote.Make = makernote.CameraMakeNikon

	if err := r.parseDirectoryTagHeadersBulkTrusted(directory, 1); err != nil {
		t.Fatalf("parseDirectoryTagHeadersBulkTrusted() error = %v", err)
	}
	if got, want := r.Exif.MakerNote.Nikon.MakerNoteVersion, "0211"; got != want {
		t.Fatalf("MakerNoteVersion = %q, want %q", got, want)
	}
}

func TestParseSubSecTimeEmbedded(t *testing.T) {
	t.Parallel()

	r := NewReader(zerolog.Nop())
	defer r.Close()

	var raw [4]byte
	copy(raw[:], []byte{'1', '2', '3', 0})
	valueOffset := utils.LittleEndian.Uint32(raw[:])
	e := tag.NewEntry(tag.TagSubSecTime, tag.TypeASCII, 4, valueOffset, tag.ExifIFD, 0, utils.LittleEndian)

	if got := r.parseSubSecTime(e); got != 123 {
		t.Fatalf("parseSubSecTime() = %d, want 123", got)
	}

	nonASCII := tag.NewEntry(tag.TagSubSecTime, tag.TypeLong, 1, valueOffset, tag.ExifIFD, 0, utils.LittleEndian)
	if got := r.parseSubSecTime(nonASCII); got != 0 {
		t.Fatalf("parseSubSecTime(non-ASCII) = %d, want 0", got)
	}
}

func TestParseSubIFDsSingleTypeIFDUsesValueOffset(t *testing.T) {
	t.Parallel()

	r := NewReader(zerolog.Nop())
	defer r.Close()

	tg := tag.NewEntry(tag.TagSubIFDs, tag.TypeIfd, 1, 0x1234, tag.IFD0, 0, utils.LittleEndian)
	r.parseSubIFDs(tg)

	if got := r.Exif.IFD0.SubIFDOffsetCount; got != 1 {
		t.Fatalf("SubIFDOffsetCount = %d, want 1", got)
	}
	if got := r.Exif.IFD0.SubIFDOffsets[0]; got != 0x1234 {
		t.Fatalf("SubIFDOffsets[0] = 0x%x, want 0x1234", got)
	}
	if got := r.state.len; got != 1 {
		t.Fatalf("state.len = %d, want 1", got)
	}

	queued := r.state.tag[0]
	if queued.IfdType != tag.SubIFD0 {
		t.Fatalf("queued.IfdType = %v, want %v", queued.IfdType, tag.SubIFD0)
	}
	if queued.Type != tag.TypeIfd {
		t.Fatalf("queued.Type = %v, want %v", queued.Type, tag.TypeIfd)
	}
	if queued.ValueOffset != 0x1234 {
		t.Fatalf("queued.ValueOffset = 0x%x, want 0x1234", queued.ValueOffset)
	}
}

func TestParseSubIFDsClampsToQueueCapacity(t *testing.T) {
	t.Parallel()

	var payload [16]byte
	utils.LittleEndian.PutUint32(payload[0:4], 0x10)
	utils.LittleEndian.PutUint32(payload[4:8], 0x20)
	utils.LittleEndian.PutUint32(payload[8:12], 0x30)
	utils.LittleEndian.PutUint32(payload[12:16], 0x40)

	r := NewReader(zerolog.Nop())
	defer r.Close()
	r.Reset(bytes.NewReader(payload[:]))
	r.state.len = tagQueueMax - 1

	tg := tag.NewEntry(tag.TagSubIFDs, tag.TypeLong, 4, 0, tag.IFD0, 0, utils.LittleEndian)
	r.parseSubIFDs(tg)

	if got := r.Exif.IFD0.SubIFDOffsetCount; got != 1 {
		t.Fatalf("SubIFDOffsetCount = %d, want 1", got)
	}
	if got := r.Exif.IFD0.SubIFDOffsets[0]; got != 0x10 {
		t.Fatalf("SubIFDOffsets[0] = 0x%x, want 0x10", got)
	}
	if got := r.state.len; got != tagQueueMax {
		t.Fatalf("state.len = %d, want %d", got, tagQueueMax)
	}
	if got := r.state.tag[tagQueueMax-1].ValueOffset; got != 0x10 {
		t.Fatalf("queued.ValueOffset = 0x%x, want 0x10", got)
	}
}

func TestParseSubIFDsClampsToOffsetCapacity(t *testing.T) {
	t.Parallel()

	var payload [16]byte
	utils.LittleEndian.PutUint32(payload[0:4], 0x10)
	utils.LittleEndian.PutUint32(payload[4:8], 0x20)
	utils.LittleEndian.PutUint32(payload[8:12], 0x30)
	utils.LittleEndian.PutUint32(payload[12:16], 0x40)

	r := NewReader(zerolog.Nop())
	defer r.Close()
	r.Reset(bytes.NewReader(payload[:]))
	r.Exif.IFD0.SubIFDOffsetCount = uint8(len(r.Exif.IFD0.SubIFDOffsets) - 1)

	tg := tag.NewEntry(tag.TagSubIFDs, tag.TypeLong, 4, 0, tag.IFD0, 0, utils.LittleEndian)
	r.parseSubIFDs(tg)

	if got := r.Exif.IFD0.SubIFDOffsetCount; got != uint8(len(r.Exif.IFD0.SubIFDOffsets)) {
		t.Fatalf("SubIFDOffsetCount = %d, want %d", got, len(r.Exif.IFD0.SubIFDOffsets))
	}
	last := len(r.Exif.IFD0.SubIFDOffsets) - 1
	if got := r.Exif.IFD0.SubIFDOffsets[last]; got != 0x10 {
		t.Fatalf("SubIFDOffsets[last] = 0x%x, want 0x10", got)
	}
	if got := r.state.len; got != 1 {
		t.Fatalf("state.len = %d, want 1", got)
	}
}

func TestApertureValueToFNumber(t *testing.T) {
	t.Parallel()

	if got := apertureValueToFNumber(0); got != 0 {
		t.Fatalf("apertureValueToFNumber(0) = %v, want 0", got)
	}
	if got, want := apertureValueToFNumber(meta.Aperture(2)), meta.Aperture(2); got != want {
		t.Fatalf("apertureValueToFNumber(2) = %v, want %v", got, want)
	}
	if got, want := apexApertureToFNumber(1), meta.Aperture(math.Sqrt2); math.Abs(float64(got-want)) > 0.0001 {
		t.Fatalf("apexApertureToFNumber(1) = %v, want %v", got, want)
	}
	if got := apexApertureToFNumber(2048); !math.IsInf(float64(got), 1) {
		t.Fatalf("apexApertureToFNumber(2048) = %v, want +Inf", got)
	}
	if got, want := apexShutterSpeedToSeconds(6), meta.ShutterSpeed(1.0/64.0); math.Abs(float64(got-want)) > 0.000001 {
		t.Fatalf("apexShutterSpeedToSeconds(6) = %v, want %v", got, want)
	}
	if got, want := apexShutterSpeedToSeconds(0), meta.ShutterSpeed(1); got != want {
		t.Fatalf("apexShutterSpeedToSeconds(0) = %v, want %v", got, want)
	}
}

func TestParseExposureBiasSignedRational(t *testing.T) {
	t.Parallel()

	var payload [8]byte
	num := int32(-1)
	utils.LittleEndian.PutUint32(payload[:4], uint32(num))
	utils.LittleEndian.PutUint32(payload[4:8], 3)

	r := NewReader(zerolog.Nop())
	defer r.Close()
	r.Reset(bytes.NewReader(payload[:]))

	tg := tag.NewEntry(
		tag.TagExposureBiasValue,
		tag.TypeSignedRational,
		1,
		0,
		tag.ExifIFD,
		0,
		utils.LittleEndian,
	)

	if got, want := r.parseExposureBias(tg), meta.NewExposureBias(-1, 3); got != want {
		t.Fatalf("parseExposureBias() = %v, want %v", got, want)
	}
}

func TestParseIFD0TagApplicationNotesSkipped(t *testing.T) {
	t.Parallel()

	const payload = "application notes payload"

	r := NewReader(zerolog.Nop())
	defer r.Close()
	r.Reset(bytes.NewReader([]byte(payload)))

	tg := tag.NewEntry(
		tag.TagApplicationNotes,
		tag.TypeUndefined,
		uint32(len(payload)),
		0,
		tag.IFD0,
		0,
		utils.LittleEndian,
	)

	if ok := r.parseIFD0Tag(tg); !ok {
		t.Fatal("parseIFD0Tag(TagApplicationNotes) = false, want true")
	}
	if got := r.po; got != 0 {
		t.Fatalf("reader offset advanced unexpectedly: got %d, want 0", got)
	}
}

func TestParseExifTagUserCommentASCII(t *testing.T) {
	t.Parallel()

	payload := append([]byte("ASCII\x00\x00\x00"), []byte("hello world\x00")...)

	r := NewReader(zerolog.Nop())
	defer r.Close()
	r.Reset(bytes.NewReader(payload))

	tg := tag.NewEntry(
		tag.TagUserComment,
		tag.TypeUndefined,
		uint32(len(payload)),
		0,
		tag.ExifIFD,
		0,
		utils.LittleEndian,
	)

	if ok := r.parseExifTag(tg); !ok {
		t.Fatal("parseExifTag(TagUserComment) = false, want true")
	}
	if got := r.Exif.ExifIFD.UserComment; got != "hello world" {
		t.Fatalf("UserComment = %q, want %q", got, "hello world")
	}
}

func TestParseExifTagUserCommentUnicode(t *testing.T) {
	t.Parallel()

	payload := append([]byte("UNICODE\x00"), []byte{
		'H', 0,
		'i', 0,
		'!', 0,
		0, 0,
	}...)

	r := NewReader(zerolog.Nop())
	defer r.Close()
	r.Reset(bytes.NewReader(payload))

	tg := tag.NewEntry(
		tag.TagUserComment,
		tag.TypeUndefined,
		uint32(len(payload)),
		0,
		tag.ExifIFD,
		0,
		utils.LittleEndian,
	)

	if ok := r.parseExifTag(tg); !ok {
		t.Fatal("parseExifTag(TagUserComment) = false, want true")
	}
	if got := r.Exif.ExifIFD.UserComment; got != "Hi!" {
		t.Fatalf("UserComment = %q, want %q", got, "Hi!")
	}
}
