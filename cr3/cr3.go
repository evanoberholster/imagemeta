// Package cr3 decodes (CR3) Canon Raw 3 Metadata using the bmff package
//
// Based on: Laurent Clévy's work on Canon CR3 file structure found at (@Lorenzo2472) (https://github.com/lclevy/canon_cr3)
package cr3

import (
	"bufio"
	"fmt"
	"io"

	"github.com/evanoberholster/imagemeta/bmff"
	"github.com/evanoberholster/imagemeta/exif"
	"github.com/evanoberholster/imagemeta/exif/ifds"
	"github.com/evanoberholster/imagemeta/imagetype"
	"github.com/evanoberholster/imagemeta/meta"
	"github.com/evanoberholster/imagemeta/xmp"
)

// Metadata is an Heic file's Metadata
type Metadata struct {
	// Reader
	r io.Reader
	// Custom Decoder interface
	ExifDecodeFn exif.DecodeFn
	XmpDecodeFn  xmp.DecodeFn
	XmpHeader    xmp.Header
	FileType     bmff.FileTypeBox
	CrxMoov      bmff.CrxMoovBox
	//Thumbnail []byte
	dim meta.Dimensions
	t   imagetype.ImageType
}

// Dimensions returns the meta.Dimensions of the
// Primary Image.
func (hm Metadata) Dimensions() meta.Dimensions {
	return hm.dim
}

// NewMetadata returns a new heic.Metadata
func NewMetadata(r io.Reader) (m *Metadata, err error) {
	m = &Metadata{r: r, t: imagetype.ImageCR3}
	err = m.getMeta()
	return m, err
}

func (hm *Metadata) getMeta() (err error) {
	bmr := bmff.NewReader(hm.r)

	hm.FileType, err = bmr.ReadFtypBox()
	if err != nil {
		return
	}

	hm.CrxMoov, err = bmr.ReadCrxMoovBox()
	if err != nil {
		return
	}

	//hm.Meta, err = bmr.ReadMetaBox()
	//if err != nil {
	//	return
	//}

	// Set Dimensions
	// TODO: fetch PitM, # images, and dimensions
	hm.dim = meta.NewDimensions(0, 0)

	return nil
}

func (hm *Metadata) DecodeExif(r meta.Reader) (err error) {
	// Don't process Exif if no Decode function is given
	if hm.ExifDecodeFn == nil {
		return
	}
	for _, cmt := range hm.CrxMoov.Meta.CMT {
		header := exif.NewHeader(cmt.ByteOrder, cmt.FirstIfdOffset, cmt.TiffHeaderOffset, cmt.ExifLength, cmt.ImageType)
		switch cmt.Bt {
		case bmff.TypeCMT1:
			header.FirstIfd = ifds.RootIFD
		case bmff.TypeCMT2:
			header.FirstIfd = ifds.ExifIFD
		case bmff.TypeCMT3:
			header.FirstIfd = ifds.MknoteIFD
		case bmff.TypeCMT4:
			header.FirstIfd = ifds.GPSIFD
		}
		if err = hm.ExifDecodeFn(r, header); err != nil {
			return err
		}
	}
	return
}

func (hm *Metadata) DecodeXMP(r meta.Reader) (err error) {
	// Don't process XMP if no Decode function is given
	if hm.XmpDecodeFn == nil {
		return
	}
	offset, length, _ := hm.CrxMoov.Meta.XPacketData()
	if _, err = r.Seek(int64(offset), 0); err != nil {
		return err
	}
	br := bufio.NewReaderSize(io.LimitReader(r, int64(length)), 1024*3/2)
	buf, err := br.Peek(24)
	if err != nil {
		return
	}
	var uuid meta.UUID
	_ = uuid.UnmarshalBinary(buf[8:24])
	if uuid != bmff.CR3XPacketUUID {
		fmt.Println("Wrong UUID")
	}
	return hm.XmpDecodeFn(br, xmp.NewHeader(24, uint32(length)))
}
