package makernote

import "github.com/evanoberholster/imagemeta/meta/utils"

// Selected Nikon maker-note tags (https://exiftool.org/TagNames/Nikon.html).
//
// Intentionally excluded:
// - LensType value decoding
const (
	TagNikonMakerNoteVersion       uint16 = 0x0001
	TagNikonISO                    uint16 = 0x0002
	TagNikonColorMode              uint16 = 0x0003
	TagNikonQuality                uint16 = 0x0004
	TagNikonWhiteBalance           uint16 = 0x0005
	TagNikonSharpness              uint16 = 0x0006
	TagNikonFocusMode              uint16 = 0x0007
	TagNikonFlashSetting           uint16 = 0x0008
	TagNikonFlashType              uint16 = 0x0009
	TagNikonISOSelection           uint16 = 0x000f
	TagNikonISOSetting             uint16 = 0x0013
	TagNikonSerialNumber           uint16 = 0x001d
	TagNikonColorSpace             uint16 = 0x001e
	TagNikonVRInfo                 uint16 = 0x001f
	TagNikonActiveDLighting        uint16 = 0x0022
	TagNikonWorldTime              uint16 = 0x0024
	TagNikonISOInfo                uint16 = 0x0025
	TagNikonVignetteControl        uint16 = 0x002a
	TagNikonShutterMode            uint16 = 0x0034
	TagNikonMechanicalShutterCount uint16 = 0x0037
	TagNikonImageSizeRAW           uint16 = 0x003e
	TagNikonColorTemperatureAuto   uint16 = 0x004f
	TagNikonLensType               uint16 = 0x0083
	TagNikonLens                   uint16 = 0x0084
	TagNikonManualFocusDistance    uint16 = 0x0085
	TagNikonDigitalZoom            uint16 = 0x0086
	TagNikonFlashMode              uint16 = 0x0087
	TagNikonAFInfo                 uint16 = 0x0088
	TagNikonShootingMode           uint16 = 0x0089
	TagNikonLensFStops             uint16 = 0x008b
	TagNikonSerialNumber2          uint16 = 0x00a0
	TagNikonImageCount             uint16 = 0x00a5
	TagNikonDeletedImageCount      uint16 = 0x00a6
	TagNikonShutterCount           uint16 = 0x00a7
	TagNikonPowerUpTime            uint16 = 0x00b6
	TagNikonAFInfo2                uint16 = 0x00b7
	TagNikonFileInfo               uint16 = 0x00b8
	TagNikonAFTune                 uint16 = 0x00b9
	TagNikonSilentPhotography      uint16 = 0x00bf
)

const nikonHeaderLength = 18

// HasNikonHeader reports whether the maker-note payload starts with a Nikon label.
func HasNikonHeader(buf []byte) bool {
	return len(buf) >= 5 && string(buf[:5]) == "Nikon"
}

// ParseNikonHeader parses Nikon maker-note TIFF info from the leading header.
func ParseNikonHeader(buf []byte) (bo utils.ByteOrder, ifdRelOffset uint32, ok bool) {
	if len(buf) < nikonHeaderLength || !HasNikonHeader(buf) {
		return utils.UnknownEndian, 0, false
	}
	bo = utils.BinaryOrder(buf[10:14])
	if bo == utils.UnknownEndian {
		return utils.UnknownEndian, 0, false
	}
	return bo, bo.Uint32(buf[14:18]), true
}
