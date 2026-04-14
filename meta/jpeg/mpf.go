package jpeg

import (
	"io"

	"github.com/evanoberholster/imagemeta/meta/utils"
)

// MPF stores Multi-Picture Format metadata from an APP2 MPF segment.
type MPF struct {
	MPFVersion     string
	NumberOfImages uint32
	Images         []MPImage
}

// MPImage stores one MPF MPImage entry.
type MPImage struct {
	MPImageFlags               uint8
	MPImageFormat              uint8
	MPImageType                uint32
	MPImageLength              uint32
	MPImageStart               uint32
	DependentImage1EntryNumber uint16
	DependentImage2EntryNumber uint16
}

func parseMPF(payload []byte, segmentOffset uint32) (*MPF, error) {
	if !isMPFPayload(payload) || len(payload) < 12 {
		return nil, errShortSegment("MPF")
	}
	tiff := payload[4:]
	order := utils.BinaryOrder(tiff)
	if order == utils.UnknownEndian || len(tiff) < 8 {
		return nil, io.ErrUnexpectedEOF
	}
	firstIFD := order.Uint32(tiff[4:8])
	if firstIFD > uint32(len(tiff)-2) {
		return nil, io.ErrUnexpectedEOF
	}
	pos := int(firstIFD)
	count := int(order.Uint16(tiff[pos : pos+2]))
	pos += 2
	if pos+count*12 > len(tiff) {
		return nil, io.ErrUnexpectedEOF
	}

	mpf := &MPF{}
	for i := 0; i < count; i++ {
		entry := tiff[pos+i*12:]
		tagID := order.Uint16(entry[0:2])
		typ := order.Uint16(entry[2:4])
		n := order.Uint32(entry[4:8])
		raw := entry[8:12]
		value, ok := tiffValue(tiff, order, typ, n, raw)
		if !ok {
			continue
		}
		switch tagID {
		case 0xb000:
			mpf.MPFVersion = trimNULString(value)
		case 0xb001:
			if len(value) >= 4 {
				mpf.NumberOfImages = order.Uint32(value)
			}
		case 0xb002:
			parseMPImages(mpf, order, value, segmentOffset)
		}
	}
	return mpf, nil
}

func parseMPImages(mpf *MPF, order utils.ByteOrder, value []byte, segmentOffset uint32) {
	n := len(value) / 16
	if mpf.NumberOfImages > 0 && int(mpf.NumberOfImages) < n {
		n = int(mpf.NumberOfImages)
	}
	mpf.Images = make([]MPImage, 0, n)
	for i := 0; i < n; i++ {
		b := value[i*16:]
		attr := order.Uint32(b[0:4])
		start := order.Uint32(b[8:12])
		if start != 0 {
			start += segmentOffset
		}
		mpf.Images = append(mpf.Images, MPImage{
			MPImageFlags:               uint8((attr & 0xf8000000) >> 27),
			MPImageFormat:              uint8((attr & 0x07000000) >> 24),
			MPImageType:                attr & 0x00ffffff,
			MPImageLength:              order.Uint32(b[4:8]),
			MPImageStart:               start,
			DependentImage1EntryNumber: order.Uint16(b[12:14]),
			DependentImage2EntryNumber: order.Uint16(b[14:16]),
		})
	}
}

func tiffValue(data []byte, order utils.ByteOrder, typ uint16, count uint32, raw []byte) ([]byte, bool) {
	size := tiffTypeSize(typ)
	if size == 0 || count == 0 {
		return nil, false
	}
	total := uint64(size) * uint64(count)
	if total <= 4 {
		return raw[:int(total)], true
	}
	offset := order.Uint32(raw)
	if uint64(offset)+total > uint64(len(data)) {
		return nil, false
	}
	start := int(offset)
	end := start + int(total)
	return data[start:end], true
}

func tiffTypeSize(typ uint16) uint32 {
	switch typ {
	case 1, 2, 6, 7:
		return 1
	case 3, 8:
		return 2
	case 4, 9, 11:
		return 4
	case 5, 10, 12:
		return 8
	default:
		return 0
	}
}
