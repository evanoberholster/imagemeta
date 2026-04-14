package jpeg

import (
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/evanoberholster/imagemeta/meta/utils"
)

func trimNULString(b []byte) string {
	if i := strings.IndexByte(string(b), 0); i >= 0 {
		b = b[:i]
	}
	return strings.TrimSpace(string(b))
}

func u16ListString(order utils.ByteOrder, b []byte) string {
	if len(b) < 2 {
		return ""
	}
	n := len(b) / 2
	parts := make([]string, 0, n)
	for i := 0; i < n; i++ {
		parts = append(parts, strconv.FormatUint(uint64(order.Uint16(b[i*2:])), 10))
	}
	return strings.Join(parts, " ")
}

func u8ListString(b []byte) string {
	parts := make([]string, 0, len(b))
	for _, v := range b {
		parts = append(parts, strconv.FormatUint(uint64(v), 10))
	}
	return strings.Join(parts, " ")
}

func u32ListString(order utils.ByteOrder, b []byte) string {
	if len(b) < 4 {
		return ""
	}
	n := len(b) / 4
	parts := make([]string, 0, n)
	for i := 0; i < n; i++ {
		parts = append(parts, strconv.FormatUint(uint64(order.Uint32(b[i*4:])), 10))
	}
	return strings.Join(parts, " ")
}

func s15Fixed16(order utils.ByteOrder, b []byte) float64 {
	if len(b) < 4 {
		return 0
	}
	return float64(int32(order.Uint32(b))) / 65536.0
}

func fixed16(order utils.ByteOrder, b []byte) float64 {
	if len(b) < 4 {
		return 0
	}
	return float64(order.Uint32(b)) / 65536.0
}

func ciffFloat32(order utils.ByteOrder, b []byte) float64 {
	if len(b) < 4 {
		return 0
	}
	return float64(math.Float32frombits(order.Uint32(b)))
}

func ciffTime(order utils.ByteOrder, b []byte) time.Time {
	if len(b) < 4 {
		return time.Time{}
	}
	sec := int64(order.Uint32(b))
	if sec == 0 {
		return time.Time{}
	}
	return time.Unix(sec, 0).UTC()
}
