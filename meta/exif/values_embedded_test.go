package exif

import (
	"testing"

	"github.com/evanoberholster/imagemeta/meta/exif/tag"
	metalog "github.com/evanoberholster/imagemeta/meta/logging"
	"github.com/evanoberholster/imagemeta/meta/utils"
)

func TestParseFirstUint32FromEmbeddedLong(t *testing.T) {
	t.Parallel()

	r := NewReader(metalog.GetLogger())
	defer r.Close()

	// a single LONG is embedded in the entry's value field; JPEGInterchangeFormat
	// offsets and lengths arrive exactly like this
	entry := tag.NewEntry(tag.TagThumbnailOffset, tag.TypeLong, 1, 166912, tag.SubIFD0, 0, utils.LittleEndian)

	var dst ImageIFD
	if !r.parseImageIFDTag(entry, &dst) {
		t.Fatal("parseImageIFDTag(TagThumbnailOffset) = false, want true")
	}
	if got, want := dst.ImageOffset, uint32(166912); got != want {
		t.Fatalf("ImageIFD.ImageOffset = %d, want %d", got, want)
	}
}
