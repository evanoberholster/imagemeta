package utils

import "encoding/binary"

type ByteOrder int8

const (
	UnknownEndian ByteOrder = iota
	LittleEndian
	BigEndian
)

func (bo ByteOrder) String() string {
	switch bo {
	case LittleEndian:
		return "LittleEndian"
	case BigEndian:
		return "BigEndian"
	default:
		return "UnknownEndian"
	}
}

func (bo ByteOrder) Uint16(buf []byte) uint16 {
	if bo == BigEndian {
		return binary.BigEndian.Uint16(buf)
	}
	return binary.LittleEndian.Uint16(buf)
}

func (bo ByteOrder) Uint32(buf []byte) uint32 {
	if bo == BigEndian {
		return binary.BigEndian.Uint32(buf)
	}
	return binary.LittleEndian.Uint32(buf)
}

func (bo ByteOrder) Uint64(buf []byte) uint64 {
	if bo == BigEndian {
		return binary.BigEndian.Uint64(buf)
	}
	return binary.LittleEndian.Uint64(buf)
}

func (bo ByteOrder) PutUint16(b []byte, v uint16) {
	if bo == BigEndian {
		binary.BigEndian.PutUint16(b, v)
		return
	}
	binary.LittleEndian.PutUint16(b, v)
}

func (bo ByteOrder) PutUint32(b []byte, v uint32) {
	if bo == BigEndian {
		binary.BigEndian.PutUint32(b, v)
		return
	}
	binary.LittleEndian.PutUint32(b, v)
}

func (bo ByteOrder) PutUint64(b []byte, v uint64) {
	if bo == BigEndian {
		binary.BigEndian.PutUint64(b, v)
		return
	}
	binary.LittleEndian.PutUint64(b, v)
}

// BinaryOrder returns the binary.ByteOrder for a Tiff Header based
// on 4 bytes from the buf.
//
// Good reference:
// CIPA DC-008-2016; JEITA CP-3451D
// -> http://www.cipa.jp/std/documents/e/DC-008-Translation-2016-E.pdf
func BinaryOrder(buf []byte) ByteOrder {
	if isTiffBigEndian(buf) {
		return BigEndian
	}
	if isTiffLittleEndian(buf) {
		return LittleEndian
	}
	return UnknownEndian
}

// IsTiffLittleEndian checks the buf for the Tiff LittleEndian Signature
func isTiffLittleEndian(buf []byte) bool {
	return string(buf[:4]) == "II*\000"
	//return buf[0] == 0x49 &&
	//	buf[1] == 0x49 &&
	//	buf[2] == 0x2a &&
	//	buf[3] == 0x00
}

// IsTiffBigEndian checks the buf for the TiffBigEndianSignature
func isTiffBigEndian(buf []byte) bool {
	return string(buf[:4]) == "MM\000*"
	//return buf[0] == 0x4d &&
	//	buf[1] == 0x4d &&
	//	buf[2] == 0x00 &&
	//	buf[3] == 0x2a
}
