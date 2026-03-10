package exif

import (
	"github.com/evanoberholster/imagemeta/meta/exif/ifd"
	"github.com/evanoberholster/imagemeta/meta/exif/makernote"
	"github.com/evanoberholster/imagemeta/meta/exif/tag"
)

// readMakerNoteDirectory parses a maker-note directory for the active camera make.
func (r *Reader) readMakerNoteDirectory(parent tag.Entry, child ifd.Directory) error {
	info := r.makerNoteInfo()
	if info.Make == makernote.CameraMakeUnknown {
		if r.Exif.CameraMakeID == makernote.CameraMakeUnknown && r.Exif.IFD0.Make != "" {
			r.Exif.CameraMakeID = makernote.IdentifyCameraMakeString(r.Exif.IFD0.Make)
		}
		info.Make = r.Exif.CameraMakeID
	}

	switch info.Make {
	case makernote.CameraMakeNikon:
		return r.readNikonMakerNoteDirectory(parent, child)
	case makernote.CameraMakeCanon, makernote.CameraMakeApple:
		// Canon and Apple maker notes are parsed as a regular IFD at the maker-note offset.
		return r.readDirectory(child, false)
	default:
		if r.debugEnabled() {
			r.debug().
				Str("make", r.Exif.CameraMake()).
				Str("makerNoteMake", info.Make.String()).
				Msg("skipping maker-note parsing for unsupported make")
		}
		return nil
	}
}

// readNikonMakerNoteDirectory parses Nikon maker-note TIFF headers and directory offsets.
func (r *Reader) readNikonMakerNoteDirectory(parent tag.Entry, child ifd.Directory) error {
	header, err := r.fastRead(makernoteNikonHeaderLength)
	if err != nil {
		return err
	}

	byteOrder, ifdRelOffset, ok := makernote.ParseNikonHeader(header)
	if !ok {
		// Nikon maker notes without the standard label are not parsed yet.
		// TODO: support unlabeled Nikon maker-note variants.
		return nil
	}

	nikonDir := ifd.New(
		byteOrder,
		ifd.MakerNoteIFD,
		child.Index,
		parent.ValueOffset+ifdRelOffset,
		parent.ValueOffset,
	)
	return r.readDirectory(nikonDir, false)
}

const makernoteNikonHeaderLength = 18
