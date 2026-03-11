package exif

import (
	"testing"

	"github.com/evanoberholster/imagemeta/meta/exif/ifd"
	"github.com/evanoberholster/imagemeta/meta/exif/tag"
	"github.com/evanoberholster/imagemeta/meta/utils"
)

var tagDecodeSink uint32

// tagFromBufferFast is a candidate decode path with manual endian loads and inlined type remap.
func tagFromBufferFast(directory ifd.Directory, buf []byte) (tag.Entry, error) {
	var id tag.ID
	var typ tag.Type
	var unitCount uint32
	var valueOffset uint32

	if directory.ByteOrder == utils.BigEndian {
		id = tag.ID(uint16(buf[0])<<8 | uint16(buf[1]))
		typ = tag.Type(uint16(buf[2])<<8 | uint16(buf[3]))
		unitCount = uint32(buf[4])<<24 | uint32(buf[5])<<16 | uint32(buf[6])<<8 | uint32(buf[7])
		valueOffset = uint32(buf[8])<<24 | uint32(buf[9])<<16 | uint32(buf[10])<<8 | uint32(buf[11])
	} else {
		id = tag.ID(uint16(buf[0]) | uint16(buf[1])<<8)
		typ = tag.Type(uint16(buf[2]) | uint16(buf[3])<<8)
		unitCount = uint32(buf[4]) | uint32(buf[5])<<8 | uint32(buf[6])<<16 | uint32(buf[7])<<24
		valueOffset = uint32(buf[8]) | uint32(buf[9])<<8 | uint32(buf[10])<<16 | uint32(buf[11])<<24
	}
	valueOffset += directory.BaseOffset

	if typ == tag.TypeLong || typ == tag.TypeUndefined {
		switch directory.Type {
		case ifd.IFD0:
			if id == tag.TagExifIFDPointer || id == tag.TagGPSIFDPointer {
				typ = tag.TypeIfd
			}
		case ifd.ExifIFD:
			if id == tag.TagMakerNote {
				typ = tag.TypeIfd
			}
		}
	}

	entry := tag.Entry{
		ValueOffset: valueOffset,
		UnitCount:   unitCount,
		ID:          id,
		Type:        typ,
		IfdType:     directory.Type,
		IfdIndex:    directory.Index,
		ByteOrder:   directory.ByteOrder,
	}
	if !entry.IsValid() {
		return entry, tag.ErrTagTypeNotValid
	}
	return entry, nil
}

func benchmarkTagDecode(b *testing.B, fn func(ifd.Directory, []byte) (tag.Entry, error), directory ifd.Directory, raw [12]byte) {
	var sum uint32
	b.ReportAllocs()
	b.SetBytes(12)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		e, err := fn(directory, raw[:])
		if err != nil {
			b.Fatal(err)
		}
		sum += uint32(e.ID) + uint32(e.Type) + e.UnitCount + e.ValueOffset
	}
	tagDecodeSink = sum
}

func BenchmarkTagFromBufferDecode(b *testing.B) {
	var leRaw [12]byte
	utils.LittleEndian.PutUint16(leRaw[0:2], uint16(tag.TagImageWidth))
	utils.LittleEndian.PutUint16(leRaw[2:4], uint16(tag.TypeLong))
	utils.LittleEndian.PutUint32(leRaw[4:8], 1)
	utils.LittleEndian.PutUint32(leRaw[8:12], 4000)
	leDir := ifd.New(utils.LittleEndian, ifd.IFD0, 0, 0, 0)

	var beRaw [12]byte
	utils.BigEndian.PutUint16(beRaw[0:2], uint16(tag.TagImageWidth))
	utils.BigEndian.PutUint16(beRaw[2:4], uint16(tag.TypeLong))
	utils.BigEndian.PutUint32(beRaw[4:8], 1)
	utils.BigEndian.PutUint32(beRaw[8:12], 4000)
	beDir := ifd.New(utils.BigEndian, ifd.IFD0, 0, 0, 0)

	b.Run("CurrentLE", func(b *testing.B) {
		benchmarkTagDecode(b, tagFromBuffer, leDir, leRaw)
	})
	b.Run("FastLE", func(b *testing.B) {
		benchmarkTagDecode(b, tagFromBufferFast, leDir, leRaw)
	})
	b.Run("CurrentBE", func(b *testing.B) {
		benchmarkTagDecode(b, tagFromBuffer, beDir, beRaw)
	})
	b.Run("FastBE", func(b *testing.B) {
		benchmarkTagDecode(b, tagFromBufferFast, beDir, beRaw)
	})
}
