package xmp

import (
	"strings"
	"testing"
	"time"

	"github.com/evanoberholster/imagemeta/meta"
)

func TestParseExifAdditionalTags(t *testing.T) {
	const src = `<x:xmpmeta xmlns:x="adobe:ns:meta/"><rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#"><rdf:Description xmlns:exif="http://ns.adobe.com/exif/1.0/" exif:CFAPattern="0 1 1 2" exif:CFAPatternColumns="2" exif:CFAPatternRows="2" exif:CFAPatternValues="0 1 1 2" exif:DeviceSettingDescription="camera settings" exif:DeviceSettingDescriptionColumns="3" exif:DeviceSettingDescriptionRows="2" exif:DeviceSettingDescriptionSettings="1 2 3" exif:ExposureIndex="160/1" exif:FlashEnergy="15/10" exif:GPSAreaInformation="zone-a" exif:GPSDestBearing="123/1" exif:GPSDestBearingRef="T" exif:GPSDestDistance="2000/1" exif:GPSDestDistanceRef="K" exif:GPSDestLatitude="33,12.5N" exif:GPSDestLongitude="151,12.5E" exif:GPSHPositioningError="3/2" exif:GPSImgDirection="90/1" exif:GPSImgDirectionRef="T" exif:GPSProcessingMethod="GPS" exif:GPSSpeed="88/1" exif:GPSSpeedRef="K" exif:GPSTrack="45/1" exif:GPSTrackRef="T" exif:ImageUniqueID="img-id" exif:MakerNote="note" exif:NativeDigest="digest" exif:Opto-ElectricConvFactor="oecf" exif:OECFColumns="3" exif:OECFNames="RGB" exif:OECFRows="2" exif:OECFValues="1 2 3 4 5 6" exif:RelatedSoundFile="AUDIO.WAV" exif:SensingMethod="2" exif:SpatialFrequencyResponse="sfr" exif:SpatialFrequencyResponseColumns="2" exif:SpatialFrequencyResponseNames="R G" exif:SpatialFrequencyResponseRows="1" exif:SpatialFrequencyResponseValues="0.2 0.3" exif:SpectralSensitivity="sensitive" exif:SubjectArea="100,200,50,50" exif:SubjectDistanceRange="3" exif:SubjectLocation="100,200" exif:FocalLengthIn35mmFormat="50" exif:GPSDateTime="2025-01-02T03:04:05Z" exif:Contrast="1"/></rdf:RDF></x:xmpmeta>`

	x, err := ParseXmp(strings.NewReader(src))
	if err != nil {
		t.Fatal(err)
	}

	exif := x.Exif
	if exif.CFAPattern != "0 1 1 2" || exif.CFAPatternColumns != 2 || exif.CFAPatternRows != 2 || exif.CFAPatternValues != "0 1 1 2" {
		t.Fatalf("cfa pattern mismatch: %+v", exif)
	}
	if exif.DeviceSettingDescription != "camera settings" || exif.DeviceSettingDescriptionColumns != 3 || exif.DeviceSettingDescriptionRows != 2 || exif.DeviceSettingDescriptionSettings != "1 2 3" {
		t.Fatalf("device setting description mismatch: %+v", exif)
	}
	assertApproxFloat64(t, exif.ExposureIndex, 160, 0.0001, "Exif.ExposureIndex")
	assertApproxFloat64(t, exif.FlashEnergy, 1.5, 0.0001, "Exif.FlashEnergy")
	if exif.GPS.AreaInformation != "zone-a" {
		t.Fatalf("GPS.AreaInformation = %q", exif.GPS.AreaInformation)
	}
	assertApproxFloat64(t, exif.GPS.DestinationBearing.Value, 123, 0.0001, "Exif.GPS.DestinationBearing.Value")
	if exif.GPS.DestinationBearing.Ref != meta.GPSRefTrue {
		t.Fatalf("GPS.DestinationBearing.Ref = %v", exif.GPS.DestinationBearing.Ref)
	}
	assertApproxFloat64(t, exif.GPS.DestinationDistance.Value, 2000, 0.0001, "Exif.GPS.DestinationDistance.Value")
	if exif.GPS.DestinationDistance.Ref != meta.GPSRefKilometers {
		t.Fatalf("GPS.DestinationDistance.Ref = %v", exif.GPS.DestinationDistance.Ref)
	}
	assertApproxFloat64(t, exif.GPS.DestinationLatitude.Signed(), 33.2083333333, 0.000001, "Exif.GPS.DestinationLatitude")
	assertApproxFloat64(t, exif.GPS.DestinationLongitude.Signed(), 151.2083333333, 0.000001, "Exif.GPS.DestinationLongitude")
	assertApproxFloat64(t, exif.GPS.HPositioningError, 1.5, 0.0001, "Exif.GPS.HPositioningError")
	assertApproxFloat64(t, exif.GPS.ImageDirection.Value, 90, 0.0001, "Exif.GPS.ImageDirection.Value")
	if exif.GPS.ImageDirection.Ref != meta.GPSRefTrue {
		t.Fatalf("GPS.ImageDirection.Ref = %v", exif.GPS.ImageDirection.Ref)
	}
	if exif.GPS.ProcessingMethod != "GPS" {
		t.Fatalf("GPS.ProcessingMethod = %q", exif.GPS.ProcessingMethod)
	}
	assertApproxFloat64(t, exif.GPS.Speed.Value, 88, 0.0001, "Exif.GPS.Speed.Value")
	if exif.GPS.Speed.Ref != meta.GPSRefKilometers {
		t.Fatalf("GPS.Speed.Ref = %v", exif.GPS.Speed.Ref)
	}
	assertApproxFloat64(t, exif.GPS.Track.Value, 45, 0.0001, "Exif.GPS.Track.Value")
	if exif.GPS.Track.Ref != meta.GPSRefTrue {
		t.Fatalf("GPS.Track.Ref = %v", exif.GPS.Track.Ref)
	}
	if exif.ImageUniqueID != "img-id" || exif.MakerNote != "note" || exif.NativeDigest != "digest" {
		t.Fatalf("image id / makernote / digest mismatch: %+v", exif)
	}
	if exif.OECF != "oecf" || exif.OECFColumns != 3 || exif.OECFRows != 2 || exif.OECFNames != "RGB" || exif.OECFValues != "1 2 3 4 5 6" {
		t.Fatalf("oecf mismatch: %+v", exif)
	}
	if exif.RelatedSoundFile != "AUDIO.WAV" || exif.SensingMethod != 2 {
		t.Fatalf("related sound file / sensing method mismatch: %+v", exif)
	}
	if exif.SpatialFrequencyResponse != "sfr" || exif.SpatialFrequencyResponseColumns != 2 || exif.SpatialFrequencyResponseRows != 1 || exif.SpatialFrequencyResponseNames != "R G" || exif.SpatialFrequencyResponseValues != "0.2 0.3" {
		t.Fatalf("spatial frequency response mismatch: %+v", exif)
	}
	if exif.SpectralSensitivity != "sensitive" {
		t.Fatalf("SpectralSensitivity = %q", exif.SpectralSensitivity)
	}
	if exif.SubjectArea != "100,200,50,50" || exif.SubjectDistanceRange != 3 || exif.SubjectLocation != "100,200" {
		t.Fatalf("subject mismatch: %+v", exif)
	}
	if exif.FocalLengthIn35mmFilm != 50 {
		t.Fatalf("FocalLengthIn35mmFilm = %d", exif.FocalLengthIn35mmFilm)
	}
	if exif.Contrast != 1 {
		t.Fatalf("Contrast = %d", exif.Contrast)
	}
	if got, want := exif.GPS.Time.Format(time.RFC3339), "2025-01-02T03:04:05Z"; got != want {
		t.Fatalf("GPS.Time = %q, want %q", got, want)
	}
}

func TestParseExifNonSpecTagsIgnored(t *testing.T) {
	const src = `<x:xmpmeta xmlns:x="adobe:ns:meta/"><rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#"><rdf:Description xmlns:exif="http://ns.adobe.com/exif/1.0/" exif:DateTime="2025-01-02T03:04:05Z" exif:SubsecTime="25" exif:DateTimeDigitized="2025-01-02T03:04:05Z" exif:SubsecTimeDigitized="7" exif:DateTimeOriginal="2025-01-02T03:04:05Z" exif:SubsecTimeOriginal="61" exif:SamplesPerPixel="3" exif:PhotometricInterpretation="2"/></rdf:RDF></x:xmpmeta>`

	x, err := ParseXmp(strings.NewReader(src))
	if err != nil {
		t.Fatal(err)
	}

	exif := x.Exif
	if got, want := exif.CreateDate.Format(time.RFC3339), "2025-01-02T03:04:05Z"; got != want {
		t.Fatalf("CreateDate = %q, want %q", got, want)
	}
	if got, want := exif.DateTimeOriginal.Format(time.RFC3339), "2025-01-02T03:04:05Z"; got != want {
		t.Fatalf("DateTimeOriginal = %q, want %q", got, want)
	}
}

func TestParseExifEXAdditionalTags(t *testing.T) {
	const src = `<x:xmpmeta xmlns:x="adobe:ns:meta/"><rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#"><rdf:Description xmlns:exifEX="http://cipa.jp/exif/1.0/" exifEX:Acceleration="1/2" exifEX:AmbientTemperature="27.5" exifEX:CameraElevationAngle="-12.5" exifEX:CameraFirmware="1.2.3" exifEX:CompImageImagesPerSequence="3" exifEX:CompImageMaxExposureAll="9.5" exifEX:CompImageMaxExposureUsed="8.5" exifEX:CompImageMinExposureAll="0.5" exifEX:CompImageMinExposureUsed="1.5" exifEX:CompImageNumSequences="5" exifEX:CompImageSumExposureAll="100.25" exifEX:CompImageSumExposureUsed="95.25" exifEX:CompImageTotalExposurePeriod="3.0" exifEX:CompImageValues="1 2 3" exifEX:CompositeImage="1" exifEX:CompositeImageCount="2" exifEX:CompositeImageExposureTimes="1/30 1/60" exifEX:Gamma="2.2" exifEX:Humidity="65.5" exifEX:ISOSpeed="400" exifEX:ISOSpeedLatitudeyyy="1" exifEX:ISOSpeedLatitudezzz="2" exifEX:ImageEditingSoftware="EditSoft" exifEX:ImageEditor="EditorName" exifEX:ImageTitle="Title A" exifEX:InteropIndex="R98" exifEX:LensMake="Canon" exifEX:MetadataEditingSoftware="MetaSoft" exifEX:OwnerName="Owner A" exifEX:Photographer="Photographer A" exifEX:Pressure="1005.5" exifEX:RAWDevelopingSoftware="RawTool" exifEX:StandardOutputSensitivity="500" exifEX:WaterDepth="12.5"/></rdf:RDF></x:xmpmeta>`

	x, err := ParseXmp(strings.NewReader(src))
	if err != nil {
		t.Fatal(err)
	}

	exif := x.Exif
	assertApproxFloat64(t, exif.Acceleration, 0.5, 0.0001, "Exif.Acceleration")
	assertApproxFloat64(t, exif.AmbientTemperature, 27.5, 0.0001, "Exif.AmbientTemperature")
	assertApproxFloat64(t, exif.CameraElevationAngle, -12.5, 0.0001, "Exif.CameraElevationAngle")
	if exif.CameraFirmware != "1.2.3" {
		t.Fatalf("CameraFirmware = %q", exif.CameraFirmware)
	}
	if exif.CompImageImagesPerSequence != 3 || exif.CompImageNumSequences != 5 {
		t.Fatalf("comp image sequence counts mismatch: %+v", exif)
	}
	assertApproxFloat64(t, exif.CompImageMaxExposureAll, 9.5, 0.0001, "Exif.CompImageMaxExposureAll")
	assertApproxFloat64(t, exif.CompImageMaxExposureUsed, 8.5, 0.0001, "Exif.CompImageMaxExposureUsed")
	assertApproxFloat64(t, exif.CompImageMinExposureAll, 0.5, 0.0001, "Exif.CompImageMinExposureAll")
	assertApproxFloat64(t, exif.CompImageMinExposureUsed, 1.5, 0.0001, "Exif.CompImageMinExposureUsed")
	assertApproxFloat64(t, exif.CompImageSumExposureAll, 100.25, 0.0001, "Exif.CompImageSumExposureAll")
	assertApproxFloat64(t, exif.CompImageSumExposureUsed, 95.25, 0.0001, "Exif.CompImageSumExposureUsed")
	assertApproxFloat64(t, exif.CompImageTotalExposurePeriod, 3.0, 0.0001, "Exif.CompImageTotalExposurePeriod")
	if exif.CompImageValues != "1 2 3" {
		t.Fatalf("CompImageValues = %q", exif.CompImageValues)
	}
	if exif.CompositeImage != 1 || exif.CompositeImageCount != 2 || exif.CompositeImageExposureTimes != "1/30 1/60" {
		t.Fatalf("composite image mismatch: %+v", exif)
	}
	assertApproxFloat64(t, exif.Gamma, 2.2, 0.0001, "Exif.Gamma")
	assertApproxFloat64(t, exif.Humidity, 65.5, 0.0001, "Exif.Humidity")
	if exif.ISOSpeed != 400 || exif.ISOSpeedLatitudeyyy != 1 || exif.ISOSpeedLatitudezzz != 2 {
		t.Fatalf("iso speed exifex mismatch: %+v", exif)
	}
	if exif.ImageEditingSoftware != "EditSoft" || exif.ImageEditor != "EditorName" || exif.ImageTitle != "Title A" {
		t.Fatalf("image editor fields mismatch: %+v", exif)
	}
	if exif.InteroperabilityIndex != "R98" || exif.LensMake != "Canon" {
		t.Fatalf("interop/lens make mismatch: %+v", exif)
	}
	if exif.MetadataEditingSoftware != "MetaSoft" || exif.OwnerName != "Owner A" || exif.Photographer != "Photographer A" {
		t.Fatalf("metadata ownership fields mismatch: %+v", exif)
	}
	assertApproxFloat64(t, exif.Pressure, 1005.5, 0.0001, "Exif.Pressure")
	if exif.RAWDevelopingSoftware != "RawTool" || exif.StandardOutputSensitivity != 500 {
		t.Fatalf("raw developing fields mismatch: %+v", exif)
	}
	assertApproxFloat64(t, exif.WaterDepth, 12.5, 0.0001, "Exif.WaterDepth")
}
