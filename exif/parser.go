package exif

import (
	"time"

	"github.com/evanoberholster/imagemeta/exif/tag"
	"github.com/pkg/errors"
)

// Parsing Errors
var (
	ErrParseBufSize   = errors.New("error parse has insufficient data")
	ErrParseGPS       = errors.New("error parsing GPS coords")
	ErrParseTimeStamp = errors.New("error parsing timestamp")
	ErrParseSubSecond = errors.New("error parsing sub second")
	ErrParseRationals = errors.New("error parsing rationals")
)

////
// Highlevel Parsers
////

// ParseTimeStamp parses a time.Time from 2 ASCII Tag's.
// ex: 1997:09:01 12:00:00
// Based on: http://www.cipa.jp/std/documents/e/DC-008-2012_E.pdf (Last checked: 24/02/2021)
func (e *Data) ParseTimeStamp(date tag.Tag, subSec tag.Tag, tz *time.Location) (t time.Time, err error) {
	if date.Type() == tag.TypeASCII {
		if tz == nil {
			tz = time.UTC
		}
		var buf []byte
		buf, err = e.er.TagValue(date)
		if err != nil {
			err = errors.Wrap(err, "ParseTimeStamp")
			return
		}
		// check recieved value
		if buf[4] == ':' && buf[7] == ':' && buf[10] == ' ' &&
			buf[13] == ':' && buf[16] == ':' {

			year := parseUint(buf[0:4])
			month := parseUint(buf[5:7])
			day := parseUint(buf[8:10])
			hour := parseUint(buf[11:13])
			min := parseUint(buf[14:16])
			sec := parseUint(buf[17:19])

			sub, _ := e.ParseSubSec(subSec)

			return time.Date(int(year), time.Month(month), int(day), int(hour), int(min), int(sec), int(sub), tz), nil
		}
	}

	return time.Time{}, ErrParseTimeStamp
}

// GPS Info Parser
// Lat, Lng, Alt, Time

// ParseGPSTimeStamp parses the GPSDateStamp, GPSTimeStamp Tags in the given Timezone in UTC.
// Optionally add subSec tag from Exif.
func (e *Data) ParseGPSTimeStamp(ds tag.Tag, ts tag.Tag, subSec tag.Tag, tz *time.Location) (t time.Time, err error) {
	byteOrder := e.er.byteOrder
	if !(ts.UnitCount == 3 && ts.Type() == tag.TypeRational && ds.Type() == tag.TypeASCII) {
		err = errors.Wrap(ErrParseTimeStamp, "ParseGPSTimeStamp TagType")
		return
	}

	// Read GPS DateStamp Tag with the format is "YYYY:MM:DD."
	var buf []byte
	buf, err = e.er.TagValue(ds)
	if err != nil {
		err = errors.Wrap(ErrParseBufSize, "ParseGPSTimeStamp DateStamp")
		return
	}

	// Parse yyyy:mm:dd from recieved value
	if buf[4] == ':' && buf[7] == ':' { //&& buf[10] == '.' {
		year := parseUint(buf[0:4])
		month := parseUint(buf[5:7])
		day := parseUint(buf[8:10])

		// Read GPS TimeStamp Tag
		buf, err = e.er.TagValue(ts)
		if err != nil {
			err = errors.Wrap(ErrParseBufSize, "ParseGPSTimeStamp TimeStamp")
			return
		}
		hour := int(byteOrder.Uint32(buf[:4]) / byteOrder.Uint32(buf[4:8]))
		min := int(byteOrder.Uint32(buf[8:12]) / byteOrder.Uint32(buf[12:16]))
		sec := int(byteOrder.Uint32(buf[16:20]) / byteOrder.Uint32(buf[20:24]))

		sub, _ := e.ParseSubSec(subSec)

		if tz == nil {
			tz = time.UTC
		}

		return time.Date(int(year), time.Month(month), int(day), int(hour), int(min), int(sec), int(sub), tz), nil
	}
	return time.Time{}, ErrParseTimeStamp
}

// ParseGPSCoord parses the GPS Coordinate (Lat or Lng) with the corresponding reference Tag.
func (e *Data) ParseGPSCoord(refTag tag.Tag, coordTag tag.Tag) (coord float64, err error) {
	if !(refTag.IsEmbedded() && coordTag.UnitCount == 3 && coordTag.Type() == tag.TypeRational) {
		return 0.0, ErrParseGPS
	}

	byteOrder := e.er.byteOrder
	// Read GPS Coord Tag
	var buf []byte
	buf, err = e.er.TagValue(coordTag)
	if err != nil {
		err = errors.Wrap(err, "ParseGPSCoord")
		return
	}
	coord = (float64(byteOrder.Uint32(buf[:4])) / float64(byteOrder.Uint32(buf[4:8])))
	coord += (float64(byteOrder.Uint32(buf[8:12])) / float64(byteOrder.Uint32(buf[12:16])) / 60.0)
	coord += (float64(byteOrder.Uint32(buf[16:20])) / float64(byteOrder.Uint32(buf[20:24])) / 3600.0)

	// Read Reference Tag
	// Coordinate is a negative value for a South or West Orientation
	buf = e.er.embeddedTagValue(refTag.ValueOffset)
	if buf[0] == 'S' || buf[0] == 'W' {
		coord *= -1
	}
	return coord, nil
}

// ParseSubSec parses a Subsecond Tag and returns an int in Nanoseconds.
// Returns ErrParseSubSecond if an err occurs
func (e *Data) ParseSubSec(subSec tag.Tag) (int, error) {
	if subSec.Type() == tag.TypeASCII && subSec.IsEmbedded() {
		buf := e.er.embeddedTagValue(subSec.ValueOffset)
		return int(parseUint(buf) * 1000000), nil
	}
	return 0, ErrParseSubSecond
}

////
// Low-Level Parsers
////

// ParseASCIIValue parses the ASCII value of the tag as a string
// and returns an error if it encounters one
func (e *Data) ParseASCIIValue(t tag.Tag) (value string, err error) {
	if t.Type() == tag.TypeASCII || t.Type() == tag.TypeASCIINoNul {
		var buf []byte
		if buf, err = e.er.TagValue(t); err != nil {
			err = errors.Wrap(err, "ParseASCIIValue")
			return
		}

		// Trim trailing spaces and null values
		return string(trim(buf)), nil
	}
	return "", tag.ErrTagTypeNotValid

}

// ParseUint16Value returns the Short value of the tag as a uint16
// and returns an error if it encounters one.
//
// Warning: it returns only the first value if there are more values
// use Uint16Values function
func (e *Data) ParseUint16Value(t tag.Tag) (uint16, error) {
	v, err := e.ParseUint32Value(t)
	return uint16(v), err
}

// ParseUint32Value returns the Short or Long value of the tag as a uint32
// and returns an error if it encounters one.
//
// Warning: it returns only the first value if there are more values
// use Uint16Values (Short) or Unit32Values (Long) function
func (e *Data) ParseUint32Value(t tag.Tag) (value uint32, err error) {
	if t.Type().IsValid() && t.UnitCount == 1 {
		var buf []byte
		if buf, err = e.er.TagValue(t); err != nil {
			return
		}
		byteOrder := e.er.ByteOrder()

		if t.Type() == tag.TypeShort {
			value = uint32(byteOrder.Uint16(buf[:2]))
			return
		}
		if t.Type() == tag.TypeLong {
			value = byteOrder.Uint32(buf[:4])
			return
		}

	}
	return 0, tag.ErrTagTypeNotValid
}

// ParseUint16Values parses the Short value of the tag as a uint16 array
// and returns an error if it encounters one.
func (e *Data) ParseUint16Values(t tag.Tag) (value []uint16, err error) {
	if t.Type() == tag.TypeShort {
		var buf []byte
		buf, err = e.er.TagValue(t)
		if err != nil {
			return
		}

		byteOrder := e.er.ByteOrder()
		count := int(t.UnitCount)

		value = make([]uint16, count)
		for i := 0; i < count; i++ {
			value[i] = byteOrder.Uint16(buf[i*2:])
		}

		return
	}
	return nil, tag.ErrTagTypeNotValid
}

// ParseUint32Values parses the Long value of the tag as a uint32 array
// and returns an error if it encounters one.
//
func (e *Data) ParseUint32Values(t tag.Tag) (value []uint32, err error) {
	if t.Type() == tag.TypeLong {
		var buf []byte
		if buf, err = e.er.TagValue(t); err != nil {
			return nil, err
		}

		byteOrder := e.er.ByteOrder()
		count := int(t.UnitCount)

		value = make([]uint32, count)
		for i := 0; i < count; i++ {
			value[i] = byteOrder.Uint32(buf[i*4:])
		}

		return
	}
	return nil, tag.ErrTagTypeNotValid
}

// ParseRationalValue parses the Rational value and returns a
// numerator and denominator for a single Unsigned Rational
func (e *Data) ParseRationalValue(t tag.Tag) (n, d uint32, err error) {
	if t.Type() == tag.TypeRational || t.Type() == tag.TypeSignedRational {
		if t.UnitCount > 1 {
			return 0, 0, ErrParseRationals
		}
		var buf []byte
		buf, err = e.er.TagValue(t)
		if err != nil {
			return
		}
		byteOrder := e.er.ByteOrder()
		n = byteOrder.Uint32(buf[:4])
		d = byteOrder.Uint32(buf[4:8])
		return
	}
	return 0, 0, tag.ErrTagTypeNotValid
}

// ParseSRationalValue returns a numerator and denominator for a single Signed Rational
func (e *Data) ParseSRationalValue(t tag.Tag) (num, denom int32, err error) {
	n, d, err := e.ParseRationalValue(t)
	return int32(n), int32(d), err
}

// ParseRationalValues returns a list of unsignedRationals
func (e *Data) ParseRationalValues(t tag.Tag) (value []tag.Rational, err error) {
	if t.Type() == tag.TypeRational || t.Type() == tag.TypeSignedRational {
		var buf []byte
		if buf, err = e.er.TagValue(t); err != nil {
			return nil, err
		}
		byteOrder := e.er.ByteOrder()
		count := int(t.UnitCount)

		value = make([]tag.Rational, count)
		for i := 0; i < count; i++ {
			value[i].Numerator = byteOrder.Uint32(buf[i*8:])
			value[i].Denominator = byteOrder.Uint32(buf[i*8+4:])
		}

		return
	}
	return nil, tag.ErrTagTypeNotValid
}

// ParseSRationalValues returns a list of unsignedRationals
func (e *Data) ParseSRationalValues(t tag.Tag) (value []tag.SRational, err error) {
	if t.Type() == tag.TypeRational || t.Type() == tag.TypeSignedRational {
		var buf []byte
		if buf, err = e.er.TagValue(t); err != nil {
			return nil, err
		}

		byteOrder := e.er.ByteOrder()
		count := int(t.UnitCount)

		value = make([]tag.SRational, count)
		for i := 0; i < count; i++ {
			value[i].Numerator = int32(byteOrder.Uint32(buf[i*8:]))
			value[i].Denominator = int32(byteOrder.Uint32(buf[i*8+4:]))
		}

		return
	}
	return nil, tag.ErrTagTypeNotValid
}

////
// Helper functions
////

// trim removes null value padding "0x00 and 0x20" from []byte
func trim(buf []byte) []byte {
	for i := len(buf); i > 0; i-- {
		if buf[i-1] == 0x0 ||
			buf[i-1] == 0x20 {
			continue
		}
		return buf[:i]
	}
	return nil
}

// parseUint parses a []byte of a string representation of a uint64 value and returns the value.
func parseUint(buf []byte) (u uint64) {
	for i := 0; i < len(buf); i++ {
		u *= 10
		u += uint64(buf[i] - '0')
	}
	return
}
