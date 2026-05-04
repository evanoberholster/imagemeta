package sony

import (
	"strings"

	"github.com/evanoberholster/imagemeta/meta"
	"github.com/evanoberholster/imagemeta/meta/utils"
)

func bytesAt(buf []byte, off, n int) []byte {
	if off < 0 || n <= 0 || off >= len(buf) {
		return nil
	}
	end := off + n
	if end > len(buf) {
		end = len(buf)
	}
	return buf[off:end]
}

func u8At(buf []byte, off int) uint8 {
	if off < 0 || off >= len(buf) {
		return 0
	}
	return buf[off]
}

func i8At(buf []byte, off int) int8 {
	return meta.SafecastUint8ToInt8Bits(u8At(buf, off))
}

func u16At(buf []byte, bo utils.ByteOrder, off int) uint16 {
	if off < 0 || off+2 > len(buf) {
		return 0
	}
	return bo.Uint16(buf[off : off+2])
}

func i16At(buf []byte, bo utils.ByteOrder, off int) int16 {
	return meta.SafecastUint16ToInt16Bits(u16At(buf, bo, off))
}

func u16RevAt(buf []byte, bo utils.ByteOrder, off int) uint16 {
	if off < 0 || off+2 > len(buf) {
		return 0
	}
	if bo == utils.BigEndian {
		return utils.LittleEndian.Uint16(buf[off : off+2])
	}
	return utils.BigEndian.Uint16(buf[off : off+2])
}

func u32At(buf []byte, bo utils.ByteOrder, off int) uint32 {
	if off < 0 || off+4 > len(buf) {
		return 0
	}
	return bo.Uint32(buf[off : off+4])
}

func DisplayText(buf []byte) string {
	if len(buf) == 0 {
		return ""
	}
	end := len(buf)
	for end > 0 {
		switch buf[end-1] {
		case 0, ' ', '\t', '\r', '\n':
			end--
		default:
			goto trimmed
		}
	}
trimmed:
	if end <= 0 {
		return ""
	}
	allPrintable := true
	for i := 0; i < end; i++ {
		ch := buf[i]
		if ch < 0x20 || ch > 0x7e {
			allPrintable = false
			break
		}
	}
	if allPrintable {
		return string(buf[:end])
	}
	var b strings.Builder
	b.Grow(end)
	for i := 0; i < end; i++ {
		ch := buf[i]
		if ch >= 0x20 && ch <= 0x7e {
			b.WriteByte(ch)
			continue
		}
		b.WriteByte('.')
	}
	return b.String()
}

func fillU16s(dst []uint16, buf []byte, bo utils.ByteOrder, off int, rev bool) {
	for i := range dst {
		if rev {
			dst[i] = u16RevAt(buf, bo, off+i*2)
		} else {
			dst[i] = u16At(buf, bo, off+i*2)
		}
	}
}

func fillI16s(dst []int16, buf []byte, bo utils.ByteOrder, off int) {
	for i := range dst {
		dst[i] = i16At(buf, bo, off+i*2)
	}
}
