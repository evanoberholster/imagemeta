package xmp

import (
	"strings"
	"testing"
)

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

func TestIdentifyNamespaceFallback(t *testing.T) {
	tests := []struct {
		in   string
		want Namespace
	}{
		{in: "apple-fi", want: AppleFiNS},
		{in: "darktable", want: DarktableNS},
		{in: "mwg-rs", want: MwgRSNS},
		{in: "pmi", want: PmiNS},
		{in: "stArea", want: StAreaNS},
		{in: "stDim", want: StDimNS},
	}
	for _, tc := range tests {
		got := IdentifyNamespace([]byte(tc.in))
		if got != tc.want {
			t.Fatalf("IdentifyNamespace(%q) = %s, want %s", tc.in, got, tc.want)
		}
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
		{in: "DateTime", want: DateTime},
		{in: "FocalLengthIn35mmFormat", want: FocalLengthIn35mmFilm},
		{in: "GPSDateTime", want: GPSTimeStamp},
		{in: "Opto-ElectricConvFactor", want: OECF},
	}

	for _, tc := range tests {
		got := IdentifyName([]byte(tc.in))
		if got != tc.want {
			t.Fatalf("IdentifyName(%q) = %s, want %s", tc.in, got, tc.want)
		}
	}
}

func TestIdentifyNameDublinCoreProperties(t *testing.T) {
	tests := []struct {
		in   string
		want Name
	}{
		{in: "contributor", want: Contributor},
		{in: "coverage", want: Coverage},
		{in: "date", want: Date},
		{in: "identifier", want: Identifier},
		{in: "language", want: Language},
		{in: "publisher", want: Publisher},
		{in: "relation", want: Relation},
		{in: "source", want: Source},
		{in: "type", want: Type},
	}

	for _, tc := range tests {
		got := IdentifyName([]byte(tc.in))
		if got != tc.want {
			t.Fatalf("IdentifyName(%q) = %s, want %s", tc.in, got, tc.want)
		}
	}
}

func TestIdentifyNameKnownButNotDecodedTags(t *testing.T) {
	crsTags := []string{
		"AutoLateralCA",
		"Blacks2012",
		"CameraProfile",
		"CameraProfileDigest",
		"Clarity2012",
		"ColorNoiseReduction",
		"ColorNoiseReductionDetail",
		"ColorNoiseReductionSmoothness",
		"Contrast2012",
		"ConvertToGrayscale",
		"DefringeGreenAmount",
		"DefringeGreenHueHi",
		"DefringeGreenHueLo",
		"DefringePurpleAmount",
		"DefringePurpleHueHi",
		"DefringePurpleHueLo",
		"Dehaze",
		"Exposure2012",
		"GrainAmount",
		"GrainFrequency",
		"GrainSeed",
		"GrainSize",
		"HasCrop",
		"HasSettings",
		"Highlights2012",
		"LensManualDistortionAmount",
		"LensProfileChromaticAberrationScale",
		"LensProfileDigest",
		"LensProfileDistortionScale",
		"LensProfileEnable",
		"LensProfileFilename",
		"LensProfileName",
		"LensProfileSetup",
		"LensProfileVignettingScale",
		"LookName",
		"LuminanceNoiseReductionContrast",
		"LuminanceNoiseReductionDetail",
		"LuminanceSmoothing",
		"OverrideLookVignette",
		"ParametricDarks",
		"ParametricHighlightSplit",
		"ParametricHighlights",
		"ParametricLights",
		"ParametricMidtoneSplit",
		"ParametricShadowSplit",
		"ParametricShadows",
		"PerspectiveAspect",
		"PerspectiveHorizontal",
		"PerspectiveRotate",
		"PerspectiveScale",
		"PerspectiveUpright",
		"PerspectiveVertical",
		"PerspectiveX",
		"PerspectiveY",
		"PostCropVignetteAmount",
		"PostCropVignetteFeather",
		"PostCropVignetteHighlightContrast",
		"PostCropVignetteMidpoint",
		"PostCropVignetteRoundness",
		"PostCropVignetteStyle",
		"ProcessVersion",
		"ShadowTint",
		"Shadows2012",
		"SharpenDetail",
		"SharpenEdgeMasking",
		"SharpenRadius",
		"SplitToningBalance",
		"SplitToningHighlightHue",
		"SplitToningHighlightSaturation",
		"SplitToningShadowHue",
		"SplitToningShadowSaturation",
		"Tint",
		"ToneCurveName",
		"ToneCurveName2012",
		"ToneMapStrength",
		"UprightCenterMode",
		"UprightCenterNormX",
		"UprightCenterNormY",
		"UprightFocalLength35mm",
		"UprightFocalMode",
		"UprightFourSegmentsCount",
		"UprightPreview",
		"UprightTransformCount",
		"UprightVersion",
		"Version",
		"Vibrance",
		"VignetteAmount",
		"Whites2012",
	}

	for _, name := range crsTags {
		if got := IdentifyName([]byte(name)); got == UnknownPropertyName {
			t.Fatalf("IdentifyName(%q) = UnknownPropertyName", name)
		}
	}

	if got := IdentifyName([]byte("DerivedFromInstanceID")); got == UnknownPropertyName {
		t.Fatalf("IdentifyName(%q) = UnknownPropertyName", "DerivedFromInstanceID")
	}
}

func TestIdentifyNameCanonicalCoverage(t *testing.T) {
	for name, s := range mapNameString {
		if name == UnknownPropertyName || s == "Unknown" {
			continue
		}

		inputs := []string{s, strings.ToLower(s), strings.ToUpper(s)}
		for _, in := range inputs {
			got := identifyName([]byte(in))
			if got != name {
				t.Fatalf("identifyName(%q)=%s, want %s", in, got, name)
			}
		}
	}
}
