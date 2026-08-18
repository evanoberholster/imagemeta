package imagemeta

import (
	"io"
	"sort"

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

// sofScanLimit bounds how many bytes of a candidate stream are inspected
// for the JPEG SOF marker; metadata segments (APPn, DQT, DHT) rarely exceed
// a few KB before the SOF appears.
const sofScanLimit = 64 * 1024

// maxJPEGSegments bounds the header segments walked before the SOF marker;
// real files reach it within a handful, this bounds work on crafted streams.
const maxJPEGSegments = 32

// TIFFPreviews lists the embedded JPEG preview images of a TIFF-based raw
// file (e.g. NEF, CR2, ARW, DNG, PEF), largest first. The offsets are
// relative to the start of the file. Candidates are validated against the
// actual stream: only renderable JPEGs qualify, which keeps raw sensor
// payloads out (DNG stores those as lossless JPEG, uncompressed or
// vendor-compressed strips).
func TIFFPreviews(r io.ReadSeeker) ([]PreviewImage, error) {
	if _, err := r.Seek(0, io.SeekStart); err != nil {
		return nil, errors.Wrapf(err, "seek to file start")
	}
	ex, err := DecodeTiff(r)
	if err != nil {
		return nil, err
	}

	var previews []PreviewImage
	add := func(name string, width, height, offset, length uint32) {
		if offset == 0 || length == 0 {
			return
		}
		if !isRenderableJPEGAt(r, offset, length) {
			return
		}
		previews = append(previews, PreviewImage{
			IFD:    name,
			Width:  width,
			Height: height,
			Offset: offset,
			Length: length,
		})
	}

	// IFD0 can carry both the JPEGInterchangeFormat pair (thumbnails) and a
	// strip-based image (CR2 stores its full-size JPEG that way)
	add("IFD0", 0, 0, ex.IFD0.ThumbnailOffset, ex.IFD0.ThumbnailLength)
	add("IFD0", ex.IFD0.ImageWidth, ex.IFD0.ImageHeight, ex.IFD0.ImageOffset, ex.IFD0.ImageLength)
	addImageIFD := func(name string, ifd *exif.ImageIFD) {
		if ifd == nil {
			return
		}
		add(name, ifd.ImageWidth, ifd.ImageHeight, ifd.ImageOffset, ifd.ImageLength)
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
	// Length is attacker-controlled; reject a preview whose extent lies
	// outside the file before allocating, so a bogus huge length cannot force
	// a large allocation
	size, err := r.Seek(0, io.SeekEnd)
	if err != nil {
		return nil, errors.Wrapf(err, "determine file size")
	}
	if int64(p.Offset)+int64(p.Length) > size {
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

// isRenderableJPEGAt reads the head of a candidate stream and reports
// whether it is a JPEG that common decoders can render.
func isRenderableJPEGAt(r io.ReadSeeker, offset, length uint32) bool {
	if _, err := r.Seek(int64(offset), io.SeekStart); err != nil {
		return false
	}
	window := min(int(length), sofScanLimit)
	buf := make([]byte, window)
	n, err := io.ReadFull(r, buf)
	if err != nil && n == 0 {
		return false
	}
	return isRenderableJPEG(buf[:n])
}

// isRenderableJPEG walks the JPEG segments until the SOF marker and accepts
// only the DCT processes common decoders implement. Raw sensor payloads in
// DNGs are lossless JPEG (SOF3) and start with the same SOI marker, so the
// SOI alone does not qualify a stream.
func isRenderableJPEG(buf []byte) bool {
	if len(buf) < 4 || buf[0] != 0xff || buf[1] != 0xd8 {
		return false
	}
	i := 2
	for segments := 0; i+4 <= len(buf) && segments < maxJPEGSegments; segments++ {
		if buf[i] != 0xff {
			return false
		}
		marker := buf[i+1]
		switch {
		case marker == 0xff: // fill byte
			i++
			continue
		case marker >= 0xd0 && marker <= 0xd7: // RST, no length field
			i += 2
			continue
		}
		switch marker {
		case 0xc0, 0xc1, 0xc2: // baseline, extended sequential, progressive
			return true
		case 0xc3, 0xc5, 0xc6, 0xc7, 0xc9, 0xca, 0xcb, 0xcd, 0xce, 0xcf:
			// lossless, differential and arithmetic processes
			return false
		case 0xd9, 0xda: // EOI or scan start without a SOF
			return false
		}
		segLen := int(buf[i+2])<<8 | int(buf[i+3])
		if segLen < 2 {
			return false
		}
		i += 2 + segLen
	}
	return false
}
