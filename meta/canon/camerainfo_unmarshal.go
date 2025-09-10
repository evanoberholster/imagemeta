package canon

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
)

// UnmarshalCameraInfo unmarshals the CanonCameraInfo tag by camera model. It acts as a factory,
// creating the appropriate camera-specific struct and then calling its UnmarshalBinary method.
func UnmarshalCameraInfo(modelID CanonModelID, b []byte) (ci CanonCameraInfo, err error) {
	switch modelID {
	case ModelEOS1D, ModelEOS1DS:
		ci = &CameraInfo1D{}
	case ModelEOS1DMarkII, ModelEOS1DSMarkII:
		ci = &CameraInfo1DmkII{}
	case ModelEOS1DMarkIIN:
		ci = &CameraInfo1DmkIIN{}
	case ModelEOS1DMarkIII, ModelEOS1DSMarkIII:
		ci = &CameraInfo1DmkIII{}
	case ModelEOS1DMarkIV:
		ci = &CameraInfo1DmkIV{}
	case ModelEOS1DX:
		ci = &CameraInfo1DX{}
	case ModelEOS5D:
		ci = &CameraInfo5D{}
	case ModelEOS5DMarkII:
		ci = &CameraInfo5DmkII{}
	case ModelEOS5DMarkIII:
		ci = &CameraInfo5DmkIII{}
	case ModelEOS6D:
		ci = &CameraInfo6D{}
	case ModelEOS7D:
		ci = &CameraInfo7D{}
	case ModelEOS40D:
		ci = &CameraInfo40D{}
	case ModelEOS50D:
		ci = &CameraInfo50D{}
	case ModelEOS60D, ModelEOS1200D:
		ci = &CameraInfo60D{}
	case ModelEOS70D:
		ci = &CameraInfo70D{}
	case ModelEOS80D:
		ci = &CameraInfo80D{}
	case ModelEOS450D:
		ci = &CameraInfo450D{}
	case ModelEOS1000D:
		ci = &CameraInfo1000D{}
	case ModelEOS500D:
		ci = &CameraInfo500D{}
	case ModelEOS550D:
		ci = &CameraInfo550D{}
	case ModelEOS600D, ModelEOS1100D:
		ci = &CameraInfo600D{}
	case ModelEOS650D, ModelEOS700D:
		ci = &CameraInfo650D{}
	case ModelEOS750D, ModelEOS760D:
		ci = &CameraInfo750D{}
	case ModelEOSR5, ModelEOSR6:
		ci = &CameraInfoR6{}
	case ModelEOSR6MarkII, ModelEOSR8, ModelEOSR50:
		ci = &CameraInfoR6m2{}
	case ModelPowerShotG5XMark2:
		ci = &CameraInfoG5XII{}
	default:
		ci = &CameraInfoUnknown{}
	}

	if unmarshaler, ok := ci.(CanonCameraInfo); ok {
		err = unmarshaler.UnmarshalBinary(b)
	} else {
		err = fmt.Errorf("unsupported camera info type for model %d", modelID)
	}

	return
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler for CameraInfo1D.
func (ci *CameraInfo1D) UnmarshalBinary(b []byte) error {
	if len(b) < 82 {
		return fmt.Errorf("incorrect length for CameraInfo1D, should be at least 82 bytes: %d", len(b))
	}
	ci.ExposureTime = CIExposureTimeFromRaw(b[4])                           // 0x04
	ci.FocalLength = FocalLength(binary.LittleEndian.Uint16(b[10:12]))      // 0x0a
	ci.LensType = CanonLensType(binary.LittleEndian.Uint16(b[13:15]))       // 0x0d, little-endian
	ci.MinFocalLength = CIFocalLength(binary.LittleEndian.Uint16(b[14:16])) // 0x0e
	ci.MaxFocalLength = CIFocalLength(binary.LittleEndian.Uint16(b[16:18])) // 0x10

	// Offsets differ between EOS-1D and EOS-1DS
	if ci.cameraModelID == ModelEOS1D {
		ci.SharpnessFrequency = b[0x41]                                      // 0x41
		ci.Sharpness = int8(b[0x42])                                         // 0x42
		ci.WhiteBalance = WhiteBalance(binary.LittleEndian.Uint16(b[0x44:])) // 0x44
		ci.ColorTemperature = binary.LittleEndian.Uint16(b[0x48:])           // 0x48
		ci.PictureStyle = PictureStyle(binary.LittleEndian.Uint16(b[0x4b:])) // 0x4b
	} else { // EOS-1DS
		ci.SharpnessFrequency = b[0x47]                                      // 0x47
		ci.Sharpness = int8(b[0x48])                                         // 0x48
		ci.WhiteBalance = WhiteBalance(binary.LittleEndian.Uint16(b[0x4a:])) // 0x4a
		ci.ColorTemperature = binary.LittleEndian.Uint16(b[0x4e:])           // 0x4e
		ci.PictureStyle = PictureStyle(binary.LittleEndian.Uint16(b[0x51:])) // 0x51
	}
	return nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler for CameraInfo1DmkII.
func (ci *CameraInfo1DmkII) UnmarshalBinary(b []byte) error {
	if len(b) < 122 {
		return fmt.Errorf("incorrect length for CameraInfo1DmkII, should be at least 122 bytes: %d", len(b))
	}
	ci.ExposureTime = CIExposureTimeFromRaw(b[0x04])                     // 0x04
	ci.FocalLength = FocalLength(binary.BigEndian.Uint16(b[0x09:]))      // 0x09, Big-endian
	ci.LensType = CanonLensType(binary.BigEndian.Uint16(b[0x0c:]))       // 0x0c, Big-endian
	ci.MinFocalLength = CIFocalLength(binary.BigEndian.Uint16(b[0x11:])) // 0x11, Big-endian
	ci.MaxFocalLength = CIFocalLength(binary.BigEndian.Uint16(b[0x13:])) // 0x13, Big-endian
	ci.WhiteBalance = WhiteBalance(binary.LittleEndian.Uint16(b[54:56])) // 0x36
	ci.ColorTemperature = binary.BigEndian.Uint16(b[55:57])              // 0x37
	ci.CanonImageSize = binary.LittleEndian.Uint16(b[57:59])             // 0x39
	ci.JPEGQuality = uint16(b[102])                                      // 0x66
	ci.PictureStyle = PictureStyle(b[108])                               // 0x6c
	ci.Saturation = int8(b[110])                                         // 0x6e
	ci.ColorTone = int8(b[111])                                          // 0x6f
	ci.Sharpness = int8(b[114])                                          // 0x72
	ci.Contrast = int8(b[115])                                           // 0x73
	ci.ISO = string(bytes.Trim(b[117:122], "\x00"))                      // 0x75
	return nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler for CameraInfo1DmkIIN.
func (ci *CameraInfo1DmkIIN) UnmarshalBinary(b []byte) error {
	if len(b) < 126 {
		return fmt.Errorf("incorrect length for CameraInfo1DmkIIN, should be at least 126 bytes: %d", len(b))
	}
	ci.ExposureTime = CIExposureTimeFromRaw(b[0x04])                     // 0x04
	ci.FocalLength = FocalLength(binary.BigEndian.Uint16(b[0x09:]))      // 0x09, Big-endian
	ci.LensType = CanonLensType(binary.BigEndian.Uint16(b[0x0c:]))       // 0x0c, Big-endian
	ci.MinFocalLength = CIFocalLength(binary.BigEndian.Uint16(b[0x11:])) // 0x11, Big-endian
	ci.MaxFocalLength = CIFocalLength(binary.BigEndian.Uint16(b[0x13:])) // 0x13, Big-endian
	ci.WhiteBalance = WhiteBalance(binary.LittleEndian.Uint16(b[0x36:])) // 0x36
	ci.ColorTemperature = binary.BigEndian.Uint16(b[0x37:])              // 0x37, Big-endian
	ci.PictureStyle = PictureStyle(b[0x73])                              // 0x73
	ci.Sharpness = int8(b[0x74])                                         // 0x74
	ci.Contrast = int8(b[0x75])                                          // 0x75
	ci.Saturation = int8(b[0x76])                                        // 0x76
	ci.ColorTone = int8(b[0x77])                                         // 0x77
	ci.ISO = string(bytes.Trim(b[0x79:0x79+5], "\x00"))                  // 0x79
	return nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler for CameraInfo1DmkIII.
func (ci *CameraInfo1DmkIII) UnmarshalBinary(b []byte) error {
	if len(b) < 0x45e+4 {
		return fmt.Errorf("incorrect length for CameraInfo1DmkIII, should be at least %d bytes: %d", 0x45e+4, len(b))
	}
	ci.FNumber = FNumberFromRaw(b[0x03])                                   // 0x03
	ci.ExposureTime = CIExposureTimeFromRaw(b[0x04])                       // 0x04
	ci.ISO = ISOFromRaw(b[0x06])                                           // 0x06
	ci.CameraTemperature = CameraTemperatureFromRaw(b[0x18])               // 0x18
	ci.MacroMagnification = b[0x1b]                                        // 0x1b
	ci.FocalLength = FocalLength(binary.BigEndian.Uint16(b[0x1d:]))        // 0x1d
	ci.CameraOrientation = CameraOrientation(b[0x30])                      // 0x30
	ci.FocusDistanceUpper = FocalLength(binary.BigEndian.Uint16(b[0x43:])) // 0x43, Big-endian
	ci.FocusDistanceLower = FocalLength(binary.BigEndian.Uint16(b[0x45:])) // 0x45, Big-endian
	ci.WhiteBalance = WhiteBalance(binary.LittleEndian.Uint16(b[0x5e:]))   // 0x5e
	ci.ColorTemperature = binary.LittleEndian.Uint16(b[0x62:])             // 0x62
	ci.PictureStyle = PictureStyle(b[0x86])                                // 0x86
	ci.LensType = CanonLensType(binary.BigEndian.Uint16(b[0x111:]))        // 0x111
	ci.MinFocalLength = CIFocalLength(binary.BigEndian.Uint16(b[0x113:]))  // 0x113
	ci.MaxFocalLength = CIFocalLength(binary.BigEndian.Uint16(b[0x115:]))  // 0x115
	ci.FirmwareVersion = string(bytes.Trim(b[0x136:0x13c], "\x00"))        // 0x136
	ci.FileIndex = binary.LittleEndian.Uint32(b[0x172:]) + 1               // 0x172
	ci.ShutterCount = binary.LittleEndian.Uint32(b[0x176:]) + 1            // 0x176
	ci.DirectoryIndex = binary.LittleEndian.Uint32(b[0x17e:]) - 1          // 0x17e
	//ci.TimeStamp1 = binary.LittleEndian.Uint32(b[0x45a : 0x45a+4])         // 0x45a
	ci.TimeStamp = binary.LittleEndian.Uint32(b[0x45e : 0x45e+4]) // 0x45e
	// PictureStyleInfo at 0x2aa
	return nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler for CameraInfo1DmkIV.
func (ci *CameraInfo1DmkIV) UnmarshalBinary(b []byte) error {
	if len(b) < 873 {
		return fmt.Errorf("incorrect length for CameraInfo1DmkIV, should be at least 873 bytes: %d", len(b))
	}
	ci.FNumber = FNumberFromRaw(b[0x03])                                   // 0x03
	ci.ExposureTime = CIExposureTimeFromRaw(b[0x04])                       // 0x04
	ci.ISO = ISOFromRaw(b[0x06])                                           // 0x06
	ci.HighlightTonePriority = b[0x07]                                     // 0x07
	ci.MeasuredEV2 = int16(b[0x08])                                        // 0x08
	ci.MeasuredEV3 = int16(b[0x09])                                        // 0x09
	ci.FlashMeteringMode = FlashMeteringMode(b[0x15])                      // 0x15
	ci.CameraTemperature = CameraTemperatureFromRaw(b[0x19])               // 0x19
	ci.FocalLength = FocalLength(binary.BigEndian.Uint16(b[0x1e:]))        // 0x1e
	ci.CameraOrientation = CameraOrientation(b[0x35])                      // 0x35
	ci.FocusDistanceUpper = FocalLength(binary.BigEndian.Uint16(b[0x54:])) // 0x54
	ci.FocusDistanceLower = FocalLength(binary.BigEndian.Uint16(b[0x56:])) // 0x56
	ci.WhiteBalance = WhiteBalance(binary.LittleEndian.Uint16(b[0x78:]))   // 0x78
	ci.ColorTemperature = binary.LittleEndian.Uint16(b[0x7c:])             // 0x7c
	ci.CameraPictureStyle = PictureStyle(b[0xaf])                          // 0xaf
	ci.HighISONoiseReduction = b[0xc9]                                     // 0xc9
	ci.LensType = CanonLensType(binary.BigEndian.Uint16(b[0x14f:]))        // 0x14f
	ci.MinFocalLength = CIFocalLength(binary.BigEndian.Uint16(b[0x151:]))  // 0x151
	ci.MaxFocalLength = CIFocalLength(binary.BigEndian.Uint16(b[0x153:]))  // 0x153
	ci.FirmwareVersion = string(bytes.Trim(b[0x1ed:0x1ed+6], "\x00"))      // 0x1ed
	ci.FileIndex = binary.LittleEndian.Uint32(b[0x22c:]) + 1               // 0x22c
	ci.DirectoryIndex = binary.LittleEndian.Uint32(b[0x238:]) - 1          // 0x238
	// PictureStyleInfo at 0x368
	return nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler for CameraInfo1DX.
func (ci *CameraInfo1DX) UnmarshalBinary(b []byte) error {
	if len(b) < 0x3f4+4 {
		return fmt.Errorf("incorrect length for CameraInfo1DX, should be at least %d bytes: %d", 0x3f4+4, len(b))
	}
	ci.FNumber = FNumberFromRaw(b[0x03])                                   // 0x03
	ci.ExposureTime = CIExposureTimeFromRaw(b[0x04])                       // 0x04
	ci.ISO = ISOFromRaw(b[0x06])                                           // 0x06
	ci.CameraTemperature = CameraTemperatureFromRaw(b[0x1b])               // 0x1b
	ci.FocalLength = FocalLength(binary.BigEndian.Uint16(b[0x23:]))        // 0x23
	ci.CameraOrientation = CameraOrientation(b[0x7d])                      // 0x7d
	ci.FocusDistanceUpper = FocalLength(binary.BigEndian.Uint16(b[0x8c:])) // 0x8c
	ci.FocusDistanceLower = FocalLength(binary.BigEndian.Uint16(b[0x8e:])) // 0x8e
	ci.WhiteBalance = WhiteBalance(binary.LittleEndian.Uint16(b[0xbc:]))   // 0xbc
	ci.ColorTemperature = binary.LittleEndian.Uint16(b[0xc0:])             // 0xc0
	ci.PictureStyle = PictureStyle(b[0xf4])                                // 0xf4
	ci.LensType = CanonLensType(binary.BigEndian.Uint16(b[0x1a7:]))        // 0x1a7
	ci.MinFocalLength = CIFocalLength(binary.BigEndian.Uint16(b[0x1a9:]))  // 0x1a9
	ci.MaxFocalLength = CIFocalLength(binary.BigEndian.Uint16(b[0x1ab:]))  // 0x1ab
	ci.FirmwareVersion = string(bytes.Trim(b[0x280:0x286], "\x00"))        // 0x280
	ci.FileIndex = binary.LittleEndian.Uint32(b[0x2d0:]) + 1               // 0x2d0
	ci.DirectoryIndex = binary.LittleEndian.Uint32(b[0x2dc:]) - 1          // 0x2dc
	// PictureStyleInfo at 0x3f4
	return nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler for CameraInfo5DmkII.
func (ci *CameraInfo5DmkII) UnmarshalBinary(b []byte) error {
	if len(b) < 0x2f7+4 {
		return fmt.Errorf("incorrect length for CameraInfo5DmkII, should be at least %d bytes: %d", 0x2f7+4, len(b))
	}
	ci.FNumber = FNumberFromRaw(b[0x03])                                   // 0x03
	ci.ExposureTime = CIExposureTimeFromRaw(b[0x04])                       // 0x04
	ci.ISO = ISOFromRaw(b[0x06])                                           // 0x06
	ci.HighlightTonePriority = b[0x07]                                     // 0x07
	ci.FlashMeteringMode = FlashMeteringMode(b[0x15])                      // 0x15
	ci.CameraTemperature = CameraTemperatureFromRaw(b[0x19])               // 0x19
	ci.MacroMagnification = b[0x1b]                                        // 0x1b
	ci.FocalLength = FocalLength(binary.BigEndian.Uint16(b[0x1e:]))        // 0x1e
	ci.CameraOrientation = CameraOrientation(b[0x31])                      // 0x31
	ci.FocusDistanceUpper = FocalLength(binary.BigEndian.Uint16(b[0x50:])) // 0x50, odd-byte big-endian
	ci.FocusDistanceLower = FocalLength(binary.BigEndian.Uint16(b[0x52:])) // 0x52, odd-byte big-endian
	ci.WhiteBalance = WhiteBalance(binary.LittleEndian.Uint16(b[0x6f:]))   // 0x6f
	ci.ColorTemperature = binary.LittleEndian.Uint16(b[0x73:])             // 0x73
	ci.PictureStyle = PictureStyle(b[0xa7])                                // 0xa7
	ci.HighISONoiseReduction = b[0xbd]                                     // 0xbd
	ci.AutoLightingOptimizer = b[0xbf]                                     // 0xbf
	ci.LensType = CanonLensType(binary.BigEndian.Uint16(b[0xe6:]))         // 0xe6
	ci.MinFocalLength = CIFocalLength(binary.BigEndian.Uint16(b[0xe8:]))   // 0xe8
	ci.MaxFocalLength = CIFocalLength(binary.BigEndian.Uint16(b[0xea:]))   // 0xea
	ci.FirmwareVersion = string(bytes.Trim(b[0x17e:0x17e+6], "\x00"))      // 0x17e
	ci.FileIndex = binary.LittleEndian.Uint32(b[0x1bb:]) + 1               // 0x1bb
	ci.DirectoryIndex = binary.LittleEndian.Uint32(b[0x1c7:]) - 1          // 0x1c7
	return nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler for CameraInfo5DmkIII.
func (ci *CameraInfo5DmkIII) UnmarshalBinary(b []byte) error {
	if len(b) < 0x3b0+4 {
		return fmt.Errorf("incorrect length for CameraInfo5DmkIII, should be at least %d bytes: %d", 0x3b0+4, len(b))
	}
	ci.FNumber = FNumberFromRaw(b[0x03])                                   // 0x03
	ci.ExposureTime = CIExposureTimeFromRaw(b[0x04])                       // 0x04
	ci.ISO = ISOFromRaw(b[0x06])                                           // 0x06
	ci.CameraTemperature = CameraTemperatureFromRaw(b[0x1b])               // 0x1b
	ci.FocalLength = FocalLength(binary.BigEndian.Uint16(b[0x23:]))        // 0x23
	ci.CameraOrientation = CameraOrientation(b[0x7d])                      // 0x7d
	ci.FocusDistanceUpper = FocalLength(binary.BigEndian.Uint16(b[0x8c:])) // 0x8c
	ci.FocusDistanceLower = FocalLength(binary.BigEndian.Uint16(b[0x8e:])) // 0x8e
	ci.WhiteBalance = WhiteBalance(binary.LittleEndian.Uint16(b[0xbc:]))   // 0xbc
	ci.ColorTemperature = binary.LittleEndian.Uint16(b[0xc0:])             // 0xc0
	ci.PictureStyle = PictureStyle(b[0xf4])                                // 0xf4
	ci.LensType = CanonLensType(binary.BigEndian.Uint16(b[0x153:]))        // 0x153
	ci.MinFocalLength = CIFocalLength(binary.BigEndian.Uint16(b[0x155:]))  // 0x155
	ci.MaxFocalLength = CIFocalLength(binary.BigEndian.Uint16(b[0x157:]))  // 0x157
	ci.FirmwareVersion = string(bytes.Trim(b[0x23c:0x242], "\x00"))        // 0x23c
	ci.FileIndex = binary.LittleEndian.Uint32(b[0x28c:]) + 1               // 0x28c
	ci.DirectoryIndex = binary.LittleEndian.Uint32(b[0x298:]) - 1          // 0x298
	// PictureStyleInfo at 0x3b0
	return nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler for CameraInfo5D.
func (ci *CameraInfo5D) UnmarshalBinary(b []byte) error {
	if len(b) < 0x11c+4 {
		return fmt.Errorf("incorrect length for CameraInfo5D, should be at least %d bytes: %d", 0x11c+4, len(b))
	}
	ci.FNumber = FNumberFromRaw(b[0x03])                                        // 0x03
	ci.ExposureTime = CIExposureTimeFromRaw(b[0x04])                            // 0x04
	ci.ISO = ISOFromRaw(b[0x06])                                                // 0x06
	ci.LensType = CanonLensType(binary.BigEndian.Uint16(b[0x0c:]))              // 0x0c
	ci.CameraTemperature = CameraTemperatureFromRaw(b[0x17])                    // 0x17
	ci.MacroMagnification = b[0x1b]                                             // 0x1b
	ci.CameraOrientation = CameraOrientation(b[0x27])                           // 0x27
	ci.FocalLength = FocalLength(binary.BigEndian.Uint16(b[0x28:]))             // 0x28
	ci.AFPointsInFocus5D = AFPointsInFocus5D(binary.BigEndian.Uint16(b[0x38:])) // 0x38
	ci.WhiteBalance = WhiteBalance(binary.LittleEndian.Uint16(b[0x54:]))        // 0x54
	ci.ColorTemperature = binary.LittleEndian.Uint16(b[0x58:])                  // 0x58
	ci.PictureStyle = PictureStyle(b[0x6c])                                     // 0x6c
	ci.MinFocalLength = CIFocalLength(binary.BigEndian.Uint16(b[0x93:]))        // 0x93
	ci.MaxFocalLength = CIFocalLength(binary.BigEndian.Uint16(b[0x95:]))        // 0x95
	ci.FirmwareRevision = string(bytes.Trim(b[0xa4:0xac], "\x00"))              // 0xa4 - 0xab
	ci.ShortOwnerName = string(bytes.Trim(b[0xac:0xbc], "\x00"))                // 0xac - 0xbb
	ci.DirectoryIndex = binary.LittleEndian.Uint32(b[0xcc:])                    // 0xcc
	ci.FileIndex = binary.LittleEndian.Uint16(b[0xd0:]) + 1                     // 0xd0
	ci.TimeStamp = binary.LittleEndian.Uint32(b[0x11c:])                        // 0x11c
	// PictureStyleInfo5D at 0xe8
	return nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler for CameraInfo6D.
func (ci *CameraInfo6D) UnmarshalBinary(b []byte) error {
	if len(b) < 0x3c6+4 {
		return fmt.Errorf("incorrect length for CameraInfo6D, should be at least %d bytes: %d", 0x3c6+4, len(b))
	}
	ci.FNumber = FNumberFromRaw(b[0x03])                                   // 0x03
	ci.ExposureTime = CIExposureTimeFromRaw(b[0x04])                       // 0x04
	ci.ISO = ISOFromRaw(b[0x06])                                           // 0x06
	ci.CameraTemperature = CameraTemperatureFromRaw(b[0x1b])               // 0x1b
	ci.FocalLength = FocalLength(binary.BigEndian.Uint16(b[0x23:]))        // 0x23
	ci.CameraOrientation = CameraOrientation(b[0x83])                      // 0x83
	ci.FocusDistanceUpper = FocalLength(binary.BigEndian.Uint16(b[0x92:])) // 0x92
	ci.FocusDistanceLower = FocalLength(binary.BigEndian.Uint16(b[0x94:])) // 0x94
	ci.WhiteBalance = WhiteBalance(binary.LittleEndian.Uint16(b[0xc2:]))   // 0xc2
	ci.ColorTemperature = binary.LittleEndian.Uint16(b[0xc6:])             // 0xc6
	ci.PictureStyle = PictureStyle(b[0xfa])                                // 0xfa
	ci.LensType = CanonLensType(binary.BigEndian.Uint16(b[0x161:]))        // 0x161
	ci.MinFocalLength = CIFocalLength(binary.BigEndian.Uint16(b[0x163:]))  // 0x163
	ci.MaxFocalLength = CIFocalLength(binary.BigEndian.Uint16(b[0x165:]))  // 0x165
	ci.FirmwareVersion = string(bytes.Trim(b[0x256:0x25c], "\x00"))        // 0x256 - 0x25b
	ci.FileIndex = binary.LittleEndian.Uint32(b[0x2aa:]) + 1               // 0x2aa
	ci.DirectoryIndex = binary.LittleEndian.Uint32(b[0x2b6:]) - 1          // 0x2b6
	// PictureStyleInfo at 0x3c6
	return nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler for CameraInfo7D.
func (ci *CameraInfo7D) UnmarshalBinary(b []byte) error {
	if len(b) < 0x327+4 {
		return fmt.Errorf("incorrect length for CameraInfo7D, should be at least %d bytes: %d", 0x327+4, len(b))
	}
	ci.FNumber = FNumberFromRaw(b[0x03])                                   // 0x03
	ci.ExposureTime = CIExposureTimeFromRaw(b[0x04])                       // 0x04
	ci.ISO = ISOFromRaw(b[0x06])                                           // 0x06
	ci.HighlightTonePriority = b[0x07]                                     // 0x07
	ci.MeasuredEV2 = int16(b[0x08])                                        // 0x08
	ci.MeasuredEV = int16(b[0x09])                                         // 0x09
	ci.FlashMeteringMode = FlashMeteringMode(b[0x15])                      // 0x15
	ci.CameraTemperature = CameraTemperatureFromRaw(b[0x19])               // 0x19
	ci.FocalLength = FocalLength(binary.BigEndian.Uint16(b[0x1e:]))        // 0x1e
	ci.CameraOrientation = CameraOrientation(b[0x35])                      // 0x35
	ci.FocusDistanceUpper = FocalLength(binary.BigEndian.Uint16(b[0x54:])) // 0x54
	ci.FocusDistanceLower = FocalLength(binary.BigEndian.Uint16(b[0x56:])) // 0x56
	ci.WhiteBalance = WhiteBalance(binary.LittleEndian.Uint16(b[0x77:]))   // 0x77
	ci.ColorTemperature = binary.LittleEndian.Uint16(b[0x7b:])             // 0x7b
	ci.CameraPictureStyle = PictureStyle(b[0xaf])                          // 0xaf
	// HighISONoiseReduction at 0xc9 is not included in the struct
	ci.LensType = CanonLensType(binary.BigEndian.Uint16(b[0x112:]))       // 0x112
	ci.MinFocalLength = CIFocalLength(binary.BigEndian.Uint16(b[0x114:])) // 0x114
	ci.MaxFocalLength = CIFocalLength(binary.BigEndian.Uint16(b[0x116:])) // 0x116
	ci.FirmwareVersion = string(bytes.Trim(b[0x1ac:0x1b2], "\x00"))       // 0x1ac - 0x1b1
	ci.FileIndex = binary.LittleEndian.Uint32(b[0x1eb:]) + 1              // 0x1eb
	ci.DirectoryIndex = binary.LittleEndian.Uint32(b[0x1f7:]) - 1         // 0x1f7
	// PictureStyleInfo at 0x327
	return nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler for CameraInfo40D.
func (ci *CameraInfo40D) UnmarshalBinary(b []byte) error {
	if len(b) < 0x92b+64 {
		return fmt.Errorf("incorrect length for CameraInfo40D, should be at least %d bytes: %d", 0x92b+64, len(b))
	}
	ci.FNumber = FNumberFromRaw(b[0x03])                                   // 0x03
	ci.ExposureTime = CIExposureTimeFromRaw(b[0x04])                       // 0x04
	ci.ISO = ISOFromRaw(b[0x06])                                           // 0x06
	ci.FlashMeteringMode = FlashMeteringMode(b[0x15])                      // 0x15
	ci.CameraTemperature = CameraTemperatureFromRaw(b[0x18])               // 0x18
	ci.MacroMagnification = b[0x1b]                                        // 0x1b
	ci.FocalLength = FocalLength(binary.BigEndian.Uint16(b[0x1d:]))        // 0x1d
	ci.CameraOrientation = CameraOrientation(b[0x30])                      // 0x30
	ci.FocusDistanceUpper = FocalLength(binary.BigEndian.Uint16(b[0x43:])) // 0x43
	ci.FocusDistanceLower = FocalLength(binary.BigEndian.Uint16(b[0x45:])) // 0x45
	ci.WhiteBalance = WhiteBalance(binary.LittleEndian.Uint16(b[0x6f:]))   // 0x6f
	ci.ColorTemperature = binary.LittleEndian.Uint16(b[0x73:])             // 0x73
	ci.LensType = CanonLensType(binary.BigEndian.Uint16(b[0xd6:]))         // 0xd6
	ci.MinFocalLength = CIFocalLength(binary.BigEndian.Uint16(b[0xd8:]))   // 0xd8
	ci.MaxFocalLength = CIFocalLength(binary.BigEndian.Uint16(b[0xda:]))   // 0xda
	ci.FirmwareVersion = string(bytes.Trim(b[0xff:0x105], "\x00"))         // 0xff
	ci.FileIndex = binary.LittleEndian.Uint32(b[0x133:]) + 1               // 0x133
	ci.DirectoryIndex = binary.LittleEndian.Uint32(b[0x13f:]) - 1          // 0x13f
	ci.LensModel = string(bytes.Trim(b[0x92b:0x96b], "\x00"))              // 0x92b
	// PictureStyleInfo at 0x25b
	return nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler for CameraInfo450D.
func (ci *CameraInfo450D) UnmarshalBinary(b []byte) error {
	if len(b) < 0x933+64 {
		return fmt.Errorf("incorrect length for CameraInfo450D, should be at least %d bytes: %d", 0x933+64, len(b))
	}
	ci.FNumber = FNumberFromRaw(b[0x03])                                   // 0x03
	ci.ExposureTime = CIExposureTimeFromRaw(b[0x04])                       // 0x04
	ci.ISO = ISOFromRaw(b[0x06])                                           // 0x06
	ci.FlashMeteringMode = FlashMeteringMode(b[0x15])                      // 0x15
	ci.CameraTemperature = CameraTemperatureFromRaw(b[0x18])               // 0x18
	ci.MacroMagnification = b[0x1b]                                        // 0x1b
	ci.FocalLength = FocalLength(binary.BigEndian.Uint16(b[0x1d:]))        // 0x1d
	ci.CameraOrientation = CameraOrientation(b[0x30])                      // 0x30
	ci.FocusDistanceUpper = FocalLength(binary.BigEndian.Uint16(b[0x43:])) // 0x43
	ci.FocusDistanceLower = FocalLength(binary.BigEndian.Uint16(b[0x45:])) // 0x45
	ci.WhiteBalance = WhiteBalance(binary.LittleEndian.Uint16(b[0x6f:]))   // 0x6f
	ci.ColorTemperature = binary.LittleEndian.Uint16(b[0x73:])             // 0x73
	ci.LensType = CanonLensType(binary.BigEndian.Uint16(b[0xde:]))         // 0xde
	ci.FirmwareVersion = string(bytes.Trim(b[0x107:0x10d], "\x00"))        // 0x107
	ci.OwnerName = string(bytes.Trim(b[0x10f:0x12f], "\x00"))              // 0x10f
	ci.DirectoryIndex = binary.LittleEndian.Uint32(b[0x133:])              // 0x133
	ci.FileIndex = binary.LittleEndian.Uint32(b[0x13f:]) + 1               // 0x13f
	ci.LensModel = string(bytes.Trim(b[0x933:0x973], "\x00"))              // 0x933
	return nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler for CameraInfo50D.
func (ci *CameraInfo50D) UnmarshalBinary(b []byte) error {
	if len(b) < 0x2d7+4 {
		return fmt.Errorf("incorrect length for CameraInfo50D, should be at least %d bytes: %d", 0x2d7+4, len(b))
	}
	ci.FNumber = FNumberFromRaw(b[0x03])                                   // 0x03
	ci.ExposureTime = CIExposureTimeFromRaw(b[0x04])                       // 0x04
	ci.ISO = ISOFromRaw(b[0x06])                                           // 0x06
	ci.HighlightTonePriority = b[0x07]                                     // 0x07
	ci.FlashMeteringMode = FlashMeteringMode(b[0x15])                      // 0x15
	ci.CameraTemperature = CameraTemperatureFromRaw(b[0x19])               // 0x19
	ci.FocalLength = FocalLength(binary.BigEndian.Uint16(b[0x1e:]))        // 0x1e
	ci.CameraOrientation = CameraOrientation(b[0x31])                      // 0x31
	ci.FocusDistanceUpper = FocalLength(binary.BigEndian.Uint16(b[0x50:])) // 0x50
	ci.FocusDistanceLower = FocalLength(binary.BigEndian.Uint16(b[0x52:])) // 0x52
	ci.WhiteBalance = WhiteBalance(binary.LittleEndian.Uint16(b[0x6f:]))   // 0x6f
	ci.ColorTemperature = binary.LittleEndian.Uint16(b[0x73:])             // 0x73
	ci.PictureStyle = PictureStyle(b[0xa7])                                // 0xa7
	ci.HighISONoiseReduction = b[0xbd]                                     // 0xbd
	ci.AutoLightingOptimizer = b[0xbf]                                     // 0xbf
	ci.LensType = CanonLensType(binary.BigEndian.Uint16(b[0xea:]))         // 0xea
	ci.MinFocalLength = CIFocalLength(binary.BigEndian.Uint16(b[0xec:]))   // 0xec
	ci.MaxFocalLength = CIFocalLength(binary.BigEndian.Uint16(b[0xee:]))   // 0xee
	ci.FirmwareVersion = string(bytes.Trim(b[0x15e:0x164], "\x00"))        // 0x15e
	ci.FileIndex = binary.LittleEndian.Uint32(b[0x19b:]) + 1               // 0x19b
	ci.DirectoryIndex = binary.LittleEndian.Uint32(b[0x1a7:]) - 1          // 0x1a7
	// PictureStyleInfo at 0x2d7
	return nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler for CameraInfo60D.
func (ci *CameraInfo60D) UnmarshalBinary(b []byte) error {
	if len(b) < 0x321+4 {
		return fmt.Errorf("incorrect length for CameraInfo60D, should be at least %d bytes: %d", 0x321+4, len(b))
	}
	ci.FNumber = FNumberFromRaw(b[0x03])                                   // 0x03
	ci.ExposureTime = CIExposureTimeFromRaw(b[0x04])                       // 0x04
	ci.ISO = ISOFromRaw(b[0x06])                                           // 0x06
	ci.CameraTemperature = CameraTemperatureFromRaw(b[0x19])               // 0x19
	ci.FocalLength = FocalLength(binary.BigEndian.Uint16(b[0x1e:]))        // 0x1e
	ci.CameraOrientation = CameraOrientation(b[0x36])                      // 0x36 (60D) or 0x3a (1200D)
	ci.FocusDistanceUpper = FocalLength(binary.BigEndian.Uint16(b[0x55:])) // 0x55
	ci.FocusDistanceLower = FocalLength(binary.BigEndian.Uint16(b[0x57:])) // 0x57
	ci.ColorTemperature = binary.LittleEndian.Uint16(b[0x7d:])             // 0x7d
	ci.LensType = CanonLensType(binary.BigEndian.Uint16(b[0xe8:]))         // 0xe8
	ci.MinFocalLength = CIFocalLength(binary.BigEndian.Uint16(b[0xea:]))   // 0xea
	ci.MaxFocalLength = CIFocalLength(binary.BigEndian.Uint16(b[0xec:]))   // 0xec
	ci.FirmwareVersion = string(bytes.Trim(b[0x199:0x19f], "\x00"))        // 0x199
	ci.FileIndex = binary.LittleEndian.Uint32(b[0x1d9:]) + 1               // 0x1d9
	ci.DirectoryIndex = binary.LittleEndian.Uint32(b[0x1e5:]) - 1          // 0x1e5
	// PictureStyleInfo at 0x2f9 or 0x321
	return nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler for CameraInfo1000D.
func (ci *CameraInfo1000D) UnmarshalBinary(b []byte) error {
	if len(b) < 0x937+64 {
		return fmt.Errorf("incorrect length for CameraInfo1000D, should be at least %d bytes: %d", 0x937+64, len(b))
	}
	ci.FNumber = FNumberFromRaw(b[3])                                             // 0x03
	ci.ExposureTime = CIExposureTimeFromRaw(b[4])                                 // 0x04
	ci.ISO = ISOFromRaw(b[6])                                                     // 0x06
	ci.FlashMeteringMode = FlashMeteringMode(b[0x15])                             // 0x15
	ci.CameraTemperature = CameraTemperatureFromRaw(b[0x18])                      // 0x18
	ci.MacroMagnification = b[0x1b]                                               // 0x1b
	ci.FocalLength = FocalLength(binary.LittleEndian.Uint16(b[0x1d:0x1f]))        // 0x1d
	ci.CameraOrientation = CameraOrientation(b[0x30])                             // 0x30
	ci.FocusDistanceUpper = FocalLength(binary.LittleEndian.Uint16(b[0x43:0x45])) // 0x43, odd-byte little-endian
	ci.FocusDistanceLower = FocalLength(binary.LittleEndian.Uint16(b[0x45:0x47])) // 0x45, odd-byte little-endian
	ci.WhiteBalance = WhiteBalance(binary.LittleEndian.Uint16(b[0x6f:0x71]))      // 0x6f
	ci.ColorTemperature = binary.LittleEndian.Uint16(b[0x73:0x75])                // 0x73
	ci.LensType = CanonLensType(binary.BigEndian.Uint16(b[0xe2:0xe4]))            // 0xe2
	ci.MinFocalLength = CIFocalLength(binary.BigEndian.Uint16(b[0xe4:0xe6]))      // 0xe4
	ci.MaxFocalLength = CIFocalLength(binary.BigEndian.Uint16(b[0xe6:0xe8]))      // 0xe6
	ci.FirmwareVersion = string(bytes.Trim(b[0x10b:0x10b+6], "\x00"))             // 0x10b
	ci.DirectoryIndex = binary.LittleEndian.Uint32(b[0x137:0x13b])                // 0x137
	ci.FileIndex = binary.LittleEndian.Uint32(b[0x143:0x147]) + 1                 // 0x143
	ci.LensModel = string(bytes.Trim(b[0x937:0x937+64], "\x00"))                  // 0x937
	return nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler for CameraInfo500D.
func (ci *CameraInfo500D) UnmarshalBinary(b []byte) error {
	if len(b) < 0x30b+4 {
		return fmt.Errorf("incorrect length for CameraInfo500D, should be at least %d bytes: %d", 0x30b+4, len(b))
	}
	ci.FNumber = FNumberFromRaw(b[0x03])                                   // 0x03
	ci.ExposureTime = CIExposureTimeFromRaw(b[0x04])                       // 0x04
	ci.ISO = ISOFromRaw(b[0x06])                                           // 0x06
	ci.HighlightTonePriority = b[0x07]                                     // 0x07
	ci.FlashMeteringMode = FlashMeteringMode(b[0x15])                      // 0x15
	ci.CameraTemperature = CameraTemperatureFromRaw(b[0x19])               // 0x19
	ci.FocalLength = FocalLength(binary.BigEndian.Uint16(b[0x1e:]))        // 0x1e
	ci.CameraOrientation = CameraOrientation(b[0x31])                      // 0x31
	ci.FocusDistanceUpper = FocalLength(binary.BigEndian.Uint16(b[0x50:])) // 0x50
	ci.FocusDistanceLower = FocalLength(binary.BigEndian.Uint16(b[0x52:])) // 0x52
	ci.WhiteBalance = WhiteBalance(binary.LittleEndian.Uint16(b[0x73:]))   // 0x73
	ci.ColorTemperature = binary.LittleEndian.Uint16(b[0x77:])             // 0x77
	ci.PictureStyle = PictureStyle(b[0xab])                                // 0xab
	ci.HighISONoiseReduction = b[0xbc]                                     // 0xbc
	ci.AutoLightingOptimizer = b[0xbe]                                     // 0xbe
	ci.LensType = CanonLensType(binary.BigEndian.Uint16(b[0xf6:]))         // 0xf6
	ci.MinFocalLength = CIFocalLength(binary.BigEndian.Uint16(b[0xf8:]))   // 0xf8
	ci.MaxFocalLength = CIFocalLength(binary.BigEndian.Uint16(b[0xfa:]))   // 0xfa
	ci.FirmwareVersion = string(bytes.Trim(b[0x190:0x196], "\x00"))        // 0x190
	ci.FileIndex = binary.LittleEndian.Uint32(b[0x1d3:]) + 1               // 0x1d3
	ci.DirectoryIndex = binary.LittleEndian.Uint32(b[0x1df:]) - 1          // 0x1df
	// PictureStyleInfo at 0x30b
	return nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler for CameraInfo550D.
func (ci *CameraInfo550D) UnmarshalBinary(b []byte) error {
	if len(b) < 0x31c+4 {
		return fmt.Errorf("incorrect length for CameraInfo550D, should be at least %d bytes: %d", 0x31c+4, len(b))
	}
	ci.FNumber = FNumberFromRaw(b[3])                                             // 0x03
	ci.ExposureTime = CIExposureTimeFromRaw(b[4])                                 // 0x04
	ci.ISO = ISOFromRaw(b[6])                                                     // 0x06
	ci.HighlightTonePriority = b[7]                                               // 0x07
	ci.FlashMeteringMode = FlashMeteringMode(b[0x15])                             // 0x15
	ci.CameraTemperature = CameraTemperatureFromRaw(b[0x19])                      // 0x19
	ci.FocalLength = FocalLength(binary.BigEndian.Uint16(b[0x1e:0x20]))           // 0x1e
	ci.CameraOrientation = CameraOrientation(b[0x35])                             // 0x35
	ci.FocusDistanceUpper = FocalLength(binary.LittleEndian.Uint16(b[0x54:0x56])) // 0x54, odd-byte little-endian
	ci.FocusDistanceLower = FocalLength(binary.LittleEndian.Uint16(b[0x56:0x58])) // 0x56, odd-byte little-endian
	ci.WhiteBalance = WhiteBalance(binary.LittleEndian.Uint16(b[0x78:0x7a]))      // 0x78
	ci.ColorTemperature = binary.LittleEndian.Uint16(b[0x7c:0x7e])                // 0x7c
	ci.LensType = CanonLensType(binary.BigEndian.Uint16(b[0xff:0x101]))           // 0xff
	ci.MinFocalLength = CIFocalLength(binary.BigEndian.Uint16(b[0x101:0x103]))    // 0x101
	ci.MaxFocalLength = CIFocalLength(binary.BigEndian.Uint16(b[0x103:0x105]))    // 0x103
	ci.FirmwareVersion = string(bytes.Trim(b[0x1a4:0x1a4+6], "\x00"))             // 0x1a4
	ci.FileIndex = binary.LittleEndian.Uint32(b[0x1e4:0x1e8]) + 1                 // 0x1e4
	ci.DirectoryIndex = binary.LittleEndian.Uint32(b[0x1f0:0x1f4]) - 1            // 0x1f0
	// PictureStyleInfo at 0x31c
	return nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler for CameraInfo600D.
func (ci *CameraInfo600D) UnmarshalBinary(b []byte) error {
	if len(b) < 0x2fb+4 {
		return fmt.Errorf("incorrect length for CameraInfo600D, should be at least %d bytes: %d", 0x2fb+4, len(b))
	}
	ci.FNumber = FNumberFromRaw(b[3])                                                 // 0x03
	ci.ExposureTime = CIExposureTimeFromRaw(b[4])                                     // 0x04
	ci.ISO = ISOFromRaw(b[6])                                                         // 0x06
	ci.HighlightTonePriority = b[7]                                                   // 0x07
	ci.FlashMeteringMode = FlashMeteringMode(b[0x15])                                 // 0x15
	ci.CameraTemperature = CameraTemperatureFromRaw(b[0x19])                          // 0x19
	ci.FocalLength = FocalLength(binary.BigEndian.Uint16(b[0x1e : 0x1e+2]))           // 0x1e
	ci.CameraOrientation = CameraOrientation(b[0x38])                                 // 0x38
	ci.FocusDistanceUpper = FocalLength(binary.LittleEndian.Uint16(b[0x57 : 0x57+2])) // 0x57, odd-byte little-endian
	ci.FocusDistanceLower = FocalLength(binary.LittleEndian.Uint16(b[0x59 : 0x59+2])) // 0x59, odd-byte little-endian
	ci.WhiteBalance = WhiteBalance(binary.LittleEndian.Uint16(b[0x7b : 0x7b+2]))      // 0x7b
	ci.ColorTemperature = binary.LittleEndian.Uint16(b[0x7f : 0x7f+2])                // 0x7f
	ci.PictureStyle = PictureStyle(b[0xb3])                                           // 0xb3
	ci.LensType = CanonLensType(binary.BigEndian.Uint16(b[0xea : 0xea+2]))            // 0xea
	ci.HighISONoiseReduction = b[0xbc]                                                // 0xbc
	ci.AutoLightingOptimizer = b[0xbe]                                                // 0xbe
	ci.MinFocalLength = CIFocalLength(binary.BigEndian.Uint16(b[0xec : 0xec+2]))      // 0xec
	ci.MaxFocalLength = CIFocalLength(binary.BigEndian.Uint16(b[0xee : 0xee+2]))      // 0xee
	ci.FirmwareVersion = string(bytes.Trim(b[0x19b:0x19b+6], "\x00"))                 // 0x19b
	ci.FileIndex = binary.LittleEndian.Uint32(b[0x1db:0x1db+4]) + 1                   // 0x1db
	ci.DirectoryIndex = binary.LittleEndian.Uint32(b[0x1e7:0x1e7+4]) - 1              // 0x1e7
	return nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler for CameraInfo650D.
func (ci *CameraInfo650D) UnmarshalBinary(b []byte) error {
	if len(b) < 0x390+4 {
		return fmt.Errorf("incorrect length for CameraInfo650D, should be at least %d bytes: %d", 0x390+4, len(b))
	}
	ci.FNumber = FNumberFromRaw(b[0x03])                                   // 0x03
	ci.ExposureTime = CIExposureTimeFromRaw(b[0x04])                       // 0x04
	ci.ISO = ISOFromRaw(b[0x06])                                           // 0x06
	ci.CameraTemperature = CameraTemperatureFromRaw(b[0x1b])               // 0x1b
	ci.FocalLength = FocalLength(binary.BigEndian.Uint16(b[0x23:]))        // 0x23
	ci.CameraOrientation = CameraOrientation(b[0x7d])                      // 0x7d
	ci.FocusDistanceUpper = FocalLength(binary.BigEndian.Uint16(b[0x8c:])) // 0x8c
	ci.FocusDistanceLower = FocalLength(binary.BigEndian.Uint16(b[0x8e:])) // 0x8e
	ci.WhiteBalance = WhiteBalance(binary.LittleEndian.Uint16(b[0xbc:]))   // 0xbc
	ci.ColorTemperature = binary.LittleEndian.Uint16(b[0xc0:])             // 0xc0
	ci.PictureStyle = PictureStyle(b[0xf4])                                // 0xf4
	ci.LensType = CanonLensType(binary.BigEndian.Uint16(b[0x127:]))        // 0x127
	ci.MinFocalLength = CIFocalLength(binary.BigEndian.Uint16(b[0x129:]))  // 0x129
	ci.MaxFocalLength = CIFocalLength(binary.BigEndian.Uint16(b[0x12b:]))  // 0x12b
	ci.FirmwareVersion = string(bytes.Trim(b[0x21b:0x221], "\x00"))        // 0x21b
	ci.FileIndex = binary.LittleEndian.Uint32(b[0x270:]) + 1               // 0x270
	ci.DirectoryIndex = binary.LittleEndian.Uint32(b[0x27c:]) - 1          // 0x27c
	// PictureStyleInfo at 0x390
	return nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler for CameraInfo700D.
func (ci *CameraInfo700D) UnmarshalBinary(b []byte) error {
	if len(b) < 0x390+4 {
		return fmt.Errorf("incorrect length for CameraInfo700D, should be at least %d bytes: %d", 0x390+4, len(b))
	}
	ci.FNumber = FNumberFromRaw(b[3])
	ci.ExposureTime = CIExposureTimeFromRaw(b[4])
	ci.ISO = ISOFromRaw(b[6])
	ci.CameraTemperature = CameraTemperatureFromRaw(b[0x1b])
	ci.FocalLength = FocalLength(binary.BigEndian.Uint16(b[0x23 : 0x23+2]))
	ci.CameraOrientation = CameraOrientation(b[0x7d])
	ci.FocusDistanceUpper = FocalLength(binary.BigEndian.Uint16(b[0x8c : 0x8c+2]))
	ci.FocusDistanceLower = FocalLength(binary.BigEndian.Uint16(b[0x8e : 0x8e+2]))
	ci.WhiteBalance = WhiteBalance(binary.LittleEndian.Uint16(b[0xbc : 0xbc+2]))
	ci.ColorTemperature = binary.LittleEndian.Uint16(b[0xc0 : 0xc0+2])
	ci.PictureStyle = PictureStyle(b[0xf4])
	ci.LensType = CanonLensType(binary.BigEndian.Uint16(b[0x127 : 0x127+2]))
	ci.MinFocalLength = CIFocalLength(binary.BigEndian.Uint16(b[0x129 : 0x129+2]))
	ci.MaxFocalLength = CIFocalLength(binary.BigEndian.Uint16(b[0x12b : 0x12b+2]))
	ci.FirmwareVersion = string(bytes.Trim(b[0x220:0x220+6], "\x00"))
	ci.FileIndex = binary.LittleEndian.Uint32(b[0x274:0x274+4]) + 1
	ci.DirectoryIndex = binary.LittleEndian.Uint32(b[0x280:0x280+4]) - 1
	// PictureStyleInfo at 0x390
	return nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler for CameraInfo750D.
func (ci *CameraInfo750D) UnmarshalBinary(b []byte) error {
	// The structure for 750D/760D is complex and appears to vary.
	// Awaiting more sample data for a stable implementation.
	return nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler for CameraInfo70D.
func (ci *CameraInfo70D) UnmarshalBinary(b []byte) error {
	if len(b) < 0x3cf+4 {
		return fmt.Errorf("incorrect length for CameraInfo70D, should be at least %d bytes: %d", 0x3cf+4, len(b))
	}
	ci.FNumber = FNumberFromRaw(b[3])
	ci.ExposureTime = CIExposureTimeFromRaw(b[4])
	ci.ISO = ISOFromRaw(b[6])
	ci.CameraTemperature = CameraTemperatureFromRaw(b[0x1b])
	ci.FocalLength = FocalLength(binary.BigEndian.Uint16(b[0x23 : 0x23+2]))
	ci.CameraOrientation = CameraOrientation(b[0x84])
	ci.FocusDistanceUpper = FocalLength(binary.BigEndian.Uint16(b[0x93 : 0x93+2]))
	ci.FocusDistanceLower = FocalLength(binary.BigEndian.Uint16(b[0x95 : 0x95+2]))
	ci.ColorTemperature = binary.LittleEndian.Uint16(b[0xc7 : 0xc7+2])
	ci.LensType = CanonLensType(binary.BigEndian.Uint16(b[0x166 : 0x166+2]))
	ci.MinFocalLength = CIFocalLength(binary.BigEndian.Uint16(b[0x168 : 0x168+2]))
	ci.MaxFocalLength = CIFocalLength(binary.BigEndian.Uint16(b[0x16a : 0x16a+2]))
	ci.FirmwareVersion = string(bytes.Trim(b[0x25e:0x25e+6], "\x00"))
	ci.FileIndex = binary.LittleEndian.Uint32(b[0x2b3:0x2b3+4]) + 1
	ci.DirectoryIndex = binary.LittleEndian.Uint32(b[0x2bf:0x2bf+4]) - 1
	// PictureStyleInfo at 0x3cf
	return nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler for CameraInfo80D.
func (ci *CameraInfo80D) UnmarshalBinary(b []byte) error {
	if len(b) < 0x45a+6 {
		return fmt.Errorf("incorrect length for CameraInfo80D, should be at least %d bytes: %d", 0x45a+6, len(b))
	}
	ci.FNumber = FNumberFromRaw(b[3])
	ci.ExposureTime = CIExposureTimeFromRaw(b[4])
	ci.ISO = ISOFromRaw(b[6])
	ci.CameraTemperature = CameraTemperatureFromRaw(b[0x1b])
	ci.FocalLength = FocalLength(binary.BigEndian.Uint16(b[0x23 : 0x23+2]))
	ci.CameraOrientation = CameraOrientation(b[0x96])
	ci.FocusDistanceUpper = FocalLength(binary.BigEndian.Uint16(b[0xa5 : 0xa5+2]))
	ci.FocusDistanceLower = FocalLength(binary.BigEndian.Uint16(b[0xa7 : 0xa7+2]))
	ci.ColorTemperature = binary.LittleEndian.Uint16(b[0x13a : 0x13a+2])
	ci.LensType = CanonLensType(binary.BigEndian.Uint16(b[0x189 : 0x189+2]))
	ci.MinFocalLength = CIFocalLength(binary.BigEndian.Uint16(b[0x18b : 0x18b+2]))
	ci.MaxFocalLength = CIFocalLength(binary.BigEndian.Uint16(b[0x18d : 0x18d+2]))
	ci.FirmwareVersion = string(bytes.Trim(b[0x45a:0x45a+6], "\x00"))
	ci.FileIndex = binary.LittleEndian.Uint32(b[0x4ae:0x4ae+4]) + 1
	ci.DirectoryIndex = binary.LittleEndian.Uint32(b[0x4ba:0x4ba+4]) - 1
	return nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler for CameraInfoR6.
func (ci *CameraInfoR6) UnmarshalBinary(b []byte) error {
	if len(b) < 0x0af1+4 {
		return fmt.Errorf("incorrect length for CameraInfoR6, should be at least %d bytes: %d", 0x0af1+4, len(b))
	}
	ci.ShutterCount = binary.LittleEndian.Uint32(b[0x0af1 : 0x0af1+4]) // 0x0af1
	return nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler for CameraInfoR6m2.
func (ci *CameraInfoR6m2) UnmarshalBinary(b []byte) error {
	if len(b) < 0x0d29+4 {
		return fmt.Errorf("incorrect length for CameraInfoR6m2, should be at least %d bytes: %d", 0x0d29+4, len(b))
	}
	ci.ShutterCount = binary.LittleEndian.Uint32(b[0x0d29 : 0x0d29+4]) // 0x0d29
	return nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler for CameraInfoG5XII.
func (ci *CameraInfoG5XII) UnmarshalBinary(b []byte) error {
	// The offsets for G5XII depend on the file type (JPEG or CR3).
	// We check for the CR3 offset first as it's larger.
	if len(b) >= 0x0a95+4 {
		ci.ShutterCount = binary.LittleEndian.Uint32(b[0x0a95 : 0x0a95+4]) // 0x0a95 (CR3)
	}

	// JPEG specific offsets
	if len(b) >= 0x0b2d+4 {
		// Check JPEG ShutterCount offset if CR3 was not available or zero
		if ci.ShutterCount == 0 {
			ci.ShutterCount = binary.LittleEndian.Uint32(b[0x0293 : 0x0293+4]) // 0x0293 (JPEG)
		}
		ci.DirectoryIndex = binary.LittleEndian.Uint32(b[0x0b21 : 0x0b21+4]) // 0x0b21 (JPEG)
		ci.FileIndex = binary.LittleEndian.Uint32(b[0x0b2d:0x0b2d+4]) + 1    // 0x0b2d (JPEG)
	} else if ci.ShutterCount == 0 { // If not long enough for JPEG offsets, check just the JPEG shutter count offset
		if len(b) >= 0x0293+4 {
			ci.ShutterCount = binary.LittleEndian.Uint32(b[0x0293 : 0x0293+4]) // 0x0293 (JPEG)
		}
	}
	return nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler for CameraInfoPowerShot.
func (ci *CameraInfoPowerShot) UnmarshalBinary(b []byte) error {
	if len(b) < 138*4 { // Minimum size for this struct type
		return fmt.Errorf("incorrect length for CameraInfoPowerShot, should be at least %d bytes: %d", 138*4, len(b))
	}
	ci.ISO = ISOFromRaw(uint8(100 * math.Exp((float64(int32(binary.LittleEndian.Uint32(b[0x00*4:0x00*4+4]))-411)/96.0)*math.Log(2))))
	ci.FNumber = FNumber(math.Exp(float64(int32(binary.LittleEndian.Uint32(b[0x05*4:0x05*4+4]))) / 192.0 * math.Log(2)))
	ci.ExposureTime = CIExposureTime(math.Exp(float64(int32(binary.LittleEndian.Uint32(b[0x06*4:0x06*4+4]))*-1) / 96.0 * math.Log(2)))
	ci.Rotation = binary.LittleEndian.Uint32(b[0x17*4 : 0x17*4+4])
	ci.CameraTemperature = int32(binary.LittleEndian.Uint32(b[len(b)-12 : len(b)-8]))
	return nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler for CameraInfoPowerShot2.
func (ci *CameraInfoPowerShot2) UnmarshalBinary(b []byte) error {
	if len(b) < 156*4 { // Minimum size for this struct type
		return fmt.Errorf("incorrect length for CameraInfoPowerShot2: %d", len(b))
	}
	ci.ISO = ISOFromRaw(uint8(100 * math.Exp((float64(int32(binary.LittleEndian.Uint32(b[0x01*4:0x01*4+4]))-411)/96.0)*math.Log(2))))
	ci.FNumber = FNumber(math.Exp(float64(int32(binary.LittleEndian.Uint32(b[0x06*4:0x06*4+4]))) / 192.0 * math.Log(2)))
	ci.ExposureTime = CIExposureTime(math.Exp(float64(int32(binary.LittleEndian.Uint32(b[0x07*4:0x07*4+4]))*-1) / 96.0 * math.Log(2)))
	ci.Rotation = binary.LittleEndian.Uint32(b[0x18*4 : 0x18*4+4])
	ci.CameraTemperature = CameraTemperatureFromRaw(uint8(binary.LittleEndian.Uint32(b[len(b)-12 : len(b)-8])))
	return nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler for CameraInfoUnknown32.
func (ci *CameraInfoUnknown32) UnmarshalBinary(b []byte) error {
	// Add unmarshaling logic if any specific fields are known for this format.
	return nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler for CameraInfoUnknown.
func (ci *CameraInfoUnknown) UnmarshalBinary(b []byte) error {
	// Add unmarshaling logic if any specific fields are known for this format.
	return nil
}
