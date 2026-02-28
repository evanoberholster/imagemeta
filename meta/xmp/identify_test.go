package xmp

import "testing"

func TestPropertyIdentity(t *testing.T) {
	p := NewProperty(XmpNS, RDF)
	if XmpNS != p.Namespace() {
		t.Fatalf("namespace = %s, want %s", p.Namespace(), XmpNS)
	}
	if RDF != p.Name() {
		t.Fatalf("name = %s, want %s", p.Name(), RDF)
	}

	p1 := IdentifyProperty([]byte("xmp"), []byte("RDF"))
	if p1 != p {
		t.Fatalf("property = %s, want %s", p1, p)
	}
	if p.String() != "xmp:RDF" {
		t.Fatalf("property string = %q", p.String())
	}
}

func TestIdentifyNameAndNamespace(t *testing.T) {
	if IdentifyName([]byte("subject")).String() != "subject" {
		t.Fatalf("IdentifyName(subject) mismatch")
	}
	if IdentifyNamespace([]byte("exif")).String() != "exif" {
		t.Fatalf("IdentifyNamespace(exif) mismatch")
	}
}

func TestIdentifyNameExifToolAliases(t *testing.T) {
	tests := []struct {
		in   string
		want Name
	}{
		{in: "XMPToolkit", want: XMPToolkit},
		{in: "HistoryAction", want: Action},
		{in: "HistoryChanged", want: Changed},
		{in: "HistoryInstanceID", want: InstanceID},
		{in: "HistoryParameters", want: Parameters},
		{in: "HistorySoftwareAgent", want: SoftwareAgent},
		{in: "HistoryWhen", want: When},
		{in: "HueAdjustmentRed", want: HueAdjustmentRed},
		{in: "SaturationAdjustmentBlue", want: SaturationAdjustmentBlue},
		{in: "LuminanceAdjustmentMagenta", want: LuminanceAdjustmentMagenta},
		{in: "RedHue", want: HueAdjustmentRed},
		{in: "GreenHue", want: HueAdjustmentGreen},
		{in: "BlueHue", want: HueAdjustmentBlue},
		{in: "RedSaturation", want: SaturationAdjustmentRed},
		{in: "GreenSaturation", want: SaturationAdjustmentGreen},
		{in: "BlueSaturation", want: SaturationAdjustmentBlue},
		{in: "ColorTemperature", want: Temperature},
	}

	for _, tc := range tests {
		got := IdentifyName([]byte(tc.in))
		if got != tc.want {
			t.Fatalf("IdentifyName(%q) = %s, want %s", tc.in, got, tc.want)
		}
	}
}
