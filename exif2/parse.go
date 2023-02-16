package exif2

import (
	"fmt"
	"math"
	"time"

	"github.com/evanoberholster/imagemeta/exif2/ifds"
	"github.com/evanoberholster/imagemeta/exif2/ifds/exififd"
	"github.com/evanoberholster/imagemeta/exif2/ifds/gpsifd"
	"github.com/evanoberholster/imagemeta/exif2/tag"
	"github.com/evanoberholster/imagemeta/meta"
)

func (ir *ifdReader) processTag(t tag.Tag) {
	switch ifds.IfdType(t.Ifd) {
	case ifds.IFD0:
		switch t.ID {
		case ifds.Make:
			ir.Exif.Make = ir.parseString(t)
		case ifds.Model:
			ir.Exif.Model = ir.parseString(t)
		case ifds.Artist:
			ir.Exif.Artist = ir.parseString(t)
		case ifds.Copyright:
			ir.Exif.Copyright = ir.parseString(t)
		case ifds.ImageWidth:
			ir.Exif.ImageWidth = uint16(ir.parseUint32(t))
		case ifds.ImageLength:
			ir.Exif.ImageHeight = uint16(ir.parseUint32(t))
		case ifds.StripOffsets:
			ir.Exif.StripOffsets = ir.parseUint32(t)
		case ifds.StripByteCounts:
			ir.Exif.StripByteCounts = ir.parseUint32(t)
		case ifds.Orientation:
			ir.Exif.Orientation = ir.parseOrientation(t)
		case ifds.Software:
			ir.Exif.Software = ir.parseString(t)
		case ifds.ImageDescription:
			ir.Exif.ImageDescription = ir.parseString(t)
		case ifds.DateTime:
			ir.Exif.modifyDate = ir.parseDate(t)
		case ifds.CameraSerialNumber:
			if ir.Exif.CameraSerial == "" {
				ir.Exif.CameraSerial = ir.parseString(t)
			}
		case ifds.ApplicationNotes:
			//ir.Exif.ApplicationNotes = ir.parseApplicationNotes(t)
			//default:
			//	fmt.Println(tagString(t))
		}
	case ifds.ExifIFD:
		switch t.ID {
		case exififd.LensMake:
			ir.Exif.LensMake = ir.parseString(t)
		case exififd.LensModel:
			ir.Exif.LensModel = ir.parseString(t)
		case exififd.LensSerialNumber:
			ir.Exif.LensSerial = ir.parseString(t)
		case exififd.CameraOwnerName:
			if ir.Exif.Artist == "" {
				ir.Exif.Artist = ir.parseString(t)
			}
		case exififd.BodySerialNumber:
			if ir.Exif.CameraSerial == "" {
				ir.Exif.CameraSerial = ir.parseString(t)
			}
		case exififd.PixelXDimension:
			if ir.Exif.ImageWidth == 0 {
				ir.Exif.ImageWidth = uint16(ir.parseUint32(t))
			}
		case exififd.PixelYDimension:
			if ir.Exif.ImageHeight == 0 {
				ir.Exif.ImageHeight = uint16(ir.parseUint32(t))
			}
		case exififd.ExposureTime:
			ir.Exif.ExposureTime = ir.parseExposureTime(t)
		case exififd.ApertureValue:
			if ir.Exif.FNumber == 0.0 {
				r := ir.parseRationalU(t)
				f := float64(r[0]) / float64(r[1])
				ir.Exif.FNumber = meta.Aperture(math.Round(math.Pow(math.Sqrt2, float64(f))*100) / 100)
			}
		case exififd.FNumber:
			ir.Exif.FNumber = ir.parseAperture(t)
		case exififd.ExposureProgram:
			ir.Exif.ExposureProgram = meta.ExposureProgram(ir.parseUint16(t))
		case exififd.ExposureBiasValue:
			ir.Exif.ExposureBias = ir.parseExposureBias(t)
		case exififd.ExposureMode:
			ir.Exif.ExposureMode = meta.ExposureMode(ir.parseUint16(t))
		case exififd.MeteringMode:
			ir.Exif.MeteringMode = meta.MeteringMode(ir.parseUint16(t))
		case exififd.ISOSpeedRatings:
			ir.Exif.ISOSpeed = ir.parseUint32(t)
		case ifds.DateTimeOriginal:
			ir.Exif.dateTimeOriginal = ir.parseDate(t)
		case ifds.DateTimeDigitized:
			ir.Exif.createDate = ir.parseDate(t)
		case ifds.Flash:
			ir.Exif.Flash = meta.Flash(ir.parseUint16(t))
		case ifds.FocalLength:
			ir.Exif.FocalLength = ir.parseFocalLength(t)
		case exififd.FocalLengthIn35mmFilm:
			ir.Exif.FocalLengthIn35mmFormat = ir.parseFocalLength(t)
		case exififd.LensSpecification:
			ir.Exif.LensInfo = ir.parseLensInfo(t)
		case exififd.SubSecTime:
			ir.Exif.subSecTime = ir.parseSubSecTime(t)
		case exififd.SubSecTimeOriginal:
			ir.Exif.subSecTimeOriginal = ir.parseSubSecTime(t)
		case exififd.SubSecTimeDigitized:
			ir.Exif.subSecTimeDigitized = ir.parseSubSecTime(t)
		case exififd.OffsetTime:
			ir.Exif.offsetTime = ir.parseOffsetTime(t)
		case exififd.OffsetTimeOriginal:
			ir.Exif.offsetTimeOriginal = ir.parseOffsetTime(t)
		case exififd.OffsetTimeDigitized:
			ir.Exif.offsetTimeDigitized = ir.parseOffsetTime(t)
			//default:
			//	fmt.Println(tagString(t))
			//	fmt.Println(string(ir.parseString(t)))
		}
	case ifds.GPSIFD:
		switch t.ID {
		case gpsifd.GPSAltitudeRef:
			ir.Exif.GPS.altitudeRef = ir.parseGPSRef(t)
		case gpsifd.GPSLatitudeRef:
			ir.Exif.GPS.latitudeRef = ir.parseGPSRef(t)
		case gpsifd.GPSLongitudeRef:
			ir.Exif.GPS.longitudeRef = ir.parseGPSRef(t)
		case gpsifd.GPSAltitude:
			ir.Exif.GPS.altitude = ir.parseGPSAltitude(t)
		case gpsifd.GPSLatitude:
			ir.Exif.GPS.latitude = ir.parseGPSCoord(t)
		case gpsifd.GPSLongitude:
			ir.Exif.GPS.longitude = ir.parseGPSCoord(t)
		case gpsifd.GPSTimeStamp:
			ir.Exif.GPS.time = ir.parseGPSTimeStamp(t)
		case gpsifd.GPSDateStamp:
			ir.Exif.GPS.date = ir.parseGPSDateStamp(t)
			//default:
			//	fmt.Println(tagString(t))
		}
	}
}

func (ir *ifdReader) parseApplicationNotes(t tag.Tag) ApplicationNotes {
	buf, err := ir.readTagValue()
	if err != nil {
		ir.logParseWarn(t, "parseApplicationNotes", "unrecognized tag type", err)
	}

	res1 := trimNULBuffer(buf)
	res2 := make([]byte, len(res1))
	copy(res2, res1)
	return res2
}

func (ir *ifdReader) parseAperture(t tag.Tag) meta.Aperture {
	if t.IsType(tag.TypeRational) || t.IsType(tag.TypeSignedRational) {
		r := ir.parseRationalU(t)
		return meta.Aperture(float32(r[0]) / float32(r[1]))
	}
	ir.logParseWarn(t, "parseAperture", "unrecognized tag type", nil)
	return 0.0
}
func (ir *ifdReader) parseExposureTime(t tag.Tag) meta.ExposureTime {
	if t.IsType(tag.TypeRational) || t.IsType(tag.TypeSignedRational) {
		r := ir.parseRationalU(t)
		return meta.ExposureTime(float32(r[0]) / float32(r[1]))
	}
	ir.logParseWarn(t, "parseExposureTime", "unrecognized tag type", nil)
	return 0.0
}

func (ir *ifdReader) parseOrientation(t tag.Tag) meta.Orientation {
	return meta.Orientation(ir.parseUint16(t))
}

func (ir *ifdReader) parseExposureBias(t tag.Tag) meta.ExposureBias {
	if !t.IsEmbedded() {
		r := ir.parseRationalU(t)
		return meta.NewExposureBias(int16(r[0]), int16(r[1]))
	}
	return meta.NewExposureBias(0, 0)
}

// parseFocalLength supports tag type Rational, SRational, Short, and Long
func (ir *ifdReader) parseFocalLength(t tag.Tag) meta.FocalLength {
	if t.IsType(tag.TypeShort) || t.IsType(tag.TypeLong) {
		return meta.NewFocalLength(ir.parseUint32(t), 1)
	}
	if !t.IsEmbedded() && (t.IsType(tag.TypeRational) || t.IsType(tag.TypeSignedRational)) {
		buf, err := ir.readTagValue()
		if err != nil {
			ir.logParseWarn(t, "parseFocalLength", "", err)
			return meta.FocalLength(0)
		}

		return meta.FocalLength(float32(t.ByteOrder.Uint32(buf[:4])) / float32(t.ByteOrder.Uint32(buf[4:8])))
	}
	ir.logParseWarn(t, "parseFocalLength", "unrecognized tag type", nil)
	return meta.FocalLength(0)
}

func (ir *ifdReader) parseSubSecTime(t tag.Tag) uint16 {
	if t.IsEmbedded() && (t.IsType(tag.TypeASCII) || t.IsType(tag.TypeASCIINoNul)) {
		t.EmbeddedValue(ir.buffer.buf[:4])
		return uint16(parseStrUint(ir.buffer.buf[:4]))
	}
	return 0
}

func (ir *ifdReader) parseLensInfo(t tag.Tag) LensInfo {
	if !t.IsEmbedded() {
		buf, err := ir.readTagValue()
		if err != nil {
			ir.logParseWarn(t, "parseLensInfo", "", err)
			return LensInfo{}
		}

		return LensInfo{
			t.ByteOrder.Uint32(buf[:4]), t.ByteOrder.Uint32(buf[4:8]),
			t.ByteOrder.Uint32(buf[8:12]), t.ByteOrder.Uint32(buf[12:16]),
			t.ByteOrder.Uint32(buf[16:20]), t.ByteOrder.Uint32(buf[20:24]),
			t.ByteOrder.Uint32(buf[24:28]), t.ByteOrder.Uint32(buf[28:32])}
	}
	return LensInfo{}
}

func (ir *ifdReader) parseRationalU(t tag.Tag) [2]uint32 {
	if !t.IsEmbedded() {
		buf, err := ir.readTagValue()
		if err != nil {
			ir.logParseWarn(t, "parseRationalU", "", err)
			return [2]uint32{}
		}
		return [2]uint32{t.ByteOrder.Uint32(buf[:4]), t.ByteOrder.Uint32(buf[4:8])}
	}
	return [2]uint32{}
}

func (ir *ifdReader) parseUint32(t tag.Tag) uint32 {
	switch t.Type() {
	case tag.TypeLong:
		return uint32(t.ValueOffset)
	case tag.TypeShort:
		t.EmbeddedValue(ir.buffer.buf[:4])
		return uint32(t.ByteOrder.Uint16(ir.buffer.buf[:4]))
	default:
		return 0
	}
}

func (ir *ifdReader) parseUint16(t tag.Tag) uint16 {
	if t.IsEmbedded() && t.IsType(tag.TypeShort) {
		t.EmbeddedValue(ir.buffer.buf[:4])
		return t.ByteOrder.Uint16(ir.buffer.buf[:2])
	}
	return 0
}

func (ir *ifdReader) parseString(t tag.Tag) string {
	if t.IsEmbedded() {
		t.EmbeddedValue(ir.buffer.buf[:4])
		return trimNULString(ir.buffer.buf[:t.Size()])
	}
	if t.IsType(tag.TypeASCII) || t.IsType(tag.TypeASCIINoNul) {
		buf, err := ir.readTagValue()
		if err != nil {
			ir.logParseWarn(t, "parseString", "", err)
		}
		// Trim function
		return trimNULString(buf)
	}
	return ""
}

func (ir *ifdReader) parseString2(t tag.Tag) []byte {
	if t.IsEmbedded() {
		t.EmbeddedValue(ir.buffer.buf[:4])
		return trimNULBuffer(ir.buffer.buf[:t.Size()])
	}
	if t.IsType(tag.TypeASCII) || t.IsType(tag.TypeASCIINoNul) {
		buf, err := ir.readTagValue()
		if err != nil {
			ir.logParseWarn(t, "parseString", "", err)
		}
		// Trim function
		return trimNULBuffer(buf)
	}
	return nil
}

func (ir *ifdReader) parseDate(t tag.Tag) time.Time {
	if t.IsType(tag.TypeASCII) {
		buf, err := ir.readTagValue()
		if err != nil {
			ir.logParseWarn(t, "parseDate", "", err)
			return time.Time{}
		}
		// check recieved value
		if buf[4] == ':' && buf[7] == ':' && buf[10] == ' ' &&
			buf[13] == ':' && buf[16] == ':' {

			year := parseStrUint(buf[0:4])
			month := parseStrUint(buf[5:7])
			day := parseStrUint(buf[8:10])
			hour := parseStrUint(buf[11:13])
			min := parseStrUint(buf[14:16])
			sec := parseStrUint(buf[17:19])
			return time.Date(int(year), time.Month(month), int(day), int(hour), int(min), int(sec), 0, time.UTC)
		}
	}
	return time.Time{}
}

func (ir *ifdReader) parseOffsetTime(t tag.Tag) int16 {
	if t.IsType(tag.TypeASCII) {
		buf, err := ir.readTagValue()
		if err != nil {
			logTag(ir.logError(err), t).Send()
			return 0
		}
		var i int
		var res int16
		if buf[3] == ':' {
			hours := parseStrUint(buf[1:3])
			min := parseStrUint(buf[i+3 : i+5])
			res = int16(hours*60 + min)
			switch buf[0] {
			case '-':
				return res * -1
			case '+':
				return res
			default:
				if ir.logLevelWarn() {
					logTag(ir.logWarn(), t).Msgf("Uknown TimeOffset: %s", string(buf))
				}
				return 0
			}
		}
	}
	return 0
}

// parseGPSCoord parses the GPS Coordinate (Lat or Lng) from the corresponding Tag.
func (ir *ifdReader) parseGPSCoord(t tag.Tag) float64 {
	// Some cameras write tag out of spec using signed rational. We accept that too.
	if t.Type() != tag.TypeRational {
		if t.Type() != tag.TypeSignedRational {
			ir.logParseWarn(t, "parseGPSCoord", "error reading GPS Coord. Tag is not Rational or SRational.", nil)
			return 0.0
		}
	}
	buf, err := ir.readTagValue()
	if err != nil {
		ir.logParseWarn(t, "parseGPSCoord", "", err)
	}

	coord := (float64(t.ByteOrder.Uint32(buf[:4])) / float64(t.ByteOrder.Uint32(buf[4:8])))
	coord += (float64(t.ByteOrder.Uint32(buf[8:12])) / float64(t.ByteOrder.Uint32(buf[12:16])) / 60.0)
	coord += (float64(t.ByteOrder.Uint32(buf[16:20])) / float64(t.ByteOrder.Uint32(buf[20:24])) / 3600.0)

	return coord
}

// parseGPSAltitude parses the GPS Altitude from the corresponding Tag.
func (ir *ifdReader) parseGPSAltitude(t tag.Tag) float32 {
	// Some cameras write tag out of spec using signed rational. We accept that too.
	if t.IsType(tag.TypeRational) || t.IsType(tag.TypeSignedRational) {
		buf, err := ir.readTagValue()
		if err != nil {
			ir.logParseWarn(t, "parseGPSAltitude", "", err)
		}
		alt := (float32(t.ByteOrder.Uint32(buf[:4])) / float32(t.ByteOrder.Uint32(buf[4:8])))
		return alt
	}
	ir.logParseWarn(t, "parseGPSAltitude", "error reading GPS Alt. Tag is not Rational or SRational", nil)
	return 0.0
}

// parseGPSTimeStamp parses the GPSTimeStamp tag in UTC.
func (ir *ifdReader) parseGPSTimeStamp(t tag.Tag) uint32 {
	if t.UnitCount == 3 && t.Type() == tag.TypeRational {
		if !t.IsEmbedded() {
			buf, err := ir.readTagValue()
			if err != nil {
				ir.logParseWarn(t, "parseGPSTimeStamp", "", err)
				return 0
			}
			var result uint32
			value := [6]uint32{
				t.ByteOrder.Uint32(buf[:4]),
				t.ByteOrder.Uint32(buf[4:8]),
				t.ByteOrder.Uint32(buf[8:12]),
				t.ByteOrder.Uint32(buf[12:16]),
				t.ByteOrder.Uint32(buf[16:20]),
				t.ByteOrder.Uint32(buf[20:24])}
			if value[1] > 0 {
				result += (value[0] / value[1]) * 3600
			}
			if value[3] > 0 {
				result += (value[2] / value[3]) * 60
			}
			if value[5] > 0 {
				result += (value[4] / value[5])
			}

			return result
		}

	}
	ir.logParseWarn(t, "parseGPSTimeStamp", "error reading GPS Time Stamp", nil)
	return 0
}

func (ir *ifdReader) parseGPSDateStamp(t tag.Tag) time.Time {
	if t.IsType(tag.TypeASCII) {
		buf, err := ir.readTagValue()
		if err != nil {
			ir.logParseWarn(t, "parseGPSDateStamp", "error reading GPS Time Stamp", err)
			return time.Time{}
		}
		fmt.Println(string(buf))
		// check recieved value
		if buf[4] == ':' && buf[7] == ':' && len(buf) < 12 {
			return time.Date(int(parseStrUint(buf[0:4])), time.Month(parseStrUint(buf[5:7])), int(parseStrUint(buf[8:10])), 0, 0, 0, 0, time.UTC)
		}
		// check recieved value
		if buf[4] == ':' && buf[7] == ':' && buf[10] == ' ' &&
			buf[13] == ':' && buf[16] == ':' && len(buf) > 19 {
			return time.Date(
				int(parseStrUint(buf[0:4])),
				time.Month(parseStrUint(buf[5:7])),
				int(parseStrUint(buf[8:10])),
				int(parseStrUint(buf[11:13])),
				int(parseStrUint(buf[14:16])),
				int(parseStrUint(buf[17:19])), 0, time.UTC)
		}
	}
	return time.Time{}
}

func (ir *ifdReader) parseGPSRef(t tag.Tag) bool {
	if t.IsEmbedded() {
		t.EmbeddedValue(ir.buffer.buf[:4])
		switch t.ID {
		case gpsifd.GPSAltitudeRef:
			return t.IsType(tag.TypeByte) && (ir.buffer.buf[0] == byte(1))
		case gpsifd.GPSLatitudeRef:
			return t.IsType(tag.TypeASCII) && (ir.buffer.buf[0] == 'S')
		case gpsifd.GPSLongitudeRef:
			return t.IsType(tag.TypeASCII) && (ir.buffer.buf[0] == 'W')
		}
	}
	return false
}
