package exif

import (
	"encoding/json"
	"sort"

	"github.com/evanoberholster/imagemeta/exif/ifds"
	"github.com/evanoberholster/imagemeta/exif/tag"
	"github.com/evanoberholster/imagemeta/imagetype"
)

// DebugJSON implements the JSONMarshaler interface that is used by encoding/json
// This is used primarily for testing and debuging.
func (e *Data) DebugJSON() ([]byte, error) {
	je := debugExif{It: e.imageType, Make: e.CameraMake(), Model: e.CameraModel(), Width: e.ImageWidth(), Height: e.ImageHeight()}
	for k, t := range e.tagMap {
		ifd, ifdIndex, _ := k.Val()
		value := e.GetTagValue(t)
		je.addTag(ifd, ifdIndex, t, value)
	}

	return json.Marshal(je)
}

func (ji *debugIfds) insertSorted(e debugTags) {
	i := sort.Search(len(ji.Tags), func(i int) bool { return ji.Tags[i].ID > e.ID })
	ji.Tags = append(ji.Tags, debugTags{})
	copy(ji.Tags[i+1:], ji.Tags[i:])
	ji.Tags[i] = e
}

func (je *debugExif) addTag(ifd ifds.IfdType, ifdIndex uint8, t tag.Tag, v interface{}) {
	if je.Ifds == nil {
		je.Ifds = make(map[string]map[uint8]debugIfds)
	}
	ji, ok := je.Ifds[ifd.String()]
	if !ok {
		je.Ifds[ifd.String()] = make(map[uint8]debugIfds)
		ji = je.Ifds[ifd.String()]
	}
	jm, ok := ji[ifdIndex]
	if !ok {
		ji[ifdIndex] = debugIfds{make([]debugTags, 0)}
		jm = ji[ifdIndex]
	}
	jm.insertSorted(debugTags{Name: ifd.TagName(t.ID), Type: t.Type(), ID: t.ID, Count: t.UnitCount, Value: v})
	je.Ifds[ifd.String()][ifdIndex] = jm
}

type debugIfds struct {
	Tags []debugTags `json:"Tags"`
}

type debugExif struct {
	Ifds   map[string]map[uint8]debugIfds `json:"Ifds"`
	It     imagetype.ImageType            `json:"ImageType"`
	Make   string
	Model  string
	Width  uint16
	Height uint16
}

type debugTags struct {
	ID    tag.ID
	Name  string
	Count uint32
	Type  tag.Type
	Value interface{}
}

func (jt debugTags) MarshalJSON() ([]byte, error) {
	st := struct {
		ID    string
		Name  string
		Count uint32
		Type  string
		Value interface{} `json:"Val"`
	}{
		ID:    jt.ID.String(),
		Name:  jt.Name,
		Count: jt.Count,
		Type:  jt.Type.String(),
		Value: jt.Value,
	}
	return json.Marshal(st)
}
