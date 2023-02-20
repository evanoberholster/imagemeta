package gpsifd

import "testing"

func TestString(t *testing.T) {
	if TagString(GPSAltitude) != "GPSAltitude" {
		t.Errorf("Expected %s got %s", "GPSAltitude", TagString(GPSAltitude))
	}
	if TagString(0x1234) != "0x1234" {
		t.Errorf("Expected %s got %s", "0x1234", TagString(0x1234))
	}
}
