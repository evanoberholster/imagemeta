package bmff

import (
	"encoding/binary"
	"fmt"
	"strings"

	"github.com/evanoberholster/imagemeta/exif"
	"github.com/evanoberholster/imagemeta/imagetype"
	"github.com/evanoberholster/imagemeta/meta"
	"github.com/evanoberholster/imagemeta/tiff"
	"github.com/pkg/errors"
)

// Values are in BigEndian
var crxBinaryOrder = binary.BigEndian

// CR3MetaBoxUUID is the uuid that corresponds with Canon CR3 Metadata.
//
// uuid = 85c0b687 820f 11e0 8111 f4ce462b6a48
var CR3MetaBoxUUID = meta.UUID{133, 192, 182, 135, 130, 15, 17, 224, 129, 17, 244, 206, 70, 43, 106, 72}

// ParseCanonCR3 is a performance focused method for parsing metadata from a .CR3 file.
func (r *Reader) ParseCanonCR3() (err error) {
	// Parse Moov Box
	moovBox, err := r.readBox()
	if err != nil {
		err = errors.Wrapf(err, "Box 'moov' (readBox) error")
		return
	}
	if moovBox.boxType != TypeMoov {
		err = errors.Wrapf(ErrWrongBoxType, "Box %s", moovBox.boxType)
		return
	}
	var inner box
	for moovBox.anyRemain() {
		if inner, err = moovBox.readInnerBox(); err != nil {
			err = errors.Wrapf(err, "Box 'moov' %s (readBox)", inner.boxType)
			return
		}
		switch inner.boxType {
		case TypeUUID:
			uuid, err := inner.readUUID()
			if err != nil {
				return err
			}
			switch uuid {
			case CR3MetaBoxUUID:
				cr3, err := parseCR3MetaBox(&inner)
				fmt.Println(cr3, err)
			}
		default:
			//fmt.Println(inner)
		}
		if err = moovBox.closeInnerBox(&inner); err != nil {
			return
		}
	}
	return moovBox.discard(moovBox.remain)
}

// CR3MetaBox is a uuidBox that contains Metadata for CR3 files
type CR3MetaBox struct {
	CNCV CNCVBox
	CCTP CCTPBox
	CTBO CTBOBox
	CMT  [4]exif.Header
	THMB THMBBox
}

// Type returns TypeUUID, CR3MetaBox's boxType.
func (cr3 CR3MetaBox) Type() BoxType {
	return TypeUUID
}

// parseCR3MetaBox parses a uuid box with the uuid of 85c0b687 820f 11e0 8111 f4ce462b6a48
func parseCR3MetaBox(outer *box) (meta CR3MetaBox, err error) {
	var inner box
	cmt := 0
	for outer.anyRemain() {
		if inner, err = outer.readInnerBox(); err != nil {
			return
		}
		switch inner.boxType {
		case TypeCNCV:
			meta.CNCV, err = inner.parseCNCVBox()
		case TypeCCTP:
			meta.CCTP, err = inner.parseCCTPBox()
		case TypeCTBO:
			meta.CTBO, err = inner.parseCTBOBox()
		case TypeTMHB:
			meta.THMB, err = inner.parseTHMBBox()
		case TypeCMT1, TypeCMT2, TypeCMT3, TypeCMT4:
			meta.CMT[cmt], err = inner.parseExifHeader(imagetype.ImageCR3)
			cmt++
		default:
			fmt.Println(inner)
		}
		if err != nil {
			return
		}

		if err = outer.closeInnerBox(&inner); err != nil {
			return
		}
	}
	return
}

func (b *box) parseExifHeader(it imagetype.ImageType) (header exif.Header, err error) {
	buf, err := b.peek(8)
	if err != nil {
		return exif.Header{}, errors.Wrap(err, "parseExifHeader")
	}
	byteOrder := tiff.BinaryOrder(buf[:4])
	firstIfdOffset := byteOrder.Uint32(buf[4:8])
	tiffHeaderOffset := uint32(b.offset)
	exifLength := uint32(b.remain)
	return exif.NewHeader(byteOrder, firstIfdOffset, tiffHeaderOffset, exifLength, it), b.discard(8)
}

// THMBBox is a Canon CR3 Thumbnail Box
type THMBBox struct {
	offset        uint32
	size          uint32
	width, height uint16
}

func (b *box) parseTHMBBox() (thmb THMBBox, err error) {
	flags, err := b.readFlags()
	if err != nil {
		return thmb, errors.Wrap(err, "parseTHMBBox")
	}
	buf, err := b.peek(8)
	if err != nil {
		return thmb, errors.Wrap(err, "parseTHMBBox")
	}
	// read width, height, jpeg image size
	thmb.offset = uint32(b.offset + 8)
	thmb.width = crxBinaryOrder.Uint16(buf[:2])
	thmb.height = crxBinaryOrder.Uint16(buf[2:4])
	thmb.size = crxBinaryOrder.Uint32(buf[4:8])
	switch flags.Version() {
	case 0:
		thmb.offset += 4
	case 1:
	}

	return thmb, b.discard(b.remain)
}

// CCDTBox is a Canon CR3 definition of tracks?
type CCDTBox struct {
	//size uint32
	it  uint64
	idx uint32
}

func (b *box) parseCCDTBox() (ccdt CCDTBox, err error) {
	if b.boxType != TypeCCDT {
		err = ErrWrongBoxType
		return
	}
	buf, err := b.peek(16)
	if err != nil {
		return
	}
	if err = b.discard(16); err != nil {
		return
	}
	// uint64 value appears to be imagetype
	ccdt.it = crxBinaryOrder.Uint64(buf[0:8])
	// uint32 value apprears to be 0 or 1 for dual pixel

	// uint32 value for the trak Index
	ccdt.idx = crxBinaryOrder.Uint32(buf[12:16])

	return ccdt, b.discard(b.remain)
}

// CNCVBox is Canon Compressor Version box
// CaNon Codec Version?
type CNCVBox struct {
	//size uint32
	//format [9]byte
	//version [6]uint8
	val [30]byte
}

func (cncv CNCVBox) String() string {
	var sb strings.Builder
	sb.WriteString("CNCV | Format: ")
	sb.Write(cncv.val[0:9])
	sb.WriteString(", Version: ")
	sb.Write(cncv.val[9:30])
	return sb.String()
}

func (b *box) parseCNCVBox() (cncv CNCVBox, err error) {
	buf, err := b.peek(30)
	if err != nil {
		return
	}
	copy(cncv.val[:], buf[0:30])
	return cncv, b.discard(30)
}

// CCTPBox is Canon Compressor Table Pointers box
// Canon CR3 trak pointers?
type CCTPBox struct {
	//size uint32
	CCDT []CCDTBox
}

func (b *box) parseCCTPBox() (cctp CCTPBox, err error) {
	if b.boxType != TypeCCTP {
		err = ErrWrongBoxType
		return
	}
	buf, err := b.peek(12)
	if err != nil {
		return
	}
	if err = b.discard(12); err != nil {
		return
	}
	// CCTP Box contains 12 bytes (3 x uint32)
	// last one is number of CCDT lines. 3, or 4 for dual pixel
	count := crxBinaryOrder.Uint32(buf[8:12])
	cctp.CCDT = make([]CCDTBox, count)

	var inner box
	for i := 0; i < int(count) && b.anyRemain(); i++ {
		if inner, err = b.readInnerBox(); err != nil {
			return
		}
		if inner.boxType == TypeCCDT {
			cctp.CCDT[i], err = inner.parseCCDTBox()
		}
		if err = b.closeInnerBox(&inner); err != nil {
			return
		}
	}
	return cctp, b.discard(b.remain)
}

// CTBOBox is a Canon tracks base offsets Box?
type CTBOBox struct {
	//size uint32
	items []IndexOffset
}

// IndexOffset has an index, an offset and a length.
type IndexOffset struct {
	offset uint64
	length uint64
	idx    uint32
}

func (b *box) parseCTBOBox() (ctbo CTBOBox, err error) {
	if b.Type() != TypeCTBO {
		err = ErrWrongBoxType
		return
	}
	buf, err := b.peek(4)
	if err != nil {
		return
	}
	if err = b.discard(4); err != nil {
		return
	}
	count := crxBinaryOrder.Uint32(buf[0:4])
	ctbo.items = make([]IndexOffset, count)

	for i := 0; i < int(count) && b.anyRemain(); i++ {
		// each item is 20 bytes in length
		buf, err = b.peek(20)
		if err != nil {
			return
		}
		if err = b.discard(20); err != nil {
			return
		}
		ctbo.items[i] = IndexOffset{
			idx:    crxBinaryOrder.Uint32(buf[0:4]),
			offset: crxBinaryOrder.Uint64(buf[4:12]),
			length: crxBinaryOrder.Uint64(buf[12:20]),
		}
	}
	return ctbo, b.discard(b.remain)
}
