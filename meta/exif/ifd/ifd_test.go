package ifd

import (
	"testing"

	"github.com/evanoberholster/imagemeta/meta/utils"
)

func TestTypeProperties(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		typ       Type
		wantName  string
		wantSub   bool
		wantRoot  bool
		wantValid bool
		next      Type
		nextOK    bool
	}{
		{
			name:      "ifd0",
			typ:       IFD0,
			wantName:  "IFD0",
			wantSub:   false,
			wantRoot:  true,
			wantValid: true,
			next:      IFD1,
			nextOK:    true,
		},
		{
			name:      "ifd1",
			typ:       IFD1,
			wantName:  "IFD1",
			wantSub:   false,
			wantRoot:  true,
			wantValid: true,
			next:      IFD2,
			nextOK:    true,
		},
		{
			name:      "ifd2",
			typ:       IFD2,
			wantName:  "IFD2",
			wantSub:   false,
			wantRoot:  true,
			wantValid: true,
			next:      Unknown,
			nextOK:    false,
		},
		{
			name:      "subifd3",
			typ:       SubIFD3,
			wantName:  "SubIFD3",
			wantSub:   true,
			wantRoot:  false,
			wantValid: true,
			next:      Unknown,
			nextOK:    false,
		},
		{
			name:      "unknown",
			typ:       Unknown,
			wantName:  "UnknownIFD",
			wantSub:   false,
			wantRoot:  false,
			wantValid: false,
			next:      Unknown,
			nextOK:    false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.typ.String(); got != tt.wantName {
				t.Fatalf("Type.String() = %q, want %q", got, tt.wantName)
			}
			if got := tt.typ.IsSubIFD(); got != tt.wantSub {
				t.Fatalf("Type.IsSubIFD() = %v, want %v", got, tt.wantSub)
			}
			if got := tt.typ.IsRootIFD(); got != tt.wantRoot {
				t.Fatalf("Type.IsRootIFD() = %v, want %v", got, tt.wantRoot)
			}
			if got := tt.typ.IsValid(); got != tt.wantValid {
				t.Fatalf("Type.IsValid() = %v, want %v", got, tt.wantValid)
			}

			next, ok := tt.typ.NextRootIFD()
			if ok != tt.nextOK {
				t.Fatalf("Type.NextRootIFD() ok = %v, want %v", ok, tt.nextOK)
			}
			if next != tt.next {
				t.Fatalf("Type.NextRootIFD() next = %v, want %v", next, tt.next)
			}
		})
	}
}

func TestDirectoryNewAndString(t *testing.T) {
	t.Parallel()

	d := New(utils.LittleEndian, ExifIFD, 2, 0x1234, 0x10)

	if d.ByteOrder != utils.LittleEndian {
		t.Fatalf("ByteOrder = %v, want %v", d.ByteOrder, utils.LittleEndian)
	}
	if d.Type != ExifIFD {
		t.Fatalf("Type = %v, want %v", d.Type, ExifIFD)
	}
	if d.Index != 2 {
		t.Fatalf("Index = %d, want 2", d.Index)
	}
	if d.Offset != 0x1234 {
		t.Fatalf("Offset = 0x%x, want 0x1234", d.Offset)
	}
	if d.BaseOffset != 0x10 {
		t.Fatalf("BaseOffset = 0x%x, want 0x10", d.BaseOffset)
	}

	const want = "IFD[ExifIFD](2)@0x1234"
	if got := d.String(); got != want {
		t.Fatalf("Directory.String() = %q, want %q", got, want)
	}
}
