package bmff

import (
	"encoding/binary"
	"fmt"
	"strings"

	"github.com/evanoberholster/imagemeta/imagetype"
	"github.com/evanoberholster/imagemeta/meta"
	"github.com/pkg/errors"
)

var (
	// crx values are in BigEndian.
	crxBinaryOrder = binary.BigEndian

	// CR3MetaBoxUUID is the uuid that corresponds with Canon CR3 Metadata.
	//
	// uuid = 85c0b687 820f 11e0 8111 f4ce462b6a48
	CR3MetaBoxUUID = meta.UUID{133, 192, 182, 135, 130, 15, 17, 224, 129, 17, 244, 206, 70, 43, 106, 72}

	// CR3XPacketUUID is the uuid that corresponds with Canon CR3 xpacket data
	//
	// uuid = be7acfcb-97a9-42e8-9c71-999491e3afac
	CR3XPacketUUID = meta.UUID{190, 122, 207, 203, 151, 169, 66, 232, 156, 113, 153, 148, 145, 227, 175, 172}
)

// CrxMoovBox is a Canon Raw Moov Box
type CrxMoovBox struct {
	Meta CR3MetaBox
}

// ReadCrxMoovBox is a performance focused method for parsing the moov box from a .CR3 file.
func (r *Reader) ReadCrxMoovBox() (cmb CrxMoovBox, err error) {
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
	if debugFlag {
		traceBox(moovBox, moovBox)
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
				return cmb, err
			}
			if debugFlag {
				traceBox(UUIDBox{uuid: uuid}, inner)
			}
			switch uuid {
			case CR3MetaBoxUUID:
				cmb.Meta, err = parseCR3MetaBox(&inner)
				if err != nil {
					return cmb, err
				}
			}
		default:
			if debugFlag {
				traceBox(inner, inner)
			}
		}
		if err = moovBox.closeInnerBox(&inner); err != nil {
			return
		}
	}
	err = moovBox.discard(moovBox.remain)
	r.br.offset = moovBox.offset
	return
}

// ReadXPacketUUIDBox -
func (r *Reader) ReadXPacketUUIDBox() (err error) {
	// Parse UUID Box
	outer, err := r.readBox()
	if err != nil {
		err = errors.Wrapf(err, "Box 'uuid' (readBox) error")
		return
	}
	fmt.Println(outer.size, outer.remain, outer.offset)
	switch outer.boxType {
	case TypeUUID:
		uuid, err := outer.readUUID()
		if err != nil {
			return err
		}
		switch uuid {
		case CR3XPacketUUID:
			//cr3, err := parseCR3MetaBox(&inner)
			//fmt.Println(cr3, err)
		}
		//fmt.Println(uuid, uuid.Bytes())
	default:
		//fmt.Println(inner)
	}

	return outer.discard(outer.remain)
}

// CR3MetaBox is a uuidBox that contains Metadata for CR3 files
type CR3MetaBox struct {
	CNCV CNCVBox
	CCTP CCTPBox
	ctbo CTBOBox
	CMT  [4]CMTBox
	THMB THMBBox
}

// Type returns TypeUUID, CR3MetaBox's boxType.
func (cr3 CR3MetaBox) Type() BoxType {
	return TypeUUID
}

// XPacketData returns CTBO[0] which corresponds to XPacket data
func (cr3 CR3MetaBox) XPacketData() (offset, length uint64, err error) {
	// First 24 bytes are a UUID box. uuid = be7acfcb-97a9-42e8-9c71-999491e3afac
	offset = cr3.ctbo.items[0].offset
	length = cr3.ctbo.items[0].length
	return offset, length, nil
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
			meta.ctbo, err = inner.parseCTBOBox()
		case TypeTHMB:
			meta.THMB, err = inner.parseTHMBBox()
		case TypeCMT1, TypeCMT2, TypeCMT3, TypeCMT4:
			// TODO: Add Ifd Types
			meta.CMT[cmt], err = inner.parseCMTBox(imagetype.ImageCR3)
			cmt++
		default:

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

// CMTBox is a CMT# box
type CMTBox struct {
	ByteOrder        binary.ByteOrder
	FirstIfdOffset   uint32
	TiffHeaderOffset uint32
	ExifLength       uint32
	ImageType        imagetype.ImageType
	Bt               BoxType
}

func (cmt CMTBox) String() string {
	return fmt.Sprintf("%s | ByteOrder:%s, Offset:%d, Size:%d, TiffHeaderOffset:%d, ImageType:%s\n", cmt.Bt, cmt.ByteOrder, cmt.FirstIfdOffset, cmt.ExifLength, cmt.TiffHeaderOffset, cmt.ImageType)
}

// Type returns TypeTHMB
func (cmt CMTBox) Type() BoxType {
	return cmt.Bt
}

func (b *box) parseCMTBox(it imagetype.ImageType) (cmt CMTBox, err error) {
	buf, err := b.peek(8)
	if err != nil {
		err = errors.Wrap(err, "parseCMTBox")
		return
	}
	binaryOrder := meta.BinaryOrder(buf[:4])
	cmt = CMTBox{
		ByteOrder:        binaryOrder,
		FirstIfdOffset:   binaryOrder.Uint32(buf[4:8]),
		TiffHeaderOffset: uint32(b.offset),
		ExifLength:       uint32(b.remain),
		ImageType:        it,
		Bt:               b.boxType,
	}
	if debugFlag {
		traceBox(cmt, *b)
	}
	return cmt, b.discard(8)
}

// THMBBox is a Canon CR3 Thumbnail Box
type THMBBox struct {
	offset        uint32
	size          uint32
	width, height uint16
}

func (thmb THMBBox) String() string {
	return fmt.Sprintf("THMB | Offset: %d, Size: %d (%dx%d)", thmb.offset, thmb.size, thmb.width, thmb.height)
}

// Type returns TypeTHMB
func (thmb THMBBox) Type() BoxType {
	return TypeTHMB
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
	if debugFlag {
		traceBox(thmb, *b)
	}
	return thmb, b.discard(b.remain)
}

// CCDTBox is a Canon CR3 definition of tracks?
type CCDTBox struct {
	//size uint32
	it  uint64
	idx uint32
}

func (ccdt CCDTBox) String() string {
	return fmt.Sprintf("CCDT | Index:%d, Item:%d ", ccdt.idx, ccdt.it)
}

// Type returns TypeCCDT
func (ccdt CCDTBox) Type() BoxType {
	return TypeCCDT
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
	if debugFlag {
		traceBox(ccdt, *b)
	}
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

// Type returns TypeCNCV
func (cncv CNCVBox) Type() BoxType {
	return TypeCNCV
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
	if debugFlag {
		traceBox(cncv, *b)
	}
	return cncv, b.discard(30)
}

// CCTPBox is Canon Compressor Table Pointers box
// Canon CR3 trak pointers?
type CCTPBox struct {
	//size uint32
	CCDT []CCDTBox
}

func (cctp CCTPBox) String() string {
	return fmt.Sprintf("CCTP | ItemCount:%d \t", len(cctp.CCDT))
}

// Type returns TypeCCTP
func (cctp CCTPBox) Type() BoxType {
	return TypeCCTP
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
	if debugFlag {
		traceBox(cctp, *b)
	}
	var inner box
	for i := 0; i < int(count) && b.anyRemain(); i++ {
		if inner, err = b.readInnerBox(); err != nil {
			break
		}
		if inner.boxType == TypeCCDT {
			cctp.CCDT[i], err = inner.parseCCDTBox()
			if err != nil {
				break
			}
		}
		if err = b.closeInnerBox(&inner); err != nil {
			break
		}
	}
	return cctp, b.discard(b.remain)
}

// CTBOBox is a Canon tracks base offsets Box?
type CTBOBox struct {
	items [5]IndexOffset
	count uint32
}

func (ctbo CTBOBox) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("CTBO | ItemCount:%d \n", ctbo.count))
	for _, item := range ctbo.items {
		sb.WriteString(fmt.Sprintf("\t | Index:%d, Offset:%d, Size:%d \n", item.idx, item.offset, item.length))
	}
	return sb.String()
}

// Type returns TypeCTBO
func (ctbo CTBOBox) Type() BoxType {
	return TypeCTBO
}

// IndexOffset has an index, an offset and a length.
type IndexOffset struct {
	offset uint64
	length uint64
	idx    uint32
}

func (b *box) parseCTBOBox() (ctbo CTBOBox, err error) {
	buf, err := b.peek(4)
	if err != nil {
		return
	}
	if err = b.discard(4); err != nil {
		return
	}
	ctbo.count = crxBinaryOrder.Uint32(buf[0:4])
	//ctbo.items = make([]IndexOffset, count)
	for i := 0; i < int(ctbo.count) && b.anyRemain(); i++ {
		// each item is 20 bytes in length
		buf, err = b.peek(20)
		if err != nil {
			return
		}
		if err = b.discard(20); err != nil {
			return
		}
		if i < len(ctbo.items) {
			ctbo.items[i] = IndexOffset{
				idx:    crxBinaryOrder.Uint32(buf[0:4]),
				offset: crxBinaryOrder.Uint64(buf[4:12]),
				length: crxBinaryOrder.Uint64(buf[12:20]),
			}
		}
	}
	if debugFlag {
		traceBox(ctbo, *b)
	}
	return ctbo, b.discard(b.remain)
}
