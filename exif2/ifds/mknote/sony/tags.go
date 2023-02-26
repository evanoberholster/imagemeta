package sony

import "github.com/evanoberholster/imagemeta/exif2/tag"

// TagSonyIDMap is a Map of tag.ID to string for the SonyMakerNote tags
var TagSonyIDMap = map[tag.ID]string{}

// TagSonyString returns the string representation of a tag.ID for Sony Makernotes
func TagSonyString(id tag.ID) string {
	if name, ok := TagSonyIDMap[id]; ok {
		return name
	}
	return id.String()
}
