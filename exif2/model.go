package exif2

import (
	"fmt"
	"strings"
	"time"

	"github.com/evanoberholster/imagemeta/exif2/ifds/exififd"
	"github.com/evanoberholster/imagemeta/exif2/ifds/gpsifd"
	"github.com/evanoberholster/imagemeta/exif2/make"
	"github.com/evanoberholster/imagemeta/exif2/model"
	"github.com/evanoberholster/imagemeta/imagetype"
	"github.com/evanoberholster/imagemeta/meta"
)

// Exif data structure
type Exif struct {
	ApplicationNotes          []byte               // 0x02bc
	GPS                       gpsifd.GPSInfo       // GPSInfoIfd / 0x8825
	Exif                      exififd.ExifIfd      // ExifIfd
	SubjectArea               SubjectArea          // ExifIFD / 0x9214
	LensInfo                  LensInfo             // ExifIFD / 0xa432	(4 rational values giving focal and aperture ranges, called LensSpecification by the EXIF spec.)
	Makernotes                MakerNotes           // ExifIFD / MakerNote
	Time                      TimeTags             // TimeTags
	ProcessingSoftware        string               // IFD0 / 0x000b
	DocumentName              string               // IFD0 / 0x010d
	ImageDescription          string               // IFD0 / 0x010e
	Software                  string               // IFD0 / 0x0131
	Artist                    string               // IFD0 / 0x013b
	Copyright                 string               // IFD0 / 0x8298
	LensMake                  string               // ExifIFD / 0xa433
	LensModel                 string               // ExifIFD / 0xa434
	LensSerial                string               // ExifIFD / 0xa435
	ImageUniqueID             string               // ExifIFD / 0xa420
	OwnerName                 string               // ExifIFD / 0xa430	(called CameraOwnerName by the EXIF spec.)
	CameraSerial              string               // ExifIFD / 0xa431	(called BodySerialNumber by the EXIF spec.)
	Make                      string               // IFD0 / 0x010f
	Model                     string               // IFD0 / 0x0110
	CameraModel               model.CameraModel    // CameraModel
	CameraMake                make.CameraMake      // CameraMake
	XResolution               uint32               // IFD0 / 0x011a rational64u
	YResolution               uint32               // IFD0 / 0x011b
	ExposureTime              meta.ExposureTime    // 0x829a
	SubjectDistance           float32              // ExifIFD / 0x9206
	FocalLength               meta.FocalLength     // ExifIFD / 0x920a
	FocalLengthIn35mmFormat   meta.FocalLength     // ExifIFD / 0xa405
	StripOffsets              uint32               // IFD0 / 0x0111 PreviewImageStart
	StripByteCounts           uint32               // IFD0 / 0x0117 PreviewImageLength
	ThumbnailOffset           uint32               // 0x0201
	ThumbnailLength           uint32               // 0x0202
	SubfileType               uint32               // IFD0 / 0x00fe
	FNumber                   meta.Aperture        // 0x829d
	ISOSpeed                  uint32               // ExifIFD / 0x8833
	ImageNumber               uint32               // ExifIFD / 0x9211
	ImageWidth                uint16               // IFD0 / 0x0100 // ExifIFD	/ 0xa002	ExifImageWidth	int16u:	(called PixelXDimension by the EXIF spec.)
	ImageHeight               uint16               // IFD0 / 0x0101 // ExifIFD	/ 0xa003	ExifImageHeight	int16u:	(called PixelYDimension by the EXIF spec.)
	Compression               meta.Compression     // IFD0 / 0x0103
	PhotometricInterpretation uint16               // IFD0 / 0x0106
	Orientation               meta.Orientation     // IFD0 / 0x0112
	ResolutionUnit            uint16               // IFD0 / 0x0128
	Rating                    uint16               // 0x4746
	ExposureProgram           meta.ExposureProgram // 0x8822
	ExposureBias              meta.ExposureBias    // 0x0d34
	ExposureMode              meta.ExposureMode    // ExifIFD / 0xa402
	ISO                       uint16               // ExifIFD / 0x8827 // FixMe
	SelfTimerMode             uint16               // ExifIFD / 0x882b
	MeteringMode              meta.MeteringMode    // ExifIFD / 0x9207
	Flash                     meta.Flash           // ExifIFD / 0x9209
	ColorSpace                ColorSpace           // ExifIFD / 0xa001
	ImageType                 imagetype.ImageType

	// 0xa20e	FocalPlaneXResolution	rational64u	ExifIFD
	// 0xa20f	FocalPlaneYResolution	rational64u	ExifIFD
}

// ModifyDate return the exif modified date with subsec offset if present
func (e Exif) ModifyDate() time.Time {
	t := e.Time.modifyDate
	if e.Time.subSecTime != 0 {
		t = t.Add(time.Duration(e.Time.subSecTime) * time.Millisecond)
	}
	if e.Time.offsetTime != nil {
		t = t.In(e.Time.offsetTime)
		_, offset := t.Zone()
		t = t.Add(time.Duration(offset) * -1 * time.Second)
	}
	return t
}

// DateTimeOriginal returns the exif Original DateTime with subsec offset if present
func (e Exif) DateTimeOriginal() time.Time {
	t := e.Time.dateTimeOriginal
	if e.Time.subSecTimeOriginal != 0 {
		t = t.Add(time.Duration(e.Time.subSecTimeOriginal) * time.Millisecond)
	}
	if e.Time.offsetTimeOriginal != nil {
		t = t.In(e.Time.offsetTimeOriginal)
		_, offset := t.Zone()
		t = t.Add(time.Duration(offset) * -1 * time.Second)
	}
	return t
}

// CreateDate reurns the CreateDate with subsec offset if present
func (e Exif) CreateDate() time.Time {
	t := e.Time.createDate
	if e.Time.subSecTimeDigitized != 0 {
		t = t.Add(time.Duration(e.Time.subSecTimeDigitized) * time.Millisecond)
	}
	if e.Time.offsetTimeDigitized != nil {
		t = t.In(e.Time.offsetTimeDigitized)
		_, offset := t.Zone()
		t = t.Add(time.Duration(offset) * -1 * time.Second)
	}
	return t
}

// Sring implements the Stringer interface for Exif
func (e Exif) String() string {
	sb := strings.Builder{}
	sb.WriteString("Exif\n")
	sb.WriteString(fmt.Sprintf("ImageType: \t%s\n", e.ImageType))
	sb.WriteString(fmt.Sprintf("Make: \t\t%s\n", e.Make))
	sb.WriteString(fmt.Sprintf("Model: \t\t%s\n", e.Model))
	sb.WriteString(fmt.Sprintf("LensMake: \t%s\n", e.LensMake))
	sb.WriteString(fmt.Sprintf("LensModel: \t%s\n", e.LensModel))
	sb.WriteString(fmt.Sprintf("CameraSerial: \t%s\n", e.CameraSerial))
	sb.WriteString(fmt.Sprintf("LensSerial: \t%s\n", e.LensSerial))
	sb.WriteString(fmt.Sprintf("Image Size: \t%dx%d\n", e.ImageWidth, e.ImageHeight))
	sb.WriteString(fmt.Sprintf("Orientation: \t%s\n", e.Orientation))
	sb.WriteString(fmt.Sprintf("ShutterSpeed: \t%s\n", e.ExposureTime))
	sb.WriteString(fmt.Sprintf("Aperture: \t%0.2f\n", e.FNumber))
	sb.WriteString(fmt.Sprintf("ISO: \t\t%d\n", e.ISOSpeed))
	sb.WriteString(fmt.Sprintf("Flash: \t\t%s\n", e.Flash))
	sb.WriteString(fmt.Sprintf("Focal Length: \t%s\n", e.FocalLength))
	sb.WriteString(fmt.Sprintf("Fl 35mm Eqv: \t%s\n", e.FocalLengthIn35mmFormat))
	sb.WriteString(fmt.Sprintf("Exposure Prgm: \t%s\n", e.ExposureProgram))
	sb.WriteString(fmt.Sprintf("Metering Mode: \t%s\n", e.MeteringMode))
	sb.WriteString(fmt.Sprintf("Exposure Mode: \t%s\n", e.ExposureMode))
	sb.WriteString(fmt.Sprintf("Date Modified: \t%s\n", e.ModifyDate()))
	sb.WriteString(fmt.Sprintf("Date Created: \t%s\n", e.CreateDate()))
	sb.WriteString(fmt.Sprintf("Date Original: \t%s\n", e.DateTimeOriginal()))
	sb.WriteString(fmt.Sprintf("Artist: \t%s\n", e.Artist))
	sb.WriteString(fmt.Sprintf("Copyright: \t%s\n", e.Copyright))
	sb.WriteString(fmt.Sprintf("Software: \t%s\n", e.Software))
	sb.WriteString(fmt.Sprintf("Image Desc: \t%s\n", e.ImageDescription))
	sb.WriteString(e.GPS.String())
	return sb.String()
}

// LensInfo struct
type LensInfo [8]uint32

// ColorSpace data
type ColorSpace uint16

// SubjectArea coordinates
type SubjectArea []uint16

// TimeTags contains time Exif tags
type TimeTags struct {
	modifyDate          time.Time      // IFD0 / 0x0132
	dateTimeOriginal    time.Time      // ExifIFD / 0x9003
	createDate          time.Time      // ExifIFD / 0x9004
	offsetTime          *time.Location // ExifIFD / 0x9010 (time zone for ModifyDate)
	offsetTimeOriginal  *time.Location // ExifIFD / 0x9011 (time zone for DateTimeOriginal)
	offsetTimeDigitized *time.Location // ExifIFD / 0x9012 (time zone for CreateDate)
	subSecTime          uint16         // ExifIFD / 0x9290 (fractional seconds for ModifyDate)
	subSecTimeOriginal  uint16         // ExifIFD / 0x9291 (fractional seconds for DateTimeOriginal)
	subSecTimeDigitized uint16         // ExifIFD / 0x9292 (fractional seconds for CreateDate)
	//timeZoneOffset      [2]int8        // ExifIFD / 0x882a // FixMe (1 or 2 values: 1. The time zone offset of DateTimeOriginal from GMT in hours, 2. If present, the time zone offset of ModifyDate)
}

// ApplicationNotes data are stil work in process
type ApplicationNotes []byte

// MakerNotes are still work in process
type MakerNotes interface {
}
