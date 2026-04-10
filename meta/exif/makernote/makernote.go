package makernote

import (
	"github.com/evanoberholster/imagemeta/meta/exif/makernote/nikon"
	"github.com/evanoberholster/imagemeta/meta/exif/makernote/panasonic"
	"github.com/evanoberholster/imagemeta/meta/exif/makernote/sony"
)

// Info contains parsed maker-note values for supported vendors.
type Info struct {
	Make      CameraMake
	Apple     *Apple
	Canon     *Canon
	Nikon     *nikon.Nikon
	Panasonic *panasonic.Panasonic
	Sony      *sony.Sony
}

// Apple contains selected Apple maker-note fields.
type Apple struct {
	RunTime           string
	BurstUUID         string
	ContentIdentifier string
	ImageUniqueID     string

	MakerNoteVersion int32
	AETarget         int32
	AEAverage        int32
	OISMode          int32
	ImageCaptureType int32
	AEStable         bool
	AFStable         bool
}
