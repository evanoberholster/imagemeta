package jpeg

import (
	"io"
	"math"
	"time"

	"github.com/evanoberholster/imagemeta/meta/utils"
)

// CIFF stores selected Canon CIFF metadata found in old APP0 HEAPJPGM records.
type CIFF struct {
	FileFormat             uint32
	TargetCompressionRatio float64
	ImageWidth             uint32
	ImageHeight            uint32
	PixelAspectRatio       float64
	Rotation               int32
	ComponentBitDepth      uint32
	ColorBitDepth          uint32
	ColorBW                uint32
	TargetImageType        uint16
	RecordID               uint32
	FileNumber             uint32
	DateTimeOriginal       time.Time
	TimeZoneCode           int32
	TimeZoneInfo           uint32
	OriginalFileName       string
	ThumbnailFileName      string
	ShutterReleaseMethod   uint16
	ShutterReleaseTiming   uint16
	FlashGuideNumber       float64
	FlashThreshold         float64
	ExposureCompensation   float64
	ShutterSpeedValue      float64
	ApertureValue          float64
	TargetDistanceSetting  float64
	MeasuredEV             float64
	CanonFileDescription   string
	CanonImageType         string
	OwnerName              string
	Make                   string
	Model                  string
	UnknownNumber          uint32
	BaseISO                uint16
	CanonFirmwareVersion   string
	ComponentVersion       string
	ROMOperationMode       string
	CanonFlashInfo         string
	FocalType              uint16
	FocalLength            float64
	FocalPlaneXSize        float64
	FocalPlaneYSize        float64
	UnknownTags            map[string]string
}

func parseCIFF(payload []byte) (*CIFF, error) {
	if !isCIFFPayload(payload) || len(payload) < 18 {
		return nil, errShortSegment("CIFF")
	}
	order := ciffByteOrder(payload)
	if order == utils.UnknownEndian {
		return nil, io.ErrUnexpectedEOF
	}
	c := &CIFF{UnknownTags: make(map[string]string)}
	headerLen := int(order.Uint32(payload[2:6]))
	if headerLen < 14 || headerLen > len(payload)-6 {
		headerLen = 0
	}
	if err := parseCIFFDirectory(c, payload, order, headerLen, len(payload)-headerLen, "CIFF", 0); err != nil {
		return nil, err
	}
	if len(c.UnknownTags) == 0 {
		c.UnknownTags = nil
	}
	return c, nil
}

func ciffByteOrder(payload []byte) utils.ByteOrder {
	if len(payload) < 2 {
		return utils.UnknownEndian
	}
	switch string(payload[:2]) {
	case "II":
		return utils.LittleEndian
	case "MM":
		return utils.BigEndian
	default:
		return utils.UnknownEndian
	}
}

func parseCIFFDirectory(c *CIFF, block []byte, order utils.ByteOrder, blockStart, blockSize int, dirName string, depth int) error {
	if depth > 16 || blockStart < 0 || blockSize < 6 || blockStart+blockSize > len(block) {
		return nil
	}
	dirPtrPos := blockStart + blockSize - 4
	dirOffset := int(order.Uint32(block[dirPtrPos:dirPtrPos+4])) + blockStart
	if dirOffset < blockStart || dirOffset+2 > blockStart+blockSize {
		return nil
	}
	entries := int(order.Uint16(block[dirOffset : dirOffset+2]))
	dirEntries := dirOffset + 2
	if entries < 0 || dirEntries+entries*10 > blockStart+blockSize {
		return nil
	}
	for i := 0; i < entries; i++ {
		entry := block[dirEntries+i*10:]
		tag := order.Uint16(entry[0:2])
		if tag&0x8000 != 0 {
			continue
		}
		size := int(order.Uint32(entry[2:6]))
		ptr := int(order.Uint32(entry[6:10])) + blockStart
		tagID := tag & 0x3fff
		tagType := (tag >> 8) & 0x38
		valueInDir := tag&0x4000 != 0
		if (tagType == 0x28 || tagType == 0x30) && !valueInDir {
			if ptr >= blockStart && size > 0 && ptr+size <= blockStart+blockSize {
				_ = parseCIFFDirectory(c, block, order, ptr, size, ciffDirName(tagID), depth+1)
			}
			continue
		}
		value := entry[2:10]
		if !valueInDir {
			if size < 0 || ptr < blockStart || ptr+size > blockStart+blockSize {
				continue
			}
			value = block[ptr : ptr+size]
		}
		parseCIFFTag(c, order, dirName, tagID, tagType, value)
	}
	return nil
}

func parseCIFFTag(c *CIFF, order utils.ByteOrder, dirName string, tagID uint16, tagType uint16, value []byte) {
	switch tagID {
	case 0x0805:
		if dirName == "ImageDescription" || c.CanonFileDescription == "" {
			c.CanonFileDescription = trimNULString(value)
		}
	case 0x0806:
		c.addUnknown(tagID, trimNULString(value))
	case 0x080a:
		parseCIFFMakeModel(c, value)
	case 0x080b:
		c.CanonFirmwareVersion = trimNULString(value)
	case 0x080c:
		c.ComponentVersion = trimNULString(value)
	case 0x080d:
		c.ROMOperationMode = trimNULString(value)
	case 0x0810:
		c.OwnerName = trimNULString(value)
	case 0x0815:
		c.CanonImageType = trimNULString(value)
	case 0x0816:
		c.OriginalFileName = trimNULString(value)
	case 0x0817:
		c.ThumbnailFileName = trimNULString(value)
	case 0x100a:
		c.TargetImageType = firstCIFFUint16(order, value)
	case 0x1010:
		c.ShutterReleaseMethod = firstCIFFUint16(order, value)
	case 0x1011:
		c.ShutterReleaseTiming = firstCIFFUint16(order, value)
	case 0x1013, 0x1026, 0x1812, 0x1819:
		c.addUnknown(tagID, ciffValueString(order, tagType, value))
	case 0x0006:
		if len(value) > 0 {
			c.addUnknown(tagID, u8ListString(value[:1]))
		}
	case 0x1014:
		c.addUnknown(tagID, u16ScalarString(order, value))
	case 0x1805, 0x1811:
		c.addUnknown(tagID, u32ScalarString(order, value))
	case 0x101c:
		c.BaseISO = firstCIFFUint16(order, value)
	case 0x1028:
		c.CanonFlashInfo = u16ListString(order, value)
	case 0x1029:
		parseCIFFFocalLength(c, order, value)
	case 0x1803:
		parseCIFFImageFormat(c, order, value)
	case 0x1804:
		c.RecordID = firstCIFFUint32(order, value)
	case 0x1807:
		c.TargetDistanceSetting = ciffFloat32(order, value)
	case 0x180b:
		c.UnknownNumber = firstCIFFUint32(order, value)
	case 0x180e:
		parseCIFFTimeStamp(c, order, value)
	case 0x1810:
		parseCIFFImageInfo(c, order, value)
	case 0x1813:
		parseCIFFFlashInfo(c, order, value)
	case 0x1814:
		c.MeasuredEV = ciffFloat32(order, value) + 5
	case 0x1817:
		c.FileNumber = firstCIFFUint32(order, value)
	case 0x1818:
		parseCIFFExposureInfo(c, order, value)
	default:
		c.addUnknown(tagID, ciffValueString(order, tagType, value))
	}
}

func parseCIFFImageFormat(c *CIFF, order utils.ByteOrder, value []byte) {
	if len(value) >= 4 {
		c.FileFormat = order.Uint32(value[0:4])
	}
	if len(value) >= 8 {
		c.TargetCompressionRatio = ciffFloat32(order, value[4:8])
	}
}

func parseCIFFImageInfo(c *CIFF, order utils.ByteOrder, value []byte) {
	if len(value) >= 4 {
		c.ImageWidth = order.Uint32(value[0:4])
	}
	if len(value) >= 8 {
		c.ImageHeight = order.Uint32(value[4:8])
	}
	if len(value) >= 12 {
		c.PixelAspectRatio = ciffFloat32(order, value[8:12])
	}
	if len(value) >= 16 {
		c.Rotation = int32(order.Uint32(value[12:16]))
	}
	if len(value) >= 20 {
		c.ComponentBitDepth = order.Uint32(value[16:20])
	}
	if len(value) >= 24 {
		c.ColorBitDepth = order.Uint32(value[20:24])
	}
	if len(value) >= 28 {
		c.ColorBW = order.Uint32(value[24:28])
	}
}

func parseCIFFTimeStamp(c *CIFF, order utils.ByteOrder, value []byte) {
	if len(value) >= 4 {
		c.DateTimeOriginal = ciffTime(order, value[0:4])
	}
	if len(value) >= 8 {
		c.TimeZoneCode = int32(order.Uint32(value[4:8])) / 3600
	}
	if len(value) >= 12 {
		c.TimeZoneInfo = order.Uint32(value[8:12])
	}
}

func parseCIFFFlashInfo(c *CIFF, order utils.ByteOrder, value []byte) {
	if len(value) >= 4 {
		c.FlashGuideNumber = ciffFloat32(order, value[0:4])
	}
	if len(value) >= 8 {
		c.FlashThreshold = ciffFloat32(order, value[4:8])
	}
}

func parseCIFFExposureInfo(c *CIFF, order utils.ByteOrder, value []byte) {
	if len(value) >= 4 {
		c.ExposureCompensation = ciffFloat32(order, value[0:4])
	}
	if len(value) >= 8 {
		apex := ciffFloat32(order, value[4:8])
		if math.Abs(apex) < 100 {
			c.ShutterSpeedValue = 1 / math.Pow(2, apex)
		}
	}
	if len(value) >= 12 {
		c.ApertureValue = math.Pow(2, ciffFloat32(order, value[8:12])/2)
	}
}

func parseCIFFFocalLength(c *CIFF, order utils.ByteOrder, value []byte) {
	if len(value) >= 2 {
		c.FocalType = order.Uint16(value[0:2])
	}
	if len(value) >= 4 {
		c.FocalLength = float64(order.Uint16(value[2:4]))
	}
	if len(value) >= 6 {
		c.FocalPlaneXSize = float64(order.Uint16(value[4:6])) * 25.4 / 1000
	}
	if len(value) >= 8 {
		c.FocalPlaneYSize = float64(order.Uint16(value[6:8])) * 25.4 / 1000
	}
}

func parseCIFFMakeModel(c *CIFF, value []byte) {
	if len(value) >= 6 {
		c.Make = trimNULString(value[:6])
		c.Model = trimNULString(value[6:])
		return
	}
	c.Make = trimNULString(value)
}

func firstCIFFUint16(order utils.ByteOrder, value []byte) uint16 {
	if len(value) < 2 {
		return 0
	}
	return order.Uint16(value)
}

func firstCIFFUint32(order utils.ByteOrder, value []byte) uint32 {
	if len(value) < 4 {
		return 0
	}
	return order.Uint32(value)
}

func ciffValueString(order utils.ByteOrder, tagType uint16, value []byte) string {
	switch tagType {
	case 0x00:
		return u8ListString(value)
	case 0x08:
		return trimNULString(value)
	case 0x10:
		return u16ListString(order, value)
	case 0x18:
		return u32ListString(order, value)
	default:
		return trimNULString(value)
	}
}

func u16ScalarString(order utils.ByteOrder, value []byte) string {
	if len(value) < 2 {
		return ""
	}
	return u16ListString(order, value[:2])
}

func u32ScalarString(order utils.ByteOrder, value []byte) string {
	if len(value) < 4 {
		return ""
	}
	return u32ListString(order, value[:4])
}

func ciffDirName(tagID uint16) string {
	switch tagID {
	case 0x2804:
		return "ImageDescription"
	case 0x2807:
		return "CameraObject"
	case 0x3002:
		return "ShootingRecord"
	case 0x3003:
		return "MeasuredInfo"
	case 0x3004:
		return "CameraSpecification"
	case 0x300a:
		return "ImageProps"
	case 0x300b:
		return "ExifInformation"
	default:
		return "CIFF"
	}
}

func (c *CIFF) addUnknown(tagID uint16, value string) {
	if value == "" {
		return
	}
	if c.UnknownTags == nil {
		c.UnknownTags = make(map[string]string)
	}
	c.UnknownTags[ciffRawTagName(tagID)] = value
}

func ciffRawTagName(tagID uint16) string {
	const digits = "0123456789abcdef"
	var b [15]byte
	copy(b[:], "CanonRaw_0x")
	b[11] = digits[(tagID>>12)&0xf]
	b[12] = digits[(tagID>>8)&0xf]
	b[13] = digits[(tagID>>4)&0xf]
	b[14] = digits[tagID&0xf]
	return string(b[:])
}
