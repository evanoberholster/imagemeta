package meta

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