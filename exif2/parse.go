package exif2

import (
	"fmt"
	"math"
	"time"

	"github.com/evanoberholster/imagemeta/exif2/ifds"
	"github.com/evanoberholster/imagemeta/exif2/ifds/exififd"
	"github.com/evanoberholster/imagemeta/exif2/ifds/gpsifd"
	"github.com/evanoberholster/imagemeta/exif2/ifds/mknote/apple"
	"github.com/evanoberholster/imagemeta/exif2/ifds/mknote/canon"
	"github.com/evanoberholster/imagemeta/exif2/make"
	"github.com/evanoberholster/imagemeta/exif2/model"
	"github.com/evanoberholster/imagemeta/exif2/tag"
	"github.com/evanoberholster/imagemeta/imagetype"
	"github.com/evanoberholster/imagemeta/meta"
)

func (ir *ifdReader) parseTag(t Tag) {
	if ir.customTagParser != nil {
		if err := ir.customTagParser(ir, t); err != nil {
			ir.logError(err).Send()
		}
		return
	}
	switch ifds.IfdType(t.Ifd) {
	case ifds.IFD0:
		switch t.ID {
		case ifds.Make:
			ir.Exif.CameraMake, ir.Exif.Make = ir.ParseCameraMake(t)
		case ifds.Model:
			ir.Exif.CameraModel, ir.Exif.Model = ir.ParseCameraModel(t)
		case ifds.Artist:
			ir.Exif.Artist = ir.ParseString(t)
		case ifds.Copyright:
			ir.Exif.Copyright = ir.ParseString(t)
		case ifds.ImageWidth:
			ir.Exif.ImageWidth = uint16(ir.ParseUint32(t))
		case ifds.ImageLength:
			ir.Exif.ImageHeight = uint16(ir.ParseUint32(t))
		case ifds.StripOffsets:
			ir.Exif.StripOffsets = ir.ParseUint32(t)
		case ifds.StripByteCounts:
			ir.Exif.StripByteCounts = ir.ParseUint32(t)
		case ifds.Orientation:
			ir.Exif.Orientation = meta.Orientation(ir.ParseUint16(t))
		case ifds.Software:
			ir.Exif.Software = ir.ParseString(t)
		case ifds.ImageDescription:
			ir.Exif.ImageDescription = ir.ParseString(t)
		case ifds.DateTime:
			ir.Exif.Time.modifyDate = ir.ParseDate(t)
		case ifds.DNGVersion:
			// TODO: If DNG version > 0 imagetype is DNG
			if ir.Exif.ImageType == imagetype.ImageTiff {
				ir.Exif.ImageType = imagetype.ImageDNG
			}

		case ifds.CameraSerialNumber:
			if ir.Exif.CameraSerial == "" {
				ir.Exif.CameraSerial = ir.ParseString(t)
			}
		//case ifds.ApplicationNotes:
		//fmt.Println(ir.parseApplicationNotes(t))
		//ir.Exif.ApplicationNotes =
		default:
			//t.logTag(ir.logWarn()).Send()
		}
	case ifds.ExifIFD:
		switch t.ID {
		case exififd.LensMake:
			ir.Exif.LensMake = ir.ParseString(t)
		case exififd.LensModel:
			ir.Exif.LensModel = ir.ParseString(t)
		case exififd.LensSerialNumber:
			ir.Exif.LensSerial = ir.ParseString(t)
		case exififd.CameraOwnerName:
			if ir.Exif.Artist == "" {
				ir.Exif.Artist = ir.ParseString(t)
			}
		case exififd.BodySerialNumber:
			if ir.Exif.CameraSerial == "" {
				ir.Exif.CameraSerial = ir.ParseString(t)
			}
		case exififd.PixelXDimension:
			if ir.Exif.ImageWidth == 0 {
				ir.Exif.ImageWidth = uint16(ir.ParseUint32(t))
			}
		case exififd.PixelYDimension:
			if ir.Exif.ImageHeight == 0 {
				ir.Exif.ImageHeight = uint16(ir.ParseUint32(t))
			}
		case exififd.ExposureTime:
			ir.Exif.ExposureTime = ir.parseExposureTime(t)
		case exififd.ApertureValue:
			if ir.Exif.FNumber == 0.0 {
				r := ir.ParseRationalU(t)
				f := float64(r[0]) / float64(r[1])
				ir.Exif.FNumber = meta.Aperture(math.Round(math.Pow(math.Sqrt2, float64(f))*100) / 100)
			}
		case exififd.FNumber:
			ir.Exif.FNumber = ir.parseAperture(t)
		case exififd.ExposureProgram:
			ir.Exif.ExposureProgram = meta.ExposureProgram(ir.ParseUint16(t))
		case exififd.ExposureBiasValue:
			ir.Exif.ExposureBias = ir.parseExposureBias(t)
		case exififd.ExposureMode:
			ir.Exif.ExposureMode = meta.ExposureMode(ir.ParseUint16(t))
		case exififd.MeteringMode:
			ir.Exif.MeteringMode = meta.MeteringMode(ir.ParseUint16(t))
		case exififd.ISOSpeed:
			ir.Exif.ISOSpeed = ir.ParseUint32(t)
		case ifds.Flash:
			ir.Exif.Flash = meta.Flash(ir.ParseUint16(t))
		case ifds.FocalLength:
			ir.Exif.FocalLength = ir.parseFocalLength(t)
		case exififd.FocalLengthIn35mmFilm:
			ir.Exif.FocalLengthIn35mmFormat = ir.parseFocalLength(t)
		case exififd.LensSpecification:
			ir.Exif.LensInfo = ir.parseLensInfo(t)
		case ifds.DateTimeOriginal:
			ir.Exif.Time.dateTimeOriginal = ir.ParseDate(t)
		case ifds.DateTimeDigitized:
			ir.Exif.Time.createDate = ir.ParseDate(t)
		case exififd.SubSecTime:
			ir.Exif.Time.subSecTime = ir.ParseSubSecTime(t)
		case exififd.SubSecTimeOriginal:
			ir.Exif.Time.subSecTimeOriginal = ir.ParseSubSecTime(t)
		case exififd.SubSecTimeDigitized:
			ir.Exif.Time.subSecTimeDigitized = ir.ParseSubSecTime(t)
		case exififd.OffsetTime:
			ir.Exif.Time.offsetTime = ir.ParseOffsetTime(t)
		case exififd.OffsetTimeOriginal:
			ir.Exif.Time.offsetTimeOriginal = ir.ParseOffsetTime(t)
		case exififd.OffsetTimeDigitized:
			ir.Exif.Time.offsetTimeDigitized = ir.ParseOffsetTime(t)
		default:
			//t.logTag(ir.logWarn()).Send()
		}
	case ifds.GPSIFD:
		ir.parseGPS(t)
	}
}

func (ir *ifdReader) parseGPS(t Tag) {
	var err error
	if ifds.IfdType(t.Ifd) == ifds.GPSIFD {
		switch t.ID {
		case gpsifd.GPSMapDatum: //string                 // 0x0012 string
			err = ir.Exif.GPS.GPSMapDatum.FromBytes(ir.TagValue(t))
		case gpsifd.GPSSatellites: //string                 // 0x0008 string
			ir.Exif.GPS.GPSSatellites = string(ir.ParseBytes(t))
		case gpsifd.GPSProcessingMethod: //[]byte                 // 0x001b undef	(values of "GPS", "CELLID", "WLAN" or "MANUAL" by the EXIF spec.)
			ir.Exif.GPS.GPSProcessingMethod = tag.NewGPSProcessingMethod(ir.ParseBytes(t))
		case gpsifd.GPSAreaInformation: //[]byte                 // 0x001c undef
			ir.Exif.GPS.GPSAreaInformation = append(ir.Exif.GPS.GPSAreaInformation, ir.ParseBytes(t)...)
		case gpsifd.GPSVersionID: //[4]uint8               // 0x0000 int8u[4]
			err = ir.Exif.GPS.GPSVersionID.FromBytes(ir.TagValue(t))
		case gpsifd.GPSDateStamp: //[11]byte               // 0x001d string[11] (time is stripped off if present, after adjusting date/time to UTC if time includes a timezone. Format is YYYY:mm:dd)
			ir.ParseGPSDateStamp(&ir.Exif.GPS.GPSDateStamp, t)
		case gpsifd.GPSLatitude:
			err = ir.Exif.GPS.GPSLatitude.FromBytes(ir.TagValue(t)) //tag.GPSCoordinate      // 0x0002 rational64u[3]
		case gpsifd.GPSLongitude: //tag.GPSCoordinate      // 0x0004 rational64u[3]
			err = ir.Exif.GPS.GPSLongitude.FromBytes(ir.TagValue(t))
		case gpsifd.GPSDestLatitude: //tag.GPSCoordinate      // 0x0014 rational64u[3]
			err = ir.Exif.GPS.GPSDestLatitude.FromBytes(ir.TagValue(t))
		case gpsifd.GPSDestLongitude: //tag.GPSCoordinate      // 0x0016 rational64u[3]
			err = ir.Exif.GPS.GPSDestLongitude.FromBytes(ir.TagValue(t))
		case gpsifd.GPSTimeStamp: //[3]tag.Rational64      // 0x0007 rational64u[3] (UTC time of GPS fix)
			ir.Exif.GPS.GPSTimeStamp = tag.GPSTimeStamp(ir.ParseGPSCoordinate(t))
		case gpsifd.GPSAltitude: //tag.Rational64         // 0x0006 rational64u
			ir.Exif.GPS.GPSAltitude = ir.ParseRational(t)
		case gpsifd.GPSDOP: //tag.Rational64         // 0x000b rational64u
			ir.Exif.GPS.GPSDOP = ir.ParseRational(t)
		case gpsifd.GPSSpeed: //tag.Rational64         // 0x000d rational64u
			ir.Exif.GPS.GPSSpeed = ir.ParseRational(t)
		case gpsifd.GPSTrack: //tag.Rational64         // 0x000f rational64u
			ir.Exif.GPS.GPSTrack = ir.ParseRational(t)
		case gpsifd.GPSImgDirection: //tag.Rational64         // 0x0011 rational64u
			ir.Exif.GPS.GPSImgDirection = ir.ParseRational(t)
		case gpsifd.GPSDestBearing: //tag.Rational64         // 0x0018 rational64u
			ir.Exif.GPS.GPSDestBearing = ir.ParseRational(t)
		case gpsifd.GPSDestDistance: //tag.Rational64         // 0x001a rational64u
			ir.Exif.GPS.GPSDestDistance = ir.ParseRational(t)
		case gpsifd.GPSHPositioningError: //tag.Rational64         // 0x001f rational64u
			ir.Exif.GPS.GPSHPositioningError = ir.ParseRational(t)
		case gpsifd.GPSLatitudeRef: //tag.GPSLatitudeRef     // 0x0001 string[2]
			ir.Exif.GPS.GPSLatitudeRef = tag.GPSLatitudeRef(ir.ParseGPSRef(t))
		case gpsifd.GPSLongitudeRef: //tag.GPSLongitudeRef    // 0x0003 string[2] (ExifTool will also accept a number when writing this tag, positive for east longitudes or negative for west, or a string containing E, East, W or West) 'E' = East 'W' = West
			ir.Exif.GPS.GPSLongitudeRef = tag.GPSLongitudeRef(ir.ParseGPSRef(t))
		case gpsifd.GPSAltitudeRef: //tag.GPSAltitudeRef     // 0x0005 int8u (ExifTool will also accept number when writing this tag, with negative numbers indicating below sea level) 0 = Above Sea Level 1 = Below Sea Level
			ir.Exif.GPS.GPSAltitudeRef = tag.GPSAltitudeRef(ir.ParseGPSRef(t))
		case gpsifd.GPSStatus: //tag.GPSStatus          // 0x0009 string[2]	'A' = Measurement Active 'V' = Measurement Void
			ir.Exif.GPS.GPSStatus = tag.GPSStatus(ir.ParseGPSRef(t))
		case gpsifd.GPSMeasureMode: //tag.GPSMeasureMode     // 0x000a string[2]	2 = 2-Dimensional Measurement 3 = 3-Dimensional Measurement
			ir.Exif.GPS.GPSMeasureMode = tag.GPSMeasureMode(ir.ParseGPSRef(t))
		case gpsifd.GPSSpeedRef: //tag.GPSSpeedRef        // 0x000c string[2]	'K' = km/h 'M' = mph 'N' = knots
			ir.Exif.GPS.GPSSpeedRef = tag.GPSSpeedRef(ir.ParseGPSRef(t))
		case gpsifd.GPSTrackRef: //tag.GPSDestBearingRef  // 0x000e string[2]	'M' = Magnetic North 'T' = True North
			ir.Exif.GPS.GPSTrackRef = tag.GPSDestBearingRef(ir.ParseGPSRef(t))
		case gpsifd.GPSImgDirectionRef: //tag.GPSDestBearingRef  // 0x0010 string[2]	'M' = Magnetic North  'T' = True North
			ir.Exif.GPS.GPSDestBearingRef = tag.GPSDestBearingRef(ir.ParseGPSRef(t))
		case gpsifd.GPSDestLatitudeRef: //tag.GPSLatitudeRef     // 0x0013 string[2] (tags 0x0013-0x001a used for subject location according to MWG 2.0)  'N' = North  'S' = South
			ir.Exif.GPS.GPSDestLatitudeRef = tag.GPSLatitudeRef(ir.ParseGPSRef(t))
		case gpsifd.GPSDestLongitudeRef: //tag.GPSLongitudeRef    // 0x0015 string[2]	'E' = East  'W' = West
			ir.Exif.GPS.GPSDestLongitudeRef = tag.GPSLongitudeRef(ir.ParseGPSRef(t))
		case gpsifd.GPSDestBearingRef: //tag.GPSDestBearingRef  // 0x0017 string[2]	'M' = Magnetic North  'T' = True North
			ir.Exif.GPS.GPSDestBearingRef = tag.GPSDestBearingRef(ir.ParseGPSRef(t))
		case gpsifd.GPSDestDistanceRef: //tag.GPSDestDistanceRef // 0x0019 string[2]	'K' = Kilometers 'M' = Miles 'N' = Nautical Miles
			ir.Exif.GPS.GPSDestDistanceRef = tag.GPSDestDistanceRef(ir.ParseGPSRef(t))
		case gpsifd.GPSDifferential: //uint16                 // 0x001e int16u	0 = No Correction 1 = Differential Corrected
			// TODO
		}
	}
	if err != nil {
		fmt.Println(err)
	}
}

func (ir *ifdReader) ParseCameraMake(t Tag) (make.CameraMake, string) {
	str := ir.ParseBuffer(t)
	if mk, ok := make.CameraMakeFromString(string(str)); ok {
		return mk, mk.String()
	}
	return make.CameraMakeUnknown, string(str)
}

func (ir *ifdReader) ParseCameraModel(t Tag) (model.CameraModel, string) {
	str := ir.ParseBuffer(t)
	switch ir.Exif.CameraMake {
	case make.Canon:
		if m, ok := canon.CameraModelFromString(string(str)); ok {
			return model.CameraModel(m), m.String()
		}
	case make.Apple:
		if m, ok := apple.CameraModelFromString(string(str)); ok {
			ir.Exif.CameraModel = model.CameraModel(m)
			return model.CameraModel(m), m.String()
		}
	}
	return model.CameraModelUnknown, string(str)
}

//func (ir *ifdReader) parseApplicationNotes(t Tag) ApplicationNotes {
//	if err := ir.discard(int(t.ValueOffset) - int(ir.po)); err != nil {
//		return nil
//	}
//	n := t.Size()
//	if t.Size() > bufferLength {
//		n = bufferLength
//	}
//	buf, err := ir.fastRead(int(n))
//	if err != nil {
//		return nil
//	}
//	res1 := trimNULBuffer(buf)
//	res2 := make([]byte, len(res1))
//	copy(res2, res1)
//	return res2
//}

func (ir *ifdReader) parseAperture(t Tag) meta.Aperture {
	if t.IsType(tag.TypeRational) || t.IsType(tag.TypeSignedRational) {
		r := ir.ParseRationalU(t)
		return meta.Aperture(float32(r[0]) / float32(r[1]))
	}
	if ir.logLevelWarn() {
		t.logTag(ir.logWarn()).Msg("Unrecognized tag type")
	}
	return 0.0
}
func (ir *ifdReader) parseExposureTime(t Tag) meta.ExposureTime {
	if t.IsType(tag.TypeRational) || t.IsType(tag.TypeSignedRational) {
		r := ir.ParseRationalU(t)
		return meta.ExposureTime(float32(r[0]) / float32(r[1]))
	}
	if ir.logLevelWarn() {
		t.logTag(ir.logWarn()).Msg("Unrecognized tag type")
	}
	return 0.0
}

func (ir *ifdReader) parseExposureBias(t Tag) meta.ExposureBias {
	if !t.IsEmbedded() {
		r := ir.ParseRationalU(t)
		return meta.NewExposureBias(int16(r[0]), int16(r[1]))
	}
	if ir.logLevelWarn() {
		t.logTag(ir.logWarn()).Msg("Unrecognized tag type")
	}
	return meta.NewExposureBias(0, 0)
}

// parseFocalLength supports tag type Rational, SRational, Short, and Long
func (ir *ifdReader) parseFocalLength(t Tag) meta.FocalLength {
	switch t.Type {
	case tag.TypeShort, tag.TypeLong:
		return meta.NewFocalLength(ir.ParseUint32(t), 1)
	case tag.TypeRational, tag.TypeSignedRational:
		val := ir.ParseRationalU(t)
		return meta.FocalLength(float32(val[0]) / float32(val[1]))
	}
	if ir.logLevelWarn() {
		t.logTag(ir.logWarn()).Msg("Unrecognized tag type")
	}
	return meta.FocalLength(0)
}

// ParseSubSecTime parses an ASCII or ASCII no Nul.
// Embedded tag with value length 4 bytes.
// Value is in milliseconds.
func (ir *ifdReader) ParseSubSecTime(t Tag) uint16 {
	if t.IsType(tag.TypeASCII) || t.IsType(tag.TypeASCIINoNul) {
		if t.IsEmbedded() {
			t.EmbeddedValue(ir.buffer.buf[:4])
			return uint16(parseStrUint(ir.buffer.buf[:4]))
		}
		buf := ir.ParseBuffer(t)
		return uint16(parseStrUint(buf) / 1000)
	}
	if ir.logLevelWarn() {
		t.logTag(ir.logWarn()).Msg("Unrecognized tag type")
	}
	return 0
}

func (ir *ifdReader) parseLensInfo(t Tag) LensInfo {
	if !t.IsEmbedded() {
		buf, err := ir.readTagValue()
		if err != nil {
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

func (ir *ifdReader) ParseRational(t Tag) tag.Rational {
	switch t.Type {
	case tag.TypeSignedRational, tag.TypeRational:
		buf, err := ir.readTagValue()
		if err != nil {
			return tag.Rational{}
		}
		return tag.Rational{t.ByteOrder.Uint32(buf[:4]), t.ByteOrder.Uint32(buf[4:8])}
	default:
		if ir.logLevelWarn() {
			t.logTag(ir.logWarn()).Msg("Unrecognized tag type")
		}
	}
	return tag.Rational{}
}

// ParseRationalU parses an Unsigned Rational value.
// Non-embedded tag with value length 8 bytes.
func (ir *ifdReader) ParseRationalU(t Tag) [2]uint32 {
	switch t.Type {
	case tag.TypeSignedRational, tag.TypeRational:
		buf, err := ir.readTagValue()
		if err != nil {
			return [2]uint32{}
		}
		return [2]uint32{t.ByteOrder.Uint32(buf[:4]), t.ByteOrder.Uint32(buf[4:8])}
	default:
		if ir.logLevelWarn() {
			t.logTag(ir.logWarn()).Msg("Unrecognized tag type")
		}
	}
	return [2]uint32{}
}

// ParseUint32 parses a Uint32 value.
// Embedded tag with value length 4 bytes.
func (ir *ifdReader) ParseUint32(t Tag) uint32 {
	switch t.Type {
	case tag.TypeLong:
		return uint32(t.ValueOffset)
	case tag.TypeShort:
		t.EmbeddedValue(ir.buffer.buf[:4])
		return uint32(t.ByteOrder.Uint16(ir.buffer.buf[:4]))
	default:
		if ir.logLevelWarn() {
			t.logTag(ir.logWarn()).Msg("Unrecognized tag type")
		}
	}
	return 0
}

// ParseUint16 parses a uint16 value.
// Embedded tag with value length 2 bytes.
func (ir *ifdReader) ParseUint16(t Tag) uint16 {
	if t.IsEmbedded() && t.IsType(tag.TypeShort) {
		t.EmbeddedValue(ir.buffer.buf[:4])
		return t.ByteOrder.Uint16(ir.buffer.buf[:4])
	}
	if ir.logLevelWarn() {
		t.logTag(ir.logWarn()).Msg("Unrecognized tag type")
	}
	return 0
}

// ParseString parses an ASCII value.
// Non-embedded or embedded tag with variable byte length.
// This function allocates.
func (ir *ifdReader) ParseString(t Tag) string {
	if t.IsEmbedded() {
		t.EmbeddedValue(ir.buffer.buf[:4])
		return string(trimNULBuffer(ir.buffer.buf[:t.Size()]))
	}
	if t.IsType(tag.TypeASCII) || t.IsType(tag.TypeASCIINoNul) {
		buf, _ := ir.readTagValue()
		return string(trimNULBuffer(buf)) // Trim function
	}
	if ir.logLevelWarn() {
		t.logTag(ir.logWarn()).Msg("Unrecognized tag type")
	}
	return ""
}

// ParseBuffer parses an ASCII value.
// Non-embedded or embedded tag with variable byte length.
// This function does not allocate.
func (ir *ifdReader) ParseBuffer(t Tag) []byte {
	if t.IsEmbedded() {
		t.EmbeddedValue(ir.buffer.buf[:4])
		return trimNULBuffer(ir.buffer.buf[:t.Size()])
	}
	if t.IsType(tag.TypeASCII) || t.IsType(tag.TypeASCIINoNul) {
		buf, err := ir.readTagValue()
		if err != nil {
			return nil
		}
		// Trim function
		return trimNULBuffer(buf)
	}
	if ir.logLevelWarn() {
		t.logTag(ir.logWarn()).Msg("Unrecognized tag type")
	}
	return nil
}

// TagValue
func (ir *ifdReader) TagValue(t Tag) tag.TagValue {
	if t.IsEmbedded() {
		t.EmbeddedValue(ir.buffer.buf[:4])
		return tag.TagValue{
			Buf:       trimNULBuffer(ir.buffer.buf[:t.Size()]),
			UnitCount: t.UnitCount,
			ID:        t.ID,
			Type:      t.Type,
			ByteOrder: t.ByteOrder,
		}
	}
	buf, err := ir.readTagValue()
	if err != nil {
		return tag.TagValue{}
	}
	return tag.TagValue{
		Buf:       buf,
		UnitCount: t.UnitCount,
		ID:        t.ID,
		Type:      t.Type,
		ByteOrder: t.ByteOrder,
	}
}

// ParseBytes parses an ASCII value, Non-Null ASCII value, or an undefined value.
// Non-embedded or embedded tag with variable byte length.
// This function does not allocate.
func (ir *ifdReader) ParseBytes(t Tag) []byte {
	if t.IsEmbedded() {
		t.EmbeddedValue(ir.buffer.buf[:4])
		return trimNULBuffer(ir.buffer.buf[:t.Size()])
	}
	if t.IsType(tag.TypeASCII) || t.IsType(tag.TypeASCIINoNul) || t.IsType(tag.TypeUndefined) {
		buf, err := ir.readTagValue()
		if err != nil {
			return nil
		}
		// Trim function
		return trimNULBuffer(buf)
	}
	if ir.logLevelWarn() {
		t.logTag(ir.logWarn()).Msg("Unrecognized tag type")
	}
	return nil
}

// ParseDate parses an ASCII value as a Date.
// Non-embedded tag with 20 byte length.
func (ir *ifdReader) ParseDate(t Tag) time.Time {
	if t.IsType(tag.TypeASCII) {
		buf, err := ir.readTagValue()
		if err != nil {
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
	if ir.logLevelWarn() {
		t.logTag(ir.logWarn()).Msg("Unrecognized tag type")
	}
	return time.Time{}
}

// ParseOffsetTime parses an ASCII value as a Timezone.
// Non-embedded tag with 6 byte length.
func (ir *ifdReader) ParseOffsetTime(t Tag) *time.Location {
	if t.IsType(tag.TypeASCII) {
		buf, err := ir.readTagValue()
		if err != nil {
			return time.UTC
		}
		if buf[3] == ':' {
			var offset int
			offset += int(parseStrUint(buf[1:3])) * hoursToSeconds
			offset += int(parseStrUint(buf[4:6])) * minutesToSeconds
			switch buf[0] {
			case '-':
				return getLocation(int32(offset*-1), buf[:6])
				//return time.FixedZone(string(buf[:6]), offset*-1)
			case '+':
				return getLocation(int32(offset), buf[:6])
				//return time.FixedZone(string(buf[:6]), offset)
			default:
				if ir.logLevelWarn() {
					t.logTag(ir.logWarn()).Msgf("Uknown TimeOffset: %s", string(buf))
				}
				return time.UTC
			}
		}
	}
	if ir.logLevelWarn() {
		t.logTag(ir.logWarn()).Msg("Unrecognized tag type")
	}
	return time.UTC
}

func (ir *ifdReader) ParseGPSCoordinate(t Tag) tag.GPSCoordinate {
	if t.UnitCount == 3 {
		switch t.Type {
		case tag.TypeRational, tag.TypeSignedRational: // Some cameras write tag out of spec using signed rational. We accept that too.
			buf, err := ir.readTagValue()
			if err != nil {
				return tag.GPSCoordinate{}
			}
			return tag.NewGPSCoordinate(
				t.ByteOrder.Uint32(buf[:4]), t.ByteOrder.Uint32(buf[4:8]),
				t.ByteOrder.Uint32(buf[8:12]), t.ByteOrder.Uint32(buf[12:16]),
				t.ByteOrder.Uint32(buf[16:20]), t.ByteOrder.Uint32(buf[20:24]))
		}
	}
	if ir.logLevelWarn() {
		t.logTag(ir.logWarn()).Msg("error reading GPS Coord. Tag is not Rational or SRational")
	}
	return tag.GPSCoordinate{}
}

// parseGPSCoord parses the GPS Coordinate (Lat or Lng) from the corresponding Tag.
//func (ir *ifdReader) parseGPSCoord(t Tag) float64 {
//	if t.UnitCount == 3 {
//		switch t.Type {
//		case tag.TypeRational, tag.TypeSignedRational: // Some cameras write tag out of spec using signed rational. We accept that too.
//			buf, err := ir.readTagValue()
//			if err != nil {
//				return 0.0
//			}
//			coord := (float64(t.ByteOrder.Uint32(buf[:4])) / float64(t.ByteOrder.Uint32(buf[4:8])))
//			coord += (float64(t.ByteOrder.Uint32(buf[8:12])) / float64(t.ByteOrder.Uint32(buf[12:16])) / 60.0)
//			coord += (float64(t.ByteOrder.Uint32(buf[16:20])) / float64(t.ByteOrder.Uint32(buf[20:24])) / 3600.0)
//			return coord
//		}
//	}
//	if ir.logLevelWarn() {
//		t.logTag(ir.logWarn()).Msg("error reading GPS Coord. Tag is not Rational or SRational")
//	}
//	return 0.0
//}

// parseGPSAltitude parses the GPS Altitude from the corresponding Tag.
func (ir *ifdReader) parseGPSAltitude(t Tag) float32 {
	if t.UnitCount == 1 {
		switch t.Type {
		case tag.TypeRational, tag.TypeSignedRational: // Some cameras write tag out of spec using signed rational. We accept that too.
			buf, err := ir.readTagValue()
			if err != nil {
				return 0.0
			}
			return (float32(t.ByteOrder.Uint32(buf[:4])) / float32(t.ByteOrder.Uint32(buf[4:8])))
		}
	}
	if ir.logLevelWarn() {
		t.logTag(ir.logWarn()).Msg("error reading GPS Alt. Tag is not Rational or SRational")
	}
	return 0.0
}

// parseGPSTimeStamp parses the GPSTimeStamp tag in UTC.
func (ir *ifdReader) parseGPSTimeStamp(t Tag) uint32 {
	if t.UnitCount == 3 && t.Type == tag.TypeRational {
		buf, err := ir.readTagValue()
		if err != nil {
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
			result += (value[0] / value[1]) * hoursToSeconds
		}
		if value[3] > 0 {
			result += (value[2] / value[3]) * minutesToSeconds
		}
		if value[5] > 0 {
			result += (value[4] / value[5])
		}
		return result
	}
	if ir.logLevelWarn() {
		t.logTag(ir.logWarn()).Msg("error reading GPS Time Stamp")
	}
	return 0
}

// parseGPSDateStamp parses a GPSDateStamp from the tag
func (ir *ifdReader) parseGPSDateStamp(t Tag) time.Time {
	if t.IsType(tag.TypeASCII) {
		buf, err := ir.readTagValue()
		if err != nil {
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
	if ir.logLevelWarn() {
		t.logTag(ir.logWarn()).Msg("error reading GPSDateStamp")
	}
	return time.Time{}
}

// ParseGPSDateStamp parses a GPSDateStamp
func (ir *ifdReader) ParseGPSDateStamp(ds *tag.GPSDateStamp, t Tag) {
	if t.IsType(tag.TypeASCII) {
		buf, err := ir.readTagValue()
		if err != nil {
			return
		}
		// check recieved value
		if buf[4] == ':' && buf[7] == ':' && len(buf) < 12 {
			*ds = tag.NewGPSDateStamp(uint16(parseStrUint(buf[0:4])), uint8(parseStrUint(buf[5:7])), uint8(parseStrUint(buf[8:10])))
			//return time.Date(int(parseStrUint(buf[0:4])), time.Month(parseStrUint(buf[5:7])), int(parseStrUint(buf[8:10])), 0, 0, 0, 0, time.UTC)
			return
		}
		// check recieved value
		if buf[4] == ':' && buf[7] == ':' && buf[10] == ' ' &&
			buf[13] == ':' && buf[16] == ':' && len(buf) > 19 {

			*ds = tag.NewGPSDateStamp(uint16(parseStrUint(buf[0:4])), uint8((parseStrUint(buf[5:7]))), uint8(parseStrUint(buf[8:10])))
			return
			//return time.Date(
			//	int(parseStrUint(buf[0:4])),
			//	time.Month(parseStrUint(buf[5:7])),
			//	int(parseStrUint(buf[8:10])),
			//	int(parseStrUint(buf[11:13])),
			//	int(parseStrUint(buf[14:16])),
			//	int(parseStrUint(buf[17:19])), 0, time.UTC)
		}
	}
	if ir.logLevelWarn() {
		t.logTag(ir.logWarn()).Msg("error reading GPSDateStamp")
	}
}

func (ir *ifdReader) ParseGPSRef(t Tag) uint8 {
	if t.IsEmbedded() {
		t.EmbeddedValue(ir.buffer.buf[:4])
		switch t.Type {
		case tag.TypeASCII, tag.TypeByte:
			return ir.buffer.buf[0]
		}
	}
	if ir.logLevelWarn() {
		t.logTag(ir.logWarn()).Msg("error reading GPS Reference")
	}
	return 0
}

// parseGPSRef parsese the GPS Reference for GPSAltitudeRef, GPSLatitudeRef, and GPSLongitudeRef.
// Returns bool, true is reprsentative of a negative value (-1 Altitude, S Latitude, or W Longitude)
func (ir *ifdReader) parseGPSRef(t Tag) bool {
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
	if ir.logLevelWarn() {
		t.logTag(ir.logWarn()).Msg("error reading GPS Reference")
	}
	return false
}

// TagParser interface is used for Custom Tag Parsers.
type TagParser interface {
	ParseCameraMake(t Tag) (make.CameraMake, string)
	ParseCameraModel(t Tag) (model.CameraModel, string)
	ParseDate(t Tag) time.Time
	ParseRationalU(t Tag) [2]uint32
	ParseString(t Tag) string
	ParseSubSecTime(t Tag) uint16
	ParseUint32(t Tag) uint32
	ParseUint16(t Tag) uint16
}

// TagParserFn function is used for Custom Tag Parsers.
type TagParserFn func(p TagParser, t Tag) error
