package exif2

import (
	"fmt"
	"strings"
	"time"

	"github.com/evanoberholster/imagemeta/meta"
)

type Exif struct {
	ProcessingSoftware        string               // IFD0 / 0x000b
	DocumentName              string               // IFD0 / 0x010d
	ImageDescription          string               // IFD0 / 0x010e
	Software                  string               // IFD0 / 0x0131
	Make                      string               // IFD0 / 0x010f
	Model                     string               // IFD0 / 0x0110
	LensMake                  string               // ExifIFD / 0xa433
	LensModel                 string               // ExifIFD / 0xa434
	LensSerial                string               // ExifIFD / 0xa435
	ImageUniqueID             string               // ExifIFD / 0xa420
	OwnerName                 string               // ExifIFD / 0xa430	(called CameraOwnerName by the EXIF spec.)
	CameraSerial              string               // ExifIFD / 0xa431	(called BodySerialNumber by the EXIF spec.)
	SubfileType               uint32               // IFD0 / 0x00fe
	ImageWidth                uint16               // IFD0 / 0x0100
	ImageHeight               uint16               // IFD0 / 0x0101
	Compression               Compression          // IFD0 / 0x0103
	PhotometricInterpretation uint16               // IFD0 / 0x0106
	Orientation               meta.Orientation     // IFD0 / 0x0112
	StripOffsets              uint32               // IFD0 / 0x0111 PreviewImageStart
	StripByteCounts           uint32               // IFD0 / 0x0117 PreviewImageLength
	XResolution               Resolution           // IFD0 / 0x011a rational64u
	YResolution               Resolution           // IFD0 / 0x011b
	ResolutionUnit            uint16               // IFD0 / 0x0128
	modifyDate                time.Time            // IFD0 / 0x0132
	Artist                    string               // IFD0 / 0x013b
	ThumbnailOffset           uint32               // 0x0201
	ThumbnailLength           uint32               // 0x0202
	ApplicationNotes          []byte               // 0x02bc
	Rating                    uint16               // 0x4746
	Copyright                 string               // 0x8298
	ExposureTime              ExposureTime         // 0x829a
	FNumber                   FNumber              // 0x829d
	ExposureProgram           meta.ExposureProgram // 0x8822
	ExposureBias              meta.ExposureBias    // 0x0d34
	ExposureMode              meta.ExposureMode    // ExifIFD / 0xa402
	GPS                       GPSInfo              // 0x8825
	ISO                       uint16               // ExifIFD / 0x8827 // FixMe
	TimeZoneOffset            [2]int8              // ExifIFD / 0x882a // FixMe (1 or 2 values: 1. The time zone offset of DateTimeOriginal from GMT in hours, 2. If present, the time zone offset of ModifyDate)
	SelfTimerMode             uint16               // ExifIFD / 0x882b
	ISOSpeed                  uint32               // ExifIFD / 0x8833
	dateTimeOriginal          time.Time            // ExifIFD / 0x9003
	createDate                time.Time            // ExifIFD / 0x9004
	subSecTime                uint16               // ExifIFD / 0x9290
	subSecTimeOriginal        uint16               // ExifIFD / 0x9291
	subSecTimeDigitized       uint16               // ExifIFD / 0x9292
	SubjectDistance           RationalU            // ExifIFD / 0x9206
	MeteringMode              meta.MeteringMode    // ExifIFD / 0x9207
	Flash                     meta.Flash           // ExifIFD / 0x9209
	FocalLength               meta.FocalLength     // ExifIFD / 0x920a
	ImageNumber               uint32               // ExifIFD / 0x9211
	SubjectArea               SubjectArea          // ExifIFD / 0x9214
	Makernote                 string               // ExifIFD / MakerNote
	ColorSpace                ColorSpace           // ExifIFD / 0xa001
	FocalLengthIn35mmFormat   meta.FocalLength     // ExifIFD / 0xa405
	LensInfo                  LensInfo             // ExifIFD / 0xa432	(4 rational values giving focal and aperture ranges, called LensSpecification by the EXIF spec.)
	// 0xa002	ExifImageWidth	int16u:	ExifIFD	(called PixelXDimension by the EXIF spec.)
	// 0xa003	ExifImageHeight	int16u:	ExifIFD	(called PixelYDimension by the EXIF spec.)
	// 0xa20e	FocalPlaneXResolution	rational64u	ExifIFD
	// 0xa20f	FocalPlaneYResolution	rational64u	ExifIFD
	// OffsetTime          string          // ExifIFD / 0x9010 (time zone for ModifyDate)
	// OffsetTimeOriginal  string          // ExifIFD / 0x9011 (time zone for DateTimeOriginal)
	// OffsetTimeDigitized string          // ExifIFD / 0x9012 (time zone for CreateDate)
	// 0x9290	SubSecTime			string	ExifIFD	(fractional seconds for ModifyDate)
	// 0x9291	SubSecTimeOriginal	string	ExifIFD	(fractional seconds for DateTimeOriginal)
	// 0x9292	SubSecTimeDigitized	string	ExifIFD	(fractional seconds for CreateDate)
}

func (e Exif) ModifyDate() time.Time {
	if e.subSecTime == 0 {
		return e.modifyDate
	}
	return e.modifyDate.Add(time.Duration(e.subSecTime) * time.Millisecond)
}

func (e Exif) DateTimeOriginal() time.Time {
	if e.subSecTimeOriginal == 0 {
		return e.dateTimeOriginal
	}
	return e.dateTimeOriginal.Add(time.Duration(e.subSecTimeOriginal) * time.Millisecond)
}

func (e Exif) CreateDate() time.Time {
	if e.subSecTimeDigitized == 0 {
		return e.createDate
	}
	return e.createDate.Add(time.Duration(e.subSecTimeDigitized) * time.Millisecond)
}

func (e Exif) String() string {
	sb := strings.Builder{}
	sb.WriteString("Exif\n")
	sb.WriteString(fmt.Sprintf("Make: \t\t%s\n", e.Make))
	sb.WriteString(fmt.Sprintf("Model: \t\t%s\n", e.Model))
	sb.WriteString(fmt.Sprintf("LensMake: \t%s\n", e.LensMake))
	sb.WriteString(fmt.Sprintf("LensModel: \t%s\n", e.LensModel))
	sb.WriteString(fmt.Sprintf("CameraSerial: \t%s\n", e.CameraSerial))
	sb.WriteString(fmt.Sprintf("LensSerial: \t%s\n", e.LensSerial))
	sb.WriteString(fmt.Sprintf("Image Size: \t%dx%d\n", e.ImageWidth, e.ImageHeight))
	sb.WriteString(fmt.Sprintf("Orientation: \t%s\n", e.Orientation))
	sb.WriteString(fmt.Sprintf("ShutterSpeed: \t%d/%d\n", e.ExposureTime[0], e.ExposureTime[1]))
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
	sb.WriteString(fmt.Sprintf("Date GPS: \t%s\n", e.GPS.Date()))
	sb.WriteString(fmt.Sprintf("Artist: \t%s\n", e.Artist))
	sb.WriteString(fmt.Sprintf("Copyright: \t%s\n", e.Copyright))
	sb.WriteString(fmt.Sprintf("GPS Altitude: \t%0.2f\n", e.GPS.Altitude()))
	sb.WriteString(fmt.Sprintf("GPS Latitude: \t%f\n", e.GPS.Latitude()))
	sb.WriteString(fmt.Sprintf("GPS Longitude: \t%f\n", e.GPS.Longitude()))

	return sb.String()
}

type LensInfo [4]RationalU

type ColorSpace uint16

type SubjectArea []uint16

type GPSInfo struct {
	latitude     float64   // Combination of GPSLatitudeRef and GPSLatitude
	longitude    float64   // Combination of GPSLongitudeRef and GPSLongitude
	altitude     float32   // Combination of GPSAltitudeRef and GPSAltitude
	date         [2]uint16 // [0]months since AD0 [1]day
	time         uint32    // time in seconds
	latitudeRef  bool
	longitudeRef bool
	altitudeRef  bool
}

func (g GPSInfo) Date() time.Time {
	year := int(g.date[0] / 12)
	month := int(g.date[0] % 12)
	day := int(g.date[1])
	hour := int(g.time / 3600)
	min := (int(g.time) - (hour * 3600)) / 60
	sec := (int(g.time) - (hour * 3600) - (min * 60))
	return time.Date(year, time.Month(month), day, hour, min, sec, 0, time.UTC)
}

func (g GPSInfo) Latitude() float64 {
	if g.latitudeRef {
		return -1 * g.latitude
	}
	return g.latitude
}

func (g GPSInfo) Longitude() float64 {
	if g.longitudeRef {
		return -1 * g.longitude
	}
	return g.longitude
}

func (g GPSInfo) Altitude() float32 {
	if g.altitudeRef {
		return -1 * g.altitude
	}
	return g.altitude
}

type Flash uint16

type FocalLength RationalU

type SonyRaw struct {
	//SonyRawFileType // 0x7000
}

type FNumber float32

type ExposureTime RationalU

type Resolution RationalU

type RationalU [2]uint32

func (rU RationalU) AsFloat() float32 {
	return float32(rU[0]) / float32(rU[1])
}

// Compression is Exif Compression.
type Compression uint16

//1		= Uncompressed
//2		= CCITT 1D
//3		= T4/Group 3 Fax
//4		= T6/Group 4 Fax
//5		= LZW
//6		= JPEG (old-style)
//7		= JPEG
//8		= Adobe Deflate
//9		= JBIG B&W
//10		= JBIG Color
//99		= JPEG
//262		= Kodak 262
//32766	= Next
//32767	= Sony ARW Compressed
//32769	= Packed RAW
//32770	= Samsung SRW Compressed
//32771	= CCIRLEW
//32772	= Samsung SRW Compressed 2
//32773	= PackBits
//32809	= Thunderscan
//32867	= Kodak KDC Compressed
//32895	= IT8CTPAD
//32896	= IT8LW
//32897	= IT8MP
//32898	= IT8BL
//32908	= PixarFilm
//32909	= PixarLog
//32946	= Deflate
//32947	= DCS
//33003	= Aperio JPEG 2000 YCbCr
//33005	= Aperio JPEG 2000 RGB
//34661	= JBIG
//34676	= SGILog
//34677	= SGILog24
//34712	= JPEG 2000
//34713	= Nikon NEF Compressed
//34715	= JBIG2 TIFF FX
//34718	= Microsoft Document Imaging (MDI) Binary Level Codec
//34719	= Microsoft Document Imaging (MDI) Progressive Transform Codec
//34720	= Microsoft Document Imaging (MDI) Vector
//34887	= ESRI Lerc
//34892	= Lossy JPEG
//34925	= LZMA2
//34926	= Zstd
//34927	= WebP
//34933	= PNG
//34934	= JPEG XR
//65000	= Kodak DCR Compressed
//65535	= Pentax PEF Compressed

// Orientation is Exif Orientation
type Orientation uint16

// 1 = Horizontal (normal)
// 2 = Mirror horizontal
// 3 = Rotate 180
// 4 = Mirror vertical
// 5 = Mirror horizontal and rotate 270 CW
// 6 = Rotate 90 CW
// 7 = Mirror horizontal and rotate 90 CW
// 8 = Rotate 270 CW
