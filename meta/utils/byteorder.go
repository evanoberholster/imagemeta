package utils

import "encoding/binary"

type ByteOrder int8

const (
	UnknownEndian ByteOrder = iota
	LittleEndian
	BigEndian
)

const (
	tiffLittleEndianSignature    = "II*\000"
	tiffBigEndianSignature       = "MM\000*"
	bigTiffLittleEndianSignature = "II+\000"
	bigTiffBigEndianSignature    = "MM\000+"
	tiffSignatureLength          = len(tiffLittleEndianSignature)
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
	if len(buf) < tiffSignatureLength {
		return UnknownEndian
	}

	switch string(buf[:tiffSignatureLength]) {
	case tiffLittleEndianSignature, bigTiffLittleEndianSignature:
		return LittleEndian
	case tiffBigEndianSignature, bigTiffBigEndianSignature:
		return BigEndian
	default:
		return UnknownEndian
	}
}

// IsBigTiffLittleEndian reports the BigTIFF (magic 43) little-endian signature.
func IsBigTiffLittleEndian(buf []byte) bool {
	return len(buf) >= tiffSignatureLength && string(buf[:tiffSignatureLength]) == bigTiffLittleEndianSignature
}

// IsBigTiffBigEndian reports the BigTIFF (magic 43) big-endian signature.
func IsBigTiffBigEndian(buf []byte) bool {
	return len(buf) >= tiffSignatureLength && string(buf[:tiffSignatureLength]) == bigTiffBigEndianSignature
}

// IsBigTiff reports whether buf carries either BigTIFF byte-order signature.
func IsBigTiff(buf []byte) bool {
	return IsBigTiffLittleEndian(buf) || IsBigTiffBigEndian(buf)
}

// IsTiffLittleEndian checks the buf for the Tiff LittleEndian Signature
func IsTiffLittleEndian(buf []byte) bool {
	return string(buf[:tiffSignatureLength]) == tiffLittleEndianSignature
}

// IsTiffBigEndian checks the buf for the TiffBigEndianSignature
func IsTiffBigEndian(buf []byte) bool {
	return string(buf[:tiffSignatureLength]) == tiffBigEndianSignature
}
