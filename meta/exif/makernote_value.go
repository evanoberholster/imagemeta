package exif

import "github.com/evanoberholster/imagemeta/meta/exif/makernote"

// makerNoteInfo returns the typed maker-note container from Exif.MakerNote.
//
// Exif.MakerNote is a pointer to the generic makernote.Makernote interface
// value; parser code normalizes it to *makernote.Info for field updates.
func (r *Reader) makerNoteInfo() *makernote.Info {
	return &r.Exif.MakerNote
}
