package jpeg

import (
	"bytes"
	"encoding/hex"
	"io"
	"time"

	"github.com/evanoberholster/imagemeta/meta/utils"
)

// ICCProfile stores selected ICC profile header and tag values.
type ICCProfile struct {
	ProfileCMMType            string
	ProfileVersion            uint32
	ProfileClass              string
	ColorSpaceData            string
	ProfileConnectionSpace    string
	ProfileDateTime           time.Time
	ProfileFileSignature      string
	PrimaryPlatform           string
	CMMFlags                  uint32
	DeviceManufacturer        string
	DeviceModel               string
	DeviceAttributes          [2]uint32
	RenderingIntent           uint32
	ConnectionSpaceIlluminant [3]float64
	ProfileCreator            string
	ProfileID                 [16]byte
	ProfileCopyright          string
	ProfileDescription        string
	MediaWhitePoint           [3]float64
	MediaBlackPoint           [3]float64
	RedMatrixColumn           [3]float64
	GreenMatrixColumn         [3]float64
	BlueMatrixColumn          [3]float64
	DeviceMfgDesc             string
	DeviceModelDesc           string
	ViewingCondDesc           string
	ViewingCondIlluminant     [3]float64
	ViewingCondSurround       [3]float64
	ViewingCondIlluminantType uint32
	Luminance                 [3]float64
	MeasurementObserver       uint32
	MeasurementBacking        [3]float64
	MeasurementGeometry       uint32
	MeasurementFlare          float64
	MeasurementIlluminant     uint32
	Technology                string
	RedTRCLength              uint32
	GreenTRCLength            uint32
	BlueTRCLength             uint32
}

func (m *Metadata) addICCChunk(payload []byte) error {
	if !isICCPayload(payload) || len(payload) < 14 {
		return errShortSegment("ICC")
	}
	seq := payload[12]
	total := payload[13]
	if seq == 0 || total == 0 || seq > total {
		return nil
	}
	if m.iccChunks == nil {
		m.iccChunks = make(map[uint8][]byte)
		m.iccTotal = total
	}
	if m.iccTotal != total {
		return nil
	}
	m.iccChunks[seq] = append(m.iccChunks[seq][:0], payload[14:]...)
	return nil
}

func (m *Metadata) finishICC() error {
	if m.iccTotal == 0 || len(m.iccChunks) == 0 {
		return nil
	}
	var assembled bytes.Buffer
	for i := uint8(1); i <= m.iccTotal; i++ {
		chunk, ok := m.iccChunks[i]
		if !ok {
			return nil
		}
		assembled.Write(chunk)
	}
	icc, err := parseICCProfile(assembled.Bytes())
	if err != nil {
		return err
	}
	m.ICC = icc
	return nil
}

func parseICCProfile(data []byte) (*ICCProfile, error) {
	if len(data) < 132 {
		return nil, errShortSegment("ICC profile")
	}
	if string(data[36:40]) != "acsp" {
		return nil, io.ErrUnexpectedEOF
	}
	icc := &ICCProfile{
		ProfileCMMType:         string(data[4:8]),
		ProfileVersion:         jpegEndian.Uint32(data[8:12]) >> 16,
		ProfileClass:           string(data[12:16]),
		ColorSpaceData:         string(data[16:20]),
		ProfileConnectionSpace: string(data[20:24]),
		ProfileFileSignature:   string(data[36:40]),
		PrimaryPlatform:        string(data[40:44]),
		CMMFlags:               jpegEndian.Uint32(data[44:48]),
		DeviceManufacturer:     string(data[48:52]),
		DeviceModel:            cleanICCSignature(data[52:56]),
		DeviceAttributes:       [2]uint32{jpegEndian.Uint32(data[56:60]), jpegEndian.Uint32(data[60:64])},
		RenderingIntent:        jpegEndian.Uint32(data[64:68]),
		ConnectionSpaceIlluminant: [3]float64{
			s15Fixed16(utils.BigEndian, data[68:72]),
			s15Fixed16(utils.BigEndian, data[72:76]),
			s15Fixed16(utils.BigEndian, data[76:80]),
		},
		ProfileCreator: string(data[80:84]),
	}
	copy(icc.ProfileID[:], data[84:100])
	icc.ProfileDateTime = parseICCTime(data[24:36])

	tagCount := int(jpegEndian.Uint32(data[128:132]))
	if 132+tagCount*12 > len(data) {
		return icc, nil
	}
	for i := 0; i < tagCount; i++ {
		entry := data[132+i*12:]
		sig := string(entry[0:4])
		offset := int(jpegEndian.Uint32(entry[4:8]))
		size := int(jpegEndian.Uint32(entry[8:12]))
		if offset < 0 || size < 0 || offset+size > len(data) {
			continue
		}
		parseICCTag(icc, sig, data[offset:offset+size])
	}
	return icc, nil
}

func cleanICCSignature(data []byte) string {
	for _, b := range data {
		if b != 0 {
			return string(data)
		}
	}
	return ""
}

func parseICCTag(icc *ICCProfile, sig string, data []byte) {
	switch sig {
	case "cprt":
		icc.ProfileCopyright = parseICCText(data)
	case "desc":
		icc.ProfileDescription = parseICCDesc(data)
	case "wtpt":
		icc.MediaWhitePoint = parseICCXYZ(data)
	case "bkpt":
		icc.MediaBlackPoint = parseICCXYZ(data)
	case "rXYZ":
		icc.RedMatrixColumn = parseICCXYZ(data)
	case "gXYZ":
		icc.GreenMatrixColumn = parseICCXYZ(data)
	case "bXYZ":
		icc.BlueMatrixColumn = parseICCXYZ(data)
	case "dmnd":
		icc.DeviceMfgDesc = parseICCDesc(data)
	case "dmdd":
		icc.DeviceModelDesc = parseICCDesc(data)
	case "vued":
		icc.ViewingCondDesc = parseICCDesc(data)
	case "view":
		if len(data) >= 36 {
			icc.ViewingCondIlluminant = parseICCXYZData(data[8:20])
			icc.ViewingCondSurround = parseICCXYZData(data[20:32])
			icc.ViewingCondIlluminantType = jpegEndian.Uint32(data[32:36])
		}
	case "lumi":
		icc.Luminance = parseICCXYZ(data)
	case "meas":
		if len(data) >= 36 {
			icc.MeasurementObserver = jpegEndian.Uint32(data[8:12])
			icc.MeasurementBacking = parseICCXYZData(data[12:24])
			icc.MeasurementGeometry = jpegEndian.Uint32(data[24:28])
			icc.MeasurementFlare = float64(jpegEndian.Uint32(data[28:32])) / 65536.0
			icc.MeasurementIlluminant = jpegEndian.Uint32(data[32:36])
		}
	case "tech":
		if len(data) >= 12 {
			icc.Technology = string(data[8:12])
		}
	case "rTRC":
		icc.RedTRCLength = uint32(len(data))
	case "gTRC":
		icc.GreenTRCLength = uint32(len(data))
	case "bTRC":
		icc.BlueTRCLength = uint32(len(data))
	}
}

func parseICCTime(data []byte) time.Time {
	if len(data) < 12 {
		return time.Time{}
	}
	year := int(jpegEndian.Uint16(data[0:2]))
	month := time.Month(jpegEndian.Uint16(data[2:4]))
	day := int(jpegEndian.Uint16(data[4:6]))
	hour := int(jpegEndian.Uint16(data[6:8]))
	minute := int(jpegEndian.Uint16(data[8:10]))
	second := int(jpegEndian.Uint16(data[10:12]))
	return time.Date(year, month, day, hour, minute, second, 0, time.UTC)
}

func parseICCText(data []byte) string {
	if len(data) < 8 {
		return ""
	}
	switch string(data[:4]) {
	case "text":
		return trimNULString(data[8:])
	case "desc":
		return parseICCDesc(data)
	default:
		return ""
	}
}

func parseICCDesc(data []byte) string {
	if len(data) < 12 || string(data[:4]) != "desc" {
		return parseICCText(data)
	}
	n := int(jpegEndian.Uint32(data[8:12]))
	if n <= 0 {
		return ""
	}
	start := 12
	end := start + n
	if end > len(data) {
		end = len(data)
	}
	return trimNULString(data[start:end])
}

func parseICCXYZ(data []byte) [3]float64 {
	if len(data) < 20 || string(data[:4]) != "XYZ " {
		return [3]float64{}
	}
	return parseICCXYZData(data[8:20])
}

func parseICCXYZData(data []byte) [3]float64 {
	if len(data) < 12 {
		return [3]float64{}
	}
	return [3]float64{
		s15Fixed16(utils.BigEndian, data[0:4]),
		s15Fixed16(utils.BigEndian, data[4:8]),
		s15Fixed16(utils.BigEndian, data[8:12]),
	}
}

func (icc ICCProfile) ProfileIDHex() string {
	return hex.EncodeToString(icc.ProfileID[:])
}
