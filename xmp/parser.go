package xmp

import (
	"fmt"
	"time"

	"github.com/evanoberholster/imagemeta/meta"
	"github.com/evanoberholster/imagemeta/xmp/xmpns"
)

func (xmp *XMP) parser(p property) (err error) {
	if len(p.Value()) == 0 {
		return
	}
	switch p.Namespace() {
	case xmpns.XMLnsNS:
		// Null operation
		return
	case xmpns.ExifNS:
		err = xmp.Exif.parse(p)
	case xmpns.AuxNS:
		err = xmp.Aux.parse(p)
	case xmpns.DcNS:
		err = xmp.DC.parse(p)
	case xmpns.XmpNS, xmpns.XapNS:
		err = xmp.Basic.parse(p)
	case xmpns.TiffNS:
		err = xmp.Tiff.parse(p)
	case xmpns.CrsNS:
		err = xmp.CRS.parse(p)
	case xmpns.XmpMMNS, xmpns.XapMMNS:
		err = xmp.MM.parse(p)
	default:
		//fmt.Println(p, ns)
		return
	}
	if err != nil {
		err = nil
		//fmt.Println(err, "\t", p)
	}

	return
}

// parseDate
func parseDate(buf []byte) (t time.Time, err error) {
	str := string(buf)
	if t, err = time.Parse("2006-01-02T15:04:05Z07:00", str); err != nil {
		if t, err = time.Parse("2006-01-02T15:04:05.00", str); err != nil {
			t, err = time.Parse("2006-01-02T15:04:05", str)
		}
	}
	return
}

func parseUUID(buf []byte) (uuid meta.UUID) {
	_, b := readUntil(buf, ':')
	if len(b) > 0 {
		buf = b
	}
	err := uuid.UnmarshalText(buf)
	if err != nil {
		if DebugMode {
			fmt.Println("Parse UUID error: ", err)
		}
	}
	return
}

//func parseUint32(s string) uint32 {
//	u64, err := strconv.ParseUint(s, 10, 32)
//	if err != nil {
//		return 0
//	}
//	return uint32(u64)
//}

// parseInt parses a []byte of a string representation of an int64 value and returns the value
func parseInt(buf []byte) (i int64) {
	var neg bool
	if buf[0] == '-' {
		buf = buf[1:]
		neg = true
	}
	for j := 0; j < len(buf); j++ {
		i *= 10
		i += int64(buf[j] - '0')
	}
	if neg {
		i *= -1
	}
	return
}

// parseUint parses a []byte of a string representation of a uint64 value and returns the value.
func parseUint(buf []byte) (u uint64) {
	for i := 0; i < len(buf); i++ {
		u *= 10
		u += uint64(buf[i] - '0')
	}
	return
}

// parseString parses a []byte and returns a string
func parseString(buf []byte) string {
	return string(buf)
}

func parseRational(buf []byte) (n uint32, d uint32) {
	// TODO: Improve parsing functionality
	for i := 0; i < len(buf); i++ {
		if buf[i] == '/' {
			if i < len(buf)+1 {
				n = uint32(parseUint(buf[:i]))
				d = uint32(parseUint(buf[i+1:]))
				return
			}
		}
	}
	return
}

func readUntil(buf []byte, delimiter byte) (a []byte, b []byte) {
	for i := 0; i < len(buf); i++ {
		if buf[i] == delimiter || buf[i] == '>' {
			return buf[:i], buf[i+1:]
		}
	}
	return buf, nil
}
