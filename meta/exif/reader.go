package exif

import (
	"bufio"
	"errors"
	"io"
	"log/slog"
	"sync"
	"time"

	"github.com/evanoberholster/imagemeta/imagetype"
	"github.com/evanoberholster/imagemeta/meta"
	"github.com/evanoberholster/imagemeta/meta/exif/tag"
	"github.com/evanoberholster/imagemeta/meta/isobmff"
	metalog "github.com/evanoberholster/imagemeta/meta/logging"
	"github.com/evanoberholster/imagemeta/meta/utils"
)

const (
	maxTagCount          = 256
	parseProbeReaderSize = 64
	parseTiffReaderSize  = 4096
)

type eofReader struct{}

func (eofReader) Read(_ []byte) (int, error) {
	return 0, io.EOF
}

var (
	pooledEOFReader = eofReader{}
	parseReaderPool = sync.Pool{
		New: func() any {
			return &Reader{Mixin: metalog.NewComponentMixin(metalog.GetLogger(), "exif")}
		},
	}
	isobmffReaderPool = sync.Pool{
		New: func() any {
			return isobmff.NewReader(nil, nil, nil, nil)
		},
	}
	tiffBufioReaderPool = utils.NewBufioReaderPool(parseTiffReaderSize, pooledEOFReader)
)

// Reader reads and parses EXIF IFD trees.
type Reader struct {
	metalog.Mixin
	reader      utils.BufferedReader
	ownedReader *bufio.Reader
	state       *state

	Exif Exif

	afInfoDecodeOptions AFInfoDecodeOptions

	po               uint32
	exifLength       uint32
	firstIFDOffset   uint32
	tiffHeaderOffset uint32
}

// NewReader creates an EXIF reader. Call Close when done.
func NewReader(l *slog.Logger, opts ...ReaderOption) *Reader {
	s, ok := statePool.Get().(*state)
	if !ok || s == nil {
		s = new(state)
	}
	s.reset()
	r := &Reader{
		Mixin:               metalog.NewComponentMixin(l, "exif"),
		state:               s,
		afInfoDecodeOptions: AFInfoDecodeAll,
	}
	applyReaderOptions(r, opts)
	return r
}

func (r *Reader) resetOffsets() {
	r.po = 0
	r.exifLength = 0
	r.firstIFDOffset = 0
	r.tiffHeaderOffset = 0
}

func (r *Reader) resetDecodeState(resetExif bool) {
	if r.state != nil {
		r.state.reset()
	}
	if resetExif {
		r.Exif = Exif{}
	}
	r.resetOffsets()
}

// AcquirePooledReader gets a Reader from the shared pool.
func AcquirePooledReader(l *slog.Logger) *Reader {
	r, ok := parseReaderPool.Get().(*Reader)
	if !ok || r == nil {
		r = &Reader{Mixin: metalog.NewComponentMixin(l, "exif")}
	}
	r.SetLogger(l)
	if r.state == nil {
		s, stateOK := statePool.Get().(*state)
		if !stateOK || s == nil {
			s = new(state)
		}
		r.state = s
	}
	r.releaseOwnedReader()
	r.reader = nil
	r.resetDecodeState(true)
	r.afInfoDecodeOptions = AFInfoDecodeAll
	return r
}

func acquirePooledReader(l *slog.Logger) *Reader {
	return AcquirePooledReader(l)
}

func releasePooledReader(r *Reader) {
	ReleasePooledReader(r)
}

// ReleasePooledReader returns a Reader to the shared pool.
func ReleasePooledReader(r *Reader) {
	if r == nil {
		return
	}
	if r.state != nil {
		r.state.reset()
		statePool.Put(r.state)
		r.state = nil
	}
	r.releaseOwnedReader()
	r.reader = nil
	r.Exif = Exif{}
	r.afInfoDecodeOptions = AFInfoDecodeAll
	r.resetOffsets()
	parseReaderPool.Put(r)
}

func acquirePooledISOBMFFReader(src io.Reader, exifReader meta.ExifReader) *isobmff.Reader {
	r, ok := isobmffReaderPool.Get().(*isobmff.Reader)
	if !ok || r == nil {
		return isobmff.NewReader(src, exifReader, nil, nil)
	}
	r.Reset(src, exifReader, nil, nil)
	return r
}

func releasePooledISOBMFFReader(r *isobmff.Reader) {
	if r == nil {
		return
	}
	r.Close()
	isobmffReaderPool.Put(r)
}

// Close returns parser state to the pool.
func (r *Reader) Close() {
	r.releaseOwnedReader()
	if r.state != nil {
		r.state.reset()
		statePool.Put(r.state)
		r.state = nil
	}
}

// Reset prepares the reader for a new decode operation.
func (r *Reader) Reset(reader io.Reader) {
	r.setReader(reader)
	r.resetDecodeState(true)
}

// Parse scans a TIFF header and parses EXIF into Exif.
func Parse(rs io.ReadSeeker) (Exif, error) {
	return ParseWithReaderOptions(rs)
}

// ParseWithReaderOptions scans a TIFF header and parses EXIF into Exif
// with optional reader-specific parse options.
func ParseWithReaderOptions(rs io.ReadSeeker, opts ...ReaderOption) (Exif, error) {
	br := tiffBufioReaderPool.Acquire(rs)
	defer tiffBufioReaderPool.Release(br)

	probe := scanParseProbe(br)
	if probe.imageType.IsISOBMFF() {
		return parseCR3FromReader(br, opts...)
	}

	reader := acquirePooledReader(metalog.GetLogger())
	defer releasePooledReader(reader)
	applyReaderOptions(reader, opts)

	header := probe.panHeader
	if probe.hasPanasonicHeader {
		// Panasonic RW2 uses an alternate TIFF signature at byte 0.
		// Use the probe-derived header directly and decode in a single pass.
	} else {
		var err error
		header, err = ScanTiffHeader(br, probe.imageType)
		if err != nil {
			return Exif{}, err
		}
	}

	return reader.Exif, reader.DecodeTiff(br, header)
}

type parseProbeInfo struct {
	imageType          imagetype.ImageType
	panHeader          meta.ExifHeader
	hasPanasonicHeader bool
}

// scanParseProbe scans input bytes once to detect the image type and optional
// Panasonic RW2 alternate TIFF header.
func scanParseProbe(br *bufio.Reader) parseProbeInfo {
	var out parseProbeInfo

	buf, err := br.Peek(parseProbeReaderSize)
	if err != nil && len(buf) == 0 {
		return out
	}

	if it, detectErr := imagetype.Buf(buf); detectErr == nil {
		out.imageType = it
	}

	if len(buf) < 8 {
		return out
	}

	switch {
	case buf[0] == 'I' && buf[1] == 'I' && buf[2] == 0x55 && buf[3] == 0x00:
		offset := utils.LittleEndian.Uint32(buf[4:8])
		if offset >= 8 {
			out.panHeader = meta.NewExifHeader(utils.LittleEndian, offset, 0, 0, imagetype.ImagePanaRAW)
			out.hasPanasonicHeader = true
		}
	case buf[0] == 'M' && buf[1] == 'M' && buf[2] == 0x00 && buf[3] == 0x55:
		offset := utils.BigEndian.Uint32(buf[4:8])
		if offset >= 8 {
			out.panHeader = meta.NewExifHeader(utils.BigEndian, offset, 0, 0, imagetype.ImagePanaRAW)
			out.hasPanasonicHeader = true
		}
	}
	return out
}

func parseCR3FromReader(src io.Reader, opts ...ReaderOption) (Exif, error) {
	reader := acquirePooledReader(metalog.GetLogger())
	defer releasePooledReader(reader)
	applyReaderOptions(reader, opts)

	bmr := acquirePooledISOBMFFReader(src, reader.DecodeIfdAppend)
	defer releasePooledISOBMFFReader(bmr)

	if err := bmr.ReadFTYP(); err != nil {
		return Exif{}, err
	}
	for {
		err := bmr.ReadMetadata()
		if err == nil {
			continue
		}
		if errors.Is(err, io.EOF) {
			break
		}
		return reader.Exif, err
	}
	return reader.Exif, nil
}

// initDecode initializes reader state for one decode operation.
func (r *Reader) initDecode(reader io.Reader, header meta.ExifHeader, resetExif bool) {
	if resetExif {
		r.Reset(reader)
	} else {
		r.setReader(reader)
		r.resetDecodeState(false)
	}

	r.Exif.ImageType = header.ImageType
	r.firstIFDOffset = header.FirstIfdOffset
	r.tiffHeaderOffset = header.TiffHeaderOffset
	// ExifLength == 0 means unknown length (for full TIFF/RAW streams).
	// Keep parsing unbounded by length in that case to avoid false short-read
	// failures on large maker-note offsets.
	r.exifLength = header.ExifLength
}

func (r *Reader) setReader(reader io.Reader) {
	r.releaseOwnedReader()
	if reader == nil {
		r.reader = nil
		return
	}
	if br, ok := reader.(utils.BufferedReader); ok {
		r.reader = br
		return
	}
	br := tiffBufioReaderPool.Acquire(reader)
	r.ownedReader = br
	r.reader = br
}

func (r *Reader) releaseOwnedReader() {
	if r.ownedReader == nil {
		return
	}
	tiffBufioReaderPool.Release(r.ownedReader)
	r.ownedReader = nil
}

// rootDirectory resolves the root IFD for a header.
func (r *Reader) rootDirectory(header meta.ExifHeader) (tag.Directory, bool) {
	rootType := header.FirstIfd
	if !rootType.IsValid() {
		if r.WarnEnabled() {
			r.readerLogContext(r.Warn(3)).
				Object("header", header).
				Str("ifd", header.FirstIfd.String()).
				Uint8("ifdID", uint8(header.FirstIfd)).
				Msg("unsupported root ifd type")
		}
		return tag.Directory{}, false
	}
	return tag.NewDirectory(header.ByteOrder, rootType, 0, header.FirstIfdOffset, 0), true
}

// decodeRootIFD decodes a root IFD with optional pre-positioned reader state.
func (r *Reader) decodeRootIFD(reader io.Reader, header meta.ExifHeader, resetExif bool, positioned bool) (err error) {
	r.initDecode(reader, header, resetExif)
	infoEnabled := r.InfoEnabled()
	var start time.Time
	if infoEnabled {
		start = time.Now()
		r.Info(3).
			Object("header", header).
			Bool("append", !resetExif).
			Bool("positioned", positioned).
			Msg("starting exif decode")
		defer func() {
			ev := r.Info(3).
				Str("imageType", r.Exif.ImageType.String()).
				Uint32("readerOffset", r.po).
				Dur("elapsed", time.Since(start))
			if err != nil {
				ev.Err(err)
			}
			ev.Msg("finished exif decode")
		}()
	}
	if positioned {
		r.po = header.FirstIfdOffset
	} else if err = r.discard(int(header.FirstIfdOffset)); err != nil {
		if r.WarnEnabled() {
			r.readerLogContext(r.Warn(3)).
				Err(err).
				Object("header", header).
				Msg("failed seeking to root ifd")
		}
		return err
	}
	root, ok := r.rootDirectory(header)
	if !ok {
		return nil
	}
	if err = r.readDirectory(root, true); err != nil {
		if r.WarnEnabled() {
			r.directoryLogContext(r.Warn(3), root).
				Err(err).
				Object("header", header).
				Msg("failed reading root ifd")
		}
		return err
	}
	return nil
}

// DecodeTiff parses EXIF from a TIFF-like stream.
func (r *Reader) DecodeTiff(reader io.Reader, header meta.ExifHeader) error {
	return r.decodeRootIFD(reader, header, true, false)
}

// DecodeJPEGIfd parses EXIF from JPEG APP1 payload bytes.
func (r *Reader) DecodeJPEGIfd(reader io.Reader, header meta.ExifHeader) error {
	return r.decodeRootIFD(reader, header, true, false)
}

// DecodeIfd parses EXIF from an already positioned TIFF IFD stream.
func (r *Reader) DecodeIfd(reader io.Reader, header meta.ExifHeader) error {
	return r.decodeRootIFD(reader, header, true, true)
}

// DecodeIfdAppend parses a TIFF IFD payload and merges results into the current Exif value.
//
// Use this for containers (for example CR3/ISOBMFF) that deliver metadata across
// multiple IFD payloads.
func (r *Reader) DecodeIfdAppend(reader io.Reader, header meta.ExifHeader) error {
	return r.decodeRootIFD(reader, header, false, true)
}

// readDirectory reads data from the underlying stream or parser buffers.
func (r *Reader) readDirectory(directory tag.Directory, drainQueue bool) error {
	if !directory.Type.IsValid() {
		return nil
	}
	infoEnabled := r.InfoEnabled()
	tagCount, err := r.readUint16(directory)
	if err != nil {
		if r.WarnEnabled() {
			r.directoryLogContext(r.Warn(3), directory).
				Err(err).
				Msg("failed reading exif tag count")
		}
		return err
	}
	if tagCount > maxTagCount {
		if r.WarnEnabled() {
			r.directoryLogContext(r.Warn(3), directory).
				Uint16("tagCount", tagCount).
				Uint16("tagCountMax", maxTagCount).
				Msg("exif tag count exceeds parser limit")
		}
		return nil
	}

	if err = r.parseDirectoryTagHeadersBulkTrusted(directory, tagCount); err != nil {
		if r.WarnEnabled() {
			r.directoryLogContext(r.Warn(3), directory).
				Err(err).
				Uint16("tagCount", tagCount).
				Msg("failed reading exif tag headers")
		}
		return err
	}

	nextIFDOffset, err := r.readUint32(directory)
	if err != nil {
		if r.WarnEnabled() {
			r.directoryLogContext(r.Warn(3), directory).
				Err(err).
				Msg("failed reading next ifd offset")
		}
		return err
	}
	if _, ok := directory.Type.NextRootIFD(); ok && nextIFDOffset != 0 {
		r.addTag(tag.NewEntry(tag.TagNextIFD, tag.TypeIfd, 1, nextIFDOffset, directory.Type, directory.Index+1, directory.ByteOrder))
	}

	if !drainQueue {
		if infoEnabled {
			r.infoDirectoryLogContext(r.Info(3), directory).
				Uint16("tagCount", tagCount).
				Uint32("nextIFDOffset", nextIFDOffset).
				Msg("queued exif directory tags")
		}
		return nil
	}

	r.state.sortAll()

	if err = r.drainQueuedTags(); err != nil {
		return err
	}

	if infoEnabled {
		r.infoDirectoryLogContext(r.Info(3), directory).
			Uint16("tagCount", tagCount).
			Uint32("nextIFDOffset", nextIFDOffset).
			Msg("parsed exif directory")
	}
	return nil
}

func (r *Reader) drainQueuedTags() error {
	for t := r.state.currentTag(); r.state.validTag(); t = r.state.advanceTag() {
		switch {
		case t.IsIfd():
			if t.IfdType == tag.IFD0 {
				switch t.ID {
				case tag.TagExifIFDPointer:
					r.Exif.IFD0.exifIfdPointer = t.ValueOffset
				case tag.TagGPSIFDPointer:
					r.Exif.IFD0.gpsIfdPointer = t.ValueOffset
				}
			}
			if err := r.seekToTag(t); err != nil {
				return err
			}
			r.state.resetPosition()
			child := t.ChildDirectory()
			if child.Type == tag.MakerNoteIFD {
				if err := r.readMakerNoteDirectory(t, child); err != nil {
					return err
				}
				r.state.sortUnread()
				continue
			}
			if child.Type.IsValid() {
				if err := r.readDirectory(child, false); err != nil {
					return err
				}
			}
			r.state.sortUnread()
		case t.ID == tag.TagSubIFDs && t.IfdType == tag.IFD0:
			r.parseSubIFDs(t)
			r.state.sortUnread()
		default:
			r.parseTag(t)
			r.state.sortUnread()
		}
	}
	return nil
}

// parseDirectoryTagHeadersPerEntry decodes tag headers by reading 12 bytes per entry.
func (r *Reader) parseDirectoryTagHeadersPerEntry(directory tag.Directory, tagCount uint16) error {
	warnEnabled := r.WarnEnabled()
	for i := 0; i < int(tagCount); i++ {
		headerBuf, readErr := r.fastRead(12)
		if readErr != nil {
			return readErr
		}
		if len(headerBuf) < 12 {
			return io.ErrUnexpectedEOF
		}
		t, parseErr := tagFromBuffer(directory, headerBuf)
		if parseErr != nil {
			if warnEnabled {
				tagID := tag.ID(directory.ByteOrder.Uint16(headerBuf[:2]))
				tagTypeValue, ok := meta.SafecastUint16ToUint8(directory.ByteOrder.Uint16(headerBuf[2:4]))
				if !ok {
					tagTypeValue = 0
				}
				tagType := tag.Type(tagTypeValue)
				unitCount := directory.ByteOrder.Uint32(headerBuf[4:8])
				valueOffset := directory.ByteOrder.Uint32(headerBuf[8:12])
				r.rawTagHeaderLogContext(r.Warn(3), directory, i, tagID, tagType, unitCount, valueOffset).
					Err(parseErr).
					Msg("invalid exif tag header")
			}
			continue
		}
		if t.IsEmbedded() {
			r.parseTag(t)
			continue
		}
		r.addTag(t)
	}
	return nil
}

// parseDirectoryTagHeadersBulk decodes all tag headers from a single contiguous read.
func (r *Reader) parseDirectoryTagHeadersBulk(directory tag.Directory, tagCount uint16) error {
	// Retained for test parity with the non-trusted parser path.
	return r.parseDirectoryTagHeadersBulkTrusted(directory, tagCount)
}

// parseDirectoryTagHeadersBulkTrusted inlines tag decode for trusted bulk buffers.
//
// This variant avoids tagFromBuffer/NewEntry/IsEmbedded call overhead in favor of
// direct field extraction and a local embedded-size switch.
func (r *Reader) parseDirectoryTagHeadersBulkTrusted(directory tag.Directory, tagCount uint16) error {
	total := int(tagCount) * 12
	if total <= 0 {
		return nil
	}
	if total > len(r.state.buf) {
		return r.parseDirectoryTagHeadersPerEntry(directory, tagCount)
	}

	raw, err := r.fastRead(total)
	if err != nil {
		return err
	}
	if len(raw) < total {
		return io.ErrUnexpectedEOF
	}

	warnEnabled := r.WarnEnabled()

	byteOrder := directory.ByteOrder
	directoryType := directory.Type
	directoryIndex := directory.Index
	baseOffset := directory.BaseOffset

	for pos := 0; pos < total; pos += 12 {
		index := pos / 12
		h := raw[pos : pos+12]
		tagID := tag.ID(byteOrder.Uint16(h[:2]))
		tagTypeValue, ok := meta.SafecastUint16ToUint8(byteOrder.Uint16(h[2:4]))
		if !ok {
			tagTypeValue = 0
		}
		tagType := tag.Type(tagTypeValue)
		unitCount := byteOrder.Uint32(h[4:8])
		valueOffset := byteOrder.Uint32(h[8:12])

		tagType = tag.NormalizeType(directoryType, tagID, tagType)

		if !tagType.IsValid() {
			if warnEnabled {
				r.rawTagHeaderLogContext(r.Warn(3), directory, index, tagID, tagType, unitCount, valueOffset).
					Err(tag.ErrTagTypeNotValid).
					Msg("invalid exif tag header")
			}
			continue
		}

		t := tag.Entry{
			ValueOffset: valueOffset,
			UnitCount:   unitCount,
			ID:          tagID,
			Type:        tagType,
			IfdType:     directoryType,
			IfdIndex:    directoryIndex,
			ByteOrder:   byteOrder,
		}

		if t.IsEmbedded() {
			r.parseTag(t)
			continue
		} else {
			t.ValueOffset += baseOffset
		}
		r.addTag(t)
	}
	return nil
}
