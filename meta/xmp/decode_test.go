package xmp

import (
	"testing"
)

func TestIdentifyPropertyFallbackToXMPNS(t *testing.T) {
	prop := identifyProperty([]byte("stEvt:action"))
	want := NewProperty(StEvtNS, Action)
	if !prop.Equals(want) {
		t.Fatalf("identifyProperty(stEvt:action) = %s, want %s", prop.String(), want.String())
	}

	prop = identifyProperty([]byte("xmpDM:videoFrameRate"))
	want = NewProperty(XmpDMNS, VideoFrameRate)
	if !prop.Equals(want) {
		t.Fatalf("identifyProperty(xmpDM:videoFrameRate) = %s, want %s", prop.String(), want.String())
	}
}
