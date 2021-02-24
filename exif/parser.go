package exif

import (
	"errors"
	"time"

	"github.com/evanoberholster/imagemeta/exif/tag"
)

// Parsing Errors
var (
	ErrParseTimeStamp = errors.New("error parsing timestamp")
	ErrParseSubSecond = errors.New("error parsing sub second")
)

////
// Highlevel Parsers
////

// ParseTimeStamp parses a time.Time from 2 ASCII Tag's.
// ex: 1997:09:01 12:00:00
// Based on: http://www.cipa.jp/std/documents/e/DC-008-2012_E.pdf (Last checked: 24/02/2021)
func (e *Data) ParseTimeStamp(date tag.Tag, subSec tag.Tag) (t time.Time, err error) {
	if date.TagType == tag.TypeASCII {
		buf := e.er.rawBuffer[:]
		if _, err = e.er.ReadAt(buf[:date.Size()], int64(date.ValueOffset)); err != nil {
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

			return time.Date(int(year), time.Month(month), int(day), int(hour), int(min), int(sec), int(sub), time.UTC), nil
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
	if ts.UnitCount != 3 || ts.TagType != tag.TypeRational || ds.TagType != tag.TypeASCII {
		// TODO: Return error
	}

	sub, _ := e.ParseSubSec(subSec)
	buf := e.er.rawBuffer[:]

	// Read GPS DateStamp Tag
	// The format is "YYYY:MM:DD."
	if _, err = e.er.ReadAt(buf[:ds.Size()], int64(ds.ValueOffset)); err != nil {
		return
	}

	// check recieved value
	if buf[4] == ':' && buf[7] == ':' { //&& buf[10] == '.' {
		year := parseUint(buf[0:4])
		month := parseUint(buf[5:7])
		day := parseUint(buf[8:10])

		// Read GPS TimeStamp Tag
		if _, err = e.er.ReadAt(buf[:ts.Size()], int64(ts.ValueOffset)); err != nil {
			return
		}
		hour := int(byteOrder.Uint32(buf[:4]) / byteOrder.Uint32(buf[4:8]))
		min := int(byteOrder.Uint32(buf[8:12]) / byteOrder.Uint32(buf[12:16]))
		sec := int(byteOrder.Uint32(buf[16:20]) / byteOrder.Uint32(buf[20:24]))

		return time.Date(int(year), time.Month(month), int(day), int(hour), int(min), int(sec), int(sub), tz), nil
	}
	return time.Time{}, ErrParseTimeStamp
}

// ParseGPSCoord parses the GPS Coordinate (Lat or Lng) with the corresponding reference Tag.
func (e *Data) ParseGPSCoord(refTag tag.Tag, coordTag tag.Tag) (coord float64, err error) {
	buf := e.er.rawBuffer[:coordTag.Size()]
	byteOrder := e.er.byteOrder

	if !refTag.IsEmbedded() || coordTag.UnitCount != 3 || coordTag.TagType != tag.TypeRational {
		// TODO: Return error
	}

	// Read GPS Coord Tag
	if _, err = e.er.ReadAt(buf, int64(coordTag.ValueOffset)); err != nil {
		return
	}
	coord = (float64(byteOrder.Uint32(buf[:4])) / float64(byteOrder.Uint32(buf[4:8])))
	coord += (float64(byteOrder.Uint32(buf[8:12])) / float64(byteOrder.Uint32(buf[12:16])) / 60.0)
	coord += (float64(byteOrder.Uint32(buf[16:20])) / float64(byteOrder.Uint32(buf[20:24])) / 3600.0)

	// Read Reference Tag
	// Coordinate is a negative value for a South or West Orientation
	e.er.byteOrder.PutUint32(buf[:4], refTag.ValueOffset)
	if buf[0] == 'S' || buf[0] == 'W' {
		coord *= -1
	}
	return coord, nil
}

// ParseSubSec parses a Subsecond Tag and returns an int in Nanoseconds.
// Returns ErrParseSubSecond if an err occurs
func (e *Data) ParseSubSec(subSec tag.Tag) (int, error) {
	if subSec.TagType == tag.TypeASCII && subSec.IsEmbedded() {
		e.er.byteOrder.PutUint32(e.er.rawBuffer[:4], subSec.ValueOffset)
		return int(parseUint(e.er.rawBuffer[:3]) * 1000000), nil
	}
	return 0, ErrParseSubSecond
}

////
// Low-Level Parsers
////

// RawEncodedBytes returns the raw encoded bytes for the value that we represent.
func rawEncodedBytes(r *reader, t tag.Tag) (buf []byte, err error) {
	// check if Value is Embedded
	if t.IsEmbedded() {
		r.ByteOrder().PutUint32(r.rawBuffer[:4], t.ValueOffset)
		return r.rawBuffer[:4], nil
	}

	byteLength := int(t.TagType.Size() * uint32(t.UnitCount))
	if byteLength <= len(r.rawBuffer) {
		return r.ReadBufferAt(byteLength, int64(t.ValueOffset))
	}

	buf = make([]byte, byteLength)
	if _, err = r.ReadAt(buf[:byteLength], int64(t.ValueOffset)); err != nil {
		return nil, err
	}
	return buf[:byteLength], nil
}

// ParseASCIIValue parses the ASCII value of the tag as a string
// and returns an error if it encounters one
func (e *Data) ParseASCIIValue(t tag.Tag) (value string, err error) {
	if t.TagType.IsValid() {
		// TODO: Needs Typecheck

		var buf []byte
		size := t.Size()
		if size <= rawBufferSize {
			buf := e.er.rawBuffer[:t.Size()]
			// check if Value is Embedded
			if t.IsEmbedded() {
				e.er.byteOrder.PutUint32(buf[:4], t.ValueOffset)
				return string(trim(buf[:4])), nil
			}
			//return r.ReadBufferAt(byteLength, int64(t.ValueOffset))
			if _, err = e.er.ReadAt(buf[:], int64(t.ValueOffset)); err != nil {
				return
			}
			return string(trim(buf[:])), nil
		}

		buf = make([]byte, size)
		if _, err = e.er.ReadAt(buf[:], int64(t.ValueOffset)); err != nil {
			return
		}

		// Trim trailing spaces and null values
		return string(trim(buf)), nil
	}
	return "", tag.ErrTagNotValid
}

// ParseUint16Value returns the Short value of the tag as a uint16
// and returns an error if it encounters one.
//
// Warning: it returns only the first value if there are more values
// use Uint16Values function
func (e *Data) ParseUint16Value(t tag.Tag) (value uint16, err error) {
	v, err := e.ParseUint32Value(t)
	return uint16(v), err
}

// ParseUint32Value returns the Short or Long value of the tag as a uint32
// and returns an error if it encounters one.
//
// Warning: it returns only the first value if there are more values
// use Uint16Values (Short) or Unit32Values (Long) function
func (e *Data) ParseUint32Value(t tag.Tag) (value uint32, err error) {
	if t.TagType == tag.TypeShort || t.TagType == tag.TypeLong {
		if t.IsEmbedded() {
			return t.ValueOffset, nil
		}
		var buf []byte
		buf, err = rawEncodedBytes(e.er, t)
		if err != nil {
			return
		}
		if t.TagType == tag.TypeShort {
			value = uint32(e.er.byteOrder.Uint16(buf[:2]))
		}
		if t.TagType == tag.TypeLong {
			value = e.er.byteOrder.Uint32(buf[:4])
		}
		return
	}
	return 0, tag.ErrTagTypeNotValid
}

// ParseUint16Values parses the Short value of the tag as a uint16 array
// and returns an error if it encounters one.
func (e *Data) ParseUint16Values(t tag.Tag) (value []uint16, err error) {
	if t.TagType == tag.TypeShort {
		var buf []byte
		if buf, err = rawEncodedBytes(e.er, t); err != nil {
			return nil, err
		}

		if len(buf) < t.Size() {
			err = tag.ErrNotEnoughData
		}

		count := int(t.UnitCount)
		value = make([]uint16, count)
		for i := 0; i < count; i++ {
			value[i] = e.er.byteOrder.Uint16(buf[i*2:])
		}

		return
	}
	return nil, tag.ErrTagTypeNotValid
}

// ParseUint32Values parses the Long value of the tag as a uint32 array
// and returns an error if it encounters one.
//
func (e *Data) ParseUint32Values(t tag.Tag) (value []uint32, err error) {
	if t.TagType == tag.TypeLong {
		if t.IsEmbedded() {
			return append(value, t.ValueOffset), nil
		}
		var buf []byte
		if buf, err = rawEncodedBytes(e.er, t); err != nil {
			return nil, err
		}

		if len(buf) < t.Size() {
			err = tag.ErrNotEnoughData
		}

		count := int(t.UnitCount)
		value = make([]uint32, count)
		for i := 0; i < count; i++ {
			value[i] = e.er.byteOrder.Uint32(buf[i*4:])
		}

		return
	}
	return nil, tag.ErrTagTypeNotValid
}

// ParseRationalValue parses the Rational value and returns a
// numerator and denominator for a single Unsigned Rational
func (e *Data) ParseRationalValue(t tag.Tag) (n, d uint32, err error) {
	if t.TagType == tag.TypeRational || t.TagType == tag.TypeSignedRational {
		if t.Size() > len(e.er.rawBuffer) {
			// Error Multiple Rationals
		}

		buf := e.er.rawBuffer[:t.Size()]
		if _, err = e.er.ReadAt(buf[:], int64(t.ValueOffset)); err != nil {
			return
		}
		n = e.er.byteOrder.Uint32(buf[:4])
		d = e.er.byteOrder.Uint32(buf[4:8])
	}
	return
}

// ParseSRationalValue returns a numerator and denominator for a single Signed Rational
func (e *Data) ParseSRationalValue(t tag.Tag) (num, denom int32, err error) {
	n, d, err := e.ParseRationalValue(t)
	return int32(n), int32(d), err
}

// ParseRationalValues returns a list of unsignedRationals
func (e *Data) ParseRationalValues(t tag.Tag) (value []tag.Rational, err error) {
	if t.TagType == tag.TypeRational || t.TagType == tag.TypeSignedRational {
		var buf []byte
		if buf, err = rawEncodedBytes(e.er, t); err != nil {
			return nil, err
		}

		if len(buf) < t.Size() {
			return nil, tag.ErrNotEnoughData
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
	if t.TagType == tag.TypeRational || t.TagType == tag.TypeSignedRational {
		var buf []byte
		if buf, err = rawEncodedBytes(e.er, t); err != nil {
			return nil, err
		}

		if len(buf) < t.Size() {
			return nil, tag.ErrNotEnoughData
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
