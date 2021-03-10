// Package heic decodes Heic Metadata using the bmff package
package heic

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	"github.com/evanoberholster/imagemeta/bmff"
	"github.com/evanoberholster/imagemeta/exif/ifds"
	"github.com/evanoberholster/imagemeta/imagetype"
	"github.com/evanoberholster/imagemeta/meta"
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
	*meta.Metadata
	// Reader
	r        io.Reader
	FileType bmff.FileTypeBox
	Meta     bmff.MetaBox
	//Thumbnail []byte
	n uint16 // Num Images
}

// NewMetadata returns a new heic.Metadata
func NewMetadata(r io.Reader, m *meta.Metadata) (hm *Metadata, err error) {
	hm = &Metadata{r: r, Metadata: m}
	err = hm.getMeta()
	return hm, err
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
	box, err := hm.Meta.Properties.ContainerByID(hm.Meta.Primary.ItemID, bmff.TypeIspe)
	if err != nil {
		return
	}
	if ispe, ok := box.(bmff.ImageSpatialExtentsProperty); ok {
		hm.Dim = meta.NewDimensions(ispe.W, ispe.H)
	}
	hm.n = 1
	return nil
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
	hm.ExifHeader, err = readExifBox(r, hm.It, item.Location.FirstExtent.Offset, item.Location.FirstExtent.Length)
	hm.ExifHeader.ImageType = hm.It
	return hm.ExifHeader, err
}

func readExifBox(r meta.Reader, imageType imagetype.ImageType, offset uint64, length uint64) (header meta.ExifHeader, err error) {
	// Seek to Exif Box position
	if _, err = r.Seek(int64(offset), 0); err != nil {
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

	//var buf [16]byte
	//r.ReadAt(buf, int64(offset)+)
	//
	// Read Tiff header
	if _, err = r.Read(buf[:8]); err != nil {
		return
	}
	byteOrder := meta.BinaryOrder(buf[:4])
	firstIfdOffset := byteOrder.Uint32(buf[4:8])
	tiffHeaderOffset := int64(offset) + size + 4
	header = meta.NewExifHeader(byteOrder, firstIfdOffset, uint32(tiffHeaderOffset), uint32(length), imageType)
	header.FirstIfd = ifds.RootIFD
	return header, nil
}

func readBox(r io.Reader, buf [8]byte, boxType string) (size int64, err error) {
	// Read size
	if _, err = r.Read(buf[:]); err != nil {
		return
	}
	size = int64(binary.BigEndian.Uint32(buf[:4]))

	// Read BoxType
	if !(bmff.IsBoxType(boxType, buf[4:8])) {
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
