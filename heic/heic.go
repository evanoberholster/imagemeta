// Package heic decodes Heic Metadata using the bmff package
package heic

import (
	"fmt"
	"io"

	"github.com/evanoberholster/imagemeta/bmff"
	"github.com/evanoberholster/imagemeta/exif/ifds"
	"github.com/evanoberholster/imagemeta/meta"
	"github.com/pkg/errors"
)

// TODO:
// Write benchmarks
// Write tests

// Common Heif Errors
var (
	ErrItemNotFound = errors.New("error item not found")
)

// Metadata is an Heic file's Metadata
type Metadata struct {
	*meta.Metadata

	FileType bmff.FileTypeBox
	Meta     bmff.MetaBox
	//Thumbnail []byte
	n uint16 // Num Images
}

// NewMetadata returns a new heic.Metadata
func NewMetadata(r io.Reader, m *meta.Metadata) (hm Metadata, err error) {
	hm = Metadata{Metadata: m}
	bmr := bmff.NewReader(r)

	hm.FileType, err = bmr.ReadFtypBox()
	if err != nil {
		return
	}

	hm.Meta, err = bmr.ReadMetaBox()
	if err != nil {
		return
	}

	// Find PITM and set Dimensions
	box, err := hm.Meta.Properties.ContainerByID(hm.Meta.Primary.ItemID, bmff.TypeIspe)
	if err != nil {
		err = errors.Wrap(err, "Heic getMeta")
	}
	if ispe, ok := box.(bmff.ImageSpatialExtentsProperty); ok {
		hm.Dim = meta.NewDimensions(ispe.W, ispe.H)
	}

	return hm, err
}

// Images returns the number of images.
func (hm Metadata) Images() uint16 {
	return hm.n
}

// Item represents an item in a HEIF file.
type Item struct {
	ID         uint16
	Info       bmff.ItemInfoEntry
	Location   bmff.ItemLocationBoxEntry
	Properties bmff.ItemPropertyAssociationItem
}

func (hm *Metadata) itemByType(it bmff.ItemType) (item Item, err error) {
	item.Info, err = hm.Meta.ItemInfo.LastItemByType(it)
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

// ReadXmpHeader reads Xmp Header from the Heic Metadata and returns an Xmp Header.
// If an error occurs returns the error.
//
func (hm *Metadata) ReadXmpHeader(r meta.Reader) (header meta.XmpHeader, err error) {
	item, err := hm.itemByType(bmff.ItemTypeMime)
	if err != nil {
		return
	}
	hm.XmpHeader = meta.NewXMPHeader(uint32(item.Location.FirstExtent.Offset), uint32(item.Location.FirstExtent.Length))
	return hm.XmpHeader, err
}

// ReadExifHeader reads Exif Header from the Heic Metadata and returns an Exif Header.
// If an error occurs returns the error.
//
func (hm *Metadata) ReadExifHeader(r meta.Reader) (header meta.ExifHeader, err error) {
	item, err := hm.itemByType(bmff.ItemTypeExif)
	if err != nil {
		return
	}
	offset := item.Location.FirstExtent.Offset
	length := item.Location.FirstExtent.Length

	var buf [18]byte
	if _, err = r.ReadAt(buf[:], int64(offset)); err != nil {
		return
	}

	// Read BoxType
	if !(bmff.IsBoxType("Exif", buf[4:8])) {
		err = fmt.Errorf("error wrong Box Type: %d", buf[4:8])
		return
	}

	// Read Tiff header
	if _, err = r.Read(buf[:8]); err != nil {
		return
	}
	byteOrder := meta.BinaryOrder(buf[10:14])
	firstIfdOffset := byteOrder.Uint32(buf[14:18])
	tiffHeaderOffset := int64(offset) + 10
	hm.ExifHeader = meta.NewExifHeader(byteOrder, firstIfdOffset, uint32(tiffHeaderOffset), uint32(length), hm.It)
	hm.ExifHeader.FirstIfd = ifds.IFD0

	return hm.ExifHeader, err
}

// ReadXmp reads XMP metadata from the meta.Reader
func (hm *Metadata) ReadXmp(r meta.Reader) (err error) {
	if hm.XmpFn == nil {
		return
	}
	if _, err = hm.ReadXmpHeader(r); err != nil {
		return
	}
	if _, err = r.Seek(int64(hm.XmpHeader.Offset), 0); err != nil {
		return
	}
	return hm.XmpFn(io.LimitReader(r, int64(hm.XmpHeader.Length)), hm.Metadata)
}

// ReadExif reads Exif metadata from the meta.Reader
func (hm *Metadata) ReadExif(r meta.Reader) (err error) {
	if hm.ExifFn == nil {
		return
	}
	if _, err = hm.ReadExifHeader(r); err != nil {
		return
	}
	return hm.ExifFn(r, hm.Metadata)
}
