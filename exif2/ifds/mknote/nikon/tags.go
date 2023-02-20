package nikon

import "github.com/evanoberholster/imagemeta/exif2/tag"

// TagNikonString returns the string representation of a tag.ID for Nikon Makernotes
func TagNikonString(id tag.ID) string {
	if name, ok := TagNikonIDMap[id]; ok {
		return name
	}
	return id.String()
}

// TagNikonIDMap is a Map of tag.ID to string for the NikonMakerNote tags
var TagNikonIDMap = map[tag.ID]string{}
