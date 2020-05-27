package mknote

import (
	"github.com/evanoberholster/exiftool/tag"
)

// CanonMKnoteIFD TagIDs
// Source: https://exiftool.org/TagNames/Canon.html on 8/05/2020
const (
	CanonCameraSettings        tag.ID = 0x0001
	CanonFocalLength           tag.ID = 0x0002
	CanonFlashInfo             tag.ID = 0x0003
	CanonShotInfo              tag.ID = 0x0004
	CanonPanorama              tag.ID = 0x0005
	CanonImageType             tag.ID = 0x0006
	CanonFirmwareVersion       tag.ID = 0x0007
	FileNumber                 tag.ID = 0x0008
	OwnerName                  tag.ID = 0x0009
	UnknownD30                 tag.ID = 0x000a
	SerialNumber               tag.ID = 0x000c
	CanonCameraInfo            tag.ID = 0x000d // WIP
	CanonFileLength            tag.ID = 0x000e // WIP
	CustomFunctions            tag.ID = 0x000f // WIP
	CanonModelID               tag.ID = 0x0010
	MovieInfo                  tag.ID = 0x0011 // WIP
	CanonAFInfo                tag.ID = 0x0012
	ThumbnailImageValidArea    tag.ID = 0x0013 // WIP
	SerialNumberFormat         tag.ID = 0x0015 // WIP
	SuperMacro                 tag.ID = 0x001a // WIP
	DateStampMode              tag.ID = 0x001c // WIP
	MyColors                   tag.ID = 0x001d // WIP
	FirmwareRevision           tag.ID = 0x001e // WIP
	Categories                 tag.ID = 0x0023 // WIP
	FaceDetect1                tag.ID = 0x0024 // WIP
	FaceDetect2                tag.ID = 0x0025 // WIP
	CanonAFInfo2               tag.ID = 0x0026
	ContrastInfo               tag.ID = 0x0027 // WIP
	ImageUniqueID              tag.ID = 0x0028 // WIP
	WBInfo                     tag.ID = 0x0029 // WIP
	FaceDetect3                tag.ID = 0x002f // WIP
	TimeInfo                   tag.ID = 0x0035
	BatteryType                tag.ID = 0x0038 // WIP
	AFInfo3                    tag.ID = 0x003c // WIP
	RawDataOffset              tag.ID = 0x0081 // WIP
	OriginalDecisionDataOffset tag.ID = 0x0083 // WIP
	CustomFunctions1D          tag.ID = 0x0090 // WIP
	PersonalFunctions          tag.ID = 0x0091 // WIP
	PersonalFunctionValues     tag.ID = 0x0092 // WIP
	CanonFileInfo              tag.ID = 0x0093
	AFPointsInFocus1D          tag.ID = 0x0094 // WIP
	LensModel                  tag.ID = 0x0095
)
