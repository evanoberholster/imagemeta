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
		{in: "2024-07-31 14:30:45", want: time.Date(2024, time.July, 31, 14, 30, 45, 0, time.UTC)},
		{in: "2024-07-31T14:30:45.123", want: time.Date(2024, time.July, 31, 14, 30, 45, 123000000, time.UTC)},
		// Exif-style date without explicit offset is interpreted in local time.
		{in: "2024:07:31 14:30:45", want: time.Date(2024, time.July, 31, 14, 30, 45, 0, time.Local)},
		{in: "2024-07-31T14:30:45+02:00", want: time.Date(2024, time.July, 31, 14, 30, 45, 0, time.FixedZone("", 2*3600))},
		{in: "2024-07-31T14:30:45+0200", want: time.Date(2024, time.July, 31, 14, 30, 45, 0, time.FixedZone("", 2*3600))},
		{in: "2024-07-31 14:30:45-05:00", want: time.Date(2024, time.July, 31, 14, 30, 45, 0, time.FixedZone("", -5*3600))},
		{in: "2024-07-31T14:30:45.123456789Z", want: time.Date(2024, time.July, 31, 14, 30, 45, 123456789, time.UTC)},
		{in: "  2024-07-31T14:30:45Z\t", want: time.Date(2024, time.July, 31, 14, 30, 45, 0, time.UTC)},
		{in: "\n2024:07:31 14:30:45  ", want: time.Date(2024, time.July, 31, 14, 30, 45, 0, time.Local)},
	}

	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			got, ok := parseDateDublinCore([]byte(tt.in))
			if !ok {
				t.Fatalf("parseDateDublinCore(%q) failed", tt.in)
			}
			if !got.Equal(tt.want) {
				t.Fatalf("parseDateDublinCore(%q) = %s, want %s", tt.in, got.Format(time.RFC3339Nano), tt.want.Format(time.RFC3339Nano))
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
