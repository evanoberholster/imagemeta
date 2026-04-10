package tag

import (
	"testing"
)

func TestNameFor(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		ifdType IfdType
		tagID   ID
		want    string
	}{
		{name: "root", ifdType: IFD0, tagID: TagMake, want: "Make"},
		{name: "exif", ifdType: ExifIFD, tagID: TagExposureTime, want: "ExposureTime"},
		{name: "gps", ifdType: GPSIFD, tagID: TagGPSLatitude, want: "GPSLatitude"},
		{name: "gps dest distance", ifdType: GPSIFD, tagID: TagGPSDestDistance, want: "GPSDestDistance"},
		{name: "unknown tag", ifdType: IFD0, tagID: ID(0xbeef), want: "0xbeef"},
		{name: "wrong ifd", ifdType: ExifIFD, tagID: TagMake, want: "0x010f"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := NameFor(tt.ifdType, tt.tagID); got != tt.want {
				t.Fatalf("NameFor(%v, %v) = %q, want %q", tt.ifdType, tt.tagID, got, tt.want)
			}
		})
	}
}

func TestValueNameFor(t *testing.T) {
	t.Parallel()

	if got, want := ValueNameFor(IFD0, TagOrientation, 1), "Horizontal (normal)"; got != want {
		t.Fatalf("ValueNameFor() = %q, want %q", got, want)
	}
	if got, want := ValueNameFor(ExifIFD, TagExposureMode, 2), "Auto bracket"; got != want {
		t.Fatalf("ValueNameFor() = %q, want %q", got, want)
	}
	if got, want := ValueNameFor(ExifIFD, TagExposureMode, 255), "255"; got != want {
		t.Fatalf("ValueNameFor() unknown = %q, want %q", got, want)
	}
}

func TestParseValueID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		ifdType IfdType
		tagID   ID
		raw     string
		want    uint32
		wantOK  bool
	}{
		{
			name:    "numeric decimal",
			ifdType: IFD0,
			tagID:   TagOrientation,
			raw:     "6",
			want:    6,
			wantOK:  true,
		},
		{
			name:    "numeric with suffix",
			ifdType: IFD0,
			tagID:   TagOrientation,
			raw:     "6 (Rotate 90 CW)",
			want:    6,
			wantOK:  true,
		},
		{
			name:    "numeric hex",
			ifdType: ExifIFD,
			tagID:   TagFlash,
			raw:     "0x1f",
			want:    0x1f,
			wantOK:  true,
		},
		{
			name:    "enum alias normalized",
			ifdType: IFD0,
			tagID:   TagOrientation,
			raw:     " horizontal ",
			want:    1,
			wantOK:  true,
		},
		{
			name:    "enum full name",
			ifdType: ExifIFD,
			tagID:   TagExposureProgram,
			raw:     "Aperture-priority AE",
			want:    3,
			wantOK:  true,
		},
		{
			name:    "unsupported tag",
			ifdType: GPSIFD,
			tagID:   TagGPSLatitude,
			raw:     "1",
			want:    1,
			wantOK:  true,
		},
		{
			name:    "unsupported enum text",
			ifdType: GPSIFD,
			tagID:   TagGPSLatitude,
			raw:     "north",
			want:    0,
			wantOK:  false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, ok := ParseValueID(tt.ifdType, tt.tagID, tt.raw)
			if ok != tt.wantOK {
				t.Fatalf("ParseValueID() ok = %v, want %v", ok, tt.wantOK)
			}
			if got != tt.want {
				t.Fatalf("ParseValueID() value = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestNormalizeValueStringAndParseUint(t *testing.T) {
	t.Parallel()

	if got, want := normalizeValueString("  Auto_Bracket  "), "auto bracket"; got != want {
		t.Fatalf("normalizeValueString() = %q, want %q", got, want)
	}

	if got, ok := parseUint("0X10"); !ok || got != 16 {
		t.Fatalf("parseUint hex = (%d,%v), want (16,true)", got, ok)
	}
	if got, ok := parseUint(" 23, "); !ok || got != 23 {
		t.Fatalf("parseUint decimal = (%d,%v), want (23,true)", got, ok)
	}
	if _, ok := parseUint("abc"); ok {
		t.Fatal("parseUint(\"abc\") should fail")
	}
}
