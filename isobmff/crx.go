package isobmff

import (
	"fmt"
	"io"
	"strings"

	"github.com/evanoberholster/imagemeta/exif2/ifds"
	"github.com/evanoberholster/imagemeta/imagetype"
	"github.com/evanoberholster/imagemeta/meta"
	"github.com/evanoberholster/imagemeta/meta/utils"
	"github.com/rs/zerolog"
)

// CrxMoovBox is a Canon Raw Moov Box
type CrxMoovBox struct {
	Meta CR3MetaBox
	Trak [5]CR3Trak
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
	ImageSize, Offset uint32
	Width, Height     uint16
	Depth, ImageType  uint16
}

func readCrxMoovBox(b *box, exifReader ExifReader) (crx CrxMoovBox, err error) {
	defer b.close()
	var inner box
	var ok bool
	for inner, ok, err = b.readInnerBox(); err == nil && ok; inner, ok, err = b.readInnerBox() {
		switch inner.boxType {
		case typeCNCV:
			crx.Meta.CNCV, err = readCNCVBox(&inner)
		//case typeCCTP:
		case typeCTBO:
			crx.Meta.CTBO, err = readCTBOBox(&inner)
		case typeCMT1:
			crx.Meta.Exif[0], err = readCMTBox(&inner, exifReader, ifds.IFD0)
		case typeCMT2:
			crx.Meta.Exif[1], err = readCMTBox(&inner, exifReader, ifds.ExifIFD)
		case typeCMT3:
			crx.Meta.Exif[2], err = readCMTBox(&inner, exifReader, ifds.MknoteIFD)
		case typeCMT4:
			crx.Meta.Exif[3], err = readCMTBox(&inner, exifReader, ifds.GPSIFD)
		//case typeTHMB:
		default:
			if logLevelInfo() {
				logInfoBox(inner)
			}
			inner.close()
		}
		if err != nil {
			return
		}
	}
	return
}

// CMT Box

func readCMTBox(b *box, exifReader func(r io.Reader, h meta.ExifHeader) error, ifdType ifds.IfdType) (header meta.ExifHeader, err error) {
	buf, err := b.Peek(16)
	if err != nil {
		return header, err
	}

	endian := utils.BinaryOrder(buf[:4])
	header = meta.NewExifHeader(endian, endian.Uint32(buf[4:8]), 0, uint32(b.size), imagetype.ImageCR3)
	header.FirstIfd = ifdType
	if logLevelInfo() {
		logCMTBox(b, header)
	}
	if err = b.Discard(8); err != nil {
		return
	}
	if exifReader != nil {
		if err = exifReader(b, header); err != nil {
			//fmt.Println(err)
		}
	}
	return header, b.close()
}

func logCMTBox(b *box, h meta.ExifHeader) {
	ev := Logger.Info()
	b.log(ev)
	ev.Uint32("FirstIfdOffset", h.FirstIfdOffset).Str("FirstIfd", h.FirstIfd.String()).Uint32("TiffHeaderOffset", h.TiffHeaderOffset).Uint32("ExifLength", h.ExifLength).Str("Endian", h.ByteOrder.String()).Str("ImageType", h.ImageType.String())
	logTraceFunction(ev)
	ev.Send()
}

// CNCV Box

// CNCVBox is Canon Compressor Version box
// CaNon Codec Version?
type CNCVBox struct {
	//format [9]byte
	//version [6]uint8
	version [30]byte
}

func readCNCVBox(b *box) (cncv CNCVBox, err error) {
	if !b.isType(typeCNCV) {
		return cncv, ErrWrongBoxType
	}
	buf, err := b.Peek(30)
	if err != nil {
		return CNCVBox{}, err
	}
	copy(cncv.version[:], buf[:30])
	if logLevelInfo() {
		cncv.log(b)
	}
	return cncv, b.close()
}

func (cncv CNCVBox) log(b *box) {
	ev := Logger.Info()
	b.log(ev)
	ev.Str("CNCV", string(cncv.version[:]))
	logTraceFunction(ev)
	ev.Send()
}

// CTBO Box

// CTBOBox is a Canon tracks base offsets Box?
// items are [2]{offset,length}
type CTBOBox struct {
	items [5]ctboItem
	count uint32
}

// ctboItem
type ctboItem struct {
	offset uint64
	length uint64
}

func readCTBOBox(b *box) (ctbo CTBOBox, err error) {
	if !b.isType(typeCTBO) {
		return ctbo, ErrWrongBoxType
	}
	buf, err := b.Peek(b.remain)
	if err != nil {
		return ctbo, err
	}
	// Item Count
	ctbo.count = crxEndian.Uint32(buf[0:4])

	// Each item is 20 bytes in length
	for i := 4; i+20 <= len(buf); i += 20 {
		idx := crxEndian.Uint32(buf[i:i+4]) - 1
		if int(idx) < len(ctbo.items) {
			ctbo.items[idx].offset = crxEndian.Uint64(buf[i+4 : i+12])
			ctbo.items[idx].length = crxEndian.Uint64(buf[i+12 : i+20])
		}
	}
	if logLevelInfo() {
		ctbo.log(b)
	}
	return ctbo, b.close()
}

// String is the stringer interface
func (ctbo CTBOBox) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("CTBO | ItemCount:%d \n", ctbo.count))
	for idx, item := range ctbo.items {
		sb.WriteString(fmt.Sprintf("\t | Index:%d, Offset:%d, Size:%d \n", idx, item.offset, item.length))
	}
	return sb.String()[:sb.Len()-1]
}

func (ctbo CTBOBox) log(b *box) {
	ev := Logger.Info()
	b.log(ev)
	ev.Array("items", ctbo)
	logTraceFunction(ev)
	ev.Send()
}

// MarshalZerologArray is a zerolog interface for logging
func (ctbo CTBOBox) MarshalZerologArray(a *zerolog.Array) {
	for i := 0; i < int(ctbo.count); i++ {
		item := ctbo.items[i]
		if item.length == 0 && item.offset == 0 {
			break
		}
		a.Object(item)
	}
}

// MarshalZerologObject is a zerolog interface for logging
func (ctboi ctboItem) MarshalZerologObject(e *zerolog.Event) {
	e.Uint64("length", ctboi.length).Uint64("offset", ctboi.offset)
}

func readCrxMdia(b *box) (err error) {
	//var inner box
	//var ok bool
	//for inner, ok, err = b.readInnerBox(); err == nil && ok; inner, ok, err = b.readInnerBox() {
	//	switch inner.boxType {
	//	case typeMdia:
	//
	//	}
	//	if logLevelInfo() {
	//		logInfoBox(inner)
	//	}
	//	inner.close()
	//}
	b.close()
	return
}

// Trak
func readCrxTrakBox(b *box) (t CR3Trak, err error) {
	//var inner box
	//var ok bool
	////for inner, ok, err = b.readInnerBox(); err == nil && ok; inner, ok, err = b.readInnerBox() {
	//	switch inner.boxType {
	//	case typeMdia:
	//		err = readCrxMdia(&inner)
	//	case typeHdlr:
	//
	//	case typeStsd:
	//
	//	case typeStsz:
	//
	//	case typeCo64:
	//
	//	}
	//	if logLevelInfo() {
	//		logInfoBox(inner)
	//	}
	//	inner.close()
	//}
	err = b.close()

	return
}

//func parseCrxTrak(outer *box) (t CR3Trak, err error) {
//	buf, err := outer.read()
//	if err != nil {
//		return
//	}
//
//	for remain := len(buf); remain > 0; {
//		bt, size := crxBoxHeader(buf)
//		remain -= size
//		switch bt {
//		case TypeMdia, TypeMinf, TypeStbl: // open box
//			size = 8
//		case TypeHdlr:
//			if logLevelDebug() {
//				logDebugBoxWithMsg(box{size: int64(size), boxType: bt, bufReader: bufReader{offset: outer.offset}}, "hdlr | type: "+string(buf[16:20]))
//			}
//			// skip any trak whose hdlr type is not "vide"
//			if string(buf[16:20]) != "vide" {
//				return
//			}
//		case TypeStsd:
//			if boxType(buf[4+stsdHeaderSize:8+stsdHeaderSize]) == TypeCRAW {
//				t.Width = crxBinaryOrder.Uint16(buf[32+stsdHeaderSize : 34+stsdHeaderSize])
//				t.Height = crxBinaryOrder.Uint16(buf[34+stsdHeaderSize : 36+stsdHeaderSize])
//				t.Depth = crxBinaryOrder.Uint16(buf[82+stsdHeaderSize : 84+stsdHeaderSize])
//				t.ImageType = crxBinaryOrder.Uint16(buf[86+stsdHeaderSize : 88+stsdHeaderSize])
//			}
//			if logLevelDebug() {
//				logDebugBoxWithMsg(box{size: int64(size), boxType: bt, bufReader: bufReader{offset: outer.offset}}, fmt.Sprintf("stsd | width:%d, height:%d, depth:%d, imagetype:%d", t.Width, t.Height, t.Depth, t.ImageType))
//			}
//		case TypeStsz:
//			if size == 20 {
//				t.ImageSize = crxBinaryOrder.Uint32(buf[12+8+4 : 16+8+4])
//			} else if size == 24 {
//				t.ImageSize = crxBinaryOrder.Uint32(buf[12+8 : 16+8])
//			}
//			if logLevelDebug() {
//				logDebugBoxWithMsg(box{size: int64(size), boxType: bt, bufReader: bufReader{offset: outer.offset}}, fmt.Sprintf("stsz | imageSize:%d", t.ImageSize))
//			}
//		case TypeCo64:
//			t.Offset = crxBinaryOrder.Uint32(buf[12+8 : 16+8])
//			if logLevelDebug() {
//				logDebugBoxWithMsg(box{size: int64(size), boxType: bt, bufReader: bufReader{offset: outer.offset}}, fmt.Sprintf("co64 | imageOffset:%d", t.Offset))
//			}
//		default: // skip other types:
//			if logLevelDebug() {
//				logDebugBoxWithMsg(box{size: int64(size), boxType: bt, bufReader: bufReader{offset: outer.offset}}, "discard")
//			}
//		}
//		buf = buf[size:]
//	}
//	return
//}
//
