package xmp

import (
	"math"
	"strings"
	"testing"
)

func TestParseMWGRegions(t *testing.T) {
	const src = `<x:xmpmeta xmlns:x="adobe:ns:meta/">
<rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">
<rdf:Description rdf:about=""
	xmlns:mwg-rs="http://www.metadataworkinggroup.com/schemas/regions/"
	xmlns:stArea="http://ns.adobe.com/xmp/sType/Area#"
	xmlns:apple-fi="http://ns.apple.com/faceinfo/1.0/"
	xmlns:stDim="http://ns.adobe.com/xap/1.0/sType/Dimensions#">
	<mwg-rs:Regions rdf:parseType="Resource">
		<mwg-rs:RegionList>
			<rdf:Seq>
				<rdf:li rdf:parseType="Resource">
					<mwg-rs:Area rdf:parseType="Resource">
						<stArea:y>0.18099999999999999</stArea:y>
						<stArea:w>0.10699999999999998</stArea:w>
						<stArea:x>0.30049999999999999</stArea:x>
						<stArea:h>0.14200000000000002</stArea:h>
						<stArea:unit>normalized</stArea:unit>
					</mwg-rs:Area>
					<mwg-rs:Type>Face</mwg-rs:Type>
					<mwg-rs:Extensions rdf:parseType="Resource">
						<apple-fi:AngleInfoYaw>0</apple-fi:AngleInfoYaw>
						<apple-fi:AngleInfoRoll>0</apple-fi:AngleInfoRoll>
						<apple-fi:ConfidenceLevel>366</apple-fi:ConfidenceLevel>
						<apple-fi:FaceID>1</apple-fi:FaceID>
					</mwg-rs:Extensions>
				</rdf:li>
				<rdf:li rdf:parseType="Resource">
					<mwg-rs:Area rdf:parseType="Resource">
						<stArea:y>0.33600000000000002</stArea:y>
						<stArea:w>0.125</stArea:w>
						<stArea:x>0.086499999999999994</stArea:x>
						<stArea:h>0.16600000000000004</stArea:h>
						<stArea:unit>normalized</stArea:unit>
					</mwg-rs:Area>
					<mwg-rs:Type>Face</mwg-rs:Type>
					<mwg-rs:Extensions rdf:parseType="Resource">
						<apple-fi:AngleInfoYaw>315</apple-fi:AngleInfoYaw>
						<apple-fi:AngleInfoRoll>0</apple-fi:AngleInfoRoll>
						<apple-fi:ConfidenceLevel>333</apple-fi:ConfidenceLevel>
						<apple-fi:FaceID>2</apple-fi:FaceID>
					</mwg-rs:Extensions>
				</rdf:li>
			</rdf:Seq>
		</mwg-rs:RegionList>
		<mwg-rs:AppliedToDimensions rdf:parseType="Resource">
			<stDim:h>2320</stDim:h>
			<stDim:w>3088</stDim:w>
			<stDim:unit>pixel</stDim:unit>
		</mwg-rs:AppliedToDimensions>
	</mwg-rs:Regions>
</rdf:Description>
</rdf:RDF>
</x:xmpmeta>`

	got, err := ParseXmp(strings.NewReader(src))
	if err != nil {
		t.Fatal(err)
	}
	if got.Regions == nil {
		t.Fatal("Regions is nil")
	}

	if got.Regions.AppliedToDimensions.H != 2320 || got.Regions.AppliedToDimensions.W != 3088 || got.Regions.AppliedToDimensions.Unit != "pixel" {
		t.Fatalf("AppliedToDimensions = %+v", got.Regions.AppliedToDimensions)
	}
	if len(got.Regions.RegionList) != 2 {
		t.Fatalf("len(RegionList) = %d", len(got.Regions.RegionList))
	}

	first := got.Regions.RegionList[0]
	if first.Type != RegionType("Face") {
		t.Fatalf("first.Type = %q", first.Type)
	}
	if !approxFloat32(first.Area.X, 0.3005, 0.0001) ||
		!approxFloat32(first.Area.Y, 0.1810, 0.0001) ||
		!approxFloat32(first.Area.W, 0.1070, 0.0001) ||
		!approxFloat32(first.Area.H, 0.1420, 0.0001) {
		t.Fatalf("first.Area = %#v", first.Area)
	}
	if first.Extensions.FaceID != "1" || !approxFloat32(first.Extensions.ConfidenceLevel, 366, 0.0001) {
		t.Fatalf("first.Extensions = %+v", first.Extensions)
	}

	second := got.Regions.RegionList[1]
	if second.Type != RegionType("Face") {
		t.Fatalf("second.Type = %q", second.Type)
	}
	if !approxFloat32(second.Area.X, 0.0865, 0.0001) ||
		!approxFloat32(second.Area.Y, 0.3360, 0.0001) ||
		!approxFloat32(second.Area.W, 0.1250, 0.0001) ||
		!approxFloat32(second.Area.H, 0.1660, 0.0001) {
		t.Fatalf("second.Area = %#v", second.Area)
	}
	if !approxFloat32(second.Extensions.AngleInfoYaw, 315, 0.0001) || second.Extensions.FaceID != "2" {
		t.Fatalf("second.Extensions = %+v", second.Extensions)
	}
}

func TestParseMWGRegionPixelAreaScaling(t *testing.T) {
	const src = `<x:xmpmeta xmlns:x="adobe:ns:meta/">
<rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">
<rdf:Description rdf:about=""
	xmlns:mwg-rs="http://www.metadataworkinggroup.com/schemas/regions/"
	xmlns:stArea="http://ns.adobe.com/xmp/sType/Area#">
	<mwg-rs:Regions rdf:parseType="Resource">
		<mwg-rs:RegionList>
			<rdf:Seq>
				<rdf:li rdf:parseType="Resource">
					<mwg-rs:Area rdf:parseType="Resource">
						<stArea:x>120</stArea:x>
						<stArea:y>250</stArea:y>
						<stArea:w>80</stArea:w>
						<stArea:h>45</stArea:h>
						<stArea:unit>pixel</stArea:unit>
					</mwg-rs:Area>
					<mwg-rs:Type>Face</mwg-rs:Type>
				</rdf:li>
			</rdf:Seq>
		</mwg-rs:RegionList>
	</mwg-rs:Regions>
</rdf:Description>
</rdf:RDF>
</x:xmpmeta>`

	got, err := ParseXmp(strings.NewReader(src))
	if err != nil {
		t.Fatal(err)
	}
	if got.Regions == nil {
		t.Fatal("Regions is nil")
	}
	if len(got.Regions.RegionList) != 1 {
		t.Fatalf("len(RegionList) = %d", len(got.Regions.RegionList))
	}
	if !approxFloat32(got.Regions.RegionList[0].Area.X, 120, 0.0001) ||
		!approxFloat32(got.Regions.RegionList[0].Area.Y, 250, 0.0001) ||
		!approxFloat32(got.Regions.RegionList[0].Area.W, 80, 0.0001) ||
		!approxFloat32(got.Regions.RegionList[0].Area.H, 45, 0.0001) {
		t.Fatalf("Area = %#v", got.Regions.RegionList[0].Area)
	}
}

func approxFloat32(got, want, epsilon float32) bool {
	return float32(math.Abs(float64(got-want))) <= epsilon
}
