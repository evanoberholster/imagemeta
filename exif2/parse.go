package exif2

import (
	"time"

	"github.com/evanoberholster/imagemeta/exif2/ifds/gpsifd"
	"github.com/evanoberholster/imagemeta/exif2/tag"
	"github.com/evanoberholster/imagemeta/meta"
)

func (ir *ifdReader) parseAperture(t tag.Tag) meta.Aperture {
	if t.IsType(tag.TypeRational) || t.IsType(tag.TypeSignedRational) {
		return meta.Aperture(ir.parseRationalU(t).AsFloat())
	}
	if ir.logWarn() {
		ir.logger.Warn().Str("func", "parseAperture").Uint32("units", t.UnitCount).Stringer("id", t.ID).Stringer("ifd", t.Type()).Uint32("size", t.Size()).Send()
	}
	return 0.0
}
func (ir *ifdReader) parseExposureTime(t tag.Tag) meta.ExposureTime {
	r := ir.parseRationalU(t)
	return meta.ExposureTime(float32(r[0]) / float32(r[1]))
}

func (ir *ifdReader) parseOrientation(t tag.Tag) meta.Orientation {
	return meta.Orientation(ir.parseUint16(t))
}

func (ir *ifdReader) parseExposureBias(t tag.Tag) meta.ExposureBias {
	eb := ir.parseRationalU(t)
	return meta.NewExposureBias(int16(eb[0]), int16(eb[1]))
}

func (ir *ifdReader) parseFocalLength(t tag.Tag) meta.FocalLength {
	if t.IsType(tag.TypeShort) || t.IsType(tag.TypeLong) {
		return meta.NewFocalLength(ir.parseUint32(t), 1)
	}
	r := ir.parseRationalU(t)
	return meta.NewFocalLength(r[0], r[1])
}

func (ir *ifdReader) parseSubSecTime(t tag.Tag) uint16 {
	if t.IsEmbedded() && (t.IsType(tag.TypeASCII) || t.IsType(tag.TypeASCIINoNul)) {
		t.ByteOrder.PutUint32(ir.buffer.buf[:4], t.ValueOffset)
		return uint16(parseStrUint(ir.buffer.buf[:4]))
	}
	return 0
}

func (ir *ifdReader) parseLensInfo(t tag.Tag) LensInfo {
	if !t.IsEmbedded() {
		buf, err := ir.readTagValue()
		if err != nil {
			if ir.logWarn() {
				ir.logger.Warn().Err(err).Str("func", "parseLensInfo").Uint32("units", t.UnitCount).Stringer("id", t.ID).Stringer("ifd", t.Type()).Uint32("size", t.Size()).Send()
			}
			return LensInfo{}
		}

		return LensInfo{
			RationalU{t.ByteOrder.Uint32(buf[:4]), t.ByteOrder.Uint32(buf[4:8])},
			RationalU{t.ByteOrder.Uint32(buf[8:12]), t.ByteOrder.Uint32(buf[12:16])},
			RationalU{t.ByteOrder.Uint32(buf[16:20]), t.ByteOrder.Uint32(buf[20:24])},
			RationalU{t.ByteOrder.Uint32(buf[24:28]), t.ByteOrder.Uint32(buf[28:32])}}
	}
	return LensInfo{}
}

func (ir *ifdReader) parseRationalU(t tag.Tag) RationalU {
	if !t.IsEmbedded() {
		buf, err := ir.readTagValue()
		if err != nil {
			if ir.logWarn() {
				ir.logger.Warn().Err(err).Str("func", "parseRationalU").Uint32("units", t.UnitCount).Stringer("id", t.ID).Stringer("ifd", t.Type()).Uint32("size", t.Size()).Send()
			}
			return RationalU{}
		}
		return RationalU{t.ByteOrder.Uint32(buf[:4]), t.ByteOrder.Uint32(buf[4:8])}
	}
	return RationalU{}
}

func (ir *ifdReader) parseUint32(t tag.Tag) uint32 {
	switch t.Type() {
	case tag.TypeLong:
		return t.ValueOffset
	case tag.TypeShort:
		t.ByteOrder.PutUint32(ir.buffer.buf[:4], t.ValueOffset)
		return uint32(t.ByteOrder.Uint16(ir.buffer.buf[:4]))
	default:
		return 0
	}
}

func (ir *ifdReader) parseUint16(t tag.Tag) uint16 {
	if t.IsEmbedded() && t.IsType(tag.TypeShort) {
		t.ByteOrder.PutUint32(ir.buffer.buf[:4], t.ValueOffset)
		return t.ByteOrder.Uint16(ir.buffer.buf[:2])
	}
	return 0
}

func (ir *ifdReader) parseString(t tag.Tag) string {
	if t.IsEmbedded() {
		t.ByteOrder.PutUint32(ir.buffer.buf[:4], t.ValueOffset)
		return trimNULString(ir.buffer.buf[:t.Size()])
	}
	if t.IsType(tag.TypeASCII) || t.IsType(tag.TypeASCIINoNul) {
		buf, err := ir.readTagValue()
		if err != nil {
			if ir.logWarn() {
				ir.logger.Warn().Err(err).Str("func", "parseString").Uint32("units", t.UnitCount).Stringer("id", t.ID).Stringer("ifd", t.Type()).Uint32("size", t.Size()).Send()
			}
		}
		// Trim function
		return trimNULString(buf)
	}
	return ""
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
		t.ByteOrder.PutUint32(ir.buffer.buf[:4], t.ValueOffset)
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
