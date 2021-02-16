package jpeg

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"os"
	"testing"
)

func TestJPEG(t *testing.T) {
	exifHeaderTests := []struct {
		filename         string
		byteOrder        binary.ByteOrder
		firstIfdOffset   uint32
		tiffHeaderOffset uint32
		width            uint16
		height           uint16
	}{
		{"../testImages/JPEG.jpg", binary.LittleEndian, 13746, 12, 1000, 563},
		{"../testImages/NoExif.jpg", binary.BigEndian, 8, 30, 50, 50},
	}
	for _, jpg := range exifHeaderTests {
		t.Run(jpg.filename, func(t *testing.T) {
			// Open file
			f, err := os.Open(jpg.filename)
			if err != nil {
				t.Fatal(err)
			}
			defer f.Close()
			// Search for Tiff header
			br := bufio.NewReader(f)
			m, err := ScanJPEG(br, nil, nil)
			if err != nil {
				t.Fatal(err)
			}
			if m.Header.ByteOrder != jpg.byteOrder {
				t.Errorf("Incorrect Byte Order wanted %s got %s", jpg.byteOrder, m.Header.ByteOrder)
			}
			if m.Header.FirstIfdOffset != jpg.firstIfdOffset {
				t.Errorf("Incorrect first Ifd Offset wanted %d got %d ", jpg.firstIfdOffset, m.Header.FirstIfdOffset)
			}
			if m.Header.TiffHeaderOffset != jpg.tiffHeaderOffset {
				t.Errorf("Incorrect tiff Header Offset wanted %d got %d ", jpg.tiffHeaderOffset, m.Header.TiffHeaderOffset)
			}
			if !m.Header.IsValid() {
				t.Errorf("Wanted valid tiff Header")
			}
			width, height := m.Size()
			if width != jpg.width || height != jpg.height {
				t.Errorf("Incorrect Jpeg Image size wanted width: %d got width: %d ", jpg.width, width)
				t.Errorf("Incorrect Jpeg Image size wanted height: %d got height: %d ", jpg.height, height)
			}
		})

		data := []byte{0, markerFirstByte, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		r := bytes.NewReader(data)
		m := newMetadata(bufio.NewReader(r), nil, nil)

		// Test discard
		m.discard(0)
		if m.discarded != 0 {
			t.Errorf("Incorrect Metadata.discard wanted %d got %d", 0, m.discarded)
		}
		// Test Scan Markers
		buf, _ := m.br.Peek(16)
		err := m.scanMarkers(buf)
		if err != nil {
			t.Errorf("Incorrect Scan Markers error wanted %s got %s", err.Error(), ErrNoJPEGMarker)
		}

		data = []byte{markerFirstByte, markerSOI, markerFirstByte, markerEOI, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		r = bytes.NewReader(data)
		m = newMetadata(bufio.NewReader(r), nil, nil)

		// Test SOI
		buf, _ = m.br.Peek(16)
		err = m.scanMarkers(buf)
		if m.discarded != 2 || m.pos != 1 || err != nil {
			t.Errorf("Incorrect JPEG Start of Image error wanted discarded %d got %d", 2, m.discarded)
		}

		// Test EOI
		buf, _ = m.br.Peek(16)
		err = m.scanMarkers(buf)
		if m.discarded != 4 || m.pos != 0 || err != nil {
			t.Errorf("Incorrect JPEG End of Image error wanted discarded %d got %d", 4, m.discarded)
		}

		// Test Scan JPEG
		m, err = ScanJPEG(bufio.NewReader(bytes.NewReader(data)), nil, nil)
		if err != ErrNoJPEGMarker {
			t.Errorf("Incorrect JPEG error at discarded %d wanted %s got %s", m.discarded, ErrNoJPEGMarker, err.Error())
		}

	}
}
