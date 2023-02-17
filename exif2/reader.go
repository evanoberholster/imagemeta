package exif2

import (
	"io"

	"github.com/evanoberholster/imagemeta/exif2/ifds"
	"github.com/evanoberholster/imagemeta/exif2/ifds/exififd"
	"github.com/evanoberholster/imagemeta/exif2/tag"
	"github.com/evanoberholster/imagemeta/imagetype"
	"github.com/evanoberholster/imagemeta/meta"
	"github.com/evanoberholster/imagemeta/tiff"
	"github.com/rs/zerolog"
)

func Decode(r io.ReadSeeker) (Exif, error) {
	header, err := tiff.ScanTiffHeader(r, imagetype.ImageUnknown)
	if err != nil {
		return Exif{}, err
	}

	ir := NewIfdReader(r)
	defer ir.Close()

	r.Seek(int64(header.TiffHeaderOffset)+int64(header.FirstIfdOffset), 0)
	if err := ir.ReadIfd0(header); err != nil {
		return ir.Exif, err
	}

	return ir.Exif, nil
}

func (ir *ifdReader) DecodeTiff(_ io.Reader, h meta.ExifHeader) error {
	ir.buffer.clear()
	// Log Header Info
	if ir.logLevelInfo() {
		ir.logger.Info().Str("imageType", h.ImageType.String()).Uint32("tiffHeader", h.TiffHeaderOffset).Uint32("firstIfdOffset", h.FirstIfdOffset).Uint32("exifLength", h.ExifLength).Send()
	}
	ir.Exif.ImageType = h.ImageType
	ir.exifLength = 4 * 1024 * 1024 // Max size is 4 MB
	if err := ir.discard(int(h.FirstIfdOffset)); err != nil {
		return err
	}
	return ir.readIfd(ifds.NewIFD(h.ByteOrder, ifds.IfdType(h.FirstIfd), 0, ir.tiffHeaderOffset))
}

func (ir *ifdReader) DecodeJPEGIfd(r io.Reader, h meta.ExifHeader) (err error) {
	// Log Header Info
	if ir.logLevelInfo() {
		ir.logger.Info().Str("imageType", h.ImageType.String()).Uint32("tiffHeader", h.TiffHeaderOffset).Uint32("firstIfdOffset", h.FirstIfdOffset).Uint32("exifLength", h.ExifLength).Send()
	}
	ir.ResetReader(r)
	ir.Exif.ImageType = h.ImageType
	ir.exifLength = h.ExifLength
	if err = ir.discard(int(h.FirstIfdOffset)); err != nil {
		if ir.logLevelError() {
			ir.logError(err).Send()
		}
	}
	if err := ir.readIfd(ifds.NewIFD(h.ByteOrder, ifds.IfdType(h.FirstIfd), 0, ir.tiffHeaderOffset)); err != nil {
		return err
	}
	err = ir.discard(int(ir.exifLength) - int(ir.po))
	return err
}

func (ir *ifdReader) DecodeIfd(r io.Reader, h meta.ExifHeader) error {
	// Log Header Info
	if ir.logLevelInfo() {
		ir.logger.Info().Str("imageType", h.ImageType.String()).Uint32("tiffHeader", h.TiffHeaderOffset).Uint32("firstIfdOffset", h.FirstIfdOffset).Uint32("exifLength", h.ExifLength).Send()
	}
	ir.ResetReader(r)
	ir.Exif.ImageType = h.ImageType
	ir.exifLength = h.ExifLength
	ir.po = h.FirstIfdOffset
	if err := ir.readIfd(ifds.NewIFD(h.ByteOrder, ifds.IfdType(h.FirstIfd), 0, ir.tiffHeaderOffset)); err != nil {
		return err
	}
	return ir.discard(int(ir.exifLength) - int(ir.po))
}

func NewIfdReader(r io.Reader) ifdReader {
	ir := ifdReader{
		reader: r,
		po:     0,
		buffer: bufferPool.Get().(*buffer),
		logger: Logger,
	}
	ir.buffer.clear()
	return ir
}

func (ir *ifdReader) ResetReader(r io.Reader) {
	ir.buffer.clear()
	ir.reader = r
}

func (ir *ifdReader) SetCustomTagParser(fn TagParserFn) {
	ir.customTagParser = fn
}

// Close closes an ifdReader. Should be called with defer following a newIfdReader
func (ir *ifdReader) Close() {
	bufferPool.Put(ir.buffer)
}

// ifdReader reads, decodes, and parses tags from an io.Reader
type ifdReader struct {
	logger           zerolog.Logger
	reader           io.Reader
	customTagParser  TagParserFn
	buffer           *buffer
	Exif             Exif
	po               uint32
	tiffHeaderOffset uint32
	firstIfdOffset   uint32
	exifLength       uint32
}

func (ir *ifdReader) ReadIfd0(header meta.ExifHeader) error {
	// Log Header Info
	if ir.logLevelInfo() {
		ir.logInfo().Object("header", header).Send()
	}

	ir.firstIfdOffset = header.FirstIfdOffset
	ir.tiffHeaderOffset = header.TiffHeaderOffset
	ir.po = header.FirstIfdOffset

	return ir.readIfd(ifds.NewIFD(header.ByteOrder, ifds.IfdType(header.FirstIfd), 0, ir.tiffHeaderOffset))
}

func (ir *ifdReader) parseIfdHeader(ifd ifds.Ifd) (err error) {
	loglevelInfo := ir.logLevelInfo()
	var tagCount uint16 // read tagCount

	if tagCount, err = ir.readUint16(ifd); err != nil || tagCount > 256 {
		// Log Ifd Reading error
		ir.logError(err).Object("ifd", ifd).Uint32("readerOffset", ir.po).Msgf("error tag count: %d for %s", tagCount, ifd.String())
		return err
	}

	// Log Ifd Info
	if loglevelInfo {
		ir.logInfo().Object("ifd", ifd).Uint16("tagCount", tagCount).Send()
	}

	// read Tag Headers
	var t tag.Tag

	br, ok := ir.reader.(BufferedReader)
	var buf []byte
	if ok {
		if buf, err = br.Peek(int(tagCount) * 12); err != nil {
			if ir.logLevelError() {
				logTag(ir.logError(err), t).Object("ifd", ifd).Uint16("tagCount", tagCount).Send()
			}
			return
		}
		if _, err = br.Discard(len(buf)); err != nil && ir.logLevelError() {
			ir.logError(err).Msg("Discard error")
			return err
		}
		ir.po += uint32(len(buf))

	}
	for i := 0; i < int(tagCount); i++ {
		if ok {
			if t, err = tagFromBuffer(ifd, buf[i*12:]); err != nil && ir.logLevelWarn() {
				logTag(ir.logWarn().Err(err), t).Send()
			}
		} else {
			if t, err = ir.readTagHeader(ifd); err != nil {
				if err == tag.ErrTagTypeNotValid {
					if ir.logLevelWarn() {
						ir.logWarn().Err(err).Object("tag", t).Send()
					}
					continue
				}
				if ir.logLevelError() {
					logTag(ir.logError(err), t).Object("ifd", ifd).Send()
				}
				return err
			}
		}
		if t.IsEmbedded() {
			ir.parseTag(t)
		} else {
			ir.addTagBuffer(t)
		}
		if loglevelInfo { // Log Tag Info
			ir.logTagInfo(t)
		}
	}

	// read Next Ifd Tag
	return ir.readNextIfdTag(ifd)
}

// BufferedReader interface represents bufio.Reader
type BufferedReader interface {
	Peek(n int) ([]byte, error)
	Discard(n int) (discarded int, err error)
	Read(p []byte) (n int, err error)
}

func (ir *ifdReader) readNextIfdTag(ifd ifds.Ifd) error {
	var err error
	if uint32(ir.buffer.nextTag().ValueOffset) <= ir.po+tag.TypeLongSize {
		var nextIfd uint32
		if nextIfd, err = ir.readUint32(ifd); err != nil {
			if ir.logLevelError() {
				ir.logError(err).Object("ifd", ifd).Uint32("offset", ir.po).Msgf("error reading nextIFD. Offset: %d Ifd: %s", ir.po, ifd.String())
			}
			return err
		}
		if ifd.IsType(ifds.IFD0) && nextIfd != 0 {
			t, _ := tag.NewTag(ifds.SubIFDs, tag.TypeIfd, 4, nextIfd, uint8(ifds.IFD0), ifd.Index+1, ifd.ByteOrder)
			ir.addTagBuffer(t)
			ir.logTagInfo(t) // Log Tag Info
		}
	}
	return nil
}

func (ir *ifdReader) readIfd(ifd ifds.Ifd) (err error) {
	if err = ir.parseIfdHeader(ifd); err != nil {
		return err
	}

	for t := ir.buffer.currentTag(); ir.buffer.validTag(); t = ir.buffer.advanceBuffer() {
		if t.IsType(tag.TypeIfd) {
			ir.seekNextTag(t)         // seek to next tag value
			ir.buffer.resetPosition() // Reset tagbuffer position to 0
			switch ifds.IfdType(t.Ifd) {
			case ifds.IFD0:
				switch t.ID {
				case ifds.SubIFDs:
					//ir.parseIfdHeader(ifd.ChildIfd(t)) // ignore errors from SubIfds
				case ifds.GPSTag, ifds.ExifTag:
					ir.parseIfdHeader(ifds.ChildIfd(t)) // ignore errors from GPSIfd and ExifIfd
				default:
					// Log Tag Info

				}
			case ifds.ExifIFD:
				if t.ID == exififd.MakerNote {
					if ir.Exif.Make == "Canon" {
						ir.parseIfdHeader(ifds.ChildIfd(t))
					}
				}
			}
		} else {
			ir.parseTag(t)
		}
	}
	return nil
}

func (ir *ifdReader) Read(buf []byte) (n int, err error) {
	if ir.exifLength != 0 && int(ir.po)+len(buf) > int(ir.exifLength) {
		return 0, imagetype.ErrDataLength
	}
	n, err = ir.reader.Read(buf)
	ir.po += uint32(n)
	return n, err
}

// ReadUint16 reads a uint16 from an ifdReader.
func (ir *ifdReader) readUint16(ifd ifds.Ifd) (uint16, error) {
	if _, err := ir.Read(ir.buffer.buf[:2]); err != nil {
		return 0, err
	}
	return ifd.ByteOrder.Uint16(ir.buffer.buf[:2]), nil
}

// ReadUint32 reads a uint32 from an ifdReader.
func (ir *ifdReader) readUint32(ifd ifds.Ifd) (uint32, error) {
	if _, err := ir.Read(ir.buffer.buf[:4]); err != nil {
		return 0, err
	}
	return ifd.ByteOrder.Uint32(ir.buffer.buf[:4]), nil
}

// readTagHeader reads the tagID uint16, tagType uint16, unitCount uint32 and valueOffset uint32
// from an Ifd. Returns Tag and error. If the tagType is unsupported, returns tag.ErrTagTypeNotValid.
func (ir *ifdReader) readTagHeader(ifd ifds.Ifd) (tag.Tag, error) {
	if _, err := ir.Read(ir.buffer.buf[:12]); err != nil {
		return tag.Tag{}, err
	}
	tagID := tag.ID(ifd.ByteOrder.Uint16(ir.buffer.buf[:2]))      // TagID
	tagType := tag.Type(ifd.ByteOrder.Uint16(ir.buffer.buf[2:4])) // TagType
	unitCount := ifd.ByteOrder.Uint32(ir.buffer.buf[4:8])         // UnitCount
	valueOffset := ifd.ByteOrder.Uint32(ir.buffer.buf[8:12])      // ValueOffset

	return tag.NewTag(tagID, tagIsIfd(ifd.Type, tagID, tagType), unitCount, valueOffset, uint8(ifd.Type), ifd.Index, ifd.ByteOrder) // NewTag
}

func tagFromBuffer(ifd ifds.Ifd, buf []byte) (t tag.Tag, err error) {
	tagID := tag.ID(ifd.ByteOrder.Uint16(buf[:2])) // TagID
	tagType := tagIsIfd(ifd.Type, tagID, tag.Type(ifd.ByteOrder.Uint16(buf[2:4])))
	t = tag.Tag{
		ID:          tagID,
		TagType:     tagType,                         // TagType
		UnitCount:   ifd.ByteOrder.Uint32(buf[4:8]),  // UnitCount
		ValueOffset: ifd.ByteOrder.Uint32(buf[8:12]), // ValueOffset
		Ifd:         uint8(ifd.Type),
		IfdIndex:    ifd.Index,
		ByteOrder:   ifd.ByteOrder,
	}
	if !t.TagType.IsValid() {
		err = tag.ErrTagTypeNotValid
	}
	return t, err
}

func tagIsIfd(ifdType ifds.IfdType, tagID tag.ID, tagType tag.Type) tag.Type {
	if tagType.Is(tag.TypeLong) || tagType.Is(tag.TypeUndefined) {
		switch ifdType {
		case ifds.IFD0: // RootIfd Children
			switch tagID {
			case ifds.ExifTag:
				return tag.TypeIfd
			case ifds.GPSTag:
				return tag.TypeIfd
			case ifds.SubIFDs:
				return tag.TypeIfd
			}
		case ifds.ExifIFD: // ExifIfd Children
			switch tagID {
			case exififd.MakerNote:
				return tag.TypeIfd
			}
		}
	}
	return tagType
}
