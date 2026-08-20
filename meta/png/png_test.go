package png

import (
	"bytes"
	"encoding/binary"
	"errors"
	"testing"

	"github.com/evanoberholster/imagemeta/meta/utils"
)

const pngSignature = "\x89PNG\r\n\x1a\n"

type oneByteReadSeeker struct {
	*bytes.Reader
}

var errSeekRecorded = errors.New("seek recorded")

type seekRecorder struct {
	*bytes.Reader
	offset int64
}

func (r *seekRecorder) Seek(offset int64, _ int) (int64, error) {
	r.offset = offset
	return 0, errSeekRecorded
}

func (r oneByteReadSeeker) Read(p []byte) (int, error) {
	if len(p) > 1 {
		p = p[:1]
	}
	return r.Reader.Read(p)
}

func buildPNGWithExif(tiff []byte) []byte {
	buf := make([]byte, 0, len(pngSignature)+15+8+len(tiff))
	buf = append(buf, pngSignature...)

	var size [4]byte
	binary.BigEndian.PutUint32(size[:], 3)
	buf = append(buf, size[:]...)
	buf = append(buf, "tEXt"...)
	buf = append(buf, "abc"...)
	buf = append(buf, 0, 0, 0, 0)

	binary.BigEndian.PutUint32(size[:], uint32(len(tiff)))
	buf = append(buf, size[:]...)
	buf = append(buf, "eXIf"...)
	buf = append(buf, tiff...)
	return buf
}

func TestScanPngHeader(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		tiff       []byte
		wantOrder  utils.ByteOrder
		wantOffset uint32
		wantBig    bool
	}{
		{
			name:       "classic little endian",
			tiff:       []byte{'I', 'I', '*', 0, 12, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			wantOrder:  utils.LittleEndian,
			wantOffset: 12,
		},
		{
			name:       "classic big endian",
			tiff:       []byte{'M', 'M', 0, '*', 0, 0, 0, 8, 0, 0},
			wantOrder:  utils.BigEndian,
			wantOffset: 8,
		},
		{
			name: "BigTIFF little endian",
			tiff: []byte{
				'I', 'I', '+', 0, 8, 0, 0, 0,
				16, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
			},
			wantOrder:  utils.LittleEndian,
			wantOffset: 16,
			wantBig:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			data := buildPNGWithExif(tt.tiff)
			r := bytes.NewReader(data)
			h, err := ScanPngHeader(r)
			if err != nil {
				t.Fatalf("ScanPngHeader() error = %v", err)
			}
			if h.ByteOrder != tt.wantOrder {
				t.Fatalf("ByteOrder = %v, want %v", h.ByteOrder, tt.wantOrder)
			}
			if h.FirstIfdOffset != tt.wantOffset {
				t.Fatalf("FirstIfdOffset = %d, want %d", h.FirstIfdOffset, tt.wantOffset)
			}
			if h.ExifLength != uint32(len(tt.tiff)) {
				t.Fatalf("ExifLength = %d, want %d", h.ExifLength, len(tt.tiff))
			}
			if h.BigTIFF != tt.wantBig {
				t.Fatalf("BigTIFF = %t, want %t", h.BigTIFF, tt.wantBig)
			}
			if got, want := r.Size()-int64(r.Len()), int64(h.TiffHeaderOffset); got != want {
				t.Fatalf("reader offset = %d, want %d", got, want)
			}
		})
	}
}

func TestScanPngHeaderHandlesShortReads(t *testing.T) {
	t.Parallel()

	data := buildPNGWithExif([]byte{'M', 'M', 0, '*', 0, 0, 0, 8, 0, 0})
	r := oneByteReadSeeker{Reader: bytes.NewReader(data)}
	h, err := ScanPngHeader(r)
	if err != nil {
		t.Fatalf("ScanPngHeader() error = %v", err)
	}
	if h.ByteOrder != utils.BigEndian || h.FirstIfdOffset != 8 {
		t.Fatalf("ScanPngHeader() = %+v", h)
	}
}

func TestScanPngHeaderRejectsInvalidExifHeader(t *testing.T) {
	t.Parallel()

	for _, tiff := range [][]byte{
		{'I', 'I', '*'},
		{'N', 'O', 'P', 'E', 0, 0, 0, 8},
		{'I', 'I', '+', 0, 4, 0, 0, 0, 16, 0, 0, 0, 0, 0, 0, 0},
	} {
		if _, err := ScanPngHeader(bytes.NewReader(buildPNGWithExif(tiff))); err == nil {
			t.Fatalf("ScanPngHeader(%v) error = nil", tiff)
		}
	}
}

func TestScanPngHeaderChunkSkipDoesNotOverflow(t *testing.T) {
	t.Parallel()

	data := append([]byte(pngSignature), 0xff, 0xff, 0xff, 0xff)
	data = append(data, "IDAT"...)
	r := &seekRecorder{Reader: bytes.NewReader(data)}
	_, err := ScanPngHeader(r)
	if !errors.Is(err, errSeekRecorded) {
		t.Fatalf("ScanPngHeader() error = %v, want %v", err, errSeekRecorded)
	}
	if want := int64(^uint32(0)) + pngChunkCRCSize; r.offset != want {
		t.Fatalf("Seek offset = %d, want %d", r.offset, want)
	}
}

func BenchmarkScanPngHeader(b *testing.B) {
	data := buildPNGWithExif([]byte{'M', 'M', 0, '*', 0, 0, 0, 8, 0, 0})
	var r bytes.Reader

	b.ReportAllocs()
	b.SetBytes(int64(len(data)))
	for i := 0; i < b.N; i++ {
		r.Reset(data)
		if _, err := ScanPngHeader(&r); err != nil {
			b.Fatal(err)
		}
	}
}
