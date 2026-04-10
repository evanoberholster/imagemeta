package exif

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/evanoberholster/imagemeta/meta/exif/tag"
)

func TestRationalUFloat64(t *testing.T) {
	t.Parallel()

	if got, want := (tag.RationalU{Numerator: 3, Denominator: 2}).Float64(), 1.5; got != want {
		t.Fatalf("tag.RationalU.Float64() = %v, want %v", got, want)
	}
	if got := (tag.RationalU{Numerator: 3, Denominator: 0}).Float64(); got != 0 {
		t.Fatalf("tag.RationalU.Float64() with zero denominator = %v, want 0", got)
	}
}

func TestLensInfoStringAndJSON(t *testing.T) {
	t.Parallel()

	info := &LensInfo{
		MinFocalLength:        tag.RationalU{Numerator: 24, Denominator: 1},
		MaxFocalLength:        tag.RationalU{Numerator: 70, Denominator: 1},
		MaxApertureAtMinFocal: tag.RationalU{Numerator: 28, Denominator: 10},
		MaxApertureAtMaxFocal: tag.RationalU{Numerator: 4, Denominator: 1},
	}

	if got, want := info.String(), "24 70 2.8 4"; got != want {
		t.Fatalf("LensInfo.String() = %q, want %q", got, want)
	}

	buf, err := json.Marshal(struct {
		LensInfo *LensInfo
	}{LensInfo: info})
	if err != nil {
		t.Fatalf("json.Marshal(LensInfo): %v", err)
	}
	if got, want := string(buf), `{"LensInfo":"24 70 2.8 4"}`; got != want {
		t.Fatalf("json.Marshal(LensInfo) = %s, want %s", got, want)
	}

	buf, err = json.Marshal(struct {
		LensInfo *LensInfo
	}{})
	if err != nil {
		t.Fatalf("json.Marshal(nil LensInfo): %v", err)
	}
	if got, want := string(buf), `{"LensInfo":null}`; got != want {
		t.Fatalf("json.Marshal(nil LensInfo) = %s, want %s", got, want)
	}
}

func TestExifTagParsedBitsetLowAndHigh(t *testing.T) {
	t.Parallel()

	var e Exif
	e.markTagParsed(10)
	e.markTagParsed(10) // duplicate
	e.markTagParsed(700)
	e.markTagParsed(700) // duplicate
	e.markTagParsed(701)

	if !e.HasTagParsed(10) {
		t.Fatal("HasTagParsed(10) should be true")
	}
	if !e.HasTagParsed(700) || !e.HasTagParsed(701) {
		t.Fatal("HasTagParsed() should be true for high IDs")
	}
	if e.HasTagParsed(11) {
		t.Fatal("HasTagParsed(11) should be false")
	}

	bits := e.TagParsedBitset()
	if bits[0]&(uint64(1)<<10) == 0 {
		t.Fatal("TagParsedBitset() missing low tag bit")
	}
	if e.highTagCount != 2 {
		t.Fatalf("highTagCount = %d, want 2", e.highTagCount)
	}
}

func TestGPSInfoSetDateAndTimeOrder(t *testing.T) {
	t.Parallel()

	base := time.Date(2024, time.March, 1, 0, 0, 0, 0, time.UTC)
	delta := 12*time.Hour + 34*time.Minute + 56*time.Second + 789*time.Millisecond
	want := base.Add(delta)

	var a GPSInfo
	a.setDate(base)
	a.setTime(delta)
	if got := a.GPSTime(); !got.Equal(want) {
		t.Fatalf("date->time GPSTime = %v, want %v", got, want)
	}

	var b GPSInfo
	b.setTime(delta)
	if pending, ok := gpsPendingDelta(b.GPSTime()); !ok || pending != delta {
		t.Fatalf("gpsPendingDelta() = (%v,%v), want (%v,true)", pending, ok, delta)
	}
	b.setDate(base)
	if got := b.GPSTime(); !got.Equal(want) {
		t.Fatalf("time->date GPSTime = %v, want %v", got, want)
	}
}

func TestGPSInfoAccessorsAndBitset(t *testing.T) {
	t.Parallel()

	g := GPSInfo{
		latitude:         33.9,
		longitude:        18.4,
		altitude:         123.4,
		latitudeRef:      tag.GPSRefSouth,
		longitudeRef:     tag.GPSRefEast,
		altitudeRef:      tag.GPSRefBelowSeaLevel,
		speedRef:         tag.GPSRefKilometersPerHour,
		speed:            tag.RationalU{Numerator: 90, Denominator: 1},
		trackRef:         tag.GPSRefTrueDirection,
		track:            tag.RationalU{Numerator: 123, Denominator: 1},
		destLatitude:     45.5,
		destLongitude:    9.1,
		destLatitudeRef:  tag.GPSRefSouth,
		destLongitudeRef: tag.GPSRefEast,
		destDistanceRef:  tag.GPSRefKilometers,
		destDistance:     tag.RationalU{Numerator: 12, Denominator: 1},
		mapDatum:         "WGS-84",
		differential:     1,
	}
	g.markTagParsed(tag.TagGPSLatitude)

	if got := g.Latitude(); got != -33.9 {
		t.Fatalf("Latitude() = %v, want -33.9", got)
	}
	if got := g.Longitude(); got != 18.4 {
		t.Fatalf("Longitude() = %v, want 18.4", got)
	}
	if got := g.Altitude(); got != -123.4 {
		t.Fatalf("Altitude() = %v, want -123.4", got)
	}
	if got := g.SpeedWithRef(); got.Ref != "K" || got.Value.Numerator != 90 || got.Value.Denominator != 1 {
		t.Fatalf("SpeedWithRef() = %+v", got)
	}
	if got := g.TrackWithRef(); got.Ref != "T" || got.Value.Numerator != 123 || got.Value.Denominator != 1 {
		t.Fatalf("TrackWithRef() = %+v", got)
	}
	if got := g.DestLatitude(); got != -45.5 {
		t.Fatalf("DestLatitude() = %v, want -45.5", got)
	}
	if got := g.DestLongitude(); got != 9.1 {
		t.Fatalf("DestLongitude() = %v, want 9.1", got)
	}
	if got := g.DestDistanceWithRef(); got.Ref != "K" || got.Value.Numerator != 12 || got.Value.Denominator != 1 {
		t.Fatalf("DestDistanceWithRef() = %+v", got)
	}
	if got := g.MapDatum(); got != "WGS-84" {
		t.Fatalf("MapDatum() = %q, want %q", got, "WGS-84")
	}
	if got := g.Differential(); got != 1 {
		t.Fatalf("Differential() = %d, want 1", got)
	}
	if !g.HasTagParsed(uint16(tag.TagGPSLatitude)) {
		t.Fatal("HasTagParsed(GPSLatitude) should be true")
	}
	if g.HasTagParsed(64) {
		t.Fatal("HasTagParsed(64) should be false for GPS bitset")
	}
}

func TestTimeTagsBitsetAndSelection(t *testing.T) {
	t.Parallel()

	base := time.Date(2024, time.January, 2, 3, 4, 5, 0, time.UTC)
	orig := time.Date(2024, time.January, 3, 4, 5, 6, 0, time.UTC)
	create := time.Date(2024, time.January, 4, 5, 6, 7, 0, time.UTC)

	var tt TimeTags
	tt.ModifyDate = base
	tt.DateTimeOriginal = orig
	tt.CreateDate = create
	tt.SubSecTime = 125
	tt.SubSecTimeOriginal = 50
	tt.SubSecTimeDigitized = 75
	tt.markTagParsed(tag.TagDateTime)
	tt.markTagParsed(tag.TagSubSecTime)
	tt.markTagParsed(tag.TagSubSecTimeOriginal)
	tt.markTagParsed(tag.TagSubSecTimeDigitized)

	if !tt.HasTagParsed(tag.TagDateTime) {
		t.Fatal("HasTagParsed(DateTime) should be true")
	}
	if !tt.HasSubSecTime() || !tt.HasSubSecTimeOriginal() || !tt.HasSubSecTimeDigitized() {
		t.Fatal("HasSubSec* methods should all be true")
	}
	if tt.HasTagParsed(tag.TagOffsetTime) {
		t.Fatal("HasTagParsed(OffsetTime) should be false")
	}

	if got, want := tt.GetModifyDate(), base.Add(125*time.Millisecond); !got.Equal(want) {
		t.Fatalf("GetModifyDate() = %v, want %v", got, want)
	}
	if got, want := tt.GetDateTimeOriginal(), orig.Add(50*time.Millisecond); !got.Equal(want) {
		t.Fatalf("GetDateTimeOriginal() = %v, want %v", got, want)
	}
	if got, want := tt.GetCreateDate(), create.Add(75*time.Millisecond); !got.Equal(want) {
		t.Fatalf("GetCreateDate() = %v, want %v", got, want)
	}

	if got, want := tt.GetSelectedDate(), orig.Add(50*time.Millisecond); !got.Equal(want) {
		t.Fatalf("GetSelectedDate() = %v, want %v", got, want)
	}
	tt.DateTimeOriginal = time.Time{}
	if got, want := tt.GetSelectedDate(), create.Add(75*time.Millisecond); !got.Equal(want) {
		t.Fatalf("GetSelectedDate() fallback create = %v, want %v", got, want)
	}
	tt.CreateDate = time.Time{}
	if got, want := tt.GetSelectedDate(), base.Add(125*time.Millisecond); !got.Equal(want) {
		t.Fatalf("GetSelectedDate() fallback modify = %v, want %v", got, want)
	}
}

func TestApplyTimeParts(t *testing.T) {
	t.Parallel()

	if got := applyTimeParts(time.Time{}, 123, nil); !got.IsZero() {
		t.Fatalf("applyTimeParts(zero) = %v, want zero", got)
	}

	base := time.Date(2024, time.January, 2, 3, 4, 5, 0, time.UTC)
	tz := time.FixedZone("+02:00", 2*60*60)
	got := applyTimeParts(base, 250, tz)
	want := base.Add(250*time.Millisecond - 2*time.Hour)
	if !got.Equal(want) {
		t.Fatalf("applyTimeParts() = %v, want %v", got, want)
	}
}

func TestTimeTagsMarshalJSONOffsetStrings(t *testing.T) {
	t.Parallel()

	tt := TimeTags{
		ModifyDate:         time.Date(2024, time.January, 2, 3, 4, 5, 0, time.UTC),
		OffsetTimeOriginal: time.FixedZone("+01:00", 60*60),
	}

	buf, err := json.Marshal(tt)
	if err != nil {
		t.Fatalf("json.Marshal(TimeTags): %v", err)
	}
	got := string(buf)
	if !strings.Contains(got, `"OffsetTimeOriginal":"+01:00"`) {
		t.Fatalf("OffsetTimeOriginal JSON = %s", got)
	}
	if !strings.Contains(got, `"OffsetTime":null`) {
		t.Fatalf("OffsetTime nil should marshal as null: %s", got)
	}
}
