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

// CR3MetaBox is a uuidBox that contains Metadata for CR3 files
type CR3MetaBox struct {
	//CCTP CCTPBox
	//THMB THMBBox
	CNCV CNCVBox
	CTBO CTBOBox
	Exif [4]meta.ExifHeader
}

// CR3Trak is a Canon CR3 Trak box
type CR3Trak struct {
	Width, Height     uint16
	Depth, ImageType  uint16
	ImageSize, Offset uint32
}

// ReadCrxMoovBox is a performance focused method for parsing the moov box from a .CR3 file.
func (r *Reader) ReadCrxMoovBox() (cmb CrxMoovBox, err error) {
	// Parse Moov Box
	moovBox, err := r.readBox()
	if err != nil {
		return cmb, errors.Wrapf(err, "Box 'moov' (readBox) error")
	}
	if moovBox.boxType != TypeMoov {
		return cmb, errors.Wrapf(ErrWrongBoxType, "Box %s", moovBox.boxType)
	}
	if debugFlag {
		tracebox(moovBox)
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
				tracebox(inner)
			}
		default:
			if debugFlag {
				traceBoxWithMsg(inner, "discard")
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

func ParseCrxTrak(trak *box) (t CR3Trak, err error) {
	var bt BoxType
	var size int
	var buf []byte
	if buf, err = trak.read(); err != nil {
		return
	}

	for remain := len(buf); remain > 0; remain -= size {
		size = int(binary.BigEndian.Uint32(buf[:4]))
		bt = boxType(buf[4:8])
		switch bt {
		case TypeMdia, TypeMinf, TypeStbl: // open box
			size = 8
		case TypeHdlr:
			if debugFlag {
				traceBoxWithMsg(box{size: int64(size), boxType: bt, bufReader: bufReader{offset: trak.offset}}, "hdlr | type: "+string(buf[16:20]))
			}
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
			if debugFlag {
				traceBoxWithMsg(box{size: int64(size), boxType: bt, bufReader: bufReader{offset: trak.offset}}, fmt.Sprintf("stsd | width:%d, height:%d, depth:%d, imagetype:%d", t.Width, t.Height, t.Depth, t.ImageType))
			}
		case TypeStsz:
			if size == 20 {
				t.ImageSize = crxBinaryOrder.Uint32(buf[12+8+4 : 16+8+4])
			} else if size == 24 {
				t.ImageSize = crxBinaryOrder.Uint32(buf[12+8 : 16+8])
			}
			if debugFlag {
				traceBoxWithMsg(box{size: int64(size), boxType: bt, bufReader: bufReader{offset: trak.offset}}, fmt.Sprintf("stsz | imageSize:%d", t.ImageSize))
			}
		case TypeCo64:
			t.Offset = crxBinaryOrder.Uint32(buf[12+8 : 16+8])
			if debugFlag {
				traceBoxWithMsg(box{size: int64(size), boxType: bt, bufReader: bufReader{offset: trak.offset}}, fmt.Sprintf("co64 | imageOffset:%d", t.Offset))
			}
		default: // skip other types:
			if debugFlag {
				traceBoxWithMsg(box{size: int64(size), boxType: bt, bufReader: bufReader{offset: trak.offset}}, "discard")
			}
		}
		buf = buf[size:]
	}
	return
}

// XPacketData returns CTBO[0] which corresponds to XPacket data
// First 24 bytes are a UUID box. uuid = be7acfcb-97a9-42e8-9c71-999491e3afac
func (cr3 CR3MetaBox) XPacketData() (offset, length uint32, err error) {
	item := cr3.CTBO.items[0]
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
			if debugFlag {
				traceBoxWithMsg(*outer, m.CNCV.String())
			}
		case TypeCTBO:
			if buf, err = outer.peek(size); err != nil {
				return
			}
			m.CTBO, err = parseCTBO(buf[8:])
			if debugFlag {
				traceBoxWithMsg(*outer, m.CTBO.String())
			}
		case TypeCMT1:
			m.Exif[0], err = parseCMT(buf[8:16], ifds.IFD0, uint32(outer.offset+8), uint32(size-8))
			if debugFlag {
				traceBoxWithMsg(box{size: int64(size), boxType: bt, bufReader: bufReader{offset: outer.offset}}, "Exif Header: "+m.Exif[0].String())
			}
		case TypeCMT2:
			m.Exif[1], err = parseCMT(buf[8:16], ifds.ExifIFD, uint32(outer.offset+8), uint32(size-8))
			if debugFlag {
				traceBoxWithMsg(box{size: int64(size), boxType: bt, bufReader: bufReader{offset: outer.offset}}, "Exif Header: "+m.Exif[1].String())
			}
		case TypeCMT3:
			m.Exif[2], err = parseCMT(buf[8:16], ifds.MknoteIFD, uint32(outer.offset+8), uint32(size-8))
			if debugFlag {
				traceBoxWithMsg(box{size: int64(size), boxType: bt, bufReader: bufReader{offset: outer.offset}}, "Exif Header: "+m.Exif[2].String())
			}
		case TypeCMT4:
			m.Exif[3], err = parseCMT(buf[8:16], ifds.GPSIFD, uint32(outer.offset+8), uint32(size-8))
			if debugFlag {
				traceBoxWithMsg(box{size: int64(size), boxType: bt, bufReader: bufReader{offset: outer.offset}}, "Exif Header: "+m.Exif[3].String())
			}
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
	return sb.String()[:sb.Len()-1]
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
		if int(idx-1) < len(ctbo.items) {
			ctbo.items[idx-1] = IndexOffset{
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
