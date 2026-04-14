package jpeg

import (
	"encoding/hex"
	"strings"

	"github.com/evanoberholster/imagemeta/meta/utils"
)

// Photoshop stores selected Photoshop Image Resource values from APP13.
type Photoshop struct {
	XResolution              float64
	DisplayedUnitsX          uint16
	YResolution              float64
	DisplayedUnitsY          uint16
	CopyrightFlag            uint8
	CopyrightFlagSet         bool
	URL                      string
	PhotoshopThumbnailLength uint32
	IPTCDigest               string
	PhotoshopQuality         int16
	PhotoshopFormat          int16
	ProgressiveScans         int16
}

func parsePhotoshop(payload []byte) (*Photoshop, *IPTC, error) {
	if !isPhotoshopPayload(payload) {
		return nil, nil, errShortSegment("Photoshop")
	}
	p := &Photoshop{}
	var iptc *IPTC
	pos := len(photoshopPrefix)
	for pos+12 <= len(payload) {
		sig := string(payload[pos : pos+4])
		if sig != "8BIM" && sig != "8B64" {
			break
		}
		pos += 4
		id := jpegEndian.Uint16(payload[pos : pos+2])
		pos += 2

		if pos >= len(payload) {
			break
		}
		nameLen := int(payload[pos])
		nameBytes := 1 + nameLen
		if nameBytes&1 != 0 {
			nameBytes++
		}
		pos += nameBytes
		if pos+4 > len(payload) {
			break
		}
		size := int(jpegEndian.Uint32(payload[pos : pos+4]))
		pos += 4
		if size < 0 || pos+size > len(payload) {
			break
		}
		data := payload[pos : pos+size]
		parsePhotoshopResource(p, &iptc, id, data)
		pos += size
		if pos&1 != 0 {
			pos++
		}
	}
	if isEmptyPhotoshop(p) {
		p = nil
	}
	return p, iptc, nil
}

func parsePhotoshopResource(p *Photoshop, iptc **IPTC, id uint16, data []byte) {
	switch id {
	case 0x03ed:
		if len(data) >= 16 {
			p.XResolution = fixed16(utils.BigEndian, data[0:4])
			p.DisplayedUnitsX = jpegEndian.Uint16(data[4:6])
			p.YResolution = fixed16(utils.BigEndian, data[8:12])
			p.DisplayedUnitsY = jpegEndian.Uint16(data[12:14])
		}
	case 0x0404:
		parsed := parseIPTC(data)
		if !parsed.empty() {
			*iptc = parsed
		}
	case 0x0406:
		if len(data) >= 2 {
			p.PhotoshopQuality = int16(jpegEndian.Uint16(data[0:2]))
		}
		if len(data) >= 4 {
			p.PhotoshopFormat = int16(jpegEndian.Uint16(data[2:4]))
		}
		if len(data) >= 6 {
			p.ProgressiveScans = int16(jpegEndian.Uint16(data[4:6]))
		}
	case 0x040a:
		if len(data) > 0 {
			p.CopyrightFlag = data[0]
			p.CopyrightFlagSet = true
		}
	case 0x040b:
		p.URL = strings.TrimRight(string(data), "\x00")
	case 0x040c:
		p.PhotoshopThumbnailLength = uint32(len(data))
		if len(data) > 28 {
			p.PhotoshopThumbnailLength = uint32(len(data) - 28)
		}
	case 0x0425:
		p.IPTCDigest = hex.EncodeToString(data)
	}
}

func isEmptyPhotoshop(p *Photoshop) bool {
	return p == nil ||
		(p.XResolution == 0 &&
			p.DisplayedUnitsX == 0 &&
			p.YResolution == 0 &&
			p.DisplayedUnitsY == 0 &&
			!p.CopyrightFlagSet &&
			p.URL == "" &&
			p.PhotoshopThumbnailLength == 0 &&
			p.IPTCDigest == "" &&
			p.PhotoshopQuality == 0 &&
			p.PhotoshopFormat == 0 &&
			p.ProgressiveScans == 0)
}

// IPTC stores selected IPTC IIM application-record values from APP13.
type IPTC struct {
	CodedCharacterSet             string
	EnvelopeRecordVersion         uint16
	ApplicationRecordVersion      uint16
	Keywords                      []string
	DateCreated                   string
	TimeCreated                   string
	DigitalCreationDate           string
	DigitalCreationTime           string
	ByLine                        string
	ProvinceState                 string
	CountryPrimaryLocationName    string
	OriginalTransmissionReference string
	OriginatingProgram            string
	Credit                        string
	CopyrightNotice               string
	CaptionAbstract               string
	Prefs                         string
}

func parseIPTC(data []byte) *IPTC {
	iptc := &IPTC{}
	for pos := 0; pos+5 <= len(data); {
		if data[pos] != 0x1c {
			pos++
			continue
		}
		record := data[pos+1]
		dataset := data[pos+2]
		size := int(jpegEndian.Uint16(data[pos+3 : pos+5]))
		pos += 5
		if size&0x8000 != 0 {
			byteCount := size & 0x7fff
			if byteCount <= 0 || byteCount > 4 || pos+byteCount > len(data) {
				break
			}
			size = 0
			for i := 0; i < byteCount; i++ {
				size = (size << 8) | int(data[pos+i])
			}
			pos += byteCount
		}
		if size < 0 || pos+size > len(data) {
			break
		}
		parseIPTCDataset(iptc, record, dataset, data[pos:pos+size])
		pos += size
	}
	return iptc
}

func parseIPTCDataset(iptc *IPTC, record, dataset uint8, value []byte) {
	if record == 1 && dataset == 90 {
		iptc.CodedCharacterSet = string(value)
		return
	}
	if record == 1 && dataset == 0 {
		if len(value) >= 2 {
			iptc.EnvelopeRecordVersion = jpegEndian.Uint16(value)
		}
		return
	}
	if record != 2 {
		return
	}
	switch dataset {
	case 0:
		if len(value) >= 2 {
			iptc.ApplicationRecordVersion = jpegEndian.Uint16(value)
		}
	case 25:
		iptc.Keywords = append(iptc.Keywords, string(value))
	case 55:
		iptc.DateCreated = iptcDate(value)
	case 60:
		iptc.TimeCreated = iptcTime(value)
	case 62:
		iptc.DigitalCreationDate = iptcDate(value)
	case 63:
		iptc.DigitalCreationTime = iptcTime(value)
	case 65:
		iptc.OriginatingProgram = string(value)
	case 80:
		iptc.ByLine = string(value)
	case 95:
		iptc.ProvinceState = string(value)
	case 101:
		iptc.CountryPrimaryLocationName = string(value)
	case 103:
		iptc.OriginalTransmissionReference = string(value)
	case 110:
		iptc.Credit = string(value)
	case 116:
		iptc.CopyrightNotice = string(value)
	case 120:
		iptc.CaptionAbstract = string(value)
	case 221:
		iptc.Prefs = string(value)
	}
}

func (iptc *IPTC) empty() bool {
	return iptc == nil ||
		(iptc.CodedCharacterSet == "" &&
			iptc.EnvelopeRecordVersion == 0 &&
			iptc.ApplicationRecordVersion == 0 &&
			len(iptc.Keywords) == 0 &&
			iptc.DateCreated == "" &&
			iptc.TimeCreated == "" &&
			iptc.DigitalCreationDate == "" &&
			iptc.DigitalCreationTime == "" &&
			iptc.ByLine == "" &&
			iptc.ProvinceState == "" &&
			iptc.CountryPrimaryLocationName == "" &&
			iptc.OriginalTransmissionReference == "" &&
			iptc.OriginatingProgram == "" &&
			iptc.Credit == "" &&
			iptc.CopyrightNotice == "" &&
			iptc.CaptionAbstract == "" &&
			iptc.Prefs == "")
}

func iptcDate(value []byte) string {
	s := string(value)
	if len(s) == 8 {
		return s[:4] + ":" + s[4:6] + ":" + s[6:8]
	}
	return s
}

func iptcTime(value []byte) string {
	s := string(value)
	if len(s) == 11 && (s[6] == '+' || s[6] == '-') {
		return s[:2] + ":" + s[2:4] + ":" + s[4:6] + s[6:9] + ":" + s[9:11]
	}
	if len(s) >= 6 {
		return s[:2] + ":" + s[2:4] + ":" + s[4:]
	}
	return s
}
