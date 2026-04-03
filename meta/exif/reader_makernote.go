package exif

import (
	"strings"

	"github.com/evanoberholster/imagemeta/meta/exif/ifd"
	"github.com/evanoberholster/imagemeta/meta/exif/makernote"
	"github.com/evanoberholster/imagemeta/meta/exif/tag"
	"github.com/evanoberholster/imagemeta/meta/utils"
)

// readMakerNoteDirectory parses a maker-note directory for the active camera make.
func (r *Reader) readMakerNoteDirectory(parent tag.Entry, child ifd.Directory) error {
	makeID := r.ensureMakerNoteMake()
	switch makeID {
	case makernote.CameraMakeNikon:
		return r.readNikonMakerNoteDirectory(parent, child)
	case makernote.CameraMakeCanon:
		return r.readCanonMakerNoteDirectory(parent, child)
	case makernote.CameraMakeApple:
		// Apple maker notes are parsed as a regular IFD at the maker-note offset.
		return r.readDirectory(child, false)
	default:
		if r.debugEnabled() {
			r.debug().
				Str("make", r.Exif.CameraMake()).
				Str("makerNoteMake", makeID.String()).
				Msg("skipping maker-note parsing for unsupported make")
		}
		return nil
	}
}

// ensureMakerNoteMake resolves and caches maker-note camera make.
func (r *Reader) ensureMakerNoteMake() makernote.CameraMake {
	info := r.makerNoteInfo()
	if info.Make != makernote.CameraMakeUnknown {
		return info.Make
	}

	if r.Exif.CameraMakeID == makernote.CameraMakeUnknown && r.Exif.IFD0.Make != "" {
		r.Exif.CameraMakeID = makernote.IdentifyCameraMakeString(r.Exif.IFD0.Make)
	}
	// Some files provide a canonical model string but omit IFD0:Make. Infer
	// make from the first model token to unlock maker-note decoding.
	if r.Exif.CameraMakeID == makernote.CameraMakeUnknown && r.Exif.IFD0.Model != "" {
		model := strings.TrimSpace(r.Exif.IFD0.Model)
		if i := strings.IndexByte(model, ' '); i > 0 {
			model = model[:i]
		}
		if model != "" {
			r.Exif.CameraMakeID = makernote.IdentifyCameraMakeString(model)
		}
	}

	info.Make = r.Exif.CameraMakeID
	return info.Make
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

// readCanonMakerNoteDirectory parses Canon maker-note variants.
//
// Canon maker notes may be one of:
//  1. Raw IFD at maker-note offset (classic JPEG/CR2 variants).
//  2. "Canon\0\0\0" prefixed IFD payload.
//  3. Embedded TIFF header ("II*\0"/"MM\0*") with a relative MakerNote IFD offset
//     used by CR3.
func (r *Reader) readCanonMakerNoteDirectory(parent tag.Entry, child ifd.Directory) error {
	header, ok := r.peekMakerNotePrefix(canonMakerNotePrefixLength)
	if !ok {
		return r.readDirectory(child, false)
	}

	// CR3: MakerNote payload starts with an embedded TIFF header.
	if byteOrder, ifdRelOffset, valid := parseMakerNoteTIFFPrefix(header); valid {
		if ifdRelOffset >= uint32(canonMakerNotePrefixLength) && (parent.UnitCount == 0 || ifdRelOffset < parent.UnitCount) {
			if err := r.discard(int(ifdRelOffset)); err != nil {
				return err
			}
			embedded := ifd.New(
				byteOrder,
				ifd.MakerNoteIFD,
				child.Index,
				parent.ValueOffset+ifdRelOffset,
				parent.ValueOffset,
			)
			return r.readDirectory(embedded, false)
		}
	}

	// Older Canon: fixed "Canon\0\0\0" prefix before IFD entries.
	if isCanonMakerNotePrefix(header) {
		if err := r.discard(canonMakerNotePrefixLength); err != nil {
			return err
		}
		prefixed := ifd.New(
			child.ByteOrder,
			ifd.MakerNoteIFD,
			child.Index,
			parent.ValueOffset+canonMakerNotePrefixLength,
			parent.ValueOffset,
		)
		return r.readDirectory(prefixed, false)
	}

	// Default Canon behavior (raw IFD at maker-note offset).
	return r.readDirectory(child, false)
}

func (r *Reader) peekMakerNotePrefix(n int) ([]byte, bool) {
	if r.reader == nil {
		return nil, false
	}
	buf, err := r.reader.Peek(n)
	if err != nil || len(buf) < n {
		return nil, false
	}
	return buf[:n], true
}

func parseMakerNoteTIFFPrefix(prefix []byte) (byteOrder utils.ByteOrder, ifdRelOffset uint32, ok bool) {
	if len(prefix) < canonMakerNotePrefixLength {
		return utils.UnknownEndian, 0, false
	}
	byteOrder = utils.BinaryOrder(prefix[:4])
	if byteOrder == utils.UnknownEndian {
		return utils.UnknownEndian, 0, false
	}
	return byteOrder, byteOrder.Uint32(prefix[4:8]), true
}

func isCanonMakerNotePrefix(prefix []byte) bool {
	return len(prefix) >= canonMakerNotePrefixLength &&
		prefix[0] == 'C' &&
		prefix[1] == 'a' &&
		prefix[2] == 'n' &&
		prefix[3] == 'o' &&
		prefix[4] == 'n' &&
		prefix[5] == 0 &&
		prefix[6] == 0 &&
		prefix[7] == 0
}

const makernoteNikonHeaderLength = 18
const canonMakerNotePrefixLength = 8
