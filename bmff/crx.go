package bmff

import (
	"encoding/binary"
	"fmt"
	"io"
	"strings"

	"github.com/evanoberholster/imagemeta/meta"
)

// Values are in BigEndian
var crxBinaryOrder = binary.BigEndian

// CR3MetaBoxUUID is the uuid that corresponds with Canon CR3 Metadata
var CR3MetaBoxUUID, _ = meta.UUIDFromBytes([]byte{133, 192, 182, 135, 130, 15, 17, 224, 129, 17, 244, 206, 70, 43, 106, 72})

// CR3MetaBox is a uuidBox that contains Metadata for CR3 files
type CR3MetaBox struct {
	CNCV CNCVBox
	CCTP CCTPBox
	CTBO CTBOBox
	//CMT1
	//CMT2
	//CMT3
	//CMT4
	//THMB
	//mvhd
}

// Type returns TypeUUID, CR3MetaBox's boxType.
func (cr3 CR3MetaBox) Type() BoxType {
	return TypeUUID
}

// parseCR3MetaBox parses a uuid box with the uuid of 85c0b687 820f 11e0 8111 f4ce462b6a48
func parseCR3MetaBox(outer *box) (meta CR3MetaBox, err error) {
	var inner box
	for outer.anyRemain() {
		if inner, err = outer.readBox(); err != nil {
			return
		}
		switch inner.boxType {
		case TypeCNCV:
			meta.CNCV, err = parseCNCVBox(&inner)
			fmt.Println("CNCV:", meta.CNCV, err)
		case TypeCCTP:
			meta.CCTP, err = parseCCTPBox(&inner)
			fmt.Println("CCTP:", meta.CCTP, err)
		case TypeCTBO:
			meta.CTBO, err = parseCTBOBox(&inner)
			fmt.Println("CTBO:", meta.CTBO, err)
		}
		fmt.Println(inner)
		outer.remain -= int(inner.size)
		if err = inner.discard(inner.remain); err != nil {
			return
		}
	}
	return
}

// CCDTBox is a Canon CR3 definition of tracks?
type CCDTBox struct {
	//size uint32
	it  uint64
	idx uint32
}

func parseCCDTBox(outer *box) (ccdt CCDTBox, err error) {
	if outer.boxType != TypeCCDT {
		err = ErrWrongBoxType
		return
	}
	buf, err := outer.bufReader.Peek(16)
	if err != nil {
		return
	}
	if err = outer.discard(16); err != nil {
		return
	}
	// uint64 value appears to be imagetype
	ccdt.it = crxBinaryOrder.Uint64(buf[0:8])
	// uint32 value apprears to be 0 or 1 for dual pixel

	// uint32 value for the trak Index
	ccdt.idx = crxBinaryOrder.Uint32(buf[12:16])

	return ccdt, outer.discard(outer.remain)
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

func parseCNCVBox(outer *box) (cncv CNCVBox, err error) {
	buf, err := outer.bufReader.Peek(30)
	if err != nil {
		return
	}
	if err = outer.discard(30); err != nil {
		return
	}
	copy(cncv.val[:], buf[0:30])
	return cncv, nil
}

// CCTPBox is Canon Compressor Table Pointers box
// Canon CR3 trak pointers?
type CCTPBox struct {
	//size uint32
	CCDT []CCDTBox
}

func parseCCTPBox(outer *box) (cctp CCTPBox, err error) {
	if outer.boxType != TypeCCTP {
		err = ErrWrongBoxType
		return
	}
	buf, err := outer.bufReader.Peek(12)
	if err != nil {
		return
	}
	if err = outer.discard(12); err != nil {
		return
	}
	// CCTP Box contains 12 bytes (3 x uint32)
	// last one is number of CCDT lines. 3, or 4 for dual pixel
	count := crxBinaryOrder.Uint32(buf[8:12])
	cctp.CCDT = make([]CCDTBox, count)

	var inner box
	for i := 0; i < int(count) && outer.anyRemain(); i++ {
		inner, err = outer.readBox()
		if err != nil {
			if err == io.EOF {
				return
			}
			return
		}
		if inner.boxType == TypeCCDT {
			cctp.CCDT[i], err = parseCCDTBox(&inner)
		}
		outer.remain -= int(inner.size)
		if err = inner.discard(inner.remain); err != nil {
			break
		}
	}
	return cctp, outer.discard(outer.remain)
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

func parseCTBOBox(outer *box) (ctbo CTBOBox, err error) {
	if outer.Type() != TypeCTBO {
		err = ErrWrongBoxType
		return
	}
	buf, err := outer.bufReader.Peek(4)
	if err != nil {
		return
	}
	if err = outer.discard(4); err != nil {
		return
	}
	count := crxBinaryOrder.Uint32(buf[0:4])
	ctbo.items = make([]IndexOffset, count)

	for i := 0; i < int(count) && outer.anyRemain(); i++ {
		// each item is 20 bytes in length
		buf, err = outer.bufReader.Peek(20)
		if err != nil {
			return
		}
		if err = outer.discard(20); err != nil {
			return
		}
		ctbo.items[i] = IndexOffset{
			idx:    crxBinaryOrder.Uint32(buf[0:4]),
			offset: crxBinaryOrder.Uint64(buf[4:12]),
			length: crxBinaryOrder.Uint64(buf[12:20]),
		}
	}
	return ctbo, outer.discard(outer.remain)
}
