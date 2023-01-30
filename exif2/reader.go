package exif2

import (
	"io"
	"os"

	"github.com/evanoberholster/imagemeta/exif2/ifds"
	"github.com/evanoberholster/imagemeta/exif2/ifds/exififd"
	"github.com/evanoberholster/imagemeta/exif2/tag"
	"github.com/evanoberholster/imagemeta/imagetype"
	"github.com/evanoberholster/imagemeta/meta"
	"github.com/evanoberholster/imagemeta/tiff"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var Logger zerolog.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout}).Level(zerolog.PanicLevel)

func Decode(r io.ReadSeeker) (Exif, error) {
	header, err := tiff.ScanTiffHeader(r, imagetype.ImageCR3)
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

func (ir *ifdReader) DecodeTiff(r io.Reader, h meta.ExifHeader) error {
	ir.buffer.clear()
	// Log Header Info
	if ir.logInfo() {
		ir.logger.Info().Str("imageType", h.ImageType.String()).Uint32("tiffHeader", h.TiffHeaderOffset).Uint32("firstIfdOffset", h.FirstIfdOffset).Uint32("exifLength", h.ExifLength).Send()
	}
	ir.Exif.ImageType = h.ImageType
	ir.exifLength = 4 * 1024 * 1024 // Max size is 4 MB
	if err := ir.discard(int(h.FirstIfdOffset)); err != nil {
		return err
	}
	return ir.readIfd(ifds.NewIFD(h.ByteOrder, ifds.IfdType(h.FirstIfd), 0, tag.Offset(ir.tiffHeaderOffset)))
}

func (ir *ifdReader) DecodeJPEGIfd(r io.Reader, h meta.ExifHeader) error {
	// Log Header Info
	if ir.logInfo() {
		ir.logger.Info().Str("imageType", h.ImageType.String()).Uint32("tiffHeader", h.TiffHeaderOffset).Uint32("firstIfdOffset", h.FirstIfdOffset).Uint32("exifLength", h.ExifLength).Send()
	}
	ir.buffer.clear()
	ir.reader = r
	ir.exifLength = h.ExifLength
	ir.discard(int(h.FirstIfdOffset))
	if err := ir.readIfd(ifds.NewIFD(h.ByteOrder, ifds.IfdType(h.FirstIfd), 0, tag.Offset(ir.tiffHeaderOffset))); err != nil {
		return err
	}
	err := ir.discard(int(ir.exifLength) - int(ir.po))
	return err
}

func (ir *ifdReader) DecodeIfd(r io.Reader, h meta.ExifHeader) error {
	// Log Header Info
	if ir.logInfo() {
		ir.logger.Info().Str("imageType", h.ImageType.String()).Uint32("tiffHeader", h.TiffHeaderOffset).Uint32("firstIfdOffset", h.FirstIfdOffset).Uint32("exifLength", h.ExifLength).Send()
	}
	ir.buffer.clear()
	ir.reader = r
	ir.Exif.ImageType = h.ImageType
	ir.exifLength = h.ExifLength
	ir.po = h.FirstIfdOffset
	if err := ir.readIfd(ifds.NewIFD(h.ByteOrder, ifds.IfdType(h.FirstIfd), 0, tag.Offset(ir.tiffHeaderOffset))); err != nil {
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
	ir.reader = r
}

// Close closes an ifdReader. Should be called with defer following a newIfdReader
func (ir *ifdReader) Close() {
	bufferPool.Put(ir.buffer)
}

// ifdReader reads, decodes, and parses tags from an io.Reader
type ifdReader struct {
	logger           zerolog.Logger
	reader           io.Reader
	buffer           *buffer
	Exif             Exif
	po               uint32
	tiffHeaderOffset uint32
	firstIfdOffset   uint32
	exifLength       uint32
}

func (ir *ifdReader) ReadIfd0(header meta.ExifHeader) error {
	// Log Header Info
	if ir.logInfo() {
		ir.logger.Info().Str("imageType", header.ImageType.String()).Uint32("tiffHeader", header.TiffHeaderOffset).Uint32("firstIfdOffset", header.FirstIfdOffset).Send()
	}

	ir.firstIfdOffset = header.FirstIfdOffset
	ir.tiffHeaderOffset = header.TiffHeaderOffset
	ir.po = header.FirstIfdOffset

	return ir.readIfd(ifds.NewIFD(header.ByteOrder, ifds.IfdType(header.FirstIfd), 0, tag.Offset(ir.tiffHeaderOffset)))
}

func (ir *ifdReader) parseIfdHeader(ifd ifds.Ifd) error {
	var err error

	// read tagCount
	var tagCount uint16
	if tagCount, err = ir.readUint16(ifd); err != nil || tagCount > 256 {
		// Log Ifd Reading error
		ir.logger.Error().Err(err).Str("Ifd", ifd.Type.String()).Uint32("offset", ir.po).Msgf("error tag count: %d for %s", tagCount, ifd.String())
		return err
	}

	// Log Ifd Info
	if ir.logInfo() {
		ir.logger.Info().Str("ifd", ifd.Type.String()).Int8("IfdIndex", ifd.Index).Stringer("offset", ifd.Offset).Uint16("tagCount", tagCount).Send()
	}

	// read Tag Headers
	var t tag.Tag
	for i := 0; i < int(tagCount); i++ {
		if t, err = ir.readTagHeader(ifd); err != nil {
			// Log Ifd Reading error
			if err == tag.ErrTagTypeNotValid {
				if ir.logInfo() {
					ir.logger.Debug().Err(err).Stringer("id", t.ID).Stringer("ifd", ifd.Type).Stringer("offset", t.ValueOffset).Uint16("type", uint16(t.Type())).Send()
				}
				continue
			}
			ir.logger.Error().Err(err).Stringer("id", t.ID).Stringer("ifd", ifd.Type).Stringer("offset", t.ValueOffset).Stringer("type", t.Type()).Send()
			return err
		}
		if t.IsEmbedded() {
			ir.processTag(t)
		} else {
			ir.addTagBuffer(t)
		}

		// Log Tag Info
		ir.logTagInfo(t)
	}

	// read Next Ifd Tag
	err = ir.readNextIfdTag(ifd)
	return err
}

func (ir *ifdReader) readNextIfdTag(ifd ifds.Ifd) error {
	var err error
	if uint32(ir.buffer.nextTag().ValueOffset) <= ir.po+tag.TypeLongSize {
		var nextIfd uint32
		if nextIfd, err = ir.readUint32(ifd); err != nil {
			ir.logger.Error().Err(err).Str("ifd", ifd.Type.String()).Uint32("offset", ir.po).Msgf("error reading nextIFD. Offset: %d Ifd: %s", ir.po, ifd.String())
			return err
		}
		if ifd.IsType(ifds.IFD0) && nextIfd != 0 {
			t, _ := tag.NewTag(ifds.SubIFDs, tag.TypeIfd, 4, nextIfd, uint8(ifds.IFD0), ifd.Index+1, ifd.ByteOrder)
			ir.addTagBuffer(t)

			// Log Tag Info
			ir.logTagInfo(t)
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
			// parse to location
			discard := int(t.ValueOffset) - int(ir.po)
			if err := ir.discard(discard); err != nil {
				ir.logger.Error().Err(err).Stringer("ifd", t.Type()).Uint16("ifdIndex", uint16(t.IfdIndex)).Uint32("discard", uint32(discard)).Send()
				return nil
			}
			// Reset tagbuffer position to 0
			ir.buffer.resetPosition()
			switch ifds.IfdType(t.Ifd) {
			case ifds.IFD0:
				switch t.ID {
				case ifds.SubIFDs:
					//ir.parseIfdHeader(ifd.ChildIfd(t)) // ignore errors from SubIfds
					//fmt.Println(ir.buffer.tag[ir.buffer.pos:ir.buffer.len])
				case ifds.GPSTag, ifds.ExifTag:
					ir.parseIfdHeader(ifds.ChildIfd(t)) // ignore errors from GPSIfd and ExifIfd
					//fmt.Println(ir.buffer.tag[ir.buffer.pos:ir.buffer.len])
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
			ir.processTag(t)
		}
	}
	return nil
}

func (ir *ifdReader) checkLength(n uint32) bool {
	return ir.po+n >= ir.exifLength
}

// ReadUint16 reads a uint16 from an ifdReader.
func (ir *ifdReader) readUint16(ifd ifds.Ifd) (uint16, error) {
	if ir.checkLength(2) {
		return 0, imagetype.ErrDataLength
	}
	n, err := ir.reader.Read(ir.buffer.buf[:2])
	ir.po += uint32(n)
	return ifd.ByteOrder.Uint16(ir.buffer.buf[:2]), err
}

// ReadUint32 reads a uint32 from an ifdReader.
func (ir *ifdReader) readUint32(ifd ifds.Ifd) (uint32, error) {
	if ir.checkLength(4) {
		return 0, imagetype.ErrDataLength
	}
	n, err := ir.reader.Read(ir.buffer.buf[:4])
	ir.po += uint32(n)
	return ifd.ByteOrder.Uint32(ir.buffer.buf[:4]), err
}

// readTagHeader reads the tagID uint16, tagType uint16, unitCount uint32 and valueOffset uint32
// from an Ifd. Returns Tag and error. If the tagType is unsupported, returns tag.ErrTagTypeNotValid.
func (ir *ifdReader) readTagHeader(ifd ifds.Ifd) (tag.Tag, error) {
	if ir.checkLength(12) {
		return tag.Tag{}, imagetype.ErrDataLength
	}
	if _, err := ir.reader.Read(ir.buffer.buf[:12]); err != nil {
		return tag.Tag{}, err
	}
	ir.po += 12
	tagID := tag.ID(ifd.ByteOrder.Uint16(ir.buffer.buf[:2]))      // TagID
	tagType := tag.Type(ifd.ByteOrder.Uint16(ir.buffer.buf[2:4])) // TagType
	unitCount := ifd.ByteOrder.Uint32(ir.buffer.buf[4:8])         // UnitCount
	valueOffset := ifd.ByteOrder.Uint32(ir.buffer.buf[8:12])      // ValueOffset

	tagType = tagIsIfd(ifd, tagID, tagType)
	return tag.NewTag(tagID, tagType, unitCount, valueOffset, uint8(ifd.Type), ifd.Index, ifd.ByteOrder) // NewTag
}

func tagIsIfd(ifd ifds.Ifd, tagID tag.ID, tagType tag.Type) tag.Type {
	if tagType.Is(tag.TypeLong) {
		// RootIfd Children
		if ifd.IsType(ifds.IFD0) {
			switch tagID {
			case ifds.ExifTag:
				return tag.TypeIfd
			case ifds.GPSTag:
				return tag.TypeIfd
			case ifds.SubIFDs:
				return tag.TypeIfd
			}
		}
	}
	if tagType.Is(tag.TypeUndefined) {
		// ExifIfd Children
		if ifd.IsType(ifds.ExifIFD) {
			switch tagID {
			case exififd.MakerNote:
				return tag.TypeIfd
			}
		}
	}
	return tagType
}
