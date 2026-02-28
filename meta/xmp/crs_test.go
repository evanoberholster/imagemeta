package xmp

import (
	"strings"
	"testing"
)

func TestParseCRSColorTemperature(t *testing.T) {
	const src = `<x:xmpmeta xmlns:x="adobe:ns:meta/"><rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#"><rdf:Description xmlns:crs="http://ns.adobe.com/camera-raw-settings/1.0/" crs:ColorTemperature="5400"/></rdf:RDF></x:xmpmeta>`

	x, err := ParseXmp(strings.NewReader(src))
	if err != nil {
		t.Fatal(err)
	}
	if x.CRS == nil {
		t.Fatal("CRS is nil")
	}
	if x.CRS.ColorTemperature != 5400 {
		t.Fatalf("ColorTemperature = %d, want 5400", x.CRS.ColorTemperature)
	}
}

func TestParseCRSHSLNoOpCompatibility(t *testing.T) {
	const src = `<x:xmpmeta xmlns:x="adobe:ns:meta/"><rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#"><rdf:Description xmlns:crs="http://ns.adobe.com/camera-raw-settings/1.0/" crs:HueAdjustmentRed="+15" crs:RedHue="+10" crs:SaturationAdjustmentBlue="-20" crs:BlueSaturation="-25" crs:LuminanceAdjustmentMagenta="+7"/></rdf:RDF></x:xmpmeta>`

	x, err := ParseXmp(strings.NewReader(src))
	if err != nil {
		t.Fatal(err)
	}
	// HSL tags are accepted but intentionally not materialized.
	// Since there are no materialized CRS fields, CRS should remain nil.
	if x.CRS != nil {
		t.Fatalf("CRS = %+v, want nil", x.CRS)
	}
}

func TestParseCRSNoOpWithMaterializedField(t *testing.T) {
	const src = `<x:xmpmeta xmlns:x="adobe:ns:meta/"><rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#"><rdf:Description xmlns:crs="http://ns.adobe.com/camera-raw-settings/1.0/" crs:HueAdjustmentRed="+15" crs:ColorTemperature="5600"/></rdf:RDF></x:xmpmeta>`

	x, err := ParseXmp(strings.NewReader(src))
	if err != nil {
		t.Fatal(err)
	}
	if x.CRS == nil {
		t.Fatal("CRS is nil")
	}
	if x.CRS.ColorTemperature != 5600 {
		t.Fatalf("ColorTemperature = %d, want 5600", x.CRS.ColorTemperature)
	}
}
