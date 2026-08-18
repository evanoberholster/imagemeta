package imagemeta

import (
	"io"
	"sort"

	"github.com/evanoberholster/imagemeta/meta"
	"github.com/evanoberholster/imagemeta/meta/exif"
	"github.com/pkg/errors"
)

// ErrNoPreview is returned when a file contains no embedded preview image.
var ErrNoPreview = errors.New("error no embedded preview image found")

// PreviewImage describes an embedded JPEG preview of a TIFF-based (raw)
// file. Offset and Length locate the JPEG stream in the file.
type PreviewImage struct {
	// IFD names the directory the preview was found in ("IFD0", "IFD1",
	// "SubIFD1", ...).
	IFD string
	// Width and Height are the dimensions declared by the directory; zero
	// when the directory does not declare them.
	Width  uint32
	Height uint32
	Offset uint32
	Length uint32
}

// TIFFPreviews lists the embedded JPEG preview images of a TIFF-based raw
// file (e.g. NEF, CR2, ARW, DNG), largest first. The offsets are relative
// to the start of the file.
func TIFFPreviews(r io.ReadSeeker) ([]PreviewImage, error) {
	if _, err := r.Seek(0, io.SeekStart); err != nil {
		return nil, errors.Wrapf(err, "seek to file start")
	}
	ex, err := DecodeTiff(r)
	if err != nil {
		return nil, err
	}

	var previews []PreviewImage
	// IFD0 carries the JPEGInterchangeFormat pair in dedicated fields; the
	// pointed-to stream is a JPEG per TIFF/EP, no compression tag needed.
	if ex.IFD0.ThumbnailOffset != 0 && ex.IFD0.ThumbnailLength != 0 {
		previews = append(previews, PreviewImage{
			IFD:    "IFD0",
			Offset: ex.IFD0.ThumbnailOffset,
			Length: ex.IFD0.ThumbnailLength,
		})
	}
	addImageIFD := func(name string, ifd *exif.ImageIFD) {
		if ifd == nil || ifd.ImageOffset == 0 || ifd.ImageLength == 0 {
			return
		}
		// ImageOffset also captures strip-based (e.g. uncompressed)
		// images; only JPEG streams qualify as extractable previews
		if !isJPEGCompression(ifd.Compression) {
			return
		}
		previews = append(previews, PreviewImage{
			IFD:    name,
			Width:  ifd.ImageWidth,
			Height: ifd.ImageHeight,
			Offset: ifd.ImageOffset,
			Length: ifd.ImageLength,
		})
	}
	addImageIFD("IFD1", ex.IFD1)
	addImageIFD("IFD2", ex.IFD2)
	subIFDNames := [8]string{
		"SubIFD0", "SubIFD1", "SubIFD2", "SubIFD3",
		"SubIFD4", "SubIFD5", "SubIFD6", "SubIFD7",
	}
	for i, sub := range ex.SubIFDs {
		addImageIFD(subIFDNames[i], sub)
	}

	sort.Slice(previews, func(i, j int) bool { return previews[i].Length > previews[j].Length })
	return previews, nil
}

// ExtractPreview reads the bytes of one preview and validates the JPEG SOI
// marker. The returned slice is newly allocated.
func ExtractPreview(r io.ReadSeeker, p PreviewImage) ([]byte, error) {
	if p.Offset == 0 || p.Length == 0 {
		return nil, ErrNoPreview
	}
	if _, err := r.Seek(int64(p.Offset), io.SeekStart); err != nil {
		return nil, errors.Wrapf(err, "seek to preview offset %d", p.Offset)
	}
	buf := make([]byte, p.Length)
	if _, err := io.ReadFull(r, buf); err != nil {
		return nil, errors.Wrapf(err, "read preview at offset %d", p.Offset)
	}
	if len(buf) < 3 || buf[0] != 0xff || buf[1] != 0xd8 {
		return nil, ErrNoPreview
	}
	return buf, nil
}

// PreviewTIFF returns the largest embedded JPEG preview of a TIFF-based
// (raw) file, e.g. NEF/CR2/ARW/DNG, analog to PreviewCR3.
func PreviewTIFF(r io.ReadSeeker) ([]byte, error) {
	previews, err := TIFFPreviews(r)
	if err != nil {
		return nil, err
	}
	for _, p := range previews {
		buf, err := ExtractPreview(r, p)
		if err == nil {
			return buf, nil
		}
	}
	return nil, ErrNoPreview
}

// isJPEGCompression reports whether the TIFF compression value denotes an
// extractable JPEG stream (old-style, new-style, and their vendor aliases).
func isJPEGCompression(c meta.Compression) bool {
	switch c {
	case 6, 7, 99, 34892:
		return true
	}
	return false
}
