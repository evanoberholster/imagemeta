package exif

import (
	"os"
	"testing"

	"github.com/evanoberholster/imagemeta/meta"
	"github.com/evanoberholster/imagemeta/meta/exif/tag"
	metalog "github.com/evanoberholster/imagemeta/meta/logging"
	"github.com/evanoberholster/imagemeta/meta/utils"
)

var (
	benchExifHandled bool
	benchExifCount   int
)

var benchExifIFDEntries = [...]tag.Entry{
	// ExifIFD text/time-friendly compact tag.
	tag.NewEntry(tag.TagExifVersion, tag.TypeUndefined, 4, packLE4('0', '2', '3', '1'), tag.ExifIFD, 0, utils.LittleEndian),

	// ExifIFD image tags.
	tag.NewEntry(tag.TagPixelXDimension, tag.TypeLong, 1, 6000, tag.ExifIFD, 0, utils.LittleEndian),
	tag.NewEntry(tag.TagPixelYDimension, tag.TypeLong, 1, 4000, tag.ExifIFD, 0, utils.LittleEndian),
	tag.NewEntry(tag.TagInteropIFDPointer, tag.TypeLong, 1, 0x3200, tag.ExifIFD, 0, utils.LittleEndian),
	tag.NewEntry(tag.TagColorSpace, tag.TypeShort, 1, uint32(1), tag.ExifIFD, 0, utils.LittleEndian),
	tag.NewEntry(tag.TagComponentsConfiguration, tag.TypeUndefined, 4, packLE4(1, 2, 3, 0), tag.ExifIFD, 0, utils.LittleEndian),
	tag.NewEntry(tag.TagFocalPlaneResolutionUnit, tag.TypeShort, 1, uint32(2), tag.ExifIFD, 0, utils.LittleEndian),
	tag.NewEntry(tag.TagSubjectArea, tag.TypeShort, 2, packLE4(120, 0, 80, 0), tag.ExifIFD, 0, utils.LittleEndian),

	// ExifIFD exposure tags.
	tag.NewEntry(tag.TagExposureProgram, tag.TypeShort, 1, uint32(3), tag.ExifIFD, 0, utils.LittleEndian),
	tag.NewEntry(tag.TagSensitivityType, tag.TypeShort, 1, uint32(2), tag.ExifIFD, 0, utils.LittleEndian),
	tag.NewEntry(tag.TagRecommendedExposureIndex, tag.TypeLong, 1, 200, tag.ExifIFD, 0, utils.LittleEndian),
	tag.NewEntry(tag.TagExposureMode, tag.TypeShort, 1, uint32(1), tag.ExifIFD, 0, utils.LittleEndian),
	tag.NewEntry(tag.TagMeteringMode, tag.TypeShort, 1, uint32(5), tag.ExifIFD, 0, utils.LittleEndian),
	tag.NewEntry(tag.TagLightSource, tag.TypeShort, 1, uint32(1), tag.ExifIFD, 0, utils.LittleEndian),
	tag.NewEntry(tag.TagISOSpeedRatings, tag.TypeLong, 1, 640, tag.ExifIFD, 0, utils.LittleEndian),
	tag.NewEntry(tag.TagFlash, tag.TypeShort, 1, uint32(0), tag.ExifIFD, 0, utils.LittleEndian),

	// ExifIFD capture tags.
	tag.NewEntry(tag.TagSensingMethod, tag.TypeShort, 1, uint32(2), tag.ExifIFD, 0, utils.LittleEndian),
	tag.NewEntry(tag.TagFileSource, tag.TypeByte, 1, uint32(3), tag.ExifIFD, 0, utils.LittleEndian),
	tag.NewEntry(tag.TagSceneType, tag.TypeByte, 1, uint32(1), tag.ExifIFD, 0, utils.LittleEndian),
	tag.NewEntry(tag.TagCustomRendered, tag.TypeShort, 1, uint32(0), tag.ExifIFD, 0, utils.LittleEndian),
	tag.NewEntry(tag.TagWhiteBalance, tag.TypeShort, 1, uint32(0), tag.ExifIFD, 0, utils.LittleEndian),
	tag.NewEntry(tag.TagSceneCaptureType, tag.TypeShort, 1, uint32(0), tag.ExifIFD, 0, utils.LittleEndian),
	tag.NewEntry(tag.TagGainControl, tag.TypeShort, 1, uint32(0), tag.ExifIFD, 0, utils.LittleEndian),
	tag.NewEntry(tag.TagContrast, tag.TypeShort, 1, uint32(0), tag.ExifIFD, 0, utils.LittleEndian),
	tag.NewEntry(tag.TagSaturation, tag.TypeShort, 1, uint32(0), tag.ExifIFD, 0, utils.LittleEndian),
	tag.NewEntry(tag.TagSharpness, tag.TypeShort, 1, uint32(0), tag.ExifIFD, 0, utils.LittleEndian),
	tag.NewEntry(tag.TagSubjectDistanceRange, tag.TypeShort, 1, uint32(1), tag.ExifIFD, 0, utils.LittleEndian),
	tag.NewEntry(tag.TagCompositeImage, tag.TypeShort, 1, uint32(1), tag.ExifIFD, 0, utils.LittleEndian),
}

type exifIFDTagHandler func(r *Reader, t tag.Entry)

var exifIFDTagMapDispatch = map[tag.ID]exifIFDTagHandler{
	tag.TagExifVersion: func(r *Reader, t tag.Entry) {
		r.Exif.ExifIFD.ExifVersion = r.parseExifVersion(t)
	},
	tag.TagPixelXDimension: func(r *Reader, t tag.Entry) {
		r.Exif.ExifIFD.PixelXDimension = r.parseUint32(t)
		if r.Exif.IFD0.ImageWidth == 0 {
			r.Exif.IFD0.ImageWidth = r.Exif.ExifIFD.PixelXDimension
		}
	},
	tag.TagPixelYDimension: func(r *Reader, t tag.Entry) {
		r.Exif.ExifIFD.PixelYDimension = r.parseUint32(t)
		if r.Exif.IFD0.ImageHeight == 0 {
			r.Exif.IFD0.ImageHeight = r.Exif.ExifIFD.PixelYDimension
		}
	},
	tag.TagInteropIFDPointer: func(r *Reader, t tag.Entry) {
		r.Exif.ExifIFD.interopIFDPointer = r.parseUint32(t)
	},
	tag.TagColorSpace: func(r *Reader, t tag.Entry) {
		r.Exif.ExifIFD.ColorSpace = r.parseUint16(t)
	},
	tag.TagComponentsConfiguration: func(r *Reader, t tag.Entry) {
		r.parseByteList(t, r.Exif.ExifIFD.ComponentsConfiguration[:])
	},
	tag.TagFocalPlaneResolutionUnit: func(r *Reader, t tag.Entry) {
		r.Exif.ExifIFD.FocalPlaneResolutionUnit = meta.ResolutionUnit(r.parseUint16(t))
	},
	tag.TagSubjectArea: func(r *Reader, t tag.Entry) {
		r.parseUint16List(t, r.Exif.ExifIFD.SubjectArea[:])
	},
	tag.TagExposureProgram: func(r *Reader, t tag.Entry) {
		r.Exif.ExifIFD.ExposureProgram = meta.ExposureProgram(r.parseUint16(t))
	},
	tag.TagSensitivityType: func(r *Reader, t tag.Entry) {
		r.Exif.ExifIFD.SensitivityType = r.parseUint16(t)
	},
	tag.TagRecommendedExposureIndex: func(r *Reader, t tag.Entry) {
		r.Exif.ExifIFD.RecommendedExposureIndex = r.parseUint32(t)
	},
	tag.TagExposureMode: func(r *Reader, t tag.Entry) {
		r.Exif.ExifIFD.ExposureMode = meta.ExposureMode(r.parseUint16(t))
	},
	tag.TagMeteringMode: func(r *Reader, t tag.Entry) {
		r.Exif.ExifIFD.MeteringMode = meta.MeteringMode(r.parseUint16(t))
	},
	tag.TagLightSource: func(r *Reader, t tag.Entry) {
		r.Exif.ExifIFD.LightSource = r.parseUint16(t)
	},
	tag.TagISOSpeedRatings: func(r *Reader, t tag.Entry) {
		r.Exif.ExifIFD.ISOSpeedRatings = r.parseUint32(t)
	},
	tag.TagFlash: func(r *Reader, t tag.Entry) {
		r.Exif.ExifIFD.Flash = meta.Flash(r.parseUint16(t))
	},
	tag.TagSensingMethod: func(r *Reader, t tag.Entry) {
		r.Exif.ExifIFD.SensingMethod = r.parseUint16(t)
	},
	tag.TagFileSource: func(r *Reader, t tag.Entry) {
		r.Exif.ExifIFD.FileSource = r.parseSceneType(t)
	},
	tag.TagSceneType: func(r *Reader, t tag.Entry) {
		r.Exif.ExifIFD.SceneType = r.parseSceneType(t)
	},
	tag.TagCustomRendered: func(r *Reader, t tag.Entry) {
		r.Exif.ExifIFD.CustomRendered = r.parseUint16(t)
	},
	tag.TagWhiteBalance: func(r *Reader, t tag.Entry) {
		r.Exif.ExifIFD.WhiteBalance = r.parseUint16(t)
	},
	tag.TagSceneCaptureType: func(r *Reader, t tag.Entry) {
		r.Exif.ExifIFD.SceneCaptureType = r.parseUint16(t)
	},
	tag.TagGainControl: func(r *Reader, t tag.Entry) {
		r.Exif.ExifIFD.GainControl = r.parseUint16(t)
	},
	tag.TagContrast: func(r *Reader, t tag.Entry) {
		r.Exif.ExifIFD.Contrast = r.parseUint16(t)
	},
	tag.TagSaturation: func(r *Reader, t tag.Entry) {
		r.Exif.ExifIFD.Saturation = r.parseUint16(t)
	},
	tag.TagSharpness: func(r *Reader, t tag.Entry) {
		r.Exif.ExifIFD.Sharpness = r.parseUint16(t)
	},
	tag.TagSubjectDistanceRange: func(r *Reader, t tag.Entry) {
		r.Exif.ExifIFD.SubjectDistanceRange = r.parseUint16(t)
	},
	tag.TagCompositeImage: func(r *Reader, t tag.Entry) {
		r.Exif.ExifIFD.CompositeImage = r.parseUint16(t)
	},
}

func parseExifTagMapOnlyBenchmark(r *Reader, t tag.Entry) bool {
	h, ok := exifIFDTagMapDispatch[t.ID]
	if !ok {
		return false
	}
	h(r, t)
	return true
}

func BenchmarkParseExifIFDDispatch(b *testing.B) {
	b.Run("switch", func(b *testing.B) {
		r := NewReader(metalog.GetLogger())
		defer r.Close()
		b.ReportAllocs()

		count := 0
		for i := 0; i < b.N; i++ {
			for _, t := range benchExifIFDEntries {
				if r.parseExifTag(t) {
					count++
				}
			}
		}
		benchExifCount = count
		benchExifHandled = count > 0
	})

	b.Run("map", func(b *testing.B) {
		r := NewReader(metalog.GetLogger())
		defer r.Close()
		b.ReportAllocs()

		count := 0
		for i := 0; i < b.N; i++ {
			for _, t := range benchExifIFDEntries {
				if parseExifTagMapOnlyBenchmark(r, t) {
					count++
				}
			}
		}
		benchExifCount = count
		benchExifHandled = count > 0
	})
}

func BenchmarkParseExifIFDDispatchCompare(b *testing.B) {
	mode := os.Getenv("IMAGEMETA_EXIF_DISPATCH_MODE")
	parseFn := (*Reader).parseExifTag
	if mode == "map" {
		parseFn = func(r *Reader, t tag.Entry) bool {
			return parseExifTagMapOnlyBenchmark(r, t)
		}
	}

	r := NewReader(metalog.GetLogger())
	defer r.Close()
	b.ReportAllocs()

	count := 0
	for i := 0; i < b.N; i++ {
		for _, t := range benchExifIFDEntries {
			if parseFn(r, t) {
				count++
			}
		}
	}
	benchExifCount = count
	benchExifHandled = count > 0
}

func packLE4(a, b, c, d byte) uint32 {
	return utils.LittleEndian.Uint32([]byte{a, b, c, d})
}
