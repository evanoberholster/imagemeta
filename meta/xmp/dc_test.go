package xmp

import (
	"strings"
	"testing"
	"time"

	"github.com/evanoberholster/imagemeta/imagetype"
)

func TestParseDublinCoreExtendedFields(t *testing.T) {
	const src = `<x:xmpmeta xmlns:x="adobe:ns:meta/">
<rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">
<rdf:Description xmlns:dc="http://purl.org/dc/elements/1.1/">
<dc:contributor><rdf:Bag><rdf:li>Alice</rdf:li><rdf:li>Bob</rdf:li></rdf:Bag></dc:contributor>
<dc:coverage>Worldwide</dc:coverage>
<dc:creator><rdf:Seq><rdf:li>Jane</rdf:li></rdf:Seq></dc:creator>
<dc:date><rdf:Seq><rdf:li>2024-01-01</rdf:li><rdf:li>2024-01-02T03:04:05Z</rdf:li></rdf:Seq></dc:date>
<dc:description><rdf:Alt><rdf:li xml:lang="x-default">Default description</rdf:li><rdf:li xml:lang="en-US">English description</rdf:li></rdf:Alt></dc:description>
<dc:format>image/jpeg</dc:format>
<dc:identifier>ID-123</dc:identifier>
<dc:language><rdf:Bag><rdf:li>en</rdf:li><rdf:li>fr</rdf:li></rdf:Bag></dc:language>
<dc:publisher><rdf:Bag><rdf:li>Publisher A</rdf:li></rdf:Bag></dc:publisher>
<dc:relation><rdf:Bag><rdf:li>related://123</rdf:li></rdf:Bag></dc:relation>
<dc:rights><rdf:Alt><rdf:li xml:lang="x-default">Copyright 2026</rdf:li></rdf:Alt></dc:rights>
<dc:source>Camera</dc:source>
<dc:subject><rdf:Bag><rdf:li>tag1</rdf:li><rdf:li>tag2</rdf:li></rdf:Bag></dc:subject>
<dc:title><rdf:Alt><rdf:li xml:lang="x-default">Default title</rdf:li><rdf:li xml:lang="fr-FR">Titre</rdf:li></rdf:Alt></dc:title>
<dc:type><rdf:Bag><rdf:li>Image</rdf:li></rdf:Bag></dc:type>
</rdf:Description>
</rdf:RDF>
</x:xmpmeta>`

	x, err := ParseXmp(strings.NewReader(src))
	if err != nil {
		t.Fatal(err)
	}

	if len(x.DC.Contributor) != 2 || x.DC.Contributor[0] != "Alice" || x.DC.Contributor[1] != "Bob" {
		t.Fatalf("DC.Contributor = %#v", x.DC.Contributor)
	}
	if x.DC.Coverage != "Worldwide" {
		t.Fatalf("DC.Coverage = %q", x.DC.Coverage)
	}
	if len(x.DC.Creator) != 1 || x.DC.Creator[0] != "Jane" {
		t.Fatalf("DC.Creator = %#v", x.DC.Creator)
	}
	if !x.DC.Date.Equal(time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC)) {
		t.Fatalf("DC.Date = %#v", x.DC.Date)
	}
	if len(x.DC.Description) != 1 || x.DC.Description[0] != "Default description" {
		t.Fatalf("DC.Description = %#v", x.DC.Description)
	}
	if len(x.DC.DescriptionLang) != 2 || x.DC.DescriptionLang[0] != "x-default" || x.DC.DescriptionLang[1] != "en-US" {
		t.Fatalf("DC.DescriptionLang = %#v", x.DC.DescriptionLang)
	}
	if x.DC.Format != imagetype.FromString("image/jpeg") {
		t.Fatalf("DC.Format = %q", x.DC.Format.String())
	}
	if x.DC.Identifier != "ID-123" {
		t.Fatalf("DC.Identifier = %q", x.DC.Identifier)
	}
	if len(x.DC.Language) != 2 || x.DC.Language[0] != "en" || x.DC.Language[1] != "fr" {
		t.Fatalf("DC.Language = %#v", x.DC.Language)
	}
	if len(x.DC.Publisher) != 1 || x.DC.Publisher[0] != "Publisher A" {
		t.Fatalf("DC.Publisher = %#v", x.DC.Publisher)
	}
	if len(x.DC.Relation) != 1 || x.DC.Relation[0] != "related://123" {
		t.Fatalf("DC.Relation = %#v", x.DC.Relation)
	}
	if len(x.DC.Rights) != 1 || x.DC.Rights[0] != "Copyright 2026" {
		t.Fatalf("DC.Rights = %#v", x.DC.Rights)
	}
	if len(x.DC.RightsLang) != 1 || x.DC.RightsLang[0] != "x-default" {
		t.Fatalf("DC.RightsLang = %#v", x.DC.RightsLang)
	}
	if x.DC.Source != "Camera" {
		t.Fatalf("DC.Source = %q", x.DC.Source)
	}
	if len(x.DC.Subject) != 2 || x.DC.Subject[0] != "tag1" || x.DC.Subject[1] != "tag2" {
		t.Fatalf("DC.Subject = %#v", x.DC.Subject)
	}
	if len(x.DC.Title) != 1 || x.DC.Title[0] != "Default title" {
		t.Fatalf("DC.Title = %#v", x.DC.Title)
	}
	if len(x.DC.TitleLang) != 2 || x.DC.TitleLang[0] != "x-default" || x.DC.TitleLang[1] != "fr-FR" {
		t.Fatalf("DC.TitleLang = %#v", x.DC.TitleLang)
	}
	if len(x.DC.Type) != 1 || x.DC.Type[0] != "Image" {
		t.Fatalf("DC.Type = %#v", x.DC.Type)
	}
}

func TestParseDublinCoreDateFormats(t *testing.T) {
	tests := []struct {
		in   string
		want time.Time
	}{
		{in: "2024", want: time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC)},
		{in: "2024-07", want: time.Date(2024, time.July, 1, 0, 0, 0, 0, time.UTC)},
		{in: "2024-07-31", want: time.Date(2024, time.July, 31, 0, 0, 0, 0, time.UTC)},
		{in: "2024-07-31T14:30", want: time.Date(2024, time.July, 31, 14, 30, 0, 0, time.UTC)},
		{in: "2024-07-31T14:30:45", want: time.Date(2024, time.July, 31, 14, 30, 45, 0, time.UTC)},
		{in: "2024-07-31T14:30:45.123", want: time.Date(2024, time.July, 31, 14, 30, 45, 123000000, time.UTC)},
		// Timezone-bearing inputs below are real values from meta/xmp/test fixtures.
		{in: "2007-08-16T11:57:04+01:00", want: time.Date(2007, time.August, 16, 11, 57, 4, 0, time.FixedZone("", 1*3600))},
		{in: "2012-10-17T13:07:01+03:00", want: time.Date(2012, time.October, 17, 13, 7, 1, 0, time.FixedZone("", 3*3600))},
		{in: "2021-02-03T17:34:04+08:00", want: time.Date(2021, time.February, 3, 17, 34, 4, 0, time.FixedZone("", 8*3600))},
		{in: "2024-11-02T12:35:44.40-04:00", want: time.Date(2024, time.November, 2, 12, 35, 44, 400000000, time.FixedZone("", -4*3600))},
		{in: "2026-02-27T16:53:33-08:00", want: time.Date(2026, time.February, 27, 16, 53, 33, 0, time.FixedZone("", -8*3600))},
		{in: "2021-01-10T09:30:34.576Z", want: time.Date(2021, time.January, 10, 9, 30, 34, 576000000, time.UTC)},
		{in: "  2003-02-04T08:06:56Z\t", want: time.Date(2003, time.February, 4, 8, 6, 56, 0, time.UTC)},
	}

	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			got, err := parseDate([]byte(tt.in))
			if err != nil {
				t.Fatalf("parseDate(%q) failed: %v", tt.in, err)
			}
			if !got.Equal(tt.want) {
				t.Fatalf("parseDate(%q) = %s, want %s", tt.in, got.Format(time.RFC3339Nano), tt.want.Format(time.RFC3339Nano))
			}
		})
	}
}

func TestParseDateRejectsNonCIPAFormats(t *testing.T) {
	tests := []string{
		"2024:07:31 14:30:45",       // Exif-style separators
		"2024-07-31 14:30:45",       // space instead of T
		"2024-07-31T14:30:45+0200",  // offset must be +hh:mm
		"2024-07-31T14:30:45+02",    // missing minute component
		"2024-07-31T14:30:45.123Z0", // trailing data
	}

	for _, in := range tests {
		t.Run(in, func(t *testing.T) {
			if _, err := parseDate([]byte(in)); err == nil {
				t.Fatalf("parseDate(%q) unexpectedly succeeded", in)
			}
		})
	}
}

func TestParseDublinCoreDateFallback(t *testing.T) {
	const src = `<x:xmpmeta xmlns:x="adobe:ns:meta/">
<rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">
<rdf:Description xmlns:dc="http://purl.org/dc/elements/1.1/">
<dc:date><rdf:Seq><rdf:li>not-a-date</rdf:li></rdf:Seq></dc:date>
</rdf:Description>
</rdf:RDF>
</x:xmpmeta>`

	x, err := ParseXmp(strings.NewReader(src))
	if err != nil {
		t.Fatal(err)
	}

	if !x.DC.Date.IsZero() {
		t.Fatalf("DC.Date = %#v, want zero value", x.DC.Date)
	}
}
