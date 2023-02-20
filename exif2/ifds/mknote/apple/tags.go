package apple

import "github.com/evanoberholster/imagemeta/exif2/tag"

// TagAppleString returns the string representation of a tag.ID for Apple Makernotes
func TagAppleString(id tag.ID) string {
	if name, ok := TagAppleIDMap[id]; ok {
		return name
	}
	return id.String()
}

// TagAppleIDMap is a Map of tag.ID to string for the AppleMakerNote tags
var TagAppleIDMap = map[tag.ID]string{}
