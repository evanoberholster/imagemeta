package tag

import "testing"

func TestIDString(t *testing.T) {
	t.Parallel()

	if got, want := ID(0x10).String(), "0x0010"; got != want {
		t.Fatalf("ID.String() = %q, want %q", got, want)
	}
	if got, want := ID(0xffff).String(), "0xffff"; got != want {
		t.Fatalf("ID.String() = %q, want %q", got, want)
	}
}

func TestTypeProperties(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		typ       Type
		wantSize  uint8
		wantName  string
		wantValid bool
	}{
		{name: "byte", typ: TypeByte, wantSize: 1, wantName: "BYTE", wantValid: true},
		{name: "ascii", typ: TypeASCII, wantSize: 1, wantName: "ASCII", wantValid: true},
		{name: "short", typ: TypeShort, wantSize: 2, wantName: "SHORT", wantValid: true},
		{name: "long", typ: TypeLong, wantSize: 4, wantName: "LONG", wantValid: true},
		{name: "rational", typ: TypeRational, wantSize: 8, wantName: "RATIONAL", wantValid: true},
		{name: "undefined", typ: TypeUndefined, wantSize: 1, wantName: "UNDEFINED", wantValid: true},
		{name: "asciinonul", typ: TypeASCIINoNul, wantSize: 1, wantName: "_ASCII_NO_NUL", wantValid: true},
		{name: "ifd", typ: TypeIfd, wantSize: 4, wantName: "IFD", wantValid: true},
		{name: "unknown", typ: TypeUnknown, wantSize: 0, wantName: "UNKNOWN", wantValid: false},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.typ.Size(); got != tt.wantSize {
				t.Fatalf("Type.Size() = %d, want %d", got, tt.wantSize)
			}
			if got := tt.typ.String(); got != tt.wantName {
				t.Fatalf("Type.String() = %q, want %q", got, tt.wantName)
			}
			if got := tt.typ.IsValid(); got != tt.wantValid {
				t.Fatalf("Type.IsValid() = %v, want %v", got, tt.wantValid)
			}
			if !tt.typ.Is(tt.typ) {
				t.Fatalf("Type.Is() should be true for same type %v", tt.typ)
			}
		})
	}
}
