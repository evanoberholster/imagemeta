package bmff

import (
	"encoding/binary"
	"fmt"
	"strings"

	"github.com/evanoberholster/imagemeta/exif/ifds"
	"github.com/evanoberholster/imagemeta/imagetype"
	"github.com/evanoberholster/imagemeta/meta"
	"github.com/pkg/errors"
)

const (
	stsdHeaderSize = 16
)

var (
	// crx values are in BigEndian.
	crxBinaryOrder = binary.BigEndian

	// CR3MetaBoxUUID is the uuid that corresponds with Canon CR3 Metadata.
	CR3MetaBoxUUID = meta.UUIDFromString("85c0b687-820f-11e0-8111-f4ce462b6a48")

	// CR3XPacketUUID is the uuid that corresponds with Canon CR3 xpacket data
	CR3XPacketUUID = meta.UUIDFromString("be7acfcb-97a9-42e8-9c71-999491e3afac")

	// CR3PreviewDataUUID is the uuid that corresponds with Canon CR3 Perview Data
	CR3PreviewDataUUID = meta.UUIDFromString("eaf42b5e-1c98-4b88-b9fb-b7dc406e4d16")

	// CR3CMTAUUID is the uuid that corresponds with Canon CR3 CMTA box
	CR3CMTAUUID = meta.UUIDFromString("5766b829-bb6a-47c5-bcfb-8b9f2260d06d")

	// CR3CNOPUUID is the uuid that corresponds with Canon CR3 CNOP "Optional data"
	CR3CNOPUUID = meta.UUIDFromString("210f1687-9149-11e4-8111-00242131fce4")
)

// CrxMoovBox is a Canon Raw Moov Box
type CrxMoovBox struct {
	Meta CR3MetaBox
	Trak [4]CR3Trak
}

// ReadCrxMoovBox is a performance focused method for parsing the moov box from a .CR3 file.
func (r *Reader) ReadCrxMoovBox() (cmb CrxMoovBox, err error) {
	// Parse Moov Box
	moovBox, err := r.readBox()
	if err != nil {
		return cmb, errors.Wrapf(err, "Box 'moov' (readBox) error")
	}
	if moovBox.boxType != TypeMoov {
		err = errors.Wrapf(ErrWrongBoxType, "Box %s", moovBox.boxType)
		return
	}
	if debugFlag {
		traceBox(moovBox, moovBox)
	}
	var inner box
	i := 0
	for moovBox.anyRemain() {
		if inner, err = moovBox.readInnerBox(); err != nil {
			return cmb, errors.Wrapf(err, "Box 'moov' %s (readBox)", inner.boxType)
		}
		switch inner.boxType {
		case TypeUUID:
			uuid, err := inner.readUUID()
			if err != nil {
				return cmb, err
			}
			if uuid == CR3MetaBoxUUID {
				if cmb.Meta, err = parseCR3MetaBox(&inner); err != nil {
					return cmb, err
				}
			}
		case TypeTrak:
			cmb.Trak[i], err = ParseCrxTrak(&inner)
			if err != nil {
				return
			}
			i++
			if debugFlag {
				traceBox(moovBox, inner)
			}
		default:
			if debugFlag {
				traceBox(moovBox, inner)
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

type CR3Trak struct {
	Width, Height     uint16
	Depth, ImageType  uint16
	ImageSize, Offset uint32
}

func ParseCrxTrak(trak *box) (t CR3Trak, err error) {
	buf, err := trak.read()
	if err != nil {
		return
	}
	var size int
	for remain := len(buf); remain > 0; remain -= size {
		size = int(binary.BigEndian.Uint32(buf[:4]))
		switch boxType(buf[4:8]) {
		case TypeMdia, TypeMinf, TypeStbl: // open box
			size = 8
		case TypeHdlr:
			// skip any trak whose hdlr type is not "vide"
			if string(buf[16:20]) != "vide" {
				return
			}
		case TypeStsd:
			if boxType(buf[4+stsdHeaderSize:8+stsdHeaderSize]) == TypeCRAW {
				t.Width = crxBinaryOrder.Uint16(buf[32+stsdHeaderSize : 34+stsdHeaderSize])
				t.Height = crxBinaryOrder.Uint16(buf[34+stsdHeaderSize : 36+stsdHeaderSize])
				t.Depth = crxBinaryOrder.Uint16(buf[82+stsdHeaderSize : 84+stsdHeaderSize])
				t.ImageType = crxBinaryOrder.Uint16(buf[86+stsdHeaderSize : 88+stsdHeaderSize])
			}
		case TypeStsz:
			if size == 20 {
				t.ImageSize = crxBinaryOrder.Uint32(buf[12+8+4 : 16+8+4])
			} else if size == 24 {
				t.ImageSize = crxBinaryOrder.Uint32(buf[12+8 : 16+8])
			}
		case TypeCo64:
			t.Offset = crxBinaryOrder.Uint32(buf[12+8 : 16+8])
		default: // skip other types:
		}
		buf = buf[size:]
	}
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
	//CCTP CCTPBox
	CTBO CTBOBox
	//THMB THMBBox
	Exif [4]meta.ExifHeader
}

// XPacketData returns CTBO[0] which corresponds to XPacket data
// First 24 bytes are a UUID box. uuid = be7acfcb-97a9-42e8-9c71-999491e3afac
func (cr3 CR3MetaBox) XPacketData() (offset, length uint32, err error) {
	item := cr3.CTBO.items[1]
	return item.offset, item.length, nil
}

// parseCR3MetaBox parses a uuid box with the uuid of 85c0b687 820f 11e0 8111 f4ce462b6a48
func parseCR3MetaBox(outer *box) (m CR3MetaBox, err error) {
	var size int
	var bt BoxType
	var buf []byte
	for remain := outer.remain; remain > 8; remain -= size {
		if buf, err = outer.peek(40); err != nil {
			return
		}
		size = int(binary.BigEndian.Uint32(buf[:4]))
		bt = boxType(buf[4:8])
		switch bt {
		case TypeCNCV:
			copy(m.CNCV.val[:], buf[8:38])
		case TypeCTBO:
			if buf, err = outer.peek(size); err != nil {
				return
			}
			m.CTBO, err = parseCTBO(buf[8:])
		case TypeCMT1:
			m.Exif[0], err = parseCMT(buf[8:16], ifds.IFD0, uint32(outer.offset+8), uint32(size-8))
		case TypeCMT2:
			m.Exif[1], err = parseCMT(buf[8:16], ifds.ExifIFD, uint32(outer.offset+8), uint32(size-8))
		case TypeCMT3:
			m.Exif[2], err = parseCMT(buf[8:16], ifds.MknoteIFD, uint32(outer.offset+8), uint32(size-8))
		case TypeCMT4:
			m.Exif[3], err = parseCMT(buf[8:16], ifds.GPSIFD, uint32(outer.offset+8), uint32(size-8))
		}
		if err != nil {
			return
		}
		err = outer.discard(size)
	}
	return
}

// CNCVBox is Canon Compressor Version box
// CaNon Codec Version?
type CNCVBox struct {
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

// CTBOBox is a Canon tracks base offsets Box?
type CTBOBox struct {
	items [5]IndexOffset
	count uint32
}

func (ctbo CTBOBox) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("CTBO | ItemCount:%d \n", ctbo.count))
	for idx, item := range ctbo.items {
		sb.WriteString(fmt.Sprintf("\t | Index:%d, Offset:%d, Size:%d \n", idx, item.offset, item.length))
	}
	return sb.String()
}

// IndexOffset has an offset and a length.
type IndexOffset struct {
	offset uint32
	length uint32
}

func parseCTBO(buf []byte) (ctbo CTBOBox, err error) {
	// Item Count
	ctbo.count = crxBinaryOrder.Uint32(buf[0:4])

	// Each item is 20 bytes in length
	for i := 4; i+20 <= len(buf); i += 20 {
		idx := crxBinaryOrder.Uint32(buf[i : i+4])
		if int(idx) < len(ctbo.items) {
			ctbo.items[idx] = IndexOffset{
				offset: uint32(crxBinaryOrder.Uint64(buf[i+4 : i+12])),
				length: uint32(crxBinaryOrder.Uint64(buf[i+12 : i+20])),
			}
		}
	}
	return ctbo, nil
}

func parseCMT(buf []byte, ifd ifds.IfdType, offset uint32, size uint32) (meta.ExifHeader, error) {
	binaryOrder := meta.BinaryOrder(buf[:4])
	header := meta.NewExifHeader(binaryOrder, binaryOrder.Uint32(buf[4:8]), offset, size, imagetype.ImageCR3)
	header.FirstIfd = ifd
	return header, nil
}
