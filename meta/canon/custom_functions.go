package canon

import (
	"encoding/binary"
	"fmt"
)

// CustomFunctions1D represents the custom functions for EOS-1D series cameras.
// Based on Phil Harvey's exiftool.
// Reference: https://exiftool.org/TagNames/CanonCustom.html#Functions1D
type CustomFunctions1D struct {
	Cfn01ShootingModeAFStop    uint32 `json:"cfn01ShootingModeAFStop"`    // 1
	Cfn02ShutterRelease        uint32 `json:"cfn02ShutterRelease"`        // 2
	Cfn03LensAFStopButton      uint32 `json:"cfn03LensAFStopButton"`      // 3
	Cfn04ShutterAELock         uint32 `json:"cfn04ShutterAELock"`         // 4
	Cfn05AELock                uint32 `json:"cfn05AELock"`                // 5
	Cfn06ExposureComp          uint32 `json:"cfn06ExposureComp"`          // 6
	Cfn07ManualTvAv            uint32 `json:"cfn07ManualTvAv"`            // 7
	Cfn08AFAssistBeam          uint32 `json:"cfn08AFAssistBeam"`          // 8
	Cfn09MirrorLockup          uint32 `json:"cfn09MirrorLockup"`          // 9
	Cfn10AFPointRegistration   uint32 `json:"cfn10AFPointRegistration"`   // 10
	Cfn11AFPointSelection      uint32 `json:"cfn11AFPointSelection"`      // 11
	Cfn12AFPointExpansion      uint32 `json:"cfn12AFPointExpansion"`      // 12
	Cfn13AELockButton          uint32 `json:"cfn13AELockButton"`          // 13
	Cfn14FillFlash             uint32 `json:"cfn14FillFlash"`             // 14
	Cfn15ShutterCurtainSync    uint32 `json:"cfn15ShutterCurtainSync"`    // 15
	Cfn16SafetyShift           uint32 `json:"cfn16SafetyShift"`           // 16
	Cfn17LensAFSearch          uint32 `json:"cfn17LensAFSearch"`          // 17
	Cfn18AutoFlash             uint32 `json:"cfn18AutoFlash"`             // 18
	Cfn19MenuButton            uint32 `json:"cfn19MenuButton"`            // 19
	Cfn20SuperimposedDisplay   uint32 `json:"cfn20SuperimposedDisplay"`   // 20
	Cfn21AddOriginalDecision   uint32 `json:"cfn21AddOriginalDecision"`   // 21
	Cfn22ShutterSpeedRange     uint32 `json:"cfn22ShutterSpeedRange"`     // 22
	Cfn23ApertureRange         uint32 `json:"cfn23ApertureRange"`         // 23
	Cfn24EraseButton           uint32 `json:"cfn24EraseButton"`           // 24
	Cfn25ExposureLevel         uint32 `json:"cfn25ExposureLevel"`         // 25
	Cfn26FlashSyncSpeed        uint32 `json:"cfn26FlashSyncSpeed"`        // 26
	Cfn27ETTLII                uint32 `json:"cfn27ETTLII"`                // 27
	Cfn28FlashExposureCompLock uint32 `json:"cfn28FlashExposureCompLock"` // 28
	Cfn29FocusingScreen        uint32 `json:"cfn29FocusingScreen"`        // 29
	Cfn30AFPointBrightness     uint32 `json:"cfn30AFPointBrightness"`     // 30
	Cfn31ShortReleaseTimeLag   uint32 `json:"cfn31ShortReleaseTimeLag"`   // 31
}

// UnmarshalBinary unmarshals the binary data into a CustomFunctions1D struct.
func (cf *CustomFunctions1D) UnmarshalBinary(b []byte) error {
	if len(b) < 128 { // 32 functions * 4 bytes
		return fmt.Errorf("invalid data length for CustomFunctions1D: got %d, want at least 128", len(b))
	}

	cf.Cfn01ShootingModeAFStop = binary.LittleEndian.Uint32(b[4:8])
	cf.Cfn02ShutterRelease = binary.LittleEndian.Uint32(b[8:12])
	cf.Cfn03LensAFStopButton = binary.LittleEndian.Uint32(b[12:16])
	cf.Cfn04ShutterAELock = binary.LittleEndian.Uint32(b[16:20])
	cf.Cfn05AELock = binary.LittleEndian.Uint32(b[20:24])
	cf.Cfn06ExposureComp = binary.LittleEndian.Uint32(b[24:28])
	cf.Cfn07ManualTvAv = binary.LittleEndian.Uint32(b[28:32])
	cf.Cfn08AFAssistBeam = binary.LittleEndian.Uint32(b[32:36])
	cf.Cfn09MirrorLockup = binary.LittleEndian.Uint32(b[36:40])
	cf.Cfn10AFPointRegistration = binary.LittleEndian.Uint32(b[40:44])
	cf.Cfn11AFPointSelection = binary.LittleEndian.Uint32(b[44:48])
	cf.Cfn12AFPointExpansion = binary.LittleEndian.Uint32(b[48:52])
	cf.Cfn13AELockButton = binary.LittleEndian.Uint32(b[52:56])
	cf.Cfn14FillFlash = binary.LittleEndian.Uint32(b[56:60])
	cf.Cfn15ShutterCurtainSync = binary.LittleEndian.Uint32(b[60:64])
	cf.Cfn16SafetyShift = binary.LittleEndian.Uint32(b[64:68])
	cf.Cfn17LensAFSearch = binary.LittleEndian.Uint32(b[68:72])
	cf.Cfn18AutoFlash = binary.LittleEndian.Uint32(b[72:76])
	cf.Cfn19MenuButton = binary.LittleEndian.Uint32(b[76:80])
	cf.Cfn20SuperimposedDisplay = binary.LittleEndian.Uint32(b[80:84])
	cf.Cfn21AddOriginalDecision = binary.LittleEndian.Uint32(b[84:88])
	cf.Cfn22ShutterSpeedRange = binary.LittleEndian.Uint32(b[88:92])
	cf.Cfn23ApertureRange = binary.LittleEndian.Uint32(b[92:96])
	cf.Cfn24EraseButton = binary.LittleEndian.Uint32(b[96:100])
	cf.Cfn25ExposureLevel = binary.LittleEndian.Uint32(b[100:104])
	cf.Cfn26FlashSyncSpeed = binary.LittleEndian.Uint32(b[104:108])
	cf.Cfn27ETTLII = binary.LittleEndian.Uint32(b[108:112])
	cf.Cfn28FlashExposureCompLock = binary.LittleEndian.Uint32(b[112:116])
	cf.Cfn29FocusingScreen = binary.LittleEndian.Uint32(b[116:120])
	cf.Cfn30AFPointBrightness = binary.LittleEndian.Uint32(b[120:124])
	cf.Cfn31ShortReleaseTimeLag = binary.LittleEndian.Uint32(b[124:128])

	return nil
}

// CustomFunctions5D represents the custom functions for the EOS 5D.
// Based on Phil Harvey's exiftool.
// Reference: https://exiftool.org/TagNames/CanonCustom.html#Functions5D
type CustomFunctions5D struct {
	Cfn01SetButtonFunc         uint32 `json:"cfn01SetButtonFunc"`         // 1
	Cfn02LongExpNoiseReduction uint32 `json:"cfn02LongExpNoiseReduction"` // 2
	Cfn03FlashSyncSpeedAv      uint32 `json:"cfn03FlashSyncSpeedAv"`      // 3
	Cfn04ShutterAELock         uint32 `json:"cfn04ShutterAELock"`         // 4
	Cfn05AFAssistBeam          uint32 `json:"cfn05AFAssistBeam"`          // 5
	Cfn06ExposureLevel         uint32 `json:"cfn06ExposureLevel"`         // 6
	Cfn07MirrorLockup          uint32 `json:"cfn07MirrorLockup"`          // 7
	Cfn08AFPointSelection      uint32 `json:"cfn08AFPointSelection"`      // 8
	Cfn09AFPointExpansion      uint32 `json:"cfn09AFPointExpansion"`      // 9
	Cfn10FocusingScreen        uint32 `json:"cfn10FocusingScreen"`        // 10
	Cfn11MenuButton            uint32 `json:"cfn11MenuButton"`            // 11
	Cfn12SuperimposedDisplay   uint32 `json:"cfn12SuperimposedDisplay"`   // 12
	Cfn13ShutterRelease        uint32 `json:"cfn13ShutterRelease"`        // 13
	Cfn14EraseButton           uint32 `json:"cfn14EraseButton"`           // 14
	Cfn15ShutterCurtainSync    uint32 `json:"cfn15ShutterCurtainSync"`    // 15
	Cfn16SafetyShift           uint32 `json:"cfn16SafetyShift"`           // 16
	Cfn17MagnifiedView         uint32 `json:"cfn17MagnifiedView"`         // 17
	Cfn18LensAFStopButton      uint32 `json:"cfn18LensAFStopButton"`      // 18
	Cfn19AddOriginalDecision   uint32 `json:"cfn19AddOriginalDecision"`   // 19
	Cfn20ETTLII                uint32 `json:"cfn20ETTLII"`                // 20
	Cfn21FlashExposureCompLock uint32 `json:"cfn21FlashExposureCompLock"` // 21
}

// UnmarshalBinary unmarshals the binary data into a CustomFunctions5D struct.
func (cf *CustomFunctions5D) UnmarshalBinary(b []byte) error {
	if len(b) < 88 { // 22 functions * 4 bytes
		return fmt.Errorf("invalid data length for CustomFunctions5D: got %d, want at least 88", len(b))
	}

	cf.Cfn01SetButtonFunc = binary.LittleEndian.Uint32(b[4:8])
	cf.Cfn02LongExpNoiseReduction = binary.LittleEndian.Uint32(b[8:12])
	cf.Cfn03FlashSyncSpeedAv = binary.LittleEndian.Uint32(b[12:16])
	cf.Cfn04ShutterAELock = binary.LittleEndian.Uint32(b[16:20])
	cf.Cfn05AFAssistBeam = binary.LittleEndian.Uint32(b[20:24])
	cf.Cfn06ExposureLevel = binary.LittleEndian.Uint32(b[24:28])
	cf.Cfn07MirrorLockup = binary.LittleEndian.Uint32(b[28:32])
	cf.Cfn08AFPointSelection = binary.LittleEndian.Uint32(b[32:36])
	cf.Cfn09AFPointExpansion = binary.LittleEndian.Uint32(b[36:40])
	cf.Cfn10FocusingScreen = binary.LittleEndian.Uint32(b[40:44])
	cf.Cfn11MenuButton = binary.LittleEndian.Uint32(b[44:48])
	cf.Cfn12SuperimposedDisplay = binary.LittleEndian.Uint32(b[48:52])
	cf.Cfn13ShutterRelease = binary.LittleEndian.Uint32(b[52:56])
	cf.Cfn14EraseButton = binary.LittleEndian.Uint32(b[56:60])
	cf.Cfn15ShutterCurtainSync = binary.LittleEndian.Uint32(b[60:64])
	cf.Cfn16SafetyShift = binary.LittleEndian.Uint32(b[64:68])
	cf.Cfn17MagnifiedView = binary.LittleEndian.Uint32(b[68:72])
	cf.Cfn18LensAFStopButton = binary.LittleEndian.Uint32(b[72:76])
	cf.Cfn19AddOriginalDecision = binary.LittleEndian.Uint32(b[76:80])
	cf.Cfn20ETTLII = binary.LittleEndian.Uint32(b[80:84])
	cf.Cfn21FlashExposureCompLock = binary.LittleEndian.Uint32(b[84:88])

	return nil
}

// CustomFunctions10D represents the custom functions for the EOS 10D.
// Based on Phil Harvey's exiftool.
// Reference: https://exiftool.org/TagNames/CanonCustom.html#Functions10D
type CustomFunctions10D struct {
	Cfn01SetButtonFunc         uint32 `json:"cfn01SetButtonFunc"`         // 1
	Cfn02ShutterRelease        uint32 `json:"cfn02ShutterRelease"`        // 2
	Cfn03MirrorLockup          uint32 `json:"cfn03MirrorLockup"`          // 3
	Cfn04ShutterAELock         uint32 `json:"cfn04ShutterAELock"`         // 4
	Cfn05AFAssistBeam          uint32 `json:"cfn05AFAssistBeam"`          // 5
	Cfn06FlashSyncSpeedAv      uint32 `json:"cfn06FlashSyncSpeedAv"`      // 6
	Cfn07AELockButton          uint32 `json:"cfn07AELockButton"`          // 7
	Cfn08ShutterCurtainSync    uint32 `json:"cfn08ShutterCurtainSync"`    // 8
	Cfn09LensAFSearch          uint32 `json:"cfn09LensAFSearch"`          // 9
	Cfn10EraseButton           uint32 `json:"cfn10EraseButton"`           // 10
	Cfn11MenuButton            uint32 `json:"cfn11MenuButton"`            // 11
	Cfn12SuperimposedDisplay   uint32 `json:"cfn12SuperimposedDisplay"`   // 12
	Cfn13AFPointRegistration   uint32 `json:"cfn13AFPointRegistration"`   // 13
	Cfn14FillFlash             uint32 `json:"cfn14FillFlash"`             // 14
	Cfn15AutoFlash             uint32 `json:"cfn15AutoFlash"`             // 15
	Cfn16ShutterSpeedRange     uint32 `json:"cfn16ShutterSpeedRange"`     // 16
	Cfn17AFPointBrightness     uint32 `json:"cfn17AFPointBrightness"`     // 17
	Cfn18AddOriginalDecision   uint32 `json:"cfn18AddOriginalDecision"`   // 18
	Cfn19FlashExposureCompLock uint32 `json:"cfn19FlashExposureCompLock"` // 19
}

// UnmarshalBinary unmarshals the binary data into a CustomFunctions10D struct.
func (cf *CustomFunctions10D) UnmarshalBinary(b []byte) error {
	if len(b) < 80 { // 20 functions * 4 bytes
		return fmt.Errorf("invalid data length for CustomFunctions10D: got %d, want at least 80", len(b))
	}
	cf.Cfn01SetButtonFunc = binary.LittleEndian.Uint32(b[4:8])
	cf.Cfn02ShutterRelease = binary.LittleEndian.Uint32(b[8:12])
	cf.Cfn03MirrorLockup = binary.LittleEndian.Uint32(b[12:16])
	cf.Cfn04ShutterAELock = binary.LittleEndian.Uint32(b[16:20])
	cf.Cfn05AFAssistBeam = binary.LittleEndian.Uint32(b[20:24])
	cf.Cfn06FlashSyncSpeedAv = binary.LittleEndian.Uint32(b[24:28])
	cf.Cfn07AELockButton = binary.LittleEndian.Uint32(b[28:32])
	cf.Cfn08ShutterCurtainSync = binary.LittleEndian.Uint32(b[32:36])
	cf.Cfn09LensAFSearch = binary.LittleEndian.Uint32(b[36:40])
	cf.Cfn10EraseButton = binary.LittleEndian.Uint32(b[40:44])
	cf.Cfn11MenuButton = binary.LittleEndian.Uint32(b[44:48])
	cf.Cfn12SuperimposedDisplay = binary.LittleEndian.Uint32(b[48:52])
	cf.Cfn13AFPointRegistration = binary.LittleEndian.Uint32(b[52:56])
	cf.Cfn14FillFlash = binary.LittleEndian.Uint32(b[56:60])
	cf.Cfn15AutoFlash = binary.LittleEndian.Uint32(b[60:64])
	cf.Cfn16ShutterSpeedRange = binary.LittleEndian.Uint32(b[64:68])
	cf.Cfn17AFPointBrightness = binary.LittleEndian.Uint32(b[68:72])
	cf.Cfn18AddOriginalDecision = binary.LittleEndian.Uint32(b[72:76])
	cf.Cfn19FlashExposureCompLock = binary.LittleEndian.Uint32(b[76:80])
	return nil
}

// CustomFunctions20D represents the custom functions for the EOS 20D.
// Based on Phil Harvey's exiftool.
// Reference: https://exiftool.org/TagNames/CanonCustom.html#Functions20D
type CustomFunctions20D struct {
	Cfn01SetButtonFunc         uint32 `json:"cfn01SetButtonFunc"`         // 1
	Cfn02LongExpNoiseReduction uint32 `json:"cfn02LongExpNoiseReduction"` // 2
	Cfn03FlashSyncSpeedAv      uint32 `json:"cfn03FlashSyncSpeedAv"`      // 3
	Cfn04ShutterAELock         uint32 `json:"cfn04ShutterAELock"`         // 4
	Cfn05AFAssistBeam          uint32 `json:"cfn05AFAssistBeam"`          // 5
	Cfn06ExposureLevel         uint32 `json:"cfn06ExposureLevel"`         // 6
	Cfn07MirrorLockup          uint32 `json:"cfn07MirrorLockup"`          // 7
	Cfn08AFPointSelection      uint32 `json:"cfn08AFPointSelection"`      // 8
	Cfn09AFPointExpansion      uint32 `json:"cfn09AFPointExpansion"`      // 9
	Cfn10MenuButton            uint32 `json:"cfn10MenuButton"`            // 10
	Cfn11SuperimposedDisplay   uint32 `json:"cfn11SuperimposedDisplay"`   // 11
	Cfn12ShutterRelease        uint32 `json:"cfn12ShutterRelease"`        // 12
	Cfn13EraseButton           uint32 `json:"cfn13EraseButton"`           // 13
	Cfn14ShutterCurtainSync    uint32 `json:"cfn14ShutterCurtainSync"`    // 14
	Cfn15SafetyShift           uint32 `json:"cfn15SafetyShift"`           // 15
	Cfn16LensAFStopButton      uint32 `json:"cfn16LensAFStopButton"`      // 16
	Cfn17AddOriginalDecision   uint32 `json:"cfn17AddOriginalDecision"`   // 17
	Cfn18ETTLII                uint32 `json:"cfn18ETTLII"`                // 18
	Cfn19FlashExposureCompLock uint32 `json:"cfn19FlashExposureCompLock"` // 19
}

// UnmarshalBinary unmarshals the binary data into a CustomFunctions20D struct.
func (cf *CustomFunctions20D) UnmarshalBinary(b []byte) error {
	if len(b) < 80 { // 20 functions * 4 bytes
		return fmt.Errorf("invalid data length for CustomFunctions20D: got %d, want at least 80", len(b))
	}
	cf.Cfn01SetButtonFunc = binary.LittleEndian.Uint32(b[4:8])
	cf.Cfn02LongExpNoiseReduction = binary.LittleEndian.Uint32(b[8:12])
	cf.Cfn03FlashSyncSpeedAv = binary.LittleEndian.Uint32(b[12:16])
	cf.Cfn04ShutterAELock = binary.LittleEndian.Uint32(b[16:20])
	cf.Cfn05AFAssistBeam = binary.LittleEndian.Uint32(b[20:24])
	cf.Cfn06ExposureLevel = binary.LittleEndian.Uint32(b[24:28])
	cf.Cfn07MirrorLockup = binary.LittleEndian.Uint32(b[28:32])
	cf.Cfn08AFPointSelection = binary.LittleEndian.Uint32(b[32:36])
	cf.Cfn09AFPointExpansion = binary.LittleEndian.Uint32(b[36:40])
	cf.Cfn10MenuButton = binary.LittleEndian.Uint32(b[40:44])
	cf.Cfn11SuperimposedDisplay = binary.LittleEndian.Uint32(b[44:48])
	cf.Cfn12ShutterRelease = binary.LittleEndian.Uint32(b[48:52])
	cf.Cfn13EraseButton = binary.LittleEndian.Uint32(b[52:56])
	cf.Cfn14ShutterCurtainSync = binary.LittleEndian.Uint32(b[56:60])
	cf.Cfn15SafetyShift = binary.LittleEndian.Uint32(b[60:64])
	cf.Cfn16LensAFStopButton = binary.LittleEndian.Uint32(b[64:68])
	cf.Cfn17AddOriginalDecision = binary.LittleEndian.Uint32(b[68:72])
	cf.Cfn18ETTLII = binary.LittleEndian.Uint32(b[72:76])
	cf.Cfn19FlashExposureCompLock = binary.LittleEndian.Uint32(b[76:80])
	return nil
}

// CustomFunctions30D represents the custom functions for the EOS 30D.
type CustomFunctions30D = CustomFunctions20D

// CustomFunctions350D represents the custom functions for the EOS 350D.
// Based on Phil Harvey's exiftool.
// Reference: https://exiftool.org/TagNames/CanonCustom.html#Functions350D
type CustomFunctions350D struct {
	Cfn01SetButtonFunc         uint32 `json:"cfn01SetButtonFunc"`         // 1
	Cfn02LongExpNoiseReduction uint32 `json:"cfn02LongExpNoiseReduction"` // 2
	Cfn03FlashSyncSpeedAv      uint32 `json:"cfn03FlashSyncSpeedAv"`      // 3
	Cfn04ShutterAELock         uint32 `json:"cfn04ShutterAELock"`         // 4
	Cfn05AFAssistBeam          uint32 `json:"cfn05AFAssistBeam"`          // 5
	Cfn06ExposureLevel         uint32 `json:"cfn06ExposureLevel"`         // 6
	Cfn07MirrorLockup          uint32 `json:"cfn07MirrorLockup"`          // 7
	Cfn08ShutterCurtainSync    uint32 `json:"cfn08ShutterCurtainSync"`    // 8
	Cfn09LensAFSearch          uint32 `json:"cfn09LensAFSearch"`          // 9
}

// UnmarshalBinary unmarshals the binary data into a CustomFunctions350D struct.
func (cf *CustomFunctions350D) UnmarshalBinary(b []byte) error {
	if len(b) < 40 { // 10 functions * 4 bytes
		return fmt.Errorf("invalid data length for CustomFunctions350D: got %d, want at least 40", len(b))
	}
	cf.Cfn01SetButtonFunc = binary.LittleEndian.Uint32(b[4:8])
	cf.Cfn02LongExpNoiseReduction = binary.LittleEndian.Uint32(b[8:12])
	cf.Cfn03FlashSyncSpeedAv = binary.LittleEndian.Uint32(b[12:16])
	cf.Cfn04ShutterAELock = binary.LittleEndian.Uint32(b[16:20])
	cf.Cfn05AFAssistBeam = binary.LittleEndian.Uint32(b[20:24])
	cf.Cfn06ExposureLevel = binary.LittleEndian.Uint32(b[24:28])
	cf.Cfn07MirrorLockup = binary.LittleEndian.Uint32(b[28:32])
	cf.Cfn08ShutterCurtainSync = binary.LittleEndian.Uint32(b[32:36])
	cf.Cfn09LensAFSearch = binary.LittleEndian.Uint32(b[36:40])
	return nil
}

// CustomFunctions400D represents the custom functions for the EOS 400D.
// Based on Phil Harvey's exiftool.
// Reference: https://exiftool.org/TagNames/CanonCustom.html#Functions400D
type CustomFunctions400D struct {
	Cfn01SetButtonFunc         uint32 `json:"cfn01SetButtonFunc"`         // 1
	Cfn02LongExpNoiseReduction uint32 `json:"cfn02LongExpNoiseReduction"` // 2
	Cfn03FlashSyncSpeedAv      uint32 `json:"cfn03FlashSyncSpeedAv"`      // 3
	Cfn04ShutterAELock         uint32 `json:"cfn04ShutterAELock"`         // 4
	Cfn05AFAssistBeam          uint32 `json:"cfn05AFAssistBeam"`          // 5
	Cfn06ExposureLevel         uint32 `json:"cfn06ExposureLevel"`         // 6
	Cfn07MirrorLockup          uint32 `json:"cfn07MirrorLockup"`          // 7
	Cfn08ShutterCurtainSync    uint32 `json:"cfn08ShutterCurtainSync"`    // 8
	Cfn09LensAFSearch          uint32 `json:"cfn09LensAFSearch"`          // 9
	Cfn10ETTLII                uint32 `json:"cfn10ETTLII"`                // 10
	Cfn11FlashExposureCompLock uint32 `json:"cfn11FlashExposureCompLock"` // 11
}

// UnmarshalBinary unmarshals the binary data into a CustomFunctions400D struct.
func (cf *CustomFunctions400D) UnmarshalBinary(b []byte) error {
	if len(b) < 48 { // 12 functions * 4 bytes
		return fmt.Errorf("invalid data length for CustomFunctions400D: got %d, want at least 48", len(b))
	}
	cf.Cfn01SetButtonFunc = binary.LittleEndian.Uint32(b[4:8])
	cf.Cfn02LongExpNoiseReduction = binary.LittleEndian.Uint32(b[8:12])
	cf.Cfn03FlashSyncSpeedAv = binary.LittleEndian.Uint32(b[12:16])
	cf.Cfn04ShutterAELock = binary.LittleEndian.Uint32(b[16:20])
	cf.Cfn05AFAssistBeam = binary.LittleEndian.Uint32(b[20:24])
	cf.Cfn06ExposureLevel = binary.LittleEndian.Uint32(b[24:28])
	cf.Cfn07MirrorLockup = binary.LittleEndian.Uint32(b[28:32])
	cf.Cfn08ShutterCurtainSync = binary.LittleEndian.Uint32(b[32:36])
	cf.Cfn09LensAFSearch = binary.LittleEndian.Uint32(b[36:40])
	cf.Cfn10ETTLII = binary.LittleEndian.Uint32(b[40:44])
	cf.Cfn11FlashExposureCompLock = binary.LittleEndian.Uint32(b[44:48])
	return nil
}

// CustomFunctionsD30 represents the custom functions for the EOS D30 and D60.
// Based on Phil Harvey's exiftool.
// Reference: https://exiftool.org/TagNames/CanonCustom.html#FunctionsD30
type CustomFunctionsD30 struct {
	Cfn01SetButtonFunc       uint32 `json:"cfn01SetButtonFunc"`       // 1
	Cfn02ShutterRelease      uint32 `json:"cfn02ShutterRelease"`      // 2
	Cfn03MirrorLockup        uint32 `json:"cfn03MirrorLockup"`        // 3
	Cfn04ShutterAELock       uint32 `json:"cfn04ShutterAELock"`       // 4
	Cfn05AFAssistBeam        uint32 `json:"cfn05AFAssistBeam"`        // 5
	Cfn06FlashSyncSpeedAv    uint32 `json:"cfn06FlashSyncSpeedAv"`    // 6
	Cfn07AELockButton        uint32 `json:"cfn07AELockButton"`        // 7
	Cfn08ShutterCurtainSync  uint32 `json:"cfn08ShutterCurtainSync"`  // 8
	Cfn09LensAFSearch        uint32 `json:"cfn09LensAFSearch"`        // 9
	Cfn10EraseButton         uint32 `json:"cfn10EraseButton"`         // 10
	Cfn11MenuButton          uint32 `json:"cfn11MenuButton"`          // 11
	Cfn12SuperimposedDisplay uint32 `json:"cfn12SuperimposedDisplay"` // 12
}

// UnmarshalBinary unmarshals the binary data into a CustomFunctionsD30 struct.
func (cf *CustomFunctionsD30) UnmarshalBinary(b []byte) error {
	if len(b) < 52 { // 13 functions * 4 bytes
		return fmt.Errorf("invalid data length for CustomFunctionsD30: got %d, want at least 52", len(b))
	}
	cf.Cfn01SetButtonFunc = binary.LittleEndian.Uint32(b[4:8])
	cf.Cfn02ShutterRelease = binary.LittleEndian.Uint32(b[8:12])
	cf.Cfn03MirrorLockup = binary.LittleEndian.Uint32(b[12:16])
	cf.Cfn04ShutterAELock = binary.LittleEndian.Uint32(b[16:20])
	cf.Cfn05AFAssistBeam = binary.LittleEndian.Uint32(b[20:24])
	cf.Cfn06FlashSyncSpeedAv = binary.LittleEndian.Uint32(b[24:28])
	cf.Cfn07AELockButton = binary.LittleEndian.Uint32(b[28:32])
	cf.Cfn08ShutterCurtainSync = binary.LittleEndian.Uint32(b[32:36])
	cf.Cfn09LensAFSearch = binary.LittleEndian.Uint32(b[36:40])
	cf.Cfn10EraseButton = binary.LittleEndian.Uint32(b[40:44])
	cf.Cfn11MenuButton = binary.LittleEndian.Uint32(b[44:48])
	cf.Cfn12SuperimposedDisplay = binary.LittleEndian.Uint32(b[48:52])
	return nil
}

// CustomFunctions is an interface for Canon custom functions.
type CustomFunctions interface {
	UnmarshalBinary(b []byte) error
}

// UnmarshalCustomFunctions unmarshals the CanonCustomFunctions tag by camera model.
func UnmarshalCustomFunctions(modelID CanonModelID, b []byte) (cf CustomFunctions, err error) {
	switch modelID {
	case ModelEOS1D, ModelEOS1DS, ModelEOS1DMarkII, ModelEOS1DSMarkII, ModelEOS1DMarkIIN:
		cf = &CustomFunctions1D{}
	case ModelEOS5D:
		cf = &CustomFunctions5D{}
	case ModelEOS10D:
		cf = &CustomFunctions10D{}
	case ModelEOS20D:
		cf = &CustomFunctions20D{}
	case ModelEOS30D:
		cf = &CustomFunctions30D{}
	case ModelEOS350D: // ModelEOSRebelXT, ModelEOSKissNDigital
		cf = &CustomFunctions350D{}
	case ModelEOS400D: // ModelEOSRebelXTi, ModelEOSKissXDigital
		cf = &CustomFunctions400D{}
	case ModelEOSD30:
		cf = &CustomFunctionsD30{}
	case ModelEOSD60:
		cf = &CustomFunctionsD30{} // D60 uses D30 functions
	default:
		return nil, fmt.Errorf("unsupported custom functions for model %s", modelID)
	}

	if unmarshaler, ok := cf.(CustomFunctions); ok {
		err = unmarshaler.UnmarshalBinary(b)
	} else {
		err = fmt.Errorf("unsupported custom functions type for model %s", modelID)
	}

	return
}
