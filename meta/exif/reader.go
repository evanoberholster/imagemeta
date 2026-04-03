package exif

import (
	"bufio"
	"errors"
	"io"
	"sync"

	oldifd "github.com/evanoberholster/imagemeta/exif2/ifds"
	"github.com/evanoberholster/imagemeta/imagetype"
	"github.com/evanoberholster/imagemeta/meta"
	"github.com/evanoberholster/imagemeta/meta/exif/ifd"
	"github.com/evanoberholster/imagemeta/meta/exif/tag"
	"github.com/evanoberholster/imagemeta/meta/isobmff"
	"github.com/evanoberholster/imagemeta/meta/tiff"
	"github.com/evanoberholster/imagemeta/meta/utils"
	"github.com/rs/zerolog"
)

const (
	defaultExifLength    = 4 * 1024 * 1024
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
			return &Reader{loggerMixin: newLoggerMixin(Logger)}
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
	loggerMixin
	reader      utils.BufferedReader
	ownedReader *bufio.Reader
	state       *state

	Exif Exif

	afInfoDecodeOptions AFInfoDecodeOptions

	po             uint32
	exifLength     uint32
	firstIFDOffset uint32
}

// NewReader creates an EXIF reader. Call Close when done.
func NewReader(l zerolog.Logger, opts ...ReaderOption) *Reader {
	s, ok := statePool.Get().(*state)
	if !ok || s == nil {
		s = new(state)
	}
	s.reset()
	r := &Reader{
		loggerMixin:         newLoggerMixin(l),
		state:               s,
		afInfoDecodeOptions: AFInfoDecodeAll,
	}
	applyReaderOptions(r, opts)
	return r
}

func acquirePooledReader(l zerolog.Logger) *Reader {
	r, ok := parseReaderPool.Get().(*Reader)
	if !ok || r == nil {
		r = &Reader{}
	}
	r.loggerMixin = newLoggerMixin(l)
	if r.state == nil {
		s, stateOK := statePool.Get().(*state)
		if !stateOK || s == nil {
			s = new(state)
		}
		r.state = s
	}
	r.state.reset()
	r.releaseOwnedReader()
	r.reader = nil
	r.Exif = Exif{}
	r.afInfoDecodeOptions = AFInfoDecodeAll
	r.po = 0
	r.exifLength = 0
	r.firstIFDOffset = 0
	return r
}

func releasePooledReader(r *Reader) {
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
	r.po = 0
	r.exifLength = 0
	r.firstIFDOffset = 0
	parseReaderPool.Put(r)
}

func acquirePooledISOBMFFReader(src io.Reader, exifReader isobmff.ExifReader) *isobmff.Reader {
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
	r.state.reset()
	r.Exif = Exif{}
	r.po = 0
	r.exifLength = 0
	r.firstIFDOffset = 0
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

	reader := acquirePooledReader(Logger)
	defer releasePooledReader(reader)
	applyReaderOptions(reader, opts)

	header := probe.panHeader
	if probe.hasPanasonicHeader {
		// Panasonic RW2 uses an alternate TIFF signature at byte 0.
		// Use the probe-derived header directly and decode in a single pass.
	} else {
		var err error
		header, err = tiff.ScanTiffHeader(br, probe.imageType)
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
	reader := acquirePooledReader(Logger)
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

// mergeIFD0CoreFields merges parsed values into the destination EXIF model.
func mergeIFD0CoreFields(dst *Exif, src Exif) {
	if dst == nil {
		return
	}
	if dst.IFD0.SubfileType == 0 && src.IFD0.SubfileType != 0 {
		dst.IFD0.SubfileType = src.IFD0.SubfileType
	}
	if dst.IFD0.ImageWidth == 0 && src.IFD0.ImageWidth != 0 {
		dst.IFD0.ImageWidth = src.IFD0.ImageWidth
	}
	if dst.IFD0.ImageHeight == 0 && src.IFD0.ImageHeight != 0 {
		dst.IFD0.ImageHeight = src.IFD0.ImageHeight
	}
	if dst.IFD0.StripOffsets == 0 && src.IFD0.StripOffsets != 0 {
		dst.IFD0.StripOffsets = src.IFD0.StripOffsets
	}
	if dst.IFD0.RowsPerStrip == 0 && src.IFD0.RowsPerStrip != 0 {
		dst.IFD0.RowsPerStrip = src.IFD0.RowsPerStrip
	}
	if dst.IFD0.StripByteCounts == 0 && src.IFD0.StripByteCounts != 0 {
		dst.IFD0.StripByteCounts = src.IFD0.StripByteCounts
	}
	if dst.IFD0.TileWidth == 0 && src.IFD0.TileWidth != 0 {
		dst.IFD0.TileWidth = src.IFD0.TileWidth
	}
	if dst.IFD0.TileLength == 0 && src.IFD0.TileLength != 0 {
		dst.IFD0.TileLength = src.IFD0.TileLength
	}
	if dst.IFD0.TileOffsets == 0 && src.IFD0.TileOffsets != 0 {
		dst.IFD0.TileOffsets = src.IFD0.TileOffsets
	}
	if dst.IFD0.TileByteCounts == 0 && src.IFD0.TileByteCounts != 0 {
		dst.IFD0.TileByteCounts = src.IFD0.TileByteCounts
	}

	dst.PanasonicRaw = src.PanasonicRaw

	for i := range dst.ifdBitset {
		dst.ifdBitset[i] |= src.ifdBitset[i]
	}
	n := int(src.highTagCount)
	if n > len(src.highTagIDs) {
		n = len(src.highTagIDs)
	}
	for i := 0; i < n; i++ {
		dst.markTagParsed(src.highTagIDs[i])
	}
	dst.MakerNote.MergeParsedTags(src.MakerNote)
}

// mapIfdType maps one metadata representation to another.
func mapIfdType(first oldifd.IfdType) ifd.Type {
	switch first {
	case oldifd.IFD0:
		return ifd.IFD0
	case oldifd.SubIFD:
		return ifd.SubIFD0
	case oldifd.ExifIFD:
		return ifd.ExifIFD
	case oldifd.GPSIFD:
		return ifd.GPSIFD
	case oldifd.MknoteIFD:
		return ifd.MakerNoteIFD
	case oldifd.SubIfd0:
		return ifd.SubIFD0
	case oldifd.SubIfd1:
		return ifd.SubIFD1
	case oldifd.SubIfd2:
		return ifd.SubIFD2
	case oldifd.SubIfd3:
		return ifd.SubIFD3
	case oldifd.SubIfd4:
		return ifd.SubIFD4
	case oldifd.SubIfd5:
		return ifd.SubIFD5
	case oldifd.SubIfd6:
		return ifd.SubIFD6
	case oldifd.SubIfd7:
		return ifd.SubIFD7
	}

	// Fall back to direct casting for already-aligned values.
	mapped := ifd.Type(first)
	if mapped >= ifd.IFD0 && mapped <= ifd.SubIFD7 {
		return mapped
	}
	return ifd.Unknown
}

// initDecode initializes reader state for one decode operation.
func (r *Reader) initDecode(reader io.Reader, header meta.ExifHeader, resetExif bool) {
	if resetExif {
		r.Reset(reader)
	} else {
		r.setReader(reader)
		r.state.reset()
		r.po = 0
		r.exifLength = 0
		r.firstIFDOffset = 0
	}

	r.Exif.ImageType = header.ImageType
	r.firstIFDOffset = header.FirstIfdOffset
	if header.ExifLength == 0 {
		r.exifLength = defaultExifLength
	} else {
		r.exifLength = header.ExifLength
	}
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
func (r *Reader) rootDirectory(header meta.ExifHeader) (ifd.Directory, bool) {
	rootType := mapIfdType(header.FirstIfd)
	if !rootType.IsValid() {
		if r.warnEnabled() {
			r.warn().Str("ifd", header.FirstIfd.String()).Uint8("ifdID", uint8(header.FirstIfd)).Msg("unsupported root ifd type")
		}
		return ifd.Directory{}, false
	}
	return ifd.New(header.ByteOrder, rootType, 0, header.FirstIfdOffset, 0), true
}

// decodeRootIFD decodes a root IFD with optional pre-positioned reader state.
func (r *Reader) decodeRootIFD(reader io.Reader, header meta.ExifHeader, resetExif bool, positioned bool) error {
	r.initDecode(reader, header, resetExif)
	if positioned {
		r.po = header.FirstIfdOffset
	} else if err := r.discard(int(header.FirstIfdOffset)); err != nil {
		return err
	}
	root, ok := r.rootDirectory(header)
	if !ok {
		return nil
	}
	return r.readDirectory(root, true)
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
func (r *Reader) readDirectory(directory ifd.Directory, drainQueue bool) error {
	if !directory.Type.IsValid() {
		return nil
	}
	tagCount, err := r.readUint16(directory)
	if err != nil {
		return err
	}
	if tagCount > maxTagCount {
		if r.warnEnabled() {
			r.warn().Str("ifd", directory.String()).Uint16("tagCount", tagCount).Msg("exif tag count exceeds parser limit")
		}
		return nil
	}

	if err = r.parseDirectoryTagHeadersBulkTrusted(directory, tagCount); err != nil {
		return err
	}

	nextIFDOffset, err := r.readUint32(directory)
	if err != nil {
		return err
	}
	if _, ok := directory.Type.NextRootIFD(); ok && nextIFDOffset != 0 {
		r.addTag(tag.NewEntry(tag.TagNextIFD, tag.TypeIfd, 1, nextIFDOffset, directory.Type, directory.Index+1, directory.ByteOrder))
	}

	if !drainQueue {
		return nil
	}

	r.state.sortAll()

	for t := r.state.currentTag(); r.state.validTag(); t = r.state.advanceTag() {
		switch {
		case t.IsIfd():
			if t.IfdType == ifd.IFD0 {
				switch t.ID {
				case tag.TagExifIFDPointer:
					r.Exif.IFD0.ExifIFDPointer = t.ValueOffset
					r.Exif.markTagParsed(uint16(t.ID))
				case tag.TagGPSIFDPointer:
					r.Exif.IFD0.GPSIFDPointer = t.ValueOffset
					r.Exif.markTagParsed(uint16(t.ID))
				}
			}
			if err = r.seekToTag(t); err != nil {
				return err
			}
			r.state.resetPosition()
			child := t.ChildDirectory()
			if child.Type == ifd.MakerNoteIFD {
				if err = r.readMakerNoteDirectory(t, child); err != nil && r.warnEnabled() {
					r.warn().Err(err).Str("ifd", child.String()).Msg("failed parsing maker-note ifd")
				}
				r.Exif.markTagParsed(uint16(t.ID))
				r.state.sortUnread()
				continue
			}
			if child.Type.IsValid() {
				if err = r.readDirectory(child, false); err != nil {
					if r.warnEnabled() {
						r.warn().Err(err).Str("ifd", child.String()).Msg("failed parsing child ifd")
					}
				}
			}
			r.state.sortUnread()
		case t.ID == tag.TagSubIFDs && t.IfdType == ifd.IFD0:
			r.parseSubIFDs(t)
			r.Exif.markTagParsed(uint16(t.ID))
			r.state.sortUnread()
		default:
			r.parseTag(t)
			r.state.sortUnread()
		}
	}
	return nil
}

// parseDirectoryTagHeadersPerEntry decodes tag headers by reading 12 bytes per entry.
func (r *Reader) parseDirectoryTagHeadersPerEntry(directory ifd.Directory, tagCount uint16) error {
	warnEnabled := r.warnEnabled()
	ifdName := ""
	if warnEnabled {
		ifdName = directory.String()
	}
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
				r.warn().Err(parseErr).Str("ifd", ifdName).Send()
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
func (r *Reader) parseDirectoryTagHeadersBulk(directory ifd.Directory, tagCount uint16) error {
	total := int(tagCount) * 12
	if total <= 0 {
		return nil
	}
	// Fallback preserves behavior if parser limits are increased above read buffer capacity.
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

	warnEnabled := r.warnEnabled()
	ifdName := ""
	if warnEnabled {
		ifdName = directory.String()
	}
	for pos := 0; pos < total; pos += 12 {
		t, parseErr := tagFromBuffer(directory, raw[pos:pos+12])
		if parseErr != nil {
			if warnEnabled {
				r.warn().Err(parseErr).Str("ifd", ifdName).Send()
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

// parseDirectoryTagHeadersBulkTrusted inlines tag decode for trusted bulk buffers.
//
// This variant avoids tagFromBuffer/NewEntry/IsEmbedded call overhead in favor of
// direct field extraction and a local embedded-size switch.
func (r *Reader) parseDirectoryTagHeadersBulkTrusted(directory ifd.Directory, tagCount uint16) error {
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

	warnEnabled := r.warnEnabled()
	ifdName := ""
	if warnEnabled {
		ifdName = directory.String()
	}

	byteOrder := directory.ByteOrder
	directoryType := directory.Type
	directoryIndex := directory.Index
	baseOffset := directory.BaseOffset

	for pos := 0; pos < total; pos += 12 {
		h := raw[pos : pos+12]
		tagID := tag.ID(byteOrder.Uint16(h[:2]))
		tagType := tag.Type(byteOrder.Uint16(h[2:4]))
		unitCount := byteOrder.Uint32(h[4:8])
		valueOffset := byteOrder.Uint32(h[8:12]) + baseOffset

		if tagType.Is(tag.TypeLong) || tagType.Is(tag.TypeUndefined) {
			switch directoryType {
			case ifd.IFD0:
				if tagID == tag.TagExifIFDPointer || tagID == tag.TagGPSIFDPointer {
					tagType = tag.TypeIfd
				}
			case ifd.ExifIFD:
				if tagID == tag.TagMakerNote {
					tagType = tag.TypeIfd
				}
			}
		}

		if !tagType.IsValid() {
			if warnEnabled {
				r.warn().Err(tag.ErrTagTypeNotValid).Str("ifd", ifdName).Send()
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
		}
		r.addTag(t)
	}
	return nil
}
