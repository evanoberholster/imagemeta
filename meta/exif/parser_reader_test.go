package exif

import (
	"bufio"
	"bytes"
	"errors"
	"math"
	"os"
	"path/filepath"
	"testing"

	"github.com/evanoberholster/imagemeta/meta"
	"github.com/evanoberholster/imagemeta/meta/exif/makernote"
	"github.com/evanoberholster/imagemeta/meta/exif/makernote/nikon"
	"github.com/evanoberholster/imagemeta/meta/exif/tag"
	"github.com/evanoberholster/imagemeta/meta/utils"
	"github.com/rs/zerolog"
)

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
	dEmbedded.ByteOrder.PutUint16(embedded[0:2], uint16(nikon.MakerNoteVersion))
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
	utils.LittleEndian.PutUint16(payload[0:2], uint16(nikon.MakerNoteVersion))
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

func TestParseDNGAdobeDataSample(t *testing.T) {
	benchDir := os.Getenv("IMAGEMETA_BENCH_IMAGE_DIR")
	if benchDir == "" {
		benchDir = defaultBenchImageDir
	}

	samplePath := filepath.Join(benchDir, "1.dng")
	if _, err := os.Stat(samplePath); err != nil {
		t.Skipf("sample not found: %s", samplePath)
	}

	f, err := os.Open(samplePath)
	if err != nil {
		t.Fatalf("open %s: %v", samplePath, err)
	}
	defer func() { _ = f.Close() }()

	parsed, err := Parse(f)
	if err != nil {
		t.Fatalf("parse %s: %v", samplePath, err)
	}

	if got := parsed.DNG.AdobeData.RecordCount; got != 1 {
		t.Fatalf("DNG.AdobeData.RecordCount = %d, want 1", got)
	}
	if got := parsed.DNG.AdobeData.MakerNoteOriginalOffset; got != 0x03e4 {
		t.Fatalf("DNG.AdobeData.MakerNoteOriginalOffset = 0x%x, want 0x3e4", got)
	}
	if got := parsed.DNG.AdobeData.MakerNoteRecordLength; got != 68242 {
		t.Fatalf("DNG.AdobeData.MakerNoteRecordLength = %d, want 68242", got)
	}
	if parsed.MakerNote.Canon == nil {
		t.Fatalf("Canon maker-note missing for %s", samplePath)
	}
	if got := parsed.MakerNote.Canon.ImageType; got != "Canon EOS 6D" {
		t.Fatalf("Canon.ImageType = %q, want %q", got, "Canon EOS 6D")
	}
	if got := parsed.MakerNote.Canon.LensModel; got != "EF70-200mm f/2.8L IS II USM" {
		t.Fatalf("Canon.LensModel = %q, want %q", got, "EF70-200mm f/2.8L IS II USM")
	}
	if got := parsed.MakerNote.Canon.TimeInfo.TimeZone; got != 180 {
		t.Fatalf("Canon.TimeInfo.TimeZone = %d, want 180", got)
	}
}

func TestParseIFD0MakeTagTrimsAndNormalizesKnownMake(t *testing.T) {
	t.Parallel()

	r := NewReader(zerolog.Nop())
	defer r.Close()

	raw := append([]byte("NIKON CORPORATION\x00"), bytes.Repeat([]byte{'x'}, 48)...)
	r.Reset(bytes.NewReader(raw))

	e := tag.NewEntry(tag.TagMake, tag.TypeASCII, 18, 0, tag.IFD0, 0, utils.LittleEndian)
	r.Exif.CameraMakeID, r.Exif.IFD0.Make = r.parseMakeTag(e)

	if got := r.Exif.CameraMakeID; got != makernote.CameraMakeNikon {
		t.Fatalf("CameraMakeID = %v, want %v", got, makernote.CameraMakeNikon)
	}
	if got := r.Exif.IFD0.Make; got != makernote.CameraMakeNikon.String() {
		t.Fatalf("IFD0.Make = %q, want %q", got, makernote.CameraMakeNikon.String())
	}
}

func TestParseSubIFDsSingleTypeIFDUsesValueOffset(t *testing.T) {
	t.Parallel()

	r := NewReader(zerolog.Nop())
	defer r.Close()

	tg := tag.NewEntry(tag.TagSubIFDs, tag.TypeIfd, 1, 0x1234, tag.IFD0, 0, utils.LittleEndian)
	r.parseSubIFDs(tg)

	if got := r.Exif.IFD0.subIFDOffsetCount; got != 1 {
		t.Fatalf("SubIFDOffsetCount = %d, want 1", got)
	}
	if got := r.Exif.IFD0.subIFDOffsets[0]; got != 0x1234 {
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

	if got := r.Exif.IFD0.subIFDOffsetCount; got != 1 {
		t.Fatalf("SubIFDOffsetCount = %d, want 1", got)
	}
	if got := r.Exif.IFD0.subIFDOffsets[0]; got != 0x10 {
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
	r.Exif.IFD0.subIFDOffsetCount = uint8(len(r.Exif.IFD0.subIFDOffsets) - 1)

	tg := tag.NewEntry(tag.TagSubIFDs, tag.TypeLong, 4, 0, tag.IFD0, 0, utils.LittleEndian)
	r.parseSubIFDs(tg)

	if got := r.Exif.IFD0.subIFDOffsetCount; got != uint8(len(r.Exif.IFD0.subIFDOffsets)) {
		t.Fatalf("SubIFDOffsetCount = %d, want %d", got, len(r.Exif.IFD0.subIFDOffsets))
	}
	last := len(r.Exif.IFD0.subIFDOffsets) - 1
	if got := r.Exif.IFD0.subIFDOffsets[last]; got != 0x10 {
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

func TestParseExifTagUserCommentHeaderSpacePadded(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		header  []byte
		payload []byte
		bo      utils.ByteOrder
		want    string
	}{
		{
			name:    "ASCII",
			header:  []byte{'A', 'S', 'C', 'I', 'I', ' ', ' ', 0},
			payload: []byte("hello world\x00"),
			bo:      utils.LittleEndian,
			want:    "hello world",
		},
		{
			name:   "Unicode",
			header: []byte{'U', 'N', 'I', 'C', 'O', 'D', 'E', ' '},
			payload: []byte{
				'H', 0,
				'i', 0,
				'!', 0,
				0, 0,
			},
			bo:   utils.LittleEndian,
			want: "Hi!",
		},
		{
			name:    "JIS",
			header:  []byte{'J', 'I', 'S', ' ', 0, ' ', 0, 0},
			payload: []byte("JIS text\x00"),
			bo:      utils.LittleEndian,
			want:    "JIS text",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			r := NewReader(zerolog.Nop())
			defer r.Close()

			data := append(append([]byte{}, tc.header...), tc.payload...)
			r.Reset(bytes.NewReader(data))
			tg := tag.NewEntry(
				tag.TagUserComment,
				tag.TypeUndefined,
				uint32(len(data)),
				0,
				tag.ExifIFD,
				0,
				tc.bo,
			)

			if ok := r.parseExifTag(tg); !ok {
				t.Fatal("parseExifTag(TagUserComment) = false, want true")
			}
			if got := r.Exif.ExifIFD.UserComment; got != tc.want {
				t.Fatalf("UserComment = %q, want %q", got, tc.want)
			}
		})
	}
}
