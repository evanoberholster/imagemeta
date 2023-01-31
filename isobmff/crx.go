package isobmff

import (
	"io"

	"github.com/evanoberholster/imagemeta/exif2/ifds"
	"github.com/evanoberholster/imagemeta/imagetype"
	"github.com/evanoberholster/imagemeta/meta"
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
			if logLevelDebug() {
				logBoxExt(&inner, zerolog.DebugLevel).Send()
			}
			err = inner.close()
		}
		if err != nil && logLevelError() {
			logBoxExt(&inner, zerolog.ErrorLevel).Send()
			return
		}
	}
	return crx, b.close()
}

// CMT Box

// readCMTBox reads a ISOBMFF Box "CMT1","CMT2","CMT3", or "CMT4" from CR3
func readCMTBox(b *box, exifReader func(r io.Reader, h meta.ExifHeader) error, ifdType ifds.IfdType) (header meta.ExifHeader, err error) {
	header, err = readExifHeader(b, ifdType, imagetype.ImageCR3)
	if err != nil {
		return
	}

	// TODO: implement Limited Reader to reduce allocations
	if exifReader != nil {
		if err = exifReader(b, header); err != nil && logLevelError() {
			logBoxExt(b, zerolog.ErrorLevel).Object("exifReader", header).Err(err).Send()
		}
	}
	return header, b.close()
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
		logBoxExt(b, zerolog.InfoLevel).Str("CNCV", string(cncv.version[:])).Send()
	}
	return cncv, b.close()
}

// CTBO Box

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
		logBoxExt(b, zerolog.InfoLevel).Array("items", ctbo).Send()
	}
	return ctbo, b.close()
}

// CTBOBox is a Canon tracks base offsets Box?
// items are [2]{offset,length}
type CTBOBox struct {
	items [5]offsetLength
	count uint32
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
	if logLevelInfo() {
		logBoxExt(b, zerolog.InfoLevel).Send()
	}
	return b.close()
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
	if logLevelInfo() {
		logBoxExt(b, zerolog.InfoLevel).Send()
	}
	return t, b.close()
}
