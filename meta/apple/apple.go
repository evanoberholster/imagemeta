// Package apple provides data types and functions for representing Apple Exif and Makernote values
package apple

import "fmt"

// AppleMakerNote is Apple Makernote data
// Based on Phil Harvey's exiftool
// Reference:
// - https://exiftool.org/TagNames/Apple.html
// Note: This is a partial implementation of AppleMakerNote
type AppleMakerNote struct {
	MakerNoteVersion    uint16             `json:"makerNoteVersion"`    // [0x0001]
	AEStable            AFState            `json:"aeStable"`            // [0x0004]
	AETarget            int32              `json:"aeTarget"`            // [0x0005]
	AEAverage           int32              `json:"aeAverage"`           // [0x0006]
	AEStable2           AFState            `json:"aeStable2"`           // [0x0007]
	AccelerationVector  AccelerationVector `json:"accelerationVector"`  // [0x0008]
	HDRImageType        HDRImageType       `json:"hdrImageType"`        // [0x000a]
	BurstUUID           string             `json:"burstUUID"`           // [0x000b]
	FocusDistanceRange  FocusDistanceRange `json:"focusDistanceRange"`  // [0x000c]
	OISMOde             int32              `json:"oisMode"`             // [0x000f]
	ContentIdentifier   string             `json:"contentIdentifier"`   // [0x0011] // MediaGroupUUID
	ImageCaptureType    ImageCaptureType   `json:"imageCaptureType"`    // [0x0012]
	ImageUniqueID       string             `json:"imageUniqueID"`       // [0x0015] // ImageGroupIdentifier
	BracketedCaptureSeq uint32             `json:"bracketedCaptureSeq"` // [0x001c] // (BracketedCaptureSequenceNumber, ref 2)
	PhotosAppFeatFlags  int32              `json:"photosAppFeatFlags"`  // [0x001f] // 'set if person or pet detected in image'
	ImageCaptureReqID   string             `json:"imageCaptureReqID"`   // [0x0020]
	PhotoIdentifier     string             `json:"photoIdentifier"`     // [0x002b]
	ColorTemperature    int32              `json:"colorTemperature"`    // [0x002d]
	CameraType          CameraType         `json:"cameraType"`          // [0x002e]
}

// ImageCaptureType is Apple Image Capture Type
type ImageCaptureType int

const (
	ImageCaptureUnknown ImageCaptureType = iota
	ImageCaptureProRAW
	ImageCapturePortrait
	ImageCapturePhoto       ImageCaptureType = 10
	ImageCaptureManualFocus ImageCaptureType = 11
	ImageCaptureScene       ImageCaptureType = 12
)

// String returns a string representation of ImageCaptureType
func (i ImageCaptureType) String() string {
	switch i {
	case ImageCaptureProRAW:
		return "ProRAW"
	case ImageCapturePortrait:
		return "Portrait"
	case ImageCapturePhoto:
		return "Photo"
	case ImageCaptureManualFocus:
		return "Manual Focus"
	case ImageCaptureScene:
		return "Scene"
	default:
		return "Unknown"
	}
}

// CameraType is Apple Camera Type
type CameraType int

const (
	CameraBackWideAngle CameraType = iota
	CameraBackNormal
	CameraFront CameraType = 6
)

// String returns a string representation of CameraType
func (c CameraType) String() string {
	switch c {
	case CameraBackWideAngle:
		return "Back Wide Angle"
	case CameraBackNormal:
		return "Back Normal"
	case CameraFront:
		return "Front"
	default:
		return "Unknown"
	}
}

// AFState represents whether auto-focus is stable
// 0 = Not Stable, 1 = Stable
type AFState bool

// String returns a string representation of AFState
func (a AFState) String() string {
	if a {
		return "Yes"
	}
	return "No"
}

// AFStateFromRAW returns an AFState from a raw value
func AFStateFromRAW(raw int32) AFState {
	return raw == 1
}

// AccelerationVector is a 3x2 matrix of int32 values
//
// Notes: XYZ coordinates of the acceleration vector in units of g.  As viewed from
// the front of the phone, positive X is toward the left side, positive Y is
// toward the bottom, and positive Z points into the face of the phone.
type AccelerationVector [3][2]int32

// AccelerationVectorfromRaw returns a new AccelerationVector
//
// Note: the directions are contrary to the Apple documentation (which have the
// signs of all axes reversed -- apparently the Apple geeks aren't very good
// with basic physics, and don't understand the concept of acceleration.  See
// http://nscookbook.com/2013/03/ios-programming-recipe-19-using-core-motion-to-access-gyro-and-accelerometer/
// for one of the few correct descriptions of this).  Note that this leads to
// a left-handed coordinate system for acceleration.
func AccelerationVectorfromRaw(raw []int32) AccelerationVector {
	return AccelerationVector{
		{raw[0], raw[1]},
		{raw[2], raw[3]},
		{raw[4], raw[5]},
	}
}

// String returns a string representation of AccelerationVector
func (a AccelerationVector) String() string {
	return fmt.Sprintf("X:%.2f, Y:%.2f, Z:%.2f", float32(a[0][0])/float32(a[0][1]), float32(a[1][0])/float32(a[1][1]), float32(a[2][0])/float32(a[2][1]))
}

// HDRImageType represents the HDR processing state of an image
type HDRImageType int32

const (
	HDRUnknown   HDRImageType = iota
	HDRProcessed HDRImageType = 3 // HDR Image
	HDROriginal  HDRImageType = 4 // Original Image
)

// String returns a string representation of HDRImageType
func (h HDRImageType) String() string {
	switch h {
	case HDRProcessed:
		return "HDR Image"
	case HDROriginal:
		return "Original Image"
	default:
		return fmt.Sprintf("Unknown HDRImage Type %d", h)
	}
}

// FocusDistanceRange represents near and far focus distances in meters
type FocusDistanceRange [2][2]int64 // [2]Rational64s

// String formats the focus range as "near - far m" with ordered values
func (f FocusDistanceRange) String() string {
	near, far := float64(f[0][0])/float64(f[0][1]), float64(f[1][0])/float64(f[1][1])
	if near > far {
		near, far = far, near // swap to ensure correct order
	}
	return fmt.Sprintf("%.2f - %.2f m", near, far)
}
