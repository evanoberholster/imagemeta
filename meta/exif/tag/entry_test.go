package tag

import (
	"testing"

	"github.com/evanoberholster/imagemeta/meta/utils"
)

func TestEntryBasics(t *testing.T) {
	t.Parallel()

	e := NewEntry(TagMake, TypeASCII, 4, 0x31323300, IFD0, 0, utils.LittleEndian)

	if e.ID != TagMake || e.Type != TypeASCII || e.UnitCount != 4 {
		t.Fatalf("NewEntry() produced unexpected fields: %+v", e)
	}
	if got, want := e.Name(), "Make"; got != want {
		t.Fatalf("Entry.Name() = %q, want %q", got, want)
	}
	if got, want := e.Size(), uint32(4); got != want {
		t.Fatalf("Entry.Size() = %d, want %d", got, want)
	}
	if !e.IsEmbedded() {
		t.Fatal("Entry.IsEmbedded() should be true for 4-byte ASCII")
	}
	if !e.IsType(TypeASCII) {
		t.Fatal("Entry.IsType(TypeASCII) should be true")
	}
	if e.IsIfd() {
		t.Fatal("Entry.IsIfd() should be false for ASCII")
	}
	if !e.IsValid() {
		t.Fatal("Entry.IsValid() should be true for ASCII")
	}
}

func TestEntryEmbeddedValue(t *testing.T) {
	t.Parallel()

	t.Run("little-endian", func(t *testing.T) {
		t.Parallel()
		e := NewEntry(TagMake, TypeASCII, 4, 0x04030201, IFD0, 0, utils.LittleEndian)
		var dst [4]byte
		e.EmbeddedValue(dst[:])
		if got := dst; got != [4]byte{0x01, 0x02, 0x03, 0x04} {
			t.Fatalf("EmbeddedValue() = %#v", got)
		}
	})

	t.Run("big-endian", func(t *testing.T) {
		t.Parallel()
		e := NewEntry(TagMake, TypeASCII, 4, 0x04030201, IFD0, 0, utils.BigEndian)
		var dst [4]byte
		e.EmbeddedValue(dst[:])
		if got := dst; got != [4]byte{0x04, 0x03, 0x02, 0x01} {
			t.Fatalf("EmbeddedValue() = %#v", got)
		}
	})
}

func TestEntryEmbeddedShort(t *testing.T) {
	t.Parallel()

	t.Run("little-endian", func(t *testing.T) {
		t.Parallel()
		e := NewEntry(TagMake, TypeShort, 1, 0x04030201, IFD0, 0, utils.LittleEndian)
		if got, want := e.EmbeddedShort(), uint16(0x0201); got != want {
			t.Fatalf("EmbeddedShort() = %#x, want %#x", got, want)
		}
	})

	t.Run("big-endian", func(t *testing.T) {
		t.Parallel()
		e := NewEntry(TagMake, TypeShort, 1, 0x04030201, IFD0, 0, utils.BigEndian)
		if got, want := e.EmbeddedShort(), uint16(0x0403); got != want {
			t.Fatalf("EmbeddedShort() = %#x, want %#x", got, want)
		}
	})
}

func TestEntryEmbeddedShorts(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		byteOrder utils.ByteOrder
		unitCount uint32
		dstLen    int
		want      []uint16
	}{
		{
			name:      "little-endian-two-values",
			byteOrder: utils.LittleEndian,
			unitCount: 2,
			dstLen:    2,
			want:      []uint16{0x0201, 0x0403},
		},
		{
			name:      "big-endian-two-values",
			byteOrder: utils.BigEndian,
			unitCount: 2,
			dstLen:    2,
			want:      []uint16{0x0403, 0x0201},
		},
		{
			name:      "unit-count-one",
			byteOrder: utils.LittleEndian,
			unitCount: 1,
			dstLen:    2,
			want:      []uint16{0x0201},
		},
		{
			name:      "dst-len-one",
			byteOrder: utils.LittleEndian,
			unitCount: 2,
			dstLen:    1,
			want:      []uint16{0x0201},
		},
		{
			name:      "zero-unit-count",
			byteOrder: utils.LittleEndian,
			unitCount: 0,
			dstLen:    2,
			want:      nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			e := NewEntry(TagMake, TypeShort, tt.unitCount, 0x04030201, IFD0, 0, tt.byteOrder)
			dst := make([]uint16, tt.dstLen)
			gotN := e.EmbeddedShorts(dst)
			if gotN != len(tt.want) {
				t.Fatalf("EmbeddedShorts() count = %d, want %d", gotN, len(tt.want))
			}
			for i := 0; i < gotN; i++ {
				if dst[i] != tt.want[i] {
					t.Fatalf("EmbeddedShorts()[%d] = %#x, want %#x", i, dst[i], tt.want[i])
				}
			}
		})
	}
}

func TestEntryEmbeddedLong(t *testing.T) {
	t.Parallel()

	t.Run("little-endian", func(t *testing.T) {
		t.Parallel()
		e := NewEntry(TagMake, TypeLong, 1, 0x04030201, IFD0, 0, utils.LittleEndian)
		if got, want := e.EmbeddedLong(), uint32(0x04030201); got != want {
			t.Fatalf("EmbeddedLong() = %#x, want %#x", got, want)
		}
	})

	t.Run("big-endian", func(t *testing.T) {
		t.Parallel()
		e := NewEntry(TagMake, TypeLong, 1, 0x01020304, IFD0, 0, utils.BigEndian)
		if got, want := e.EmbeddedLong(), uint32(0x01020304); got != want {
			t.Fatalf("EmbeddedLong() = %#x, want %#x", got, want)
		}
	})
}

func TestEntryIsEmbeddedAndIsIfd(t *testing.T) {
	t.Parallel()

	embedded := NewEntry(TagDateTime, TypeShort, 2, 0, IFD0, 0, utils.LittleEndian) // 2*2 = 4
	if !embedded.IsEmbedded() {
		t.Fatal("Entry.IsEmbedded() should be true when size <= 4")
	}

	pointer := NewEntry(TagExifIFDPointer, TypeIfd, 1, 0, IFD0, 0, utils.LittleEndian)
	if pointer.IsEmbedded() {
		t.Fatal("Entry.IsEmbedded() should be false for TypeIfd")
	}
	if !pointer.IsIfd() {
		t.Fatal("Entry.IsIfd() should be true for TypeIfd")
	}
}

func TestEntryIsEmbeddedBoundaries(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		typ  Type
		u    uint32
		want bool
	}{
		{name: "ascii_no_nul_4", typ: TypeASCIINoNul, u: 4, want: true},
		{name: "ascii_no_nul_5", typ: TypeASCIINoNul, u: 5, want: false},
		{name: "short_2", typ: TypeShort, u: 2, want: true},
		{name: "short_3", typ: TypeShort, u: 3, want: false},
		{name: "float_1", typ: TypeFloat, u: 1, want: true},
		{name: "float_2", typ: TypeFloat, u: 2, want: false},
		{name: "ifd_1", typ: TypeIfd, u: 1, want: false},
		{name: "double_1", typ: TypeDouble, u: 1, want: false},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			e := NewEntry(TagMake, tt.typ, tt.u, 0, IFD0, 0, utils.LittleEndian)
			if got := e.IsEmbedded(); got != tt.want {
				t.Fatalf("IsEmbedded() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEntryChildDirectory(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		entry         Entry
		wantType      IfdType
		wantIndex     int8
		wantOffset    uint32
		wantBase      uint32
		wantByteOrder utils.ByteOrder
	}{
		{
			name:          "ifd0 exif pointer",
			entry:         NewEntry(TagExifIFDPointer, TypeIfd, 1, 0x120, IFD0, 0, utils.LittleEndian),
			wantType:      ExifIFD,
			wantIndex:     0,
			wantOffset:    0x120,
			wantBase:      0,
			wantByteOrder: utils.LittleEndian,
		},
		{
			name:          "ifd0 gps pointer",
			entry:         NewEntry(TagGPSIFDPointer, TypeIfd, 1, 0x220, IFD0, 0, utils.BigEndian),
			wantType:      GPSIFD,
			wantIndex:     0,
			wantOffset:    0x220,
			wantBase:      0,
			wantByteOrder: utils.BigEndian,
		},
		{
			name:          "ifd0 next ifd",
			entry:         NewEntry(TagNextIFD, TypeIfd, 1, 0x320, IFD0, 0, utils.LittleEndian),
			wantType:      IFD1,
			wantIndex:     1,
			wantOffset:    0x320,
			wantBase:      0,
			wantByteOrder: utils.LittleEndian,
		},
		{
			name:          "ifd1 next ifd",
			entry:         NewEntry(TagNextIFD, TypeIfd, 1, 0x420, IFD1, 5, utils.LittleEndian),
			wantType:      IFD2,
			wantIndex:     6,
			wantOffset:    0x420,
			wantBase:      0,
			wantByteOrder: utils.LittleEndian,
		},
		{
			name:          "exif maker note",
			entry:         NewEntry(TagMakerNote, TypeIfd, 1, 0x520, ExifIFD, 2, utils.LittleEndian),
			wantType:      MakerNoteIFD,
			wantIndex:     2,
			wantOffset:    0x520,
			wantBase:      0,
			wantByteOrder: utils.LittleEndian,
		},
		{
			name:          "subifd passthrough",
			entry:         NewEntry(TagExposureTime, TypeLong, 1, 0x620, SubIFD3, 3, utils.LittleEndian),
			wantType:      SubIFD3,
			wantIndex:     3,
			wantOffset:    0x620,
			wantBase:      0,
			wantByteOrder: utils.LittleEndian,
		},
		{
			name:          "unknown mapping",
			entry:         NewEntry(TagExposureTime, TypeLong, 1, 0x720, GPSIFD, 1, utils.LittleEndian),
			wantType:      Unknown,
			wantIndex:     1,
			wantOffset:    0x720,
			wantBase:      0,
			wantByteOrder: utils.LittleEndian,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.entry.ChildDirectory()
			if got.Type != tt.wantType {
				t.Fatalf("ChildDirectory().Type = %v, want %v", got.Type, tt.wantType)
			}
			if got.Index != tt.wantIndex {
				t.Fatalf("ChildDirectory().Index = %d, want %d", got.Index, tt.wantIndex)
			}
			if got.Offset != tt.wantOffset {
				t.Fatalf("ChildDirectory().Offset = 0x%x, want 0x%x", got.Offset, tt.wantOffset)
			}
			if got.BaseOffset != tt.wantBase {
				t.Fatalf("ChildDirectory().BaseOffset = 0x%x, want 0x%x", got.BaseOffset, tt.wantBase)
			}
			if got.ByteOrder != tt.wantByteOrder {
				t.Fatalf("ChildDirectory().ByteOrder = %v, want %v", got.ByteOrder, tt.wantByteOrder)
			}
		})
	}
}
