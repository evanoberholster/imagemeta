package makernote

import (
	"github.com/evanoberholster/imagemeta/meta/exif/makernote/canon"
	"github.com/evanoberholster/imagemeta/meta/exif/makernote/nikon"
	"github.com/evanoberholster/imagemeta/meta/exif/makernote/panasonic"
	"github.com/evanoberholster/imagemeta/meta/exif/makernote/sony"
)

// Info contains parsed maker-note values for supported vendors.
type Info struct {
	Make      CameraMake
	Apple     *Apple
	Canon     *canon.Canon
	Nikon     *nikon.Nikon
	Panasonic *panasonic.Panasonic
	Sony      *sony.Sony
}
