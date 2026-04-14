package exif

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/evanoberholster/imagemeta/meta/exif/tag"
	"github.com/evanoberholster/imagemeta/meta/utils"
	"github.com/rs/zerolog"
)

func TestParseIFD0ImageReferenceFromStripArrays(t *testing.T) {
	t.Parallel()

	var payload [16]byte
	utils.LittleEndian.PutUint32(payload[0:4], 0x2000)
	utils.LittleEndian.PutUint32(payload[4:8], 0x3000)
	utils.LittleEndian.PutUint32(payload[8:12], 4096)
	utils.LittleEndian.PutUint32(payload[12:16], 8192)

	r := NewReader(zerolog.Nop())
	defer r.Close()
	r.Reset(bufio.NewReaderSize(bytes.NewReader(payload[:]), len(payload)))

	offsetTag := tag.NewEntry(tag.TagStripOffsets, tag.TypeLong, 2, 0, tag.IFD0, 0, utils.LittleEndian)
	if !r.parseIFD0ImageTag(offsetTag) {
		t.Fatal("parseIFD0ImageTag(TagStripOffsets) = false, want true")
	}
	lengthTag := tag.NewEntry(tag.TagStripByteCounts, tag.TypeLong, 2, 8, tag.IFD0, 0, utils.LittleEndian)
	if !r.parseIFD0ImageTag(lengthTag) {
		t.Fatal("parseIFD0ImageTag(TagStripByteCounts) = false, want true")
	}

	if got, want := r.Exif.IFD0.ImageOffset, uint32(0x2000); got != want {
		t.Fatalf("IFD0.ImageOffset = %#x, want %#x", got, want)
	}
	if got, want := r.Exif.IFD0.ImageLength, uint32(4096); got != want {
		t.Fatalf("IFD0.ImageLength = %d, want %d", got, want)
	}
}

func TestParseImageIFDReferenceFromStripArrays(t *testing.T) {
	t.Parallel()

	var payload [16]byte
	utils.LittleEndian.PutUint32(payload[0:4], 0x4000)
	utils.LittleEndian.PutUint32(payload[4:8], 0x5000)
	utils.LittleEndian.PutUint32(payload[8:12], 1024)
	utils.LittleEndian.PutUint32(payload[12:16], 2048)

	r := NewReader(zerolog.Nop())
	defer r.Close()
	r.Reset(bufio.NewReaderSize(bytes.NewReader(payload[:]), len(payload)))

	var dst ImageIFD
	offsetTag := tag.NewEntry(tag.TagStripOffsets, tag.TypeLong, 2, 0, tag.IFD1, 0, utils.LittleEndian)
	if !r.parseImageIFDTag(offsetTag, &dst) {
		t.Fatal("parseImageIFDTag(TagStripOffsets) = false, want true")
	}
	lengthTag := tag.NewEntry(tag.TagStripByteCounts, tag.TypeLong, 2, 8, tag.IFD1, 0, utils.LittleEndian)
	if !r.parseImageIFDTag(lengthTag, &dst) {
		t.Fatal("parseImageIFDTag(TagStripByteCounts) = false, want true")
	}

	if got, want := dst.ImageOffset, uint32(0x4000); got != want {
		t.Fatalf("ImageIFD.ImageOffset = %#x, want %#x", got, want)
	}
	if got, want := dst.ImageLength, uint32(1024); got != want {
		t.Fatalf("ImageIFD.ImageLength = %d, want %d", got, want)
	}
}
