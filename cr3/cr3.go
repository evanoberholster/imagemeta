// Package cr3 decodes (CR3) Canon Raw 3 Metadata using the bmff package
//
// Based on: Laurent Clévy's work on Canon CR3 file structure found at (@Lorenzo2472) (https://github.com/lclevy/canon_cr3)
package cr3

import (
	"fmt"
	"io"

	"github.com/evanoberholster/imagemeta/bmff"
	"github.com/evanoberholster/imagemeta/exif"
	"github.com/evanoberholster/imagemeta/imagetype"
	"github.com/evanoberholster/imagemeta/meta"
	"github.com/evanoberholster/imagemeta/xmp"
	"github.com/pkg/errors"
)

// Metadata is an Heic file's Metadata
type Metadata struct {
	mr          meta.Reader
	ExifHeader  meta.ExifHeader
	XmpHeader   meta.XmpHeader
	jpegOffsets [2]uint32

	FileType bmff.FileTypeBox
	CrxMoov  bmff.CrxMoovBox

	e *exif.Data
	// Decode Functions for EXIF and XMP metadata
	exifFn func(r io.Reader, header meta.ExifHeader) error
	xmpFn  func(r io.Reader, header meta.XmpHeader) error
}

// Dimensions returns the dimensions (width and height) of the image
func (m Metadata) Dimensions() meta.Dimensions {
	return m.e.Dimensions()
}

// ImageType returns imagetype.ImageCR3 for Canon CR3 image
func (m Metadata) ImageType() imagetype.ImageType {
	return imagetype.ImageCR3
}

// PreviewImage returns a JPEG preview image
func (m Metadata) PreviewImage() io.Reader {
	return io.NewSectionReader(m.mr, int64(m.CrxMoov.Trak[0].Offset), int64(m.CrxMoov.Trak[0].ImageSize))
}

// Exif returns parsed Exif data from JPEG
func (m Metadata) Exif() (exif.Exif, error) {
	return m.e, nil
}

// Xmp returns parsed Xmp data from JPEG
func (m Metadata) Xmp() (xmp.XMP, error) {
	sr := io.NewSectionReader(m.mr, int64(m.XmpHeader.Offset), int64(m.XmpHeader.Length))
	return xmp.ParseXmp(sr)
}

func Parse(mr meta.Reader) (Metadata, error) {
	m := Metadata{mr: mr}
	err := m.getMeta()
	return m, err
}

func (m *Metadata) getMeta() (err error) {
	bmr := bmff.NewReader(m.mr)

	m.FileType, err = bmr.ReadFtypBox()
	if err != nil {
		return errors.Wrapf(err, "ReadFtypBox")
	}

	m.CrxMoov, err = bmr.ReadCrxMoovBox()
	if err != nil {
		return errors.Wrapf(err, "ReadCrxMoovBox")
	}
	for i, header := range m.CrxMoov.Meta.Exif {
		if i == 0 {
			m.e, err = exif.ParseExif(m.mr, header)
			if err != nil {
				fmt.Println(err)
			}
			continue
		}
		if err = m.e.ParseIfd(header); err != nil {
			return err
		}
	}
	return nil
}

//func (hm *Metadata) DecodeExif(r io.Reader) (err error) {
//	// Don't process Exif if no Decode function is given
//	if hm.ExifFn == nil {
//		return
//	}
//	for _, cmt := range hm.CrxMoov.Meta.CMT {
//		hm.ExifHeader = meta.NewExifHeader(cmt.ByteOrder, cmt.FirstIfdOffset, cmt.TiffHeaderOffset, cmt.ExifLength, cmt.ImageType)
//		switch cmt.Bt {
//		case bmff.TypeCMT1:
//			hm.ExifHeader.FirstIfd = ifds.IFD0
//		case bmff.TypeCMT2:
//			hm.ExifHeader.FirstIfd = ifds.ExifIFD
//		case bmff.TypeCMT3:
//			hm.ExifHeader.FirstIfd = ifds.MknoteIFD
//		case bmff.TypeCMT4:
//			hm.ExifHeader.FirstIfd = ifds.GPSIFD
//		}
//		if err = hm.ExifFn(r, hm.Metadata); err != nil {
//			return err
//		}
//	}
//	return
//}

// DecodeXMP decodes XMP from the underlying CR3 Image.
//func (hm *Metadata) DecodeXMP(r io.Reader) (err error) {
//	// Don't process XMP if no Decode function is given
//	if hm.XmpFn == nil {
//		return
//	}
//	offset, length, _ := hm.CrxMoov.Meta.XPacketData()
//	//if _, err = r.Seek(int64(offset), 0); err != nil {
//	//	return err
//	//}
//	//br := bufio.NewReaderSize(io.LimitReader(r, int64(length)), 1024*3/2)
//	//buf, err := br.Peek(24)
//	//if err != nil {
//	//	return
//	//}
//	//var uuid meta.UUID
//	//_ = uuid.UnmarshalBinary(buf[8:24])
//	//if uuid != bmff.CR3XPacketUUID {
//	//	fmt.Println("Wrong UUID")
//	//}
//	hm.XmpHeader = meta.NewXMPHeader(uint32(offset), uint32(length))
//	return hm.XmpFn(r, hm.Metadata)
//}
