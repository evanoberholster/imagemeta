package exif

import (
	"bytes"
	"encoding/binary"
	"testing"
	"time"

	"github.com/evanoberholster/imagemeta/exif/ifds"
	"github.com/evanoberholster/imagemeta/exif/ifds/gpsifd"
	"github.com/evanoberholster/imagemeta/exif/tag"
	"github.com/evanoberholster/imagemeta/imagetype"
	"github.com/stretchr/testify/assert"
)

func newMockReader(buf []byte) *reader {
	return newExifReader(bytes.NewReader(buf), binary.BigEndian, 0x0000)
}

func TestParseTimeStamp(t *testing.T) {
	dateTag := tag.NewTag(ifds.DateTimeDigitized, tag.TypeASCII, 20, 0)
	wrongTag := tag.NewTag(ifds.DateTimeDigitized, tag.TypeByte, 20, 0)
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
		assert.ErrorIs(t, err, ErrParseBufSize)
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
		ds := tag.NewTag(gpsifd.GPSDateStamp, tag.TypeASCII, 11, 0)
		ts := tag.NewTag(gpsifd.GPSTimeStamp, tag.TypeRational, 0x0003, 11)

		d := newData(newMockReader(buf), imagetype.ImageUnknown)
		if i == 5 {
			ts.TagType = tag.TypeByte
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
		{[]byte{255, 128, 255, 255, 10, 0, 0, 0, 120, 240, 0, 0, 1, 0, 0, 0, 250, 0, 0, 0, 1}, 'E', 0, ErrParseBufSize},
		{[]byte("255, "), 'S', 0, ErrParseGPS},
	}

	for i, v := range parseGPSCoordTests {

		lat := tag.NewTag(gpsifd.GPSLatitude, tag.TypeRational, 3, 0)
		latRef := tag.NewTag(gpsifd.GPSLatitudeRef, tag.TypeASCII, 2, binary.BigEndian.Uint32([]byte{v.ref, 0, 0, 0}))
		d := newData(newMockReader(v.buf), imagetype.ImageUnknown)
		if v.err == ErrParseGPS {
			lat.TagType = tag.TypeByte
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
