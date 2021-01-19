package xml

import (
	"encoding/xml"
	"strconv"
	"time"
)

func decodeRDF(decoder *xml.Decoder, start *xml.StartElement) (strs []string) {
	var t xml.Token
	var err error
ReadLbl:
	if t, err = decoder.RawToken(); err != nil {
		panic(err)
	}
	switch x := t.(type) {
	case xml.StartElement:
		goto ReadLbl
	case xml.EndElement:
		if x.Name != start.Name {
			goto ReadLbl
		}
	case xml.CharData:
		if x[0] == 10 || x[0] == 32 {
			goto ReadLbl
		}
		strs = append(strs, string(x))
		goto ReadLbl
	}
	return strs
}

func parseDate(str string) (t time.Time, err error) {
	if t, err = time.Parse("2006-01-02T15:04:05Z07:00", str); err != nil {
		if t, err = time.Parse("2006-01-02T15:04:05.00", str); err != nil {
			t, err = time.Parse("2006-01-02T15:04:05", str)
		}
	}
	return
}

func parseUint32(s string) uint32 {
	u64, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return 0
	}
	return uint32(u64)
}
