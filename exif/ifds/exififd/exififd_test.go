package exififd

import "testing"

func TestString(t *testing.T) {
	if TagString(FNumber) != "FNumber" {
		t.Errorf("Expected %s got %s", "FNumber", TagString(FNumber))
	}
	if TagString(0x1234) != "0x1234" {
		t.Errorf("Expected %s got %s", "0x1234", TagString(0x1234))
	}
}
