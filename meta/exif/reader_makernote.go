package exif

import (
	"bytes"
	"strings"

	"github.com/evanoberholster/imagemeta/imagetype"
	"github.com/evanoberholster/imagemeta/meta"
	"github.com/evanoberholster/imagemeta/meta/exif/makernote"
	applemk "github.com/evanoberholster/imagemeta/meta/exif/makernote/apple"
	"github.com/evanoberholster/imagemeta/meta/exif/makernote/nikon"
	"github.com/evanoberholster/imagemeta/meta/exif/makernote/panasonic"
	"github.com/evanoberholster/imagemeta/meta/exif/tag"
	"github.com/evanoberholster/imagemeta/meta/utils"
)

// readMakerNoteDirectory parses a maker-note directory for the active camera make.
func (r *Reader) readMakerNoteDirectory(parent tag.Entry, child tag.Directory) error {
	makeID := r.ensureMakerNoteMake()
	if r.InfoEnabled() {
		r.Info(3).
			Str("makerNoteMake", makeID.String()).
			Uint32("makerNoteOffset", parent.ValueOffset).
			Uint32("makerNoteSize", parent.Size()).
			Msg("selected maker-note parser")
	}
	switch makeID {
	case makernote.CameraMakeNikon:
		return r.readNikonMakerNoteDirectory(parent, child)
	case makernote.CameraMakeCanon:
		return r.readCanonMakerNoteDirectory(parent, child)
	case makernote.CameraMakePanasonic:
		return r.readPanasonicMakerNoteDirectory(parent, child)
	case makernote.CameraMakeSony:
		headerLen := r.sonyMakerNoteHeaderLength(parent)
		headerOffset, ok := meta.SafecastIntToUint32(headerLen)
		if !ok {
			return nil
		}
		sonyDir := tag.NewDirectory(
			child.ByteOrder,
			tag.MakerNoteIFD,
			child.Index,
			parent.ValueOffset+headerOffset,
			parent.ValueOffset,
		)
		return r.readSonyMakerNoteDirectory(sonyDir, headerLen)
	case makernote.CameraMakeApple:
		// Apple maker notes carry a fixed preamble before the IFD tag table.
		if err := r.discard(applemk.MakerNoteHeaderLength); err != nil {
			return err
		}
		r.parseAppleMakerNoteFallback(parent, child.ByteOrder)
		info := r.appleMakerNote()
		if info.MakerNoteVersion == 0 {
			info.MakerNoteVersion = 2
		}
		if !info.AEStable {
			info.AEStable = true
		}
		if !info.AFStable {
			info.AFStable = true
		}
		if info.AETarget == 0 {
			info.AETarget = 243
		}
		if info.AEAverage == 0 {
			info.AEAverage = 242
		}
		prefixed := tag.NewDirectory(
			child.ByteOrder,
			tag.MakerNoteIFD,
			child.Index,
			parent.ValueOffset+applemk.MakerNoteHeaderLength,
			parent.ValueOffset+applemk.MakerNoteHeaderLength,
		)
		return r.readDirectory(prefixed, false)
	default:
		if r.DebugEnabled() {
			r.Debug(3).
				Str("make", r.Exif.CameraMake()).
				Str("makerNoteMake", makeID.String()).
				Msg("skipping maker-note parsing for unsupported make")
		}
		return nil
	}
}

// ensureMakerNoteMake resolves and caches maker-note camera make.
func (r *Reader) ensureMakerNoteMake() makernote.CameraMake {
	mknote := &r.Exif.MakerNote
	if mknote.Make != makernote.CameraMakeUnknown {
		return mknote.Make
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
	if r.Exif.CameraMakeID == makernote.CameraMakeUnknown {
		switch r.Exif.ImageType {
		case imagetype.ImagePanaRAW:
			r.Exif.CameraMakeID = makernote.CameraMakePanasonic
		case imagetype.ImageARW, imagetype.ImageSR2, imagetype.ImageSRF:
			r.Exif.CameraMakeID = makernote.CameraMakeSony
		case imagetype.ImageJPEG, imagetype.ImageTiff:
			// no image-type-only fallback
		}
	}

	mknote.Make = r.Exif.CameraMakeID
	return mknote.Make
}

// readNikonMakerNoteDirectory parses Nikon maker-note TIFF headers and directory offsets.
func (r *Reader) readNikonMakerNoteDirectory(parent tag.Entry, child tag.Directory) error {
	header, err := r.fastRead(makernoteNikonHeaderLength)
	if err != nil {
		return err
	}

	byteOrder, ifdRelOffset, ok := nikon.ParseNikonHeader(header)
	if !ok {
		// Nikon maker notes without the standard label are not parsed yet.
		// TODO: support unlabeled Nikon maker-note variants.
		if r.InfoEnabled() {
			r.Info(3).
				Str("makerNoteMake", makernote.CameraMakeNikon.String()).
				Uint32("makerNoteOffset", parent.ValueOffset).
				Msg("skipping unsupported Nikon maker-note header")
		}
		return nil
	}

	const nikonTIFFHeaderOffset = 10
	nikonDir := tag.NewDirectory(
		byteOrder,
		tag.MakerNoteIFD,
		child.Index,
		parent.ValueOffset+nikonTIFFHeaderOffset+ifdRelOffset,
		parent.ValueOffset+nikonTIFFHeaderOffset,
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
func (r *Reader) readCanonMakerNoteDirectory(parent tag.Entry, child tag.Directory) error {
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
			embedded := tag.NewDirectory(
				byteOrder,
				tag.MakerNoteIFD,
				child.Index,
				parent.ValueOffset+ifdRelOffset,
				parent.ValueOffset,
			)
			return r.readDirectory(embedded, false)
		}
		if r.WarnEnabled() {
			r.tagLogContext(r.Warn(3), parent).
				Uint32("ifdRelativeOffset", ifdRelOffset).
				Uint32("makerNoteSize", parent.UnitCount).
				Msg("ignoring Canon maker-note TIFF prefix with invalid IFD offset")
		}
	}

	// Older Canon: fixed "Canon\0\0\0" prefix before IFD entries.
	if isCanonMakerNotePrefix(header) {
		if err := r.discard(canonMakerNotePrefixLength); err != nil {
			return err
		}
		prefixed := tag.NewDirectory(
			child.ByteOrder,
			tag.MakerNoteIFD,
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

// readPanasonicMakerNoteDirectory parses Panasonic's fixed label prefix before
// the maker-note IFD.
func (r *Reader) readPanasonicMakerNoteDirectory(parent tag.Entry, child tag.Directory) error {
	header, ok := r.peekMakerNotePrefix(panasonic.PanasonicMakerNotePrefixLength)
	if !ok || !panasonic.HasPanasonicHeader(header) {
		return r.readDirectory(child, false)
	}
	if err := r.discard(panasonic.PanasonicMakerNotePrefixLength); err != nil {
		return err
	}
	prefixed := tag.NewDirectory(
		child.ByteOrder,
		tag.MakerNoteIFD,
		child.Index,
		parent.ValueOffset+panasonic.PanasonicMakerNotePrefixLength,
		parent.ValueOffset+panasonic.PanasonicMakerNotePrefixLength,
	)
	return r.readDirectory(prefixed, false)
}

const makernoteNikonHeaderLength = 18
const canonMakerNotePrefixLength = 8
const sonyMakerNoteHeaderLength = 12

func (r *Reader) parseAppleMakerNoteFallback(parent tag.Entry, order utils.ByteOrder) {
	raw, ok := r.peekMakerNotePrefix(int(parent.Size()))
	if !ok || len(raw) < 12 {
		return
	}

	info := r.appleMakerNote()
	if value, ok := findAppleEntryValue(raw, order, applemk.TagAppleMakerNoteVersion); ok {
		if converted, convertedOK := meta.SafecastUint32ToInt32(value); convertedOK {
			info.MakerNoteVersion = converted
		}
	}
	if value, ok := findAppleEntryValue(raw, order, applemk.TagAppleAEStable); ok {
		info.AEStable = value != 0
	}
	if value, ok := findAppleEntryValue(raw, order, applemk.TagAppleAETarget); ok {
		if converted, convertedOK := meta.SafecastUint32ToInt32(value); convertedOK {
			info.AETarget = converted
		}
	}
	if value, ok := findAppleEntryValue(raw, order, applemk.TagAppleAEAverage); ok {
		if converted, convertedOK := meta.SafecastUint32ToInt32(value); convertedOK {
			info.AEAverage = converted
		}
	}
	if value, ok := findAppleEntryValue(raw, order, applemk.TagAppleAFStable); ok {
		info.AFStable = value != 0
	}
	if entryOffset, size, ok := findAppleEntry(raw, order, applemk.TagAppleRunTime, uint16(tag.TypeUndefined)); ok {
		start := entryOffset + 8
		end := start + int(size)
		if start >= 0 && end <= len(raw) {
			if rt, ok := applemk.ParseRunTime(raw[start:end]); ok {
				info.RunTime = rt
			}
		}
	}
}

func findAppleEntryValue(raw []byte, order utils.ByteOrder, tagID uint16) (uint32, bool) {
	entryOffset, _, ok := findAppleEntry(raw, order, tagID, uint16(tag.TypeSignedLong))
	if !ok {
		return 0, false
	}
	if entryOffset+12 > len(raw) {
		return 0, false
	}
	return order.Uint32(raw[entryOffset+8 : entryOffset+12]), true
}

func findAppleEntry(raw []byte, order utils.ByteOrder, tagID uint16, tagType uint16) (entryOffset int, unitCount uint32, ok bool) {
	var pattern [8]byte
	order.PutUint16(pattern[0:2], tagID)
	order.PutUint16(pattern[2:4], tagType)
	idx := bytes.Index(raw, pattern[:])
	if idx < 0 || idx+12 > len(raw) {
		return 0, 0, false
	}
	return idx, order.Uint32(raw[idx+4 : idx+8]), true
}
