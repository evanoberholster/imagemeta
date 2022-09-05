package exif2

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/evanoberholster/imagemeta/exif2/ifds"
	"github.com/evanoberholster/imagemeta/exif2/ifds/exififd"
	"github.com/evanoberholster/imagemeta/exif2/ifds/gpsifd"
	"github.com/evanoberholster/imagemeta/exif2/tag"
	"github.com/evanoberholster/imagemeta/imagetype"
	"github.com/evanoberholster/imagemeta/meta"
	"github.com/evanoberholster/imagemeta/tiff"
	"github.com/pkg/errors"
)

func Decode(r io.ReadSeeker) (Exif, error) {
	header, err := tiff.ScanTiffHeader(r, imagetype.ImageCR3)
	if err != nil {
		return Exif{}, err
	}

	p := bufferPool.Get().(*buffer)
	p.len = 0
	p.pos = 0
	defer bufferPool.Put(p)

	ir := ifdReader{
		reader: r,
		po:     0,
		buffer: p,
	}
	r.Seek(0, 0)
	if err := ir.ReadIfd0(header); err != nil {
		return ir.exif, err
	}

	return ir.exif, nil
}

// ifdReader reads, decodes, and parses tags from am io.Reader
type ifdReader struct {
	reader           io.ReadSeeker
	buffer           *buffer
	exif             Exif
	po               uint32
	tiffHeaderOffset uint32
	firstIfdOffset   uint32
}

func (ir *ifdReader) ReadIfd0(header meta.ExifHeader) error {
	//fmt.Println(header)
	ir.firstIfdOffset = header.FirstIfdOffset
	ir.tiffHeaderOffset = header.TiffHeaderOffset
	ir.reader.Seek(int64(header.TiffHeaderOffset)+int64(header.FirstIfdOffset), 0)
	ir.po = header.FirstIfdOffset
	if header.ByteOrder == binary.LittleEndian {
		return ir.readIfd(ifds.NewIFD(tag.LittleEndian, ifds.IFD0, 0, ir.tiffHeaderOffset))
	}
	return ir.readIfd(ifds.NewIFD(tag.BigEndian, ifds.IFD0, 0, ir.tiffHeaderOffset))
}

func (ir *ifdReader) parseIfdHeader(ifd ifds.Ifd) error {
	var err error

	// read tagCount
	var tagCount uint16
	if tagCount, err = ir.readUint16(ifd); err != nil || tagCount > 256 {
		return errors.Wrapf(err, "error tag count: %d for %s", tagCount, ifd.String())
	}

	// LogInfo Entering Ifd
	//fmt.Printf("%s Tag Count: %d\n", ifd.String(), tagCount)

	// read Tag Headers
	var t tag.Tag
	for i := 0; i < int(tagCount); i++ {
		if t, err = ir.readTagHeader(ifd); err != nil {
			if err == tag.ErrTagTypeNotValid {
				//if errors.Is(err, tag.ErrTagTypeNotValid) {
				// Log TagNotValid Error
				//ifdEnumerateLogger.Warningf(nil, "Tag in IFD [%s] at position (%d) has invalid type and will be skipped.", fqIfdPath, i)
				continue
			}
			return err
		}
		//fmt.Println(tagString(t))
		if t.IsEmbedded() {
			ir.processTag(t)
		} else {
			ir.buffer.addTag(t)
		}
	}

	// read Next Ifd
	var nextIfd uint32
	if nextIfd, err = ir.readUint32(ifd); err != nil {
		return errors.Wrapf(err, "error reading nextIFD. Offset: %d Ifd: %s", ir.po, ifd.String())
	}
	if nextIfd != 0 {
		ifdTag, _ := tag.NewTag(ifds.SubIFDs, tag.TypeIfd, 4, nextIfd, uint8(ifds.IFD0), ifd.Index+1, ifd.ByteOrder)
		ir.buffer.addTag(ifdTag)
	}
	return nil
}

func (ir *ifdReader) readSubIfd(t tag.Tag) bool {
	if t.IsType(tag.TypeIfd) {
		// parse to location
		if err := ir.discard(int(t.ValueOffset) - int(ir.po)); err != nil {
			fmt.Println(err)
			return false
		}
		// Reset tagbuffer position to 0
		ir.buffer.resetPosition()
		return true
	}
	return false
}

func (ir *ifdReader) readIfd(ifd ifds.Ifd) (err error) {
	if err = ir.parseIfdHeader(ifd); err != nil {
		return err
	}

	for t := ir.buffer.currentTag(); ir.buffer.validTag(); t = ir.buffer.nextTag() {
		//fmt.Println(ir.tagBuf, ir.tagBuf.pos, ir.tagBuf.len)
		if t.ID == ifds.XMLPacket {
			ir.discard(int(t.Size()))
			continue
		}
		if ir.readSubIfd(t) {
			if ifd.IsType(ifds.IFD0) {
				switch t.ID {
				case ifds.SubIFDs:
					if err = ir.parseIfdHeader(ifd.ChildIfd(t)); err != nil {
						return err
					}
				case ifds.GPSTag, ifds.ExifTag:
					if err = ir.parseIfdHeader(ifd.ChildIfd(t)); err != nil {
						return err
					}
				}
			}
		} else {
			ir.processTag(t)
		}
	}
	return nil
}

func (ir *ifdReader) processTag(t tag.Tag) {
	switch ifds.IfdType(t.Ifd) {
	case ifds.IFD0:
		switch t.ID {
		case ifds.Make:
			ir.exif.Make = ir.parseString(t)
		case ifds.Model:
			ir.exif.Model = ir.parseString(t)
		case ifds.Artist:
			ir.exif.Artist = ir.parseString(t)
		case ifds.Copyright:
			ir.exif.Copyright = ir.parseString(t)
		case ifds.ImageWidth:
			ir.exif.ImageWidth = ir.parseUint16(t)
		case ifds.ImageLength:
			ir.exif.ImageHeight = ir.parseUint16(t)
		case ifds.StripOffsets:
			ir.exif.StripOffsets = ir.parseUint32(t)
		case ifds.StripByteCounts:
			ir.exif.StripByteCounts = ir.parseUint32(t)
		case ifds.Orientation:
			ir.exif.Orientation = ir.parseOrientation(t)
		case ifds.Software:
			ir.exif.Software = ir.parseString(t)
		case ifds.DateTime:
			ir.exif.modifyDate = ir.parseDate(t)
			//default:
			//	fmt.Println(tagString(t))
		}
	case ifds.ExifIFD:
		switch t.ID {
		case exififd.LensMake:
			ir.exif.LensMake = ir.parseString(t)
		case exififd.LensModel:
			ir.exif.LensModel = ir.parseString(t)
		case exififd.LensSerialNumber:
			ir.exif.LensSerial = ir.parseString(t)
		case exififd.BodySerialNumber:
			ir.exif.CameraSerial = ir.parseString(t)
		case exififd.PixelXDimension:
			if ir.exif.ImageWidth == 0 {
				ir.exif.ImageWidth = uint16(ir.parseUint32(t))
			}
		case exififd.PixelYDimension:
			if ir.exif.ImageHeight == 0 {
				ir.exif.ImageHeight = uint16(ir.parseUint32(t))
			}
		case exififd.ExposureTime:
			ir.exif.ExposureTime = ir.parseExposureTime(t)
		case exififd.FNumber:
			ir.exif.FNumber = ir.parseFNumber(t)
		case exififd.ExposureProgram:
			ir.exif.ExposureProgram = meta.ExposureProgram(ir.parseUint16(t))
		case exififd.ExposureBiasValue:
			ir.exif.ExposureBias = ir.parseExposureBias(t)
		case exififd.ExposureMode:
			ir.exif.ExposureMode = meta.ExposureMode(ir.parseUint16(t))
		case exififd.MeteringMode:
			ir.exif.MeteringMode = meta.MeteringMode(ir.parseUint16(t))
		case exififd.ISOSpeedRatings:
			ir.exif.ISOSpeed = ir.parseUint32(t)
		case ifds.DateTimeOriginal:
			ir.exif.dateTimeOriginal = ir.parseDate(t)
		case ifds.DateTimeDigitized:
			ir.exif.createDate = ir.parseDate(t)
		case ifds.Flash:
			ir.exif.Flash = meta.Flash(ir.parseUint16(t))
		case ifds.FocalLength:
			ir.exif.FocalLength = ir.parseFocalLength(t)
		case exififd.FocalLengthIn35mmFilm:
			ir.exif.FocalLengthIn35mmFormat = ir.parseFocalLength(t)
		case exififd.LensSpecification:
			ir.exif.LensInfo = ir.parseLensInfo(t)
		case exififd.SubSecTime:
			ir.exif.subSecTime = ir.parseSubSecTime(t)
		case exififd.SubSecTimeOriginal:
			ir.exif.subSecTimeOriginal = ir.parseSubSecTime(t)
		case exififd.SubSecTimeDigitized:
			ir.exif.subSecTimeDigitized = ir.parseSubSecTime(t)
		}
	case ifds.GPSIFD:
		switch t.ID {
		case gpsifd.GPSAltitudeRef:
			ir.exif.GPS.altitudeRef = ir.parseGPSRef(t)
		case gpsifd.GPSLatitudeRef:
			ir.exif.GPS.latitudeRef = ir.parseGPSRef(t)
		case gpsifd.GPSLongitudeRef:
			ir.exif.GPS.longitudeRef = ir.parseGPSRef(t)
		case gpsifd.GPSAltitude:
			ir.exif.GPS.altitude = ir.parseGPSAltitude(t)
		case gpsifd.GPSLatitude:
			ir.exif.GPS.latitude = ir.parseGPSCoord(t)
		case gpsifd.GPSLongitude:
			ir.exif.GPS.longitude = ir.parseGPSCoord(t)
		case gpsifd.GPSTimeStamp:
			ir.exif.GPS.time = ir.parseGPSTimeStamp(t)
		case gpsifd.GPSDateStamp:
			ir.exif.GPS.date = ir.parseGPSDateStamp(t)
		}
	}
}

// ReadUint16 reads a uint16 from an ifdReader.
func (ir *ifdReader) readUint16(ifd ifds.Ifd) (uint16, error) {
	n, err := ir.reader.Read(ir.buffer.buf[:2])
	ir.po += uint32(n)
	return ifd.ByteOrder.Uint16(ir.buffer.buf[:2]), err
}

// ReadUint32 reads a uint32 from an ifdReader.
func (ir *ifdReader) readUint32(ifd ifds.Ifd) (uint32, error) {
	n, err := ir.reader.Read(ir.buffer.buf[:4])
	ir.po += uint32(n)
	return ifd.ByteOrder.Uint32(ir.buffer.buf[:4]), err
}

// ReadTagHeader reads the tagID uint16, tagType uint16, unitCount uint32 and valueOffset uint32
// from an Ifd. Returns Tag and error. If the tagType is unsupported, returns tag.ErrTagTypeNotValid.
func (ir *ifdReader) readTagHeader(ifd ifds.Ifd) (tag.Tag, error) {
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
