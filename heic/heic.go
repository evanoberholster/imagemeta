// Package heic decodes Heic Metadata using the bmff package
package heic

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	"github.com/evanoberholster/imagemeta/bmff"
	"github.com/evanoberholster/imagemeta/exif"
	"github.com/evanoberholster/imagemeta/imagetype"
	"github.com/evanoberholster/imagemeta/meta"
	"github.com/evanoberholster/imagemeta/xmp"
)

// TODO:
// Add support for XMP
// Write benchmarks
// Write tests

// Common Heif Errors
var (
	ErrItemNotFound = errors.New("error item not found")
)

// Metadata is an Heic file's Metadata
type Metadata struct {
	// Reader
	r io.Reader
	// Custom Decoder interface
	ExifDecodeFn exif.DecodeFn
	ExifHeader   exif.Header
	XmpDecodeFn  xmp.DecodeFn
	XmpHeader    xmp.Header
	FileType     bmff.FileTypeBox
	Meta         bmff.MetaBox
	//Thumbnail []byte
	dim meta.Dimensions
	n   uint16 // Num Images
	t   imagetype.ImageType
}

// NewMetadata returns a new heic.Metadata
func NewMetadata(r io.Reader) (m *Metadata, err error) {
	m = &Metadata{r: r, t: imagetype.ImageHEIF}
	err = m.getMeta()
	return m, err
}

// Dimensions returns the meta.Dimensions of the
// Primary Image.
func (hm Metadata) Dimensions() meta.Dimensions {
	return hm.dim
}

// Images returns the number of images.
func (hm Metadata) Images() uint16 {
	return hm.n
}

func (hm *Metadata) getMeta() (err error) {
	bmr := bmff.NewReader(hm.r)

	hm.FileType, err = bmr.ReadFtypBox()
	if err != nil {
		return
	}

	hm.Meta, err = bmr.ReadMetaBox()
	if err != nil {
		return
	}

	// Find PITM and set Dimensions
	// TODO: fetch PitM, # images, and dimensions
	hm.dim = meta.NewDimensions(0, 0)
	hm.n = 1

	//fmt.Println(mb)
	return nil
}

// Item represents an item in a HEIF file.
type Item struct {
	ID         uint16
	Info       bmff.ItemInfoEntry
	Location   bmff.ItemLocationBoxEntry
	Properties bmff.ItemPropertyAssociationItem
}

func (hm *Metadata) ExifItem() (item Item, err error) {
	item.Info, err = hm.Meta.ItemInfo.LastItemByType(bmff.ItemTypeExif)
	if err == bmff.ErrItemNotFound {
		err = meta.ErrNoExif
		return
	}
	item.ID = item.Info.ItemID
	item.Location, err = hm.Meta.Location.EntryByID(item.Info.ItemID)
	if err == bmff.ErrItemNotFound {
		err = meta.ErrNoExif
		return
	}
	return
}

// DecodeExif reads Exif Metadata from the underlying reader interface and returns Exif.
// If an error occurs returns the error.
//
// Utilizes the custom decoder ExifDecodeFn if it is not nil.
func (hm *Metadata) DecodeExif(r meta.Reader) error {
	item, err := hm.ExifItem()
	if err != nil {
		return err
	}
	header, err := readExifBox(r, hm.t, item.Location.FirstExtent.Offset, item.Location.FirstExtent.Length)
	if err != nil {
		return err
	}
	if hm.ExifDecodeFn == nil {
		return nil
	}
	return hm.ExifDecodeFn(r, header)
}

func readExifBox(r meta.Reader, imageType imagetype.ImageType, offset uint64, length uint64) (header exif.Header, err error) {
	// Seek to Exif Box position
	_, err = r.Seek(int64(offset), 0)
	if err != nil {
		return
	}

	// Read Exif Box
	var buf [8]byte
	size, err := readBox(r, buf, "Exif")
	if err != nil {
		return
	}
	_, err = r.Seek(int64(size-4), io.SeekCurrent)
	if err != nil {
		return
	}

	// Read Tiff header
	if _, err = r.Read(buf[:8]); err != nil {
		return
	}
	byteOrder := meta.BinaryOrder(buf[:4])
	firstIfdOffset := byteOrder.Uint32(buf[4:8])
	tiffHeaderOffset := int64(offset) + size + 4
	return exif.NewHeader(byteOrder, firstIfdOffset, uint32(tiffHeaderOffset), uint32(length), imageType), nil
}

func readBox(r io.Reader, buf [8]byte, boxType string) (size int64, err error) {
	// Read size
	if _, err = r.Read(buf[:]); err != nil {
		return
	}
	size = int64(binary.BigEndian.Uint32(buf[:4]))

	// Read BoxType
	if !(isBoxType(boxType, buf[4:8])) {
		err = fmt.Errorf("error wrong Box Type: %d", buf)
		return
	}

	switch size {
	case 1:
		// Read 64 bit size, after the type
		if _, err = r.Read(buf[:]); err != nil {
			return
		}
		size = int64(binary.BigEndian.Uint64(buf[:8]))
		if size < 0 {
			return //Error
		}
	case 0:
		return // Error
	}
	return size, nil
}

func isBoxType(bt1 string, bt2 []byte) bool {
	return bt1[0] == bt2[0] && bt1[1] == bt2[1] && bt1[2] == bt2[2] && bt1[3] == bt2[3]
}
