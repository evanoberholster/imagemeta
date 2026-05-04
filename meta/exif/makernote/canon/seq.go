package canon

import "github.com/evanoberholster/imagemeta/meta"

// Seq16 is a 1-based indexed view over a uint16 slice, matching ExifTool's
// FIRST_ENTRY=1 convention used by Canon CameraSettings and ShotInfo payloads.
type Seq16 []uint16

// Present reports whether sequence index n (1-based) exists.
func (s Seq16) Present(n int) bool {
	return n > 0 && n <= len(s)
}

// U16 returns the uint16 at sequence index n (1-based), or 0 if absent.
func (s Seq16) U16(n int) uint16 {
	if n <= 0 || n > len(s) {
		return 0
	}
	return s[n-1]
}

// I16 returns the int16 at sequence index n (1-based), or 0 if absent.
func (s Seq16) I16(n int) int16 {
	return meta.SafecastUint16ToInt16Bits(s.U16(n))
}

// Seq32 is a 1-based indexed view over an int32 slice, matching ExifTool's
// FIRST_ENTRY=1 convention used by Canon int32-based payloads.
type Seq32 []int32

// Present reports whether sequence index n (1-based) exists.
func (s Seq32) Present(n int) bool {
	return n > 0 && n <= len(s)
}

// I32 returns the int32 at sequence index n (1-based), or 0 if absent.
func (s Seq32) I32(n int) int32 {
	if n <= 0 || n > len(s) {
		return 0
	}
	return s[n-1]
}

// U32 returns the uint32 at sequence index n (1-based), or 0 if absent.
func (s Seq32) U32(n int) uint32 {
	return meta.SafecastInt32ToUint32Bits(s.I32(n))
}
