package exif

import (
	"github.com/evanoberholster/imagemeta/imagetype"
	"github.com/evanoberholster/imagemeta/meta"
	"github.com/evanoberholster/imagemeta/meta/exif/makernote"
	"github.com/evanoberholster/imagemeta/meta/jpeg"
)

// FromCIFF maps parsed Canon CIFF metadata into the public EXIF model.
func FromCIFF(c *jpeg.CIFF, it imagetype.ImageType) Exif {
	out := Exif{
		CameraMakeID: makernote.CameraMakeCanon,
		ImageType:    it,
	}
	out.MakerNote.Make = makernote.CameraMakeCanon
	out.IFD0.Make = c.Make
	out.IFD0.Model = c.Model
	out.IFD0.ImageWidth = c.ImageWidth
	out.IFD0.ImageHeight = c.ImageHeight
	out.IFD0.ImageDescription = c.CanonFileDescription
	out.IFD0.ModifyDate = c.DateTimeOriginal
	out.Time.DateTimeOriginal = c.DateTimeOriginal
	out.ExifIFD.PixelXDimension = c.ImageWidth
	out.ExifIFD.PixelYDimension = c.ImageHeight
	out.ExifIFD.ISOSpeedRatings = uint32(c.BaseISO)
	out.ExifIFD.CameraOwnerName = c.OwnerName
	out.IFD0.Artist = c.OwnerName
	out.ExifIFD.FocalLength = meta.FocalLength(c.FocalLength)
	out.ExifIFD.FNumber = meta.Aperture(c.ApertureValue)
	out.ExifIFD.ApertureValue = meta.Aperture(c.ApertureValue)
	out.ExifIFD.ShutterSpeedValue = meta.ShutterSpeed(c.ShutterSpeedValue)
	out.ExifIFD.ExposureTime = meta.ExposureTime(c.ShutterSpeedValue)
	return out
}
