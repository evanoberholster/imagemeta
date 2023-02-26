package canon

import "testing"

func TestString(t *testing.T) {
	if TagCanonString(CanonAFInfo) != "CanonAFInfo" {
		t.Errorf("Expected %s got %s", "CanonAFInfo", TagCanonString(CanonAFInfo))
	}
	if TagCanonString(0x1234) != "0x1234" {
		t.Errorf("Expected %s got %s", "0x1234", TagCanonString(0x1234))
	}
}
