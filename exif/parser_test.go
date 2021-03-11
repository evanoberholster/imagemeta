package exif

import (
	"bytes"
	"encoding/binary"
	"io"
	"testing"
	"time"

	"github.com/evanoberholster/imagemeta/exif/ifds"
	"github.com/evanoberholster/imagemeta/exif/ifds/gpsifd"
	"github.com/evanoberholster/imagemeta/exif/tag"
	"github.com/evanoberholster/imagemeta/imagetype"
	"github.com/stretchr/testify/assert"
)

func newMockReader(buf []byte) *reader {
	r := newExifReader(bytes.NewReader(buf), binary.BigEndian, 0x0000, 0)
	r.ifdExifOffset[ifds.RootIFD] = 0
	return r
}

func TestParseTimeStamp(t *testing.T) {
	dateTag := tag.NewTag(ifds.DateTimeDigitized, tag.TypeASCII, 20, 0, 0)
	wrongTag := tag.NewTag(ifds.DateTimeDigitized, tag.TypeByte, 20, 0, 0)
	buf := []byte("1997:09:01 12:00:00  ")
	d := newData(newMockReader(buf), imagetype.ImageUnknown)

	ts, err := d.ParseTimeStamp(dateTag, tag.Tag{}, nil)
	if err != nil {
		assert.Error(t, err, "Parse Timestamp")
	}
	expected := time.Unix(873115200, 0).UTC()
	assert.Equal(t, expected.Unix(), ts.Unix(), "Parse Timestamp")

	//
	buf = []byte("1997:09:01")
	d = newData(newMockReader(buf), imagetype.ImageUnknown)
	ts, err = d.ParseTimeStamp(wrongTag, tag.Tag{}, nil)
	if assert.Error(t, err) {
		assert.Equal(t, ErrParseTimeStamp, err)
	}

	ts, err = d.ParseTimeStamp(dateTag, tag.Tag{}, nil)
	if assert.Error(t, err) {
		assert.ErrorIs(t, err, io.EOF)
	}

}
func TestParseGPSTimeStamp(t *testing.T) {
	parseGPSTimeStampTests := []struct {
		ds  []byte
		ts  []byte
		val time.Time
		err error
	}{
		{[]byte("1992:03:01. "), []byte{255, 128, 255, 255, 10, 0, 0, 0, 120, 240, 0, 0, 0, 255, 128, 255, 255, 10, 0, 0, 0, 120, 240, 0, 0, 0}, time.Unix(699553984, 0), nil},
		{[]byte("1992:03."), []byte{255}, time.Unix(0, 0), ErrParseBufSize},
		{[]byte("1992:03."), []byte{255, 128, 255, 255, 10, 0, 0, 0, 120, 240, 0, 0}, time.Unix(0, 0), ErrParseTimeStamp},
		{[]byte("1992:03:"), []byte{255, 128, 255, 255, 10, 0, 0, 0, 120, 240, 0, 0}, time.Unix(0, 0), ErrParseBufSize},
		{[]byte("1992:03:01."), []byte("255, "), time.Unix(0, 0), ErrParseBufSize},
		{[]byte("1992:03:01."), []byte("255, "), time.Unix(0, 0), ErrParseTimeStamp},
	}

	for i, v := range parseGPSTimeStampTests {
		buf := append(v.ds, v.ts...)
		ds := tag.NewTag(gpsifd.GPSDateStamp, tag.TypeASCII, 11, 0, 0)
		ts := tag.NewTag(gpsifd.GPSTimeStamp, tag.TypeRational, 0x0003, 11, 0)

		d := newData(newMockReader(buf), imagetype.ImageUnknown)
		if i == 5 {
			ts = tag.NewTag(gpsifd.GPSTimeStamp, tag.TypeByte, 0x0003, 11, 0)
		}
		ti, err := d.ParseGPSTimeStamp(ds, ts, tag.Tag{}, nil)
		if err != nil {
			assert.ErrorIs(t, err, v.err, "Test: %d", i)
		}
		if v.err == nil {
			assert.Equal(t, v.val.Unix(), ti.Unix(), "Test: %d", i)
		}

	}

}
func TestParseGPSCoord(t *testing.T) {
	parseGPSCoordTests := []struct {
		buf []byte
		ref byte
		val float64
		err error
	}{
		{[]byte{255, 128, 255, 255, 10, 0, 0, 0, 120, 240, 0, 0, 1, 0, 0, 0, 250, 0, 0, 0, 1, 0, 0, 0}, 'S', -27.63546006348398, nil},
		{[]byte{255, 128, 255, 255, 10, 0, 0, 0, 120, 240, 0, 0, 1, 0, 0, 0, 250, 0, 0, 0, 1, 0, 0, 0}, 'N', 27.63546006348398, nil},
		{[]byte{255, 128, 255, 255, 10, 0, 0, 0, 120, 240, 0, 0, 1, 0, 0, 0, 250, 0, 0, 0, 1, 0, 0, 0}, 'W', -27.63546006348398, nil},
		{[]byte{255, 128, 255, 255, 10, 0, 0, 0, 120, 240, 0, 0, 1, 0, 0, 0, 250, 0, 0, 0, 1, 0, 0, 0}, 'E', 27.63546006348398, nil},
		{[]byte{255, 128, 255, 255, 10, 0, 0, 0, 120, 240, 0, 0, 1, 0, 0, 0, 250, 0, 0, 0, 1}, 'E', 0, io.EOF},
		{[]byte("255, "), 'S', 0, ErrParseGPS},
	}

	for i, v := range parseGPSCoordTests {

		lat := tag.NewTag(gpsifd.GPSLatitude, tag.TypeRational, 3, 0, 0)
		latRef := tag.NewTag(gpsifd.GPSLatitudeRef, tag.TypeASCII, 2, binary.BigEndian.Uint32([]byte{v.ref, 0, 0, 0}), 0)
		d := newData(newMockReader(v.buf), imagetype.ImageUnknown)
		if v.err == ErrParseGPS {
			lat = tag.NewTag(gpsifd.GPSLatitude, tag.TypeByte, 3, 0, 0)
		}
		coord, err := d.ParseGPSCoord(latRef, lat)
		if v.err != nil {
			if assert.Error(t, err) {
				assert.ErrorIs(t, err, v.err, "Test: %d", i)
			}
		}
		assert.Equal(t, v.val, coord, "Test: %d", i)
	}
}

func TestParseASCIIValue(t *testing.T) {
	tests := []struct {
		name string
		data []byte
		tag  tag.Tag
		val  string
		err  error
	}{
		{"1", []byte("  HelloWorld"), tag.NewTag(ifds.ActiveArea, tag.TypeASCII, 10, 2, 0), "HelloWorld", nil},
		{"2", []byte(""), tag.NewTag(ifds.ActiveArea, tag.TypeASCII, 2, 538986601, 0), "  Hi", nil},
		{"3", []byte{}, tag.NewTag(ifds.ActiveArea, tag.TypeRational, 1, 0, 0), "", tag.ErrTagTypeNotValid},
		{"4", []byte{}, tag.NewTag(ifds.ActiveArea, tag.TypeShort, 1, 12345773, 0), "", tag.ErrTagTypeNotValid},
		{"5", []byte{}, tag.NewTag(ifds.ActiveArea, tag.TypeLong, 1, 0, 0), "", tag.ErrTagTypeNotValid},
		{"2", []byte(""), tag.NewTag(ifds.ActiveArea, tag.TypeASCII, 6, 538986601, 0), "", io.EOF},
	}
	for _, v := range tests {
		d := newData(newMockReader(v.data), imagetype.ImageUnknown)
		val, err := d.ParseASCIIValue(v.tag)
		assert.ErrorIs(t, err, v.err, v.name)
		assert.Equal(t, v.val, val, v.name)
	}
}
func TestParseUint32Value(t *testing.T) {
	tests := []struct {
		name string
		data []byte
		tag  tag.Tag
		val  uint32
		err  error
	}{
		{"1", []byte{}, tag.NewTag(ifds.ActiveArea, tag.TypeASCII, 10, 2, 0), 0, tag.ErrTagTypeNotValid},
		{"2", []byte{}, tag.NewTag(ifds.ActiveArea, tag.TypeRational, 1, 0, 0), 0, io.EOF},
		{"3", []byte{}, tag.NewTag(ifds.ActiveArea, tag.TypeLong, 1, 12345773, 0), 12345773, nil},
		{"4", []byte{}, tag.NewTag(ifds.ActiveArea, tag.TypeShort, 1, 12345773, 0), 188, nil},
		{"5", []byte{}, tag.NewTag(ifds.ActiveArea, tag.TypeLong, 1, 0, 0), 0, nil},
	}
	for _, v := range tests {
		d := newData(newMockReader(v.data), imagetype.ImageUnknown)
		val, err := d.ParseUint32Value(v.tag)
		assert.ErrorIs(t, err, v.err, v.name)
		assert.Equal(t, int(v.val), int(val), v.name)
	}
}
func TestParseUint16Values(t *testing.T) {
	tests := []struct {
		name string
		data []byte
		tag  tag.Tag
		val  []uint16
		err  error
	}{
		{"1", []byte{}, tag.NewTag(ifds.ActiveArea, tag.TypeASCII, 10, 2, 0), nil, tag.ErrTagTypeNotValid},
		{"2", []byte{}, tag.NewTag(ifds.ActiveArea, tag.TypeRational, 1, 0, 0), nil, tag.ErrTagTypeNotValid},
		{"3", []byte{}, tag.NewTag(ifds.ActiveArea, tag.TypeLong, 1, 12345773, 0), nil, tag.ErrTagTypeNotValid},
		{"4", []byte{}, tag.NewTag(ifds.ActiveArea, tag.TypeShort, 1, 1024232342, 0), []uint16{15628}, nil},
		{"5", []byte{}, tag.NewTag(ifds.ActiveArea, tag.TypeShort, 2, 1024232342, 0), []uint16{15628, 35734}, nil},
		{"6", []byte{}, tag.NewTag(ifds.ActiveArea, tag.TypeShort, 3, 1024232342, 0), nil, io.EOF},
		{"7", []byte{0, 0, 250, 250, 125, 125, 125, 234}, tag.NewTag(ifds.ActiveArea, tag.TypeShort, 3, 2, uint8(ifds.RootIFD)), []uint16{64250, 32125, 32234}, nil},
	}
	for _, v := range tests {
		d := newData(newMockReader(v.data), imagetype.ImageUnknown)
		val, err := d.ParseUint16Values(v.tag)
		assert.ErrorIs(t, err, v.err, v.name)
		assert.Equal(t, v.val, val, v.name)
	}
}
func TestParseUint32Values(t *testing.T) {
	tests := []struct {
		name string
		data []byte
		tag  tag.Tag
		val  []uint32
		err  error
	}{
		{"1", []byte{}, tag.NewTag(ifds.ActiveArea, tag.TypeASCII, 10, 2, 0), nil, tag.ErrTagTypeNotValid},
		{"2", []byte{}, tag.NewTag(ifds.ActiveArea, tag.TypeRational, 1, 0, 0), nil, tag.ErrTagTypeNotValid},
		{"3", []byte{}, tag.NewTag(ifds.ActiveArea, tag.TypeShort, 1, 1024232342, 0), nil, tag.ErrTagTypeNotValid},
		{"4", []byte{}, tag.NewTag(ifds.ActiveArea, tag.TypeLong, 1, 12345773, 0), []uint32{12345773}, nil},
		{"5", []byte{0, 0, 250, 250, 125, 125, 0, 0, 25, 25}, tag.NewTag(ifds.ActiveArea, tag.TypeLong, 2, 2, 0), []uint32{4210720125, 6425}, nil},
		{"6", []byte{0, 0, 250, 250, 125, 125, 0, 0, 25, 25, 0, 0, 0, 25}, tag.NewTag(ifds.ActiveArea, tag.TypeLong, 3, 2, 0), []uint32{4210720125, 6425, 25}, nil},
		{"7", []byte{0, 0, 250, 250, 125, 125, 125, 234}, tag.NewTag(ifds.ActiveArea, tag.TypeLong, 3, 2, 0), nil, io.EOF},
	}
	for _, v := range tests {
		d := newData(newMockReader(v.data), imagetype.ImageUnknown)
		val, err := d.ParseUint32Values(v.tag)
		assert.ErrorIs(t, err, v.err, v.name)
		assert.Equal(t, v.val, val, v.name)
	}
}
func TestParseRationalValue(t *testing.T) {
	tests := []struct {
		name string
		data []byte
		tag  tag.Tag
		val  []uint32
		err  error
	}{
		{"1", []byte{}, tag.NewTag(ifds.ActiveArea, tag.TypeASCII, 10, 2, 0), []uint32{0, 0}, tag.ErrTagTypeNotValid},
		{"2", []byte{}, tag.NewTag(ifds.ActiveArea, tag.TypeLong, 1, 1024232342, 0), []uint32{0, 0}, tag.ErrTagTypeNotValid},
		{"3", []byte{0, 0, 0, 0, 12, 23, 0, 0, 12, 24}, tag.NewTag(ifds.ApertureValue, tag.TypeSignedRational, 1, 2, 0), []uint32{3095, 3096}, nil},
		{"4", []byte{0, 0, 0, 0, 12, 25, 0, 0, 12, 26}, tag.NewTag(ifds.ApertureValue, tag.TypeRational, 1, 2, 0), []uint32{3097, 3098}, nil},
		{"5", []byte{0, 0, 0, 0, 12, 23, 0, 0, 12, 24}, tag.NewTag(ifds.ApertureValue, tag.TypeRational, 2, 2, 0), []uint32{0, 0}, ErrParseRationals},
		{"6", []byte{}, tag.NewTag(ifds.ActiveArea, tag.TypeRational, 1, 10, 0), []uint32{0, 0}, io.EOF},
	}

	// Rational
	for _, v := range tests {
		d := newData(newMockReader(v.data), imagetype.ImageUnknown)
		a, b, err := d.ParseRationalValue(v.tag)
		assert.ErrorIs(t, err, v.err, v.name)
		assert.Equal(t, v.val, []uint32{a, b}, v.name)
	}

	// SignedRational
	for _, v := range tests {
		d := newData(newMockReader(v.data), imagetype.ImageUnknown)
		a, b, err := d.ParseSRationalValue(v.tag)
		assert.ErrorIs(t, err, v.err, v.name)
		val2 := []int32{int32(v.val[0]), int32(v.val[1])}
		assert.Equal(t, val2, []int32{a, b}, v.name)
	}
}

func TestParseRationalValues(t *testing.T) {
	tests := []struct {
		name string
		data []byte
		tag  tag.Tag
		val  []tag.Rational
		err  error
	}{
		{"1", []byte{}, tag.NewTag(ifds.ActiveArea, tag.TypeASCII, 10, 2, 0), nil, tag.ErrTagTypeNotValid},
		{"2", []byte{}, tag.NewTag(ifds.ActiveArea, tag.TypeLong, 1, 1024232342, 0), nil, tag.ErrTagTypeNotValid},
		{"3", []byte{0, 0, 0, 0, 12, 23, 0, 0, 12, 24}, tag.NewTag(ifds.ApertureValue, tag.TypeSignedRational, 1, 2, 0), []tag.Rational{{Numerator: 3095, Denominator: 3096}}, nil},
		{"4", []byte{0, 0, 0, 0, 12, 25, 0, 0, 12, 26}, tag.NewTag(ifds.ApertureValue, tag.TypeRational, 1, 2, 0), []tag.Rational{{Numerator: 3097, Denominator: 3098}}, nil},
		{"5", []byte{0, 0, 0, 0, 12, 23, 0, 0, 12, 24, 0, 0, 12, 23, 0, 0, 12, 24}, tag.NewTag(ifds.ApertureValue, tag.TypeRational, 2, 2, 0), []tag.Rational{{Numerator: 3095, Denominator: 3096}, {Numerator: 3095, Denominator: 3096}}, nil},
		{"6", []byte{}, tag.NewTag(ifds.ActiveArea, tag.TypeRational, 1, 10, 0), nil, io.EOF},
	}

	// Rational
	for _, v := range tests {
		d := newData(newMockReader(v.data), imagetype.ImageUnknown)
		val, err := d.ParseRationalValues(v.tag)
		assert.ErrorIs(t, err, v.err, v.name)
		assert.Equal(t, v.val, val, v.name)
	}
}
func TestParseSRationalValues(t *testing.T) {
	tests := []struct {
		name string
		data []byte
		tag  tag.Tag
		val  []tag.SRational
		err  error
	}{
		{"1", []byte{}, tag.NewTag(ifds.ActiveArea, tag.TypeASCII, 10, 2, 0), nil, tag.ErrTagTypeNotValid},
		{"2", []byte{}, tag.NewTag(ifds.ActiveArea, tag.TypeLong, 1, 1024232342, 0), nil, tag.ErrTagTypeNotValid},
		{"3", []byte{0, 0, 0, 0, 12, 23, 0, 0, 12, 24}, tag.NewTag(ifds.ApertureValue, tag.TypeSignedRational, 1, 2, 0), []tag.SRational{{Numerator: 3095, Denominator: 3096}}, nil},
		{"4", []byte{0, 0, 0, 0, 12, 25, 0, 0, 12, 26}, tag.NewTag(ifds.ApertureValue, tag.TypeSignedRational, 1, 2, 0), []tag.SRational{{Numerator: 3097, Denominator: 3098}}, nil},
		{"5", []byte{0, 0, 0, 0, 12, 23, 0, 0, 12, 24, 0, 0, 12, 23, 0, 0, 12, 24}, tag.NewTag(ifds.ApertureValue, tag.TypeSignedRational, 2, 2, 0), []tag.SRational{{Numerator: 3095, Denominator: 3096}, {Numerator: 3095, Denominator: 3096}}, nil},
		{"6", []byte{}, tag.NewTag(ifds.ActiveArea, tag.TypeRational, 1, 10, 0), nil, io.EOF},
	}

	// Rational
	for _, v := range tests {
		d := newData(newMockReader(v.data), imagetype.ImageUnknown)
		val, err := d.ParseSRationalValues(v.tag)
		assert.ErrorIs(t, err, v.err, v.name)
		assert.Equal(t, v.val, val, v.name)
	}
}

func TestTrim(t *testing.T) {
	// Test Trim
	a := []byte{'a', 'b', 'c', 'd', '.', ' '}
	if !bytes.Equal(trim(a), a[:len(a)-1]) {
		t.Errorf("Trim should remove trailing spaces: expected %s got %s", a[:len(a)-1], trim(a))
	}
	a = []byte{' ', ' ', ' ', ' ', ' ', ' '}
	if len(trim(a)) != 0 {
		t.Errorf("Trim should remove trailing spaces: expected %d got %d", 0, len(trim(a)))
	}
}
