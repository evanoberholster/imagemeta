package exif2

import (
	"io"

	"github.com/evanoberholster/imagemeta/exif2/ifds"
	"github.com/evanoberholster/imagemeta/exif2/ifds/exififd"
	"github.com/evanoberholster/imagemeta/exif2/ifds/mknote/nikon"
	"github.com/evanoberholster/imagemeta/exif2/tag"
	"github.com/evanoberholster/imagemeta/imagetype"
	"github.com/evanoberholster/imagemeta/meta"
	"github.com/evanoberholster/imagemeta/meta/tiff"
	"github.com/evanoberholster/imagemeta/meta/utils"
	"github.com/rs/zerolog"
)

func Parse(r io.ReadSeeker) (Exif, error) {
	h, err := tiff.ScanTiffHeader(r, imagetype.ImageUnknown)
	if err != nil {
		return Exif{}, err
	}

	ir := NewIfdReader(Logger)
	defer ir.Close()

	if _, err = r.Seek(int64(h.TiffHeaderOffset), 0); err != nil {
		return ir.Exif, err
	}
	if err := ir.DecodeTiff(r, h); err != nil {
		return ir.Exif, err
	}
	return ir.Exif, nil
}

func (ir *ifdReader) DecodeTiff(r io.Reader, h meta.ExifHeader) error {
	// Log Header Info
	if ir.logLevelInfo() {
		ir.logger.Info().Str("imageType", h.ImageType.String()).Uint32("tiffHeader", h.TiffHeaderOffset).Uint32("firstIfdOffset", h.FirstIfdOffset).Uint32("exifLength", h.ExifLength).Send()
	}
	ir.ResetReader(r)

	ir.Exif.ImageType = h.ImageType
	ir.firstIfdOffset = h.FirstIfdOffset
	ir.exifLength = 4 * 1024 * 1024 // Max size is 4 MB
	if err := ir.discard(int(h.FirstIfdOffset)); err != nil {
		return err
	}
	err := ir.readIfd(ifds.NewIFD(h.ByteOrder, ifds.IfdType(h.FirstIfd), 0, ir.tiffHeaderOffset, 0))
	return err
}

func (ir *ifdReader) DecodeJPEGIfd(r io.Reader, h meta.ExifHeader) (err error) {
	// Log Header Info
	if ir.logLevelInfo() {
		ir.logger.Info().Str("imageType", h.ImageType.String()).Uint32("tiffHeader", h.TiffHeaderOffset).Uint32("firstIfdOffset", h.FirstIfdOffset).Uint32("exifLength", h.ExifLength).Send()
	}
	ir.ResetReader(r)
	ir.Exif.ImageType = h.ImageType
	ir.firstIfdOffset = h.FirstIfdOffset
	ir.exifLength = h.ExifLength
	if err = ir.discard(int(h.FirstIfdOffset)); err != nil {
		if ir.logLevelError() {
			ir.logError(err).Send()
		}
	}
	if err := ir.readIfd(ifds.NewIFD(h.ByteOrder, ifds.IfdType(h.FirstIfd), 0, ir.tiffHeaderOffset, 0)); err != nil {
		return err
	}
	err = ir.discard(int(ir.exifLength) - int(ir.po))
	return err
}

func (ir *ifdReader) DecodeIfd(r io.Reader, h meta.ExifHeader) (err error) {
	// Log Header Info
	if ir.logLevelInfo() {
		ir.logger.Info().Str("imageType", h.ImageType.String()).Uint32("tiffHeader", h.TiffHeaderOffset).Uint32("firstIfdOffset", h.FirstIfdOffset).Uint32("exifLength", h.ExifLength).Send()
	}
	ir.ResetReader(r)
	ir.Exif.ImageType = h.ImageType
	ir.exifLength = h.ExifLength
	ir.firstIfdOffset = h.FirstIfdOffset
	ir.po = h.FirstIfdOffset
	err = ir.readIfd(ifds.NewIFD(h.ByteOrder, ifds.IfdType(h.FirstIfd), 0, ir.tiffHeaderOffset, 0))
	return err
}

// NewIfdReader creates a new IfdReader with the given io.Reader
// Need to call defer IfdReader.Close() when complete
func NewIfdReader(l zerolog.Logger) ifdReader {
	ir := ifdReader{
		buffer: bufferPool.Get().(*buffer),
		logger: l,
	}
	ir.buffer.clear()
	return ir
}

// ResetReader resets the reader and clears the buffer
func (ir *ifdReader) ResetReader(r io.Reader) {
	ir.buffer.clear()
	ir.reader = r
}

// SetCustomTagParser sets a custom tag parser
func (ir *ifdReader) SetCustomTagParser(fn TagParserFn) {
	ir.customTagParser = fn
}

// Close closes an ifdReader. Should be called with defer following a newIfdReader
func (ir *ifdReader) Close() {
	bufferPool.Put(ir.buffer)
}

// ifdReader reads, decodes, and parses tags from an io.Reader
type ifdReader struct {
	logger zerolog.Logger
	reader io.Reader
	//bufReader        BufferedReader
	customTagParser  TagParserFn
	buffer           *buffer
	Exif             Exif
	po               uint32
	tiffHeaderOffset uint32
	firstIfdOffset   uint32
	exifLength       uint32
}

func (ir *ifdReader) readIfdHeader(ifd ifds.Ifd) (err error) {
	loglevelInfo := ir.logLevelInfo()
	var tagCount uint16 // read tagCount

	if tagCount, err = ir.readUint16(ifd); err != nil || tagCount > 128 {
		// Log Ifd Reading error
		ir.logError(err).Object("ifd", ifd).Uint32("readerOffset", ir.po).Msgf("error tag count: %d for %s", tagCount, ifd.String())
		return err
	}

	if loglevelInfo { // Log Ifd Info
		ir.logInfo().Object("ifd", ifd).Uint16("tagCount", tagCount).Send()
	}

	buf, err := ir.fastRead(int(tagCount) * 12) // read Tag Headers
	if err != nil {
		if ir.logLevelError() {
			ir.logError(err).Object("ifd", ifd).Uint16("tagCount", tagCount).Send()
		}
		return err
	}
	var t Tag
	for i := 0; i < int(tagCount); i++ {
		if t, err = tagFromBuffer(ifd, buf[i*12:]); err != nil {
			if ir.logLevelWarn() {
				t.logTag(ir.logWarn().Err(err)).Send()
			}
			continue
		}
		if loglevelInfo { // Log Tag Info
			t.logTag(ir.logDebug()).Send()
		}
		if t.IsEmbedded() {
			ir.parseTag(t)
		} else {
			ir.addTagBuffer(t)
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
	if uint32(ir.buffer.nextTag().ValueOffset) <= ir.po {
		var nextIfd uint32
		if nextIfd, err = ir.readUint32(ifd); err != nil {
			if ir.logLevelError() {
				ir.logError(err).Object("ifd", ifd).Uint32("offset", ir.po).Msgf("error reading nextIFD. Offset: %d Ifd: %s", ir.po, ifd.String())
			}
			return err
		}
		if ifd.IsType(ifds.IFD0) && nextIfd != 0 {
			t := NewTag(ifds.SubIFDs, tag.TypeIfd, 4, nextIfd, ifds.IFD0, ifd.Index+1, ifd.ByteOrder)
			ir.addTagBuffer(t)
			if ir.logLevelDebug() { // Log Tag Info
				t.logTag(ir.logDebug()).Send()
			}
		}
	}
	return nil
}

func (ir *ifdReader) readIfd(ifd ifds.Ifd) (err error) {
	if err = ir.readIfdHeader(ifd); err != nil {
		return err
	}

	for t := ir.buffer.currentTag(); ir.buffer.validTag(); t = ir.buffer.advanceBuffer() {

		if t.IsType(tag.TypeIfd) {
			if err = ir.seekToTag(t); err != nil { // seek to next tag value
				ir.logError(err).Send()
			}
			ir.buffer.resetPosition() // Reset tagbuffer position to 0
			switch t.Ifd {
			case ifds.IFD0:
				switch t.ID {
				case ifds.GPSTag, ifds.ExifTag:
					if err = ir.readIfdHeader(t.childIfd()); err != nil { // ignore errors from GPSIfd and ExifIfd
						ir.logError(err).Send()
					}
				}
			case ifds.SubIfd0, ifds.SubIfd1, ifds.SubIfd2, ifds.SubIfd3, ifds.SubIfd4, ifds.SubIfd5:
				if err = ir.readIfdHeader(t.childIfd()); err != nil { // ignore errors from SubIfd0
					ir.logError(err).Send()
				}
			case ifds.ExifIFD:
				if t.ID == exififd.MakerNote {
					ir.readMakerNotes(t)
				}
			}
			continue
		}
		if t.ID == ifds.SubIFDs && t.Ifd == ifds.IFD0 {
			ir.readSubIfds(t)
			continue
		}
		// Parse all other tags
		ir.parseTag(t)

	}
	return nil
}

// readSubIfds from SubIfd Tag and add them to the TagBuffer.
// Limited to total 6 SubIfds, can be increased.
func (ir *ifdReader) readSubIfds(t Tag) {
	if t.IsType(tag.TypeLong) {
		buf, err := ir.readTagValue()
		if err != nil {
			if ir.logLevelError() {
				t.logTag(ir.logError(err)).Send()
			}
			return
		}
		for i := 0; i < int(t.UnitCount); i++ {
			ir.addTagBuffer(NewTag(t.ID, tag.TypeIfd, tag.TypeIfdSize, t.ByteOrder.Uint32(buf[4*i:]), ifds.SubIfd0+ifds.IfdType(i), 0, t.ByteOrder))
		}
	}
}

func (ir *ifdReader) readMakerNotes(t Tag) {
	switch ir.Exif.CameraMake {
	case ifds.Canon:
		if err := ir.readIfdHeader(t.childIfd()); err != nil {
			ir.logError(err).Send()
		}
	case ifds.Nikon:
		if t.Size() > 18 { // read Nikon Makernotes header 18 bytes
			buf, err := ir.fastRead(18)
			if err != nil {
				t.logTag(ir.logError(err)).Send()
			}
			if nikon.IsNikonMkNoteHeaderBytes(buf[:5]) {
				ir.Exif.ImageType = imagetype.ImageNEF
				if byteOrder := utils.BinaryOrder(buf[10:14]); byteOrder != utils.UnknownEndian {
					err = ir.readIfdHeader(ifds.NewIFD(byteOrder, ifds.MknoteIFD, t.IfdIndex, t.ValueOffset, t.ValueOffset+byteOrder.Uint32(buf[14:18])))
					if err != nil {
						ir.logError(err).Send()
					}
				}
			}
		}
	}
}

func (ir *ifdReader) fastRead(n int) (buf []byte, err error) {
	if ir.exifLength != 0 && int(ir.po)+n > int(ir.exifLength) {
		return nil, imagetype.ErrDataLength
	}
	if br, ok := ir.reader.(BufferedReader); ok {
		if buf, err = br.Peek(n); err != nil {
			if ir.logLevelError() {
				ir.logError(err).Msg("Peek error")
			}
			return
		}
		if n, err = br.Discard(len(buf)); err != nil {
			if ir.logLevelError() {
				ir.logError(err).Msg("Discard error")
			}
			return
		}
		ir.po += uint32(n)
		return
	}
	if n, err = ir.reader.Read(ir.buffer.buf[:n]); err != nil {
		if ir.logLevelError() {
			ir.logError(err).Msg("Read error")
		}
		return
	}
	ir.po += uint32(n)
	return ir.buffer.buf[:n], err
}

// ReadUint16 reads a uint16 from an ifdReader.
func (ir *ifdReader) readUint16(ifd ifds.Ifd) (uint16, error) {
	buf, err := ir.fastRead(2)
	return ifd.ByteOrder.Uint16(buf), err
}

// ReadUint32 reads a uint32 from an ifdReader.
func (ir *ifdReader) readUint32(ifd ifds.Ifd) (uint32, error) {
	buf, err := ir.fastRead(4)
	return ifd.ByteOrder.Uint32(buf), err
}

func tagFromBuffer(ifd ifds.Ifd, buf []byte) (t Tag, err error) {
	tagID := tag.ID(ifd.ByteOrder.Uint16(buf[:2]))                  // TagID
	tagType := tag.Type(ifd.ByteOrder.Uint16(buf[2:4]))             // TagType
	unitCount := ifd.ByteOrder.Uint32(buf[4:8])                     // UnitCount
	valueOffset := ifd.ByteOrder.Uint32(buf[8:12]) + ifd.BaseOffset // ValueOffset

	t = NewTag(tagID, tagIsIfd(ifd.Type, tagID, tagType), unitCount, valueOffset, ifd.Type, ifd.Index, ifd.ByteOrder) // NewTag
	if !t.IsValid() {
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
