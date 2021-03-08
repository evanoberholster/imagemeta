// Package cr3 decodes (CR3) Canon Raw 3 Metadata using the bmff package
//
// Based on: Laurent Clévy's work on Canon CR3 file structure found at (@Lorenzo2472) (https://github.com/lclevy/canon_cr3)
package cr3

import (
	"io"

	"github.com/evanoberholster/imagemeta/bmff"
	"github.com/evanoberholster/imagemeta/exif/ifds"
	"github.com/evanoberholster/imagemeta/imagetype"
	"github.com/evanoberholster/imagemeta/meta"
)

// Metadata is an Heic file's Metadata
type Metadata struct {
	meta.Metadata
	// Reader
	r io.Reader

	FileType bmff.FileTypeBox
	CrxMoov  bmff.CrxMoovBox
	//Thumbnail []byte
}

// NewMetadata returns a new heic.Metadata
func NewMetadata(r io.Reader, m meta.Metadata) (cm *Metadata, err error) {
	cm = &Metadata{r: r, Metadata: m}
	m.It = imagetype.ImageCR3
	err = cm.getMeta()
	return cm, err
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

	// Find PITM and set Dimensions
	//if err != nil {
	//	b, err := hm.CrxMoov.Meta.Properties.ContainerByID(hm.Meta.Primary.ItemID, bmff.TypeIspe)
	//	return
	//}
	//if ispe, ok := b.(bmff.ImageSpatialExtentsProperty); ok {
	//	hm.Dim = meta.NewDimensions(ispe.W, ispe.H)
	//}

	//hm.Meta, err = bmr.ReadMetaBox()
	//if err != nil {
	//	return
	//}

	// Set Dimensions
	// TODO: fetch PitM, # images, and dimensions
	hm.Dim = meta.NewDimensions(0, 0)
	return nil
}

func (hm *Metadata) DecodeExif(r io.Reader) (err error) {
	// Don't process Exif if no Decode function is given
	if hm.ExifDecodeFn == nil {
		return
	}
	for _, cmt := range hm.CrxMoov.Meta.CMT {
		hm.ExifHeader = meta.NewExifHeader(cmt.ByteOrder, cmt.FirstIfdOffset, cmt.TiffHeaderOffset, cmt.ExifLength, cmt.ImageType)
		switch cmt.Bt {
		case bmff.TypeCMT1:
			hm.ExifHeader.FirstIfd = ifds.RootIFD
		case bmff.TypeCMT2:
			hm.ExifHeader.FirstIfd = ifds.ExifIFD
		case bmff.TypeCMT3:
			hm.ExifHeader.FirstIfd = ifds.MknoteIFD
		case bmff.TypeCMT4:
			hm.ExifHeader.FirstIfd = ifds.GPSIFD
		}
		if err = hm.ExifDecodeFn(r, hm.ExifHeader); err != nil {
			return err
		}
	}
	return
}

// DecodeXMP decodes XMP from the underlying CR3 Image.
func (hm *Metadata) DecodeXMP(r io.Reader) (err error) {
	// Don't process XMP if no Decode function is given
	if hm.XmpDecodeFn == nil {
		return
	}
	offset, length, _ := hm.CrxMoov.Meta.XPacketData()
	//if _, err = r.Seek(int64(offset), 0); err != nil {
	//	return err
	//}
	//br := bufio.NewReaderSize(io.LimitReader(r, int64(length)), 1024*3/2)
	//buf, err := br.Peek(24)
	//if err != nil {
	//	return
	//}
	//var uuid meta.UUID
	//_ = uuid.UnmarshalBinary(buf[8:24])
	//if uuid != bmff.CR3XPacketUUID {
	//	fmt.Println("Wrong UUID")
	//}
	return hm.XmpDecodeFn(r, meta.NewXMPHeader(uint32(offset), uint32(length)))
}
