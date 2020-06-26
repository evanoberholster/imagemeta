// Package jpegmeta provides types and functions for decoding JPEG file Metadata
package jpegmeta

import (
	"bufio"
	"errors"
	"io"
)

// Errors
var (
	ErrNoExif       = errors.New("error No Exif")
	ErrNoJPEGMarker = errors.New("no JPEG Marker")
	ErrEndOfImage   = errors.New("end of Image")
)

// Scan -
// TODO: Write tests
func Scan(reader io.Reader) (m Metadata, err error) {
	r := bufio.NewReader(reader)
	return scan(r)
}

// ScanBuf -
// TODO: Write tests
func ScanBuf(r *bufio.Reader) (m Metadata, err error) {
	return scan(r)
}

func scan(r *bufio.Reader) (m Metadata, err error) {
	defer func() {
		if state := recover(); state != nil {
			err = state.(error)
		}
	}()
	m = newMetadata(r)
	var buf []byte
	for {
		if buf, err = m.br.Peek(16); err != nil {
			if err == io.EOF {
				err = ErrNoJPEGMarker
				return
			}
			panic(err)
		}

		if !isMarkerFirstByte(buf) {
			if err = m.discard(1); err != nil {
				panic(err)
			}
			continue
		}
		if isSOIMarker(buf) {
			m.pos++
			//fmt.Println("SOI:", m.discarded, m.pos)
			if err = m.discard(2); err != nil {
				panic(err)
			}
			continue
		}
		if m.pos > 0 {
			if err := m.scanMarkers(buf); err == nil {
				continue
			}
		}

		break
	}
	if !m.tiffHeader.IsValid() {
		err = ErrNoExif
		return
	}
	return
}

func (m *Metadata) scanMarkers(buf []byte) (err error) {
	switch buf[1] {
	case markerSOF0, markerSOF1,
		markerSOF2, markerSOF3,
		markerSOF5, markerSOF6,
		markerSOF7, markerSOF9,
		markerSOF10:
		return m.readSOF(buf)
	case markerDHT:
		// Artifical End Of Image for DHT Marker.
		// This is done to improve performance.
		if m.pos == 1 {
			return ErrEndOfImage
		}
		// Ignore DHT Markers
		return m.ignoreMarker(buf)
	case markerSOI:
		m.pos++
		return m.discard(2)
	case markerEOI:
		m.pos--
		// Return EndOfImage
		if m.pos == 1 {
			return ErrEndOfImage
		}
		return m.discard(2)
	case markerDQT:
		// Ignore DQT Markers
		return m.ignoreMarker(buf)
	case markerDRI:
		return m.discard(6)
	case markerAPP0:
		return m.ignoreMarker(buf)
	case markerAPP2:
		if isICCProfilePrefix(buf) {
			// Ignore ICC Profile Marker
			return m.ignoreMarker(buf)
		}
		return m.ignoreMarker(buf)
	case markerAPP7, markerAPP8,
		markerAPP9, markerAPP10:
		return m.ignoreMarker(buf)
	case markerAPP13:
		if isPhotoshopPrefix(buf) {
			// Ignore Photoshop Profile Marker
			return m.ignoreMarker(buf)
		}
		return m.ignoreMarker(buf)
	case markerAPP14:
		return m.ignoreMarker(buf)
	case markerAPP1:
		return m.readAPP1(buf)
	}
	//fmt.Println(m.discarded)
	return m.discard(1)
}
