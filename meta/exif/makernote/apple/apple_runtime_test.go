package apple

import "testing"

// buildPlist assembles a minimal binary-plist body (8-byte header, objects,
// 1-byte offset table, 32-byte trailer) using 1-byte offset and object-ref
// sizes. objects is the concatenated object data starting at offset 8;
// offsets are the absolute byte offsets of each object.
func buildPlist(objects []byte, offsets []byte, topObject byte) []byte {
	raw := []byte("bplist00")
	raw = append(raw, objects...)
	offTableOff := len(raw)
	raw = append(raw, offsets...)
	trailer := make([]byte, 32)
	trailer[6] = 1                   // offsetIntSize
	trailer[7] = 1                   // objRefSize
	trailer[15] = byte(len(offsets)) // numObjects
	trailer[23] = topObject          // topObject
	trailer[31] = byte(offTableOff)  // offsetTableOffset
	raw = append(raw, trailer...)
	return raw
}

// TestParseRunTimeCyclicReference ensures a binary plist whose array refers
// back to itself is rejected rather than recursing until the goroutine stack
// is exhausted (an unrecoverable fatal error / DoS).
func TestParseRunTimeCyclicReference(t *testing.T) {
	// obj0 at offset 8: array of length 1 whose single element refers to obj0.
	objects := []byte{0xA1, 0x00}
	offsets := []byte{0x08}
	raw := buildPlist(objects, offsets, 0)

	if _, ok := ParseRunTime(raw); ok {
		t.Fatalf("expected ParseRunTime to reject cyclic plist")
	}
}

// TestParseRunTimeValidDict confirms a well-formed binary-plist dictionary is
// still parsed correctly after the cycle guard was added.
func TestParseRunTimeValidDict(t *testing.T) {
	// obj0 dict{obj1:obj2}, obj1 string "flags", obj2 int 42.
	objects := []byte{
		0xD1, 0x01, 0x02, // off 8: dict len1 key->1 val->2
		0x55, 'f', 'l', 'a', 'g', 's', // off 11: string "flags"
		0x10, 0x2A, // off 17: int 42
	}
	offsets := []byte{0x08, 0x0B, 0x11}
	raw := buildPlist(objects, offsets, 0)

	rt, ok := ParseRunTime(raw)
	if !ok {
		t.Fatalf("expected ParseRunTime to parse valid dict")
	}
	if rt.Flags != 42 {
		t.Fatalf("Flags = %d, want 42", rt.Flags)
	}
}
