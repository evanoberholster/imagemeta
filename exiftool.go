package exiftool
import (
	"fmt"
	"log"
	"encoding/json"
	"io"
	"errors"
	"strconv"
	"time"
)

const (
	GroupExif = "EXIF"
	GroupXML = "XML"
	GroupFile = "File"
	GroupComposite = "Composite"
	GroupMakerNotes = "MakerNotes"

	gExif = "EXIF"
	gFile = "File"
	gComposite = "Composite"
	gMakerNotes = "MakerNotes"

	// Time Layout
	TimeLayoutBasic =  "2006:01:02 15:04:05"
	TimeLayoutSubSec = "2006:01:02 15:04:05.999"
	TimeLayoutGPS = "2006:01:02 15:04:05.999Z"
)

type ExifResponse struct {
	Group		map[string]ExifGroup
}

type ExifGroup struct {
	Item 	map[string]ExifItem
}

type ExifItem struct {
	desc			string
	vString			string
	nInt			int
	nFloat			float64
	vStringSlice	[]string
	error 			error
}

type ItemPath struct {
	Group 	string
	Item 	string
}

func GetWithReader(reader io.Reader) {
	resp, err := extractJson(reader)
	if err != nil {
		log.Println(err)
	}

	// Parse response keys
	data := parseResponse(resp)

	log.Println(data.ToExif())
	log.Println(data.ToComposite())
	log.Println(data.ToLocation())

	log.Println(data.ToMakerNotes())

	//log.Println(data.GetExifItem(GroupComposite, "Aperture"))
	//log.Println(data.Group[GroupFile].Item["ImageHeight"])
}

func ExtractWithReader(reader io.Reader) (ExifResponse, error) {
	resp, err := extractJson(reader)
	if err != nil {
		return ExifResponse{}, err
	}
	data := parseResponse(resp)
	return data, nil
}

func extractJson(reader io.Reader) (map[string]interface{}, error) {
	var response map[string]interface{}
	meta, err := ExtractReader("exiftool", reader, "-json", "-a", "-l", "-g")
	if err != nil {
		return response, err
	} else {
		// Remove [ brackets ]
		meta = meta[:len(meta)-2] 
		meta = meta[1:] 
		//
		if err := json.Unmarshal(meta, &response); err != nil {
			return response, err
	    }
	}
	return response, nil
}

func parseResponse(response map[string]interface{}) ExifResponse {
	var resp ExifResponse
	resp.Group = make(map[string]ExifGroup)
	for key, value := range response {
		if v, ok := value.(map[string]interface{}); ok {
			resp.Group[key], _ = parseGroup(v)
		} else {
			//log.Printf("Not Used: " + key)
		}
	}
	return resp
}

func parseGroup(group map[string]interface{}) (ExifGroup, error) {
	var item ExifGroup
	item.Item = make(map[string]ExifItem)
	for key, value := range group {
		if v, ok := value.(map[string]interface{}); ok {
			item.Item[key], _ = parseItem(v)
		} else {
			log.Printf("Unknown item: " + key)
		}
	}
	return item, nil
}

func parseItem(item map[string]interface{}) (ExifItem, error) {
	var e ExifItem
	for key, value := range item {
		switch v := value.(type) {
		case string:
			if key == "val" {
				e.setString(v)
			} else if key == "num" {
				e.setString(v)
			} 
			// Testing Purposes: Key Desription
			//if key == "desc" {
			//	e.setDesc(v)
			//}
		case int:
			if key == "val" {
				e.setInt(v)
			} else if key == "num" {
				e.setInt(v)
			}
		case float64:
			if key == "val" {
				e.setFloat(v)
				e.setInt(int(v))
				e.setString(strconv.FormatFloat(v, 'f', -1, 64))
			} else if key == "num" {
				e.setInt(int(v))
				e.setFloat(v)
			}
		case bool:
			if key == "val" {
				e.setString(strconv.FormatBool(true))
			}
		default:
			// TODO: include type []interface{} for []string
			fmt.Printf("I don't know about type %T!\n", v)
		}
	}
	return e, nil
}

func (e ExifResponse)GetExifItem(g string, i string) (ExifItem) {
	if item, ok := e.Group[g].Item[i]; ok {
		return item
	}
	return ExifItem{error: errors.New("Item \"" + g + "/" + i + "\": Does not exist")}
}

func (e ExifResponse)GetExifGroup(g string) (ExifGroup, error) {
	if group, ok := e.Group[g]; ok {
		return group, nil
	} 
	return ExifGroup{}, errors.New("Group \"" + g + "\": Does not exist")
}

func (e ExifResponse)ExifPath(paths ...ItemPath) (ExifItem) {
	for _, path := range paths {
		if item, ok := e.Group[path.Group].Item[path.Item]; ok {
			return item
		}
	}
	return ExifItem{error: errors.New("Items: Does not exist")}
}

func Path(g string, i string) ItemPath {
	return ItemPath{Group: g, Item: i}
}

// Access Path Groups
func ExifPath(i string) ItemPath { return ItemPath{Group: GroupExif, Item: i} }

func CompositePath(i string) ItemPath { return ItemPath{Group: GroupComposite, Item: i} }

func FilePath(i string) ItemPath { return ItemPath{Group: GroupFile, Item: i} }

func MakerNotesPath(i string) ItemPath { return ItemPath{Group: GroupMakerNotes, Item: i} }

func (e ExifItem) String() string {
	return fmt.Sprintf("\n Description: %s\n StringValue: %s\n NumberValue: %d\n FloatValue: %f\n", 
		e.desc, e.vString, e.nInt, e.nFloat)
}

// ExifItem Setters
func (e *ExifItem)setDesc(v string) { e.desc = v }

func (e *ExifItem)setString(v string) { e.vString = v }

func (e *ExifItem)setInt(v int) { e.nInt = v }

func (e *ExifItem)setFloat(v float64) { e.nFloat = v }

// ExifItem Getters
func (e ExifItem) ToString() (string) { return e.vString }

func (e ExifItem) ToInt() (int) { return e.nInt }

func (e ExifItem) ToFloat() (float64)  { return e.nFloat }

func (e ExifItem) ToTime() (time.Time) {
	timeLayouts := []string{TimeLayoutSubSec, TimeLayoutGPS, TimeLayoutBasic}
	for _, layout := range timeLayouts {
		t, err := time.Parse(layout, e.vString)
		if err == nil {
			return t
		}
	}
	return time.Time{} // Return nil datetime
}




